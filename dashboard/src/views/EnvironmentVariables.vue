<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import {
  useVueTable,
  getCoreRowModel,
  getPaginationRowModel,
  type ColumnDef,
} from "@tanstack/vue-table";
import type { AcceptableValue } from "reka-ui";

interface EnvironmentVariable {
  id: number;
  key: string;
  value: string;
  scope: "global" | "node";
  node_id?: number;
  node_name?: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

interface Node {
  id: number;
  name: string;
  display_name?: string | null;
}

// Mock nodes data
const mockNodes: Node[] = [
  { id: 1, name: "node-01", display_name: "Node 01" },
  { id: 2, name: "node-02", display_name: "Node 02" },
  { id: 3, name: "node-03", display_name: "Node 03" },
  { id: 4, name: "local-node", display_name: "Local Node" },
];

// Mock environment variables data
const mockEnvVars: EnvironmentVariable[] = [
  {
    id: 1,
    key: "DATABASE_URL",
    value: "postgresql://user:pass@localhost:5432/db",
    scope: "global",
    description: "Database connection URL",
    created_at: "2024-01-15 10:00:00",
    updated_at: "2024-01-20 14:30:00",
  },
  {
    id: 2,
    key: "REDIS_URL",
    value: "redis://localhost:6379",
    scope: "global",
    description: "Redis connection URL",
    created_at: "2024-01-15 10:05:00",
    updated_at: "2024-01-20 14:30:00",
  },
  {
    id: 3,
    key: "API_KEY",
    value: "sk-1234567890abcdef",
    scope: "global",
    description: "API authentication key",
    created_at: "2024-01-15 10:10:00",
    updated_at: "2024-01-20 14:30:00",
  },
  {
    id: 4,
    key: "NODE_SPECIFIC_CONFIG",
    value: "custom-value-1",
    scope: "node",
    node_id: 1,
    node_name: "Node 01",
    description: "Node-specific configuration",
    created_at: "2024-01-15 10:15:00",
    updated_at: "2024-01-20 14:30:00",
  },
  {
    id: 5,
    key: "NODE_SPECIFIC_CONFIG",
    value: "custom-value-2",
    scope: "node",
    node_id: 2,
    node_name: "Node 02",
    description: "Node-specific configuration",
    created_at: "2024-01-15 10:20:00",
    updated_at: "2024-01-20 14:30:00",
  },
  {
    id: 6,
    key: "LOG_LEVEL",
    value: "debug",
    scope: "node",
    node_id: 3,
    node_name: "Node 03",
    description: "Logging level for this node",
    created_at: "2024-01-15 10:25:00",
    updated_at: "2024-01-20 14:30:00",
  },
];

const envVars = ref<EnvironmentVariable[]>([]);
const loading = ref(false);
const nodes = ref<Node[]>([]);

// Filter
const filterScope = ref<"all" | "global" | "node">("all");
const selectedNodeId = ref<string>("all");

const filteredEnvVars = computed(() => {
  let filtered = envVars.value;

  if (filterScope.value === "global") {
    filtered = filtered.filter((v) => v.scope === "global");
  } else if (filterScope.value === "node") {
    filtered = filtered.filter((v) => v.scope === "node");
  }

  if (selectedNodeId.value !== "all") {
    filtered = filtered.filter(
      (v) => v.node_id?.toString() === selectedNodeId.value
    );
  }

  return filtered;
});

// Dialog state
const isSheetOpen = ref(false);
const isEditMode = ref(false);
const editingEnvVar = ref<EnvironmentVariable | null>(null);

// Form data
const formData = ref({
  key: "",
  value: "",
  scope: "global" as "global" | "node",
  node_id: undefined as number | undefined,
  description: "",
});

// Fetch data
const fetchEnvVars = () => {
  loading.value = true;
  setTimeout(() => {
    envVars.value = [...mockEnvVars];
    loading.value = false;
  }, 300);
};

const fetchNodes = () => {
  nodes.value = [...mockNodes];
};

// Open add dialog
const openAddDialog = () => {
  isEditMode.value = false;
  editingEnvVar.value = null;
  formData.value = {
    key: "",
    value: "",
    scope: "global",
    node_id: undefined,
    description: "",
  };
  isSheetOpen.value = true;
};

// Open edit dialog
const openEditDialog = (envVar: EnvironmentVariable) => {
  isEditMode.value = true;
  editingEnvVar.value = envVar;
  formData.value = {
    key: envVar.key,
    value: envVar.value,
    scope: envVar.scope,
    node_id: envVar.node_id,
    description: envVar.description || "",
  };
  isSheetOpen.value = true;
};

// Save environment variable
const saveEnvVar = () => {
  if (!formData.value.key || !formData.value.value) {
    return;
  }

  loading.value = true;
  setTimeout(() => {
    if (isEditMode.value && editingEnvVar.value) {
      const index = envVars.value.findIndex(
        (v) => v.id === editingEnvVar.value!.id
      );
      if (index !== -1) {
        const existing = envVars.value[index]!;
        envVars.value[index] = {
          id: existing.id,
          key: formData.value.key,
          value: formData.value.value,
          scope: formData.value.scope,
          node_id: formData.value.node_id,
          node_name:
            formData.value.scope === "node" && formData.value.node_id
              ? nodes.value.find((n) => n.id === formData.value.node_id)
                ?.display_name || nodes.value.find((n) => n.id === formData.value.node_id)?.name
              : undefined,
          description: formData.value.description || undefined,
          created_at: existing.created_at,
          updated_at: existing.updated_at,
        };
      }
    } else {
      const newId = Math.max(...envVars.value.map((v) => v.id), 0) + 1;
      const newNodeName =
        formData.value.scope === "node" && formData.value.node_id
          ? nodes.value.find((n) => n.id === formData.value.node_id)
            ?.display_name || nodes.value.find((n) => n.id === formData.value.node_id)?.name
          : undefined;
      envVars.value.push({
        id: newId,
        key: formData.value.key,
        value: formData.value.value,
        scope: formData.value.scope,
        node_id: formData.value.node_id,
        node_name: newNodeName,
        description: formData.value.description || undefined,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      });
    }
    isSheetOpen.value = false;
    loading.value = false;
  }, 300);
};

// Delete environment variable
const deleteEnvVar = (envVar: EnvironmentVariable) => {
  if (
    !confirm(
      `Are you sure you want to delete environment variable "${envVar.key}"?`
    )
  ) {
    return;
  }

  loading.value = true;
  setTimeout(() => {
    const index = envVars.value.findIndex((v) => v.id === envVar.id);
    if (index !== -1) {
      envVars.value.splice(index, 1);
    }
    loading.value = false;
  }, 300);
};

// Define table columns
const columns: ColumnDef<EnvironmentVariable>[] = [
  {
    accessorKey: "key",
    header: "Key",
  },
  {
    accessorKey: "value",
    header: "Value",
    cell: (info) => {
      const value = info.getValue() as string;
      return value.length > 50 ? value.substring(0, 50) + "..." : value;
    },
  },
  {
    accessorKey: "scope",
    header: "Scope",
    cell: (info) => {
      const scope = info.getValue() as string;
      return { scope };
    },
  },
  {
    accessorKey: "node_name",
    header: "Node",
  },
  {
    accessorKey: "description",
    header: "Description",
  },
  {
    id: "actions",
    header: "Actions",
    cell: (info) => {
      return { envVar: info.row.original };
    },
  },
];

// Create table instance
const table = useVueTable({
  get data() {
    return filteredEnvVars.value;
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

const handleNodeIdUpdate = (val: AcceptableValue) => {
  if (val === null || val === undefined) {
    formData.value.node_id = undefined;
    return;
  }
  const stringVal = typeof val === 'string' ? val : String(val);
  formData.value.node_id = stringVal ? parseInt(stringVal, 10) : undefined;
};

onMounted(() => {
  fetchEnvVars();
  fetchNodes();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold">Environment Variables</h1>
        <p class="text-muted-foreground mt-1">
          Manage global and node-specific environment variables
        </p>
      </div>
      <Button @click="openAddDialog">
        <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
        Add Variable
      </Button>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2">
        <label class="text-sm font-medium">Scope:</label>
        <Select v-model="filterScope">
          <SelectTrigger class="w-[150px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All</SelectItem>
            <SelectItem value="global">Global</SelectItem>
            <SelectItem value="node">Node</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div v-if="filterScope === 'node' || filterScope === 'all'" class="flex items-center gap-2">
        <label class="text-sm font-medium">Node:</label>
        <Select v-model="selectedNodeId">
          <SelectTrigger class="w-[200px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Nodes</SelectItem>
            <SelectItem v-for="node in nodes" :key="node.id" :value="node.id.toString()">
              {{ node.display_name || node.name }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && envVars.length === 0" class="text-center py-8 text-muted-foreground">
      Loading...
    </div>

    <!-- Environment Variables Table -->
    <div v-else-if="filteredEnvVars.length > 0" class="rounded-lg border bg-card">
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
              <template v-if="cell.column.id === 'scope'">
                <span :class="[
                  'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium',
                  (cell.getValue() as { scope: string }).scope === 'global'
                    ? 'bg-blue-500/10 text-blue-500 border-blue-500/20'
                    : 'bg-purple-500/10 text-purple-500 border-purple-500/20',
                ]">
                  {{ (cell.getValue() as { scope: string }).scope }}
                </span>
              </template>
              <template v-else-if="cell.column.id === 'actions'">
                <div class="flex items-center gap-2">
                  <Button variant="ghost" size="sm"
                    @click="openEditDialog((cell.getValue() as { envVar: EnvironmentVariable }).envVar)"
                    class="h-8 px-2">
                    <Icon icon="lucide:edit" class="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm"
                    @click="deleteEnvVar((cell.getValue() as { envVar: EnvironmentVariable }).envVar)"
                    class="h-8 px-2 text-red-500 hover:text-red-600">
                    <Icon icon="lucide:trash-2" class="h-4 w-4" />
                  </Button>
                </div>
              </template>
              <template v-else>
                <span v-if="cell.column.id === 'key'" class="font-medium font-mono text-xs">
                  {{ cell.getValue() }}
                </span>
                <span v-else-if="cell.column.id === 'value'" class="font-mono text-xs">
                  {{ cell.getValue() }}
                </span>
                <span v-else-if="cell.column.id === 'node_name'">
                  {{ cell.getValue() || "-" }}
                </span>
                <span v-else>
                  {{ cell.getValue() || "-" }}
                </span>
              </template>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <!-- Pagination Controls -->
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
              filteredEnvVars.length
            )
          }}
          of {{ filteredEnvVars.length }} variables
        </div>
        <div class="flex items-center gap-2">
          <Button variant="outline" size="sm" :disabled="!table.getCanPreviousPage()" @click="table.previousPage()">
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

    <!-- Empty State -->
    <div v-else class="rounded-lg border bg-card p-12 text-center">
      <Icon icon="lucide:key" class="h-12 w-12 mx-auto text-muted-foreground mb-4" />
      <p class="text-muted-foreground">No environment variables found</p>
      <Button class="mt-4" @click="openAddDialog">
        <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
        Add First Variable
      </Button>
    </div>

    <!-- Add/Edit Environment Variable Dialog -->
    <Sheet v-model:open="isSheetOpen">
      <SheetContent class="px-6 py-6">
        <SheetHeader class="!px-0 !pt-0">
          <SheetTitle>
            {{ isEditMode ? "Edit Environment Variable" : "Add Environment Variable" }}
          </SheetTitle>
          <SheetDescription>
            {{
              isEditMode
                ? "Update environment variable information"
                : "Fill in the information to add a new environment variable"
            }}
          </SheetDescription>
        </SheetHeader>

        <div class="space-y-4 !px-0">
          <div class="space-y-2">
            <label for="env-key" class="text-sm font-medium">Key *</label>
            <Input id="env-key" v-model="formData.key" placeholder="ENV_VARIABLE_NAME" />
          </div>

          <div class="space-y-2">
            <label for="env-value" class="text-sm font-medium">Value *</label>
            <Input id="env-value" v-model="formData.value" placeholder="variable value" />
          </div>

          <div class="space-y-2">
            <label for="env-scope" class="text-sm font-medium">Scope *</label>
            <Select v-model="formData.scope">
              <SelectTrigger id="env-scope">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="global">Global</SelectItem>
                <SelectItem value="node">Node</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div v-if="formData.scope === 'node'" class="space-y-2">
            <label for="env-node" class="text-sm font-medium">Node *</label>
            <Select :model-value="formData.node_id?.toString()" @update:model-value="handleNodeIdUpdate">
              <SelectTrigger id="env-node">
                <SelectValue placeholder="Select a node" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="node in nodes" :key="node.id" :value="node.id.toString()">
                  {{ node.display_name || node.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="space-y-2">
            <label for="env-description" class="text-sm font-medium">Description</label>
            <Input id="env-description" v-model="formData.description" placeholder="Optional description" />
          </div>
        </div>

        <SheetFooter class="!px-0 !pb-0">
          <Button variant="outline" @click="isSheetOpen = false">Cancel</Button>
          <Button @click="saveEnvVar" :disabled="loading">
            {{ loading ? "Saving..." : isEditMode ? "Update" : "Add" }}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  </div>
</template>
