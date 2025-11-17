<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
} from "@/components/ui/sheet";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    useVueTable,
    getCoreRowModel,
    getPaginationRowModel,
    type ColumnDef,
} from "@tanstack/vue-table";
import { nodeApi } from "@/api/node";

interface Node {
    id: number;
    is_local: boolean;
    name: string;
    display_name?: string | null;
    address: string;
    ssh_port: number;
    agent_port: number;
    ssh_user: string;
    ssh_password: string;
    status?: string;
}

interface CreateNodeRequest {
    name: string;
    address: string;
    ssh_port: number;
    agent_port: number;
    ssh_user: string;
    ssh_password: string;
    is_local: boolean;
    display_name?: string;
}

// 数据
const nodes = ref<Node[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

// 对话框状态
const isSheetOpen = ref(false);
const isEditMode = ref(false);
const editingNode = ref<Node | null>(null);

// 删除确认对话框状态
const isDeleteDialogOpen = ref(false);
const nodeToDelete = ref<Node | null>(null);

// 表单数据
const formData = ref<CreateNodeRequest>({
    name: "",
    address: "",
    ssh_port: 22,
    agent_port: 8080,
    ssh_user: "root",
    ssh_password: "",
    is_local: false,
    display_name: "",
});

// 获取节点列表
const fetchNodes = async () => {
    loading.value = true;
    error.value = null;
    const page = table.getState().pagination.pageIndex + 1;
    const pageSize = table.getState().pagination.pageSize;
    const response = await nodeApi.page(page, pageSize);
    nodes.value = response.data?.nodes.map((node) => ({
        id: node.id,
        name: node.name,
        address: node.address,
        is_local: node.is_local,
        display_name: node.display_name || null,
        ssh_port: 0,
        agent_port: 0,
        ssh_user: "",
        ssh_password: "",
        status: node.status || "offline",
    })) ?? [];
    error.value = response.error ?? null;
    loading.value = false;
};

// 打开添加对话框
const openAddDialog = () => {
    isEditMode.value = false;
    editingNode.value = null;
    formData.value = {
        name: "",
        address: "",
        ssh_port: 22,
        agent_port: 8080,
        ssh_user: "root",
        ssh_password: "",
        is_local: false,
        display_name: "",
    };
    isSheetOpen.value = true;
};

// 打Separator开编辑对话框
const openEditDialog = (node: Node) => {
    isEditMode.value = true;
    editingNode.value = node;
    formData.value = {
        name: node.name,
        address: node.address,
        ssh_port: node.ssh_port,
        agent_port: node.agent_port,
        ssh_user: node.ssh_user,
        ssh_password: node.ssh_password,
        is_local: node.is_local,
        display_name: node.display_name || "",
    };
    isSheetOpen.value = true;
};

// 保存节点（添加或更新）
const saveNode = async () => {
    if (!formData.value.name || !formData.value.address) {
        error.value = "Name and address are required";
        return;
    }

    loading.value = true;
    error.value = null;

    if (isEditMode.value && editingNode.value) {
        const response = await nodeApi.update({
            id: editingNode.value.id,
            name: formData.value.name,
            address: formData.value.address,
            ssh_port: formData.value.ssh_port,
            agent_port: formData.value.agent_port,
            ssh_user: formData.value.ssh_user,
            ssh_password: formData.value.ssh_password,
            is_local: formData.value.is_local,
            display_name: formData.value.display_name || "",
        });
        if (response.success) {
            isSheetOpen.value = false;
            fetchNodes();
        }
        error.value = response.error ?? null;
        loading.value = false;
    } else {
        const response = await nodeApi.create(formData.value);
        if (response.success) {
            isSheetOpen.value = false;
            fetchNodes();
        }
        error.value = response.error ?? null;
        loading.value = false;
    }

};

// 打开删除确认对话框
const openDeleteDialog = (node: Node) => {
    nodeToDelete.value = node;
    isDeleteDialogOpen.value = true;
};

// 确认删除节点
const confirmDeleteNode = async () => {
    if (!nodeToDelete.value) {
        return;
    }

    loading.value = true;
    error.value = null;

    const response = await nodeApi.delete(nodeToDelete.value.id);
    if (response.success) {
        fetchNodes();
    } else {
        error.value = response.error ?? null;
    }
    loading.value = false;
    isDeleteDialogOpen.value = false;
    nodeToDelete.value = null;
};

const getStatusColor = (status?: string) => {
    const colors: Record<string, string> = {
        online: "text-green-500 bg-green-500/10 border-green-500/20",
        offline: "text-red-500 bg-red-500/10 border-red-500/20",
        maintenance: "text-yellow-500 bg-yellow-500/10 border-yellow-500/20",
    };
    return colors[status || "offline"] || colors.offline;
};

const getStatusDot = (status?: string) => {
    const colors: Record<string, string> = {
        online: "bg-green-500",
        offline: "bg-red-500",
        maintenance: "bg-yellow-500",
    };
    return colors[status || "offline"] || colors.offline;
};

// 定义表格列
const columns: ColumnDef<Node>[] = [
    {
        accessorKey: "id",
        header: "ID",
        cell: (info) => info.getValue(),
    },
    {
        accessorKey: "name",
        header: "Name",
        cell: (info) => {
            const node = info.row.original;
            return node.display_name || node.name;
        },
    },
    {
        accessorKey: "address",
        header: "Address",
    },
    {
        accessorKey: "ssh_port",
        header: "SSH Port",
    },
    {
        accessorKey: "agent_port",
        header: "Agent Port",
    },
    {
        accessorKey: "ssh_user",
        header: "SSH User",
    },
    {
        accessorKey: "is_local",
        header: "Type",
    },
    {
        accessorKey: "status",
        header: "Status",
    },
    {
        id: "actions",
        header: "Actions",
    },
];

// 创建表格实例
const table = useVueTable({
    get data() {
        return nodes.value;
    },
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    initialState: {
        pagination: {
            pageSize: 10,
        },
    },
});

// 初始化
onMounted(() => {
    fetchNodes();
});
</script>

<template>
    <div class="space-y-6">
        <div class="flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-bold">Nodes</h1>
                <p class="text-muted-foreground mt-1">
                    Manage and monitor your server nodes
                </p>
            </div>
            <Button @click="openAddDialog">
                <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
                Add Node
            </Button>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="rounded-lg border border-red-500/20 bg-red-500/10 p-4 text-red-500">
            {{ error }}
        </div>

        <!-- 加载状态 -->
        <div v-if="loading && nodes.length === 0" class="text-center py-8 text-muted-foreground">
            Loading...
        </div>

        <!-- 节点表格 -->
        <div v-else-if="nodes.length > 0" class="rounded-lg border bg-card">
            <Table>
                <TableHeader>
                    <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
                        <TableHead v-for="header in headerGroup.headers" :key="header.id" class="px-6">
                            <div v-if="!header.isPlaceholder">
                                {{
                                    typeof header.column.columnDef.header === "string"
                                        ? header.column.columnDef.header
                                        : header.column.id
                                }}
                            </div>
                        </TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    <TableRow v-for="row in table.getRowModel().rows" :key="row.id">
                        <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id" class="px-6">
                            <template v-if="cell.column.id === 'is_local'">
                                <span :class="[
                                    'inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium',
                                    (cell.getValue() as boolean)
                                        ? 'bg-blue-500/10 text-blue-500 border-blue-500/20'
                                        : 'bg-gray-500/10 text-gray-500 border-gray-500/20',
                                ]">
                                    {{
                                        (cell.getValue() as boolean)
                                            ? "Local"
                                            : "Remote"
                                    }}
                                </span>
                            </template>
                            <template v-else-if="cell.column.id === 'status'">
                                <span :class="[
                                    'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium',
                                    getStatusColor(cell.getValue() as string),
                                ]">
                                    <span :class="[
                                        'h-1.5 w-1.5 rounded-full',
                                        getStatusDot(cell.getValue() as string),
                                    ]"></span>
                                    {{ cell.getValue() }}
                                </span>
                            </template>
                            <template v-else-if="cell.column.id === 'actions'">
                                <div class="flex items-center gap-2">
                                    <Button variant="ghost" size="sm" @click="
                                        openEditDialog(cell.row.original as Node)
                                        " class="h-8 px-2">
                                        <Icon icon="lucide:edit" class="h-4 w-4" />
                                    </Button>
                                    <Button variant="ghost" size="sm" @click="
                                        openDeleteDialog(cell.row.original as Node)
                                        " class="h-8 px-2 text-red-500 hover:text-red-600">
                                        <Icon icon="lucide:trash-2" class="h-4 w-4" />
                                    </Button>
                                </div>
                            </template>
                            <template v-else>
                                <span v-if="cell.column.id === 'name'" class="font-medium">
                                    {{ cell.getValue() }}
                                </span>
                                <span v-else>
                                    {{ cell.getValue() }}
                                </span>
                            </template>
                        </TableCell>
                    </TableRow>
                </TableBody>
            </Table>

            <!-- 分页控件 -->
            <div v-if="table.getPageCount() > 1" class="flex items-center justify-between border-t px-6 py-4">
                <div class="text-sm text-muted-foreground">
                    Showing
                    {{
                        table.getState().pagination.pageIndex *
                        table.getState().pagination.pageSize +
                        1
                    }}
                    -
                    {{
                        Math.min(
                            (table.getState().pagination.pageIndex + 1) *
                            table.getState().pagination.pageSize,
                            nodes.length
                        )
                    }}
                    of {{ nodes.length }} nodes
                </div>
                <div class="flex items-center gap-2">
                    <Button variant="outline" size="sm" :disabled="!table.getCanPreviousPage()"
                        @click="table.previousPage()">
                        <Icon icon="lucide:chevron-left" class="h-4 w-4" />
                    </Button>
                    <div class="flex items-center gap-1">
                        <Button v-for="page in table.getPageCount()" :key="page" variant="outline" size="sm" :class="[
                            'min-w-[40px]',
                            table.getState().pagination.pageIndex + 1 === page
                                ? 'bg-primary text-primary-foreground'
                                : '',
                        ]" @click="table.setPageIndex(page - 1)">
                            {{ page }}
                        </Button>
                    </div>
                    <Button variant="outline" size="sm" :disabled="!table.getCanNextPage()" @click="table.nextPage()">
                        <Icon icon="lucide:chevron-right" class="h-4 w-4" />
                    </Button>
                </div>
            </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="rounded-lg border bg-card p-12 text-center">
            <Icon icon="lucide:server" class="h-12 w-12 mx-auto text-muted-foreground mb-4" />
            <p class="text-muted-foreground">No nodes found</p>
            <Button class="mt-4" @click="openAddDialog">
                <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
                Add First Node
            </Button>
        </div>

        <!-- 添加/编辑节点对话框 -->
        <Sheet v-model:open="isSheetOpen">
            <SheetContent class="flex h-screen max-h-screen flex-col overflow-hidden px-6 py-6">
                <SheetHeader class="px-3 sm:px-5">
                    <SheetTitle>{{ isEditMode ? "Edit Node" : "Add Node" }}</SheetTitle>
                    <SheetDescription>
                        {{
                            isEditMode
                                ? "Update node information"
                                : "Fill in the node information to add a new node"
                        }}
                    </SheetDescription>
                </SheetHeader>

                <div class="flex-1 overflow-y-auto pr-6">
                    <div class="px-3 sm:px-5">
                        <div class="space-y-2">
                            <label for="node-name" class="text-sm font-medium">Name *</label>
                            <Input id="node-name" v-model="formData.name" placeholder="Node name" />
                        </div>

                        <div class="space-y-2">
                            <label for="node-display-name" class="text-sm font-medium">Display Name</label>
                            <Input id="node-display-name" v-model="formData.display_name"
                                placeholder="Display name (optional)" />
                        </div>

                        <div class="space-y-2">
                            <label for="node-address" class="text-sm font-medium">Address *</label>
                            <Input id="node-address" v-model="formData.address" placeholder="192.168.1.1" />
                        </div>

                        <div class="grid grid-cols-2 gap-4">
                            <div class="space-y-2">
                                <label for="node-ssh-port" class="text-sm font-medium">SSH Port</label>
                                <Input id="node-ssh-port" v-model.number="formData.ssh_port" type="number"
                                    placeholder="22" />
                            </div>
                            <div class="space-y-2">
                                <label for="node-agent-port" class="text-sm font-medium">Agent Port</label>
                                <Input id="node-agent-port" v-model.number="formData.agent_port" type="number"
                                    placeholder="8080" />
                            </div>
                        </div>

                        <div class="space-y-2">
                            <label for="node-ssh-user" class="text-sm font-medium">SSH User</label>
                            <Input id="node-ssh-user" v-model="formData.ssh_user" placeholder="root" />
                        </div>

                        <div class="space-y-2">
                            <label for="node-ssh-password" class="text-sm font-medium">SSH Password</label>
                            <Input id="node-ssh-password" v-model="formData.ssh_password" type="password"
                                placeholder="SSH password" />
                        </div>

                        <div class="flex items-center space-x-2">
                            <input id="is_local" v-model="formData.is_local" type="checkbox"
                                class="h-4 w-4 rounded border-gray-300" />
                            <label for="is_local" class="text-sm font-medium">Local Node</label>
                        </div>
                    </div>
                </div>

                <SheetFooter class="px-3 sm:px-5">
                    <Button variant="outline" @click="isSheetOpen = false">Cancel</Button>
                    <Button @click="saveNode" :disabled="loading">
                        {{ loading ? "Saving..." : isEditMode ? "Update" : "Add" }}
                    </Button>
                </SheetFooter>
            </SheetContent>
        </Sheet>

        <!-- Delete Confirmation Dialog -->
        <AlertDialog v-model:open="isDeleteDialogOpen">
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Confirm Delete</AlertDialogTitle>
                    <AlertDialogDescription>
                        Are you sure you want to delete node "{{ nodeToDelete?.name || nodeToDelete?.display_name }}"?
                        This action cannot be undone.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction @click="confirmDeleteNode" :disabled="loading">
                        {{ loading ? "Deleting..." : "Delete" }}
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    </div>
</template>
