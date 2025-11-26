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
  useVueTable,
  getCoreRowModel,
  getPaginationRowModel,
  type ColumnDef,
} from "@tanstack/vue-table";
import type { AcceptableValue } from "reka-ui";
import { envApi, type AddOrUpdateEnvRequest, type GetEnvPageRequest } from "@/api/env";
import { ApiResponseHelper } from "@/api/base";

interface EnvironmentVariable {
  id: number;
  key: string;
  value: string;
  scope: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

const envVars = ref<EnvironmentVariable[]>([]);
const loading = ref(false);

// Filter
const filterScope = ref<string>("all");
const scopeOptions = ref<string[]>(["all"]);

// Dialog state
const isSheetOpen = ref(false);
const isEditMode = ref(false);
const editingEnvVar = ref<EnvironmentVariable | null>(null);
const isDeleteDialogOpen = ref(false);
const envVarToDelete = ref<EnvironmentVariable | null>(null);

// Form data
const formData = ref({
  id: null as number | null,
  key: "",
  value: "",
  scope: "",
  description: "",
});

// Fetch data
const fetchEnvVars = async () => {
  loading.value = true;

  try {
    const request: GetEnvPageRequest = {
      scope: filterScope.value === "all" ? "" : filterScope.value,
      page: table.getState().pagination.pageIndex,
      page_size: table.getState().pagination.pageSize,
    };
    const response = await envApi.page(request);
    if (!ApiResponseHelper.isSuccess(response)) {
      return;
    }

    const data = response.data;
    envVars.value = data?.list || [];
    table.setPageIndex(data?.current_page || 0);
    table.setPageSize(data?.page_size || 10);
  } finally {
    loading.value = false;
  }

};

const fetchScopes = async () => {
  const response = await envApi.scopes();
  if (!ApiResponseHelper.isSuccess(response)) {
    return;
  }
  const data = response.data;
  scopeOptions.value = data || [];
  scopeOptions.value.unshift("all");
};

// Open add dialog
const openAddDialog = () => {
  isEditMode.value = false;
  editingEnvVar.value = null;
  formData.value = {
    id: null,
    key: "",
    value: "",
    scope: "",
    description: "",
  };
  isSheetOpen.value = true;
};

// Open edit dialog
const openEditDialog = (envVar: EnvironmentVariable) => {
  isEditMode.value = true;
  editingEnvVar.value = envVar;
  formData.value = {
    id: envVar.id,
    key: envVar.key,
    value: envVar.value,
    scope: envVar.scope,
    description: envVar.description || "",
  };
  isSheetOpen.value = true;
  refresh();
};

const refresh = () => {
  fetchEnvVars();
  fetchScopes();
};

// Save environment variable
const saveEnvVar = async () => {
  if (!formData.value.key || !formData.value.value) {
    return;
  }

  loading.value = true;
  // Use "global" as default if scope is empty or only whitespace
  const scope = formData.value.scope?.trim() || "global";
  let request: AddOrUpdateEnvRequest = {
    id: formData.value.id,
    key: formData.value.key,
    value: formData.value.value,
    scope: scope,
  };
  const response = await envApi.addOrUpdate(request);
  if (!ApiResponseHelper.isSuccess(response)) {
    return;
  }
  isSheetOpen.value = false;
  loading.value = false;
  refresh();
};

// Open delete dialog
const openDeleteDialog = (envVar: EnvironmentVariable) => {
  envVarToDelete.value = envVar;
  isDeleteDialogOpen.value = true;
};

// Confirm delete environment variable
const confirmDeleteEnvVar = async () => {
  if (!envVarToDelete.value) {
    return;
  }

  const id = envVarToDelete.value.id;
  if (!id) {
    return;
  }

  loading.value = true;
  const response = await envApi.delete(id);
  if (!ApiResponseHelper.isSuccess(response)) {
    loading.value = false;
    return;
  }
  loading.value = false;
  isDeleteDialogOpen.value = false;
  envVarToDelete.value = null;
  refresh();
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
  },
  {
    accessorKey: "scope",
    header: "Scope",
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
  },
];

// Create table instance
const table = useVueTable({
  get data() {
    return envVars.value;
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

onMounted(() => {
  fetchEnvVars();
  fetchScopes();
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
        <Select v-model="filterScope" @update:model-value="refresh">
          <SelectTrigger class="w-[150px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="scope in scopeOptions" :key="scope" :value="scope">
              {{ scope }}
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
    <div v-else-if="envVars.length > 0" class="rounded-lg border bg-card">
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
                  (cell.getValue() as string) === 'global'
                    ? 'bg-blue-500/10 text-blue-500 border-blue-500/20'
                    : 'bg-purple-500/10 text-purple-500 border-purple-500/20',
                ]">
                  {{ cell.getValue() }}
                </span>
              </template>
              <template v-else-if="cell.column.id === 'actions'">
                <div class="flex items-center gap-2">
                  <Button variant="ghost" size="sm" @click="openEditDialog(row.original)" class="h-8 px-2">
                    <Icon icon="lucide:edit" class="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" @click="openDeleteDialog(row.original)"
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
              envVars.length
            )
          }}
          of {{ envVars.length }} variables
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
            <label for="env-scope" class="text-sm font-medium">Scope</label>
            <Input id="env-scope" v-model="formData.scope" placeholder="global (default)" />
            <p class="text-xs text-muted-foreground">Leave empty to use default scope "global"</p>
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

    <!-- Delete Confirmation Dialog -->
    <AlertDialog v-model:open="isDeleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete Environment Variable</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to delete environment variable "{{ envVarToDelete?.key }}"?
            This action cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction @click="confirmDeleteEnvVar" :disabled="loading">
            {{ loading ? "Deleting..." : "Delete" }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
