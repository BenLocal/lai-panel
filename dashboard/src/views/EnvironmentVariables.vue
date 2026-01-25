<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Sidebar from 'primevue/sidebar'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import { envApi, type AddOrUpdateEnvRequest, type GetEnvPageRequest } from "@/api/env";
import { ApiResponseHelper } from "@/api/base";

interface EnvironmentVariable {
  id: number;
  key: string;
  value: string;
  scope: string;
  description?: string;
  node_name?: string;
}

const envVars = ref<EnvironmentVariable[]>([]);
const loading = ref(false);
const filterScope = ref<string>("all");
const scopeOptions = ref<string[]>(["all"]);
const isSheetOpen = ref(false);
const isEditMode = ref(false);
const editingEnvVar = ref<EnvironmentVariable | null>(null);
const isDeleteDialogOpen = ref(false);
const envVarToDelete = ref<EnvironmentVariable | null>(null);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)));

const formData = ref({
  id: null as number | null,
  key: "",
  value: "",
  scope: "",
  description: "",
});

const fetchEnvVars = async () => {
  loading.value = true;
  try {
    const req: GetEnvPageRequest = {
      scope: filterScope.value === "all" ? "" : filterScope.value,
      page: currentPage.value,
      page_size: pageSize.value,
    };
    const res = await envApi.page(req);
    if (!ApiResponseHelper.isSuccess(res)) return;
    const d = res.data!;
    envVars.value = (d.list ?? []) as EnvironmentVariable[];
    total.value = d.total ?? 0;
    currentPage.value = d.current_page ?? 1;
    pageSize.value = d.page_size ?? 10;
  } finally {
    loading.value = false;
  }
};

const fetchScopes = async () => {
  const res = await envApi.scopes();
  if (!ApiResponseHelper.isSuccess(res)) return;
  const list = (res.data ?? []) as string[];
  scopeOptions.value = ["all", ...list];
};

const openAddDialog = () => {
  isEditMode.value = false;
  editingEnvVar.value = null;
  formData.value = { id: null, key: "", value: "", scope: "", description: "" };
  isSheetOpen.value = true;
};

const openEditDialog = (ev: EnvironmentVariable) => {
  isEditMode.value = true;
  editingEnvVar.value = ev;
  formData.value = {
    id: ev.id,
    key: ev.key,
    value: ev.value,
    scope: ev.scope,
    description: ev.description ?? "",
  };
  isSheetOpen.value = true;
  fetchEnvVars();
  fetchScopes();
};

const saveEnvVar = async () => {
  if (!formData.value.key?.trim() || !formData.value.value?.trim()) return;
  loading.value = true;
  const scope = formData.value.scope?.trim() || "global";
  const req: AddOrUpdateEnvRequest = {
    id: formData.value.id ?? null,
    key: formData.value.key.trim(),
    value: formData.value.value.trim(),
    scope,
  };
  const res = await envApi.addOrUpdate(req);
  loading.value = false;
  if (!ApiResponseHelper.isSuccess(res)) return;
  isSheetOpen.value = false;
  fetchEnvVars();
  fetchScopes();
};

const openDeleteDialog = (ev: EnvironmentVariable) => {
  envVarToDelete.value = ev;
  isDeleteDialogOpen.value = true;
};

const confirmDeleteEnvVar = async () => {
  if (!envVarToDelete.value?.id) return;
  loading.value = true;
  const res = await envApi.delete(envVarToDelete.value.id);
  loading.value = false;
  if (!ApiResponseHelper.isSuccess(res)) return;
  isDeleteDialogOpen.value = false;
  envVarToDelete.value = null;
  fetchEnvVars();
  fetchScopes();
};

const goToPage = (p: number) => {
  if (p < 1 || p > totalPages.value) return;
  currentPage.value = p;
  fetchEnvVars();
};

onMounted(() => {
  fetchEnvVars();
  fetchScopes();
});
</script>

<template>
  <div class="page-root">
    <div class="page-header page-header-row">
      <div>
        <h1>Environment Variables</h1>
        <p class="text-muted-foreground">Manage global and node-specific environment variables</p>
      </div>
      <Button label="Add Variable" icon="pi pi-plus" @click="openAddDialog" />
    </div>

    <div class="toolbar">
      <label>Scope:</label>
      <Select v-model="filterScope" :options="scopeOptions" @update:model-value="fetchEnvVars" class="scope-select" />
    </div>

    <div v-if="loading && !envVars.length" class="loading-state">Loading...</div>

    <div v-else-if="envVars.length > 0" class="table-wrap">
      <DataTable :value="envVars" size="small" striped-rows>
        <Column field="key" header="Key">
          <template #body="{ data }">
            <span class="font-mono text-sm">{{ data.key }}</span>
          </template>
        </Column>
        <Column field="value" header="Value">
          <template #body="{ data }">
            <span class="font-mono text-sm">{{ data.value }}</span>
          </template>
        </Column>
        <Column field="scope" header="Scope">
          <template #body="{ data }">
            <Tag :value="data.scope" :severity="data.scope === 'global' ? 'info' : 'contrast'" />
          </template>
        </Column>
        <Column field="node_name" header="Node">
          <template #body="{ data }">{{ data.node_name ?? "-" }}</template>
        </Column>
        <Column field="description" header="Description">
          <template #body="{ data }">{{ data.description ?? "-" }}</template>
        </Column>
        <Column header="Actions">
          <template #body="{ data }">
            <div class="action-btns">
              <Button text rounded size="small" @click="openEditDialog(data)" v-tooltip.top="'Edit'">
                <i class="pi pi-pencil"></i>
              </Button>
              <Button text rounded size="small" severity="danger" @click="openDeleteDialog(data)" v-tooltip.top="'Delete'">
                <i class="pi pi-trash"></i>
              </Button>
            </div>
          </template>
        </Column>
      </DataTable>
      <div v-if="totalPages > 1" class="pagination-bar">
        <span class="text-muted-foreground">
          {{ (currentPage - 1) * pageSize + 1 }}–{{ Math.min(currentPage * pageSize, total) }} of {{ total }}
        </span>
        <div class="pagination-btns">
          <Button outlined size="small" :disabled="currentPage === 1" @click="goToPage(currentPage - 1)">
            <i class="pi pi-chevron-left"></i>
          </Button>
          <Button
            v-for="p in totalPages"
            :key="p"
            outlined
            size="small"
            :class="{ 'pagination-active': currentPage === p }"
            @click="goToPage(p)"
          >
            {{ p }}
          </Button>
          <Button outlined size="small" :disabled="currentPage === totalPages" @click="goToPage(currentPage + 1)">
            <i class="pi pi-chevron-right"></i>
          </Button>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <i class="pi pi-key"></i>
      <p>No environment variables found</p>
      <Button label="Add First Variable" icon="pi pi-plus" @click="openAddDialog" />
    </div>

    <Sidebar v-model:visible="isSheetOpen" position="right" :style="{ width: '90vw', maxWidth: '1200px' }" class="sheet">
      <div class="sheet-header">
        <h2>{{ isEditMode ? "Edit" : "Add" }} Environment Variable</h2>
        <p class="text-muted-foreground">{{ isEditMode ? "Update" : "Fill in" }} the variable details</p>
      </div>
      <div class="sheet-body">
        <div class="form-group">
          <label for="env-key">Key *</label>
          <InputText id="env-key" v-model="formData.key" placeholder="ENV_VARIABLE_NAME" />
        </div>
        <div class="form-group">
          <label for="env-value">Value *</label>
          <InputText id="env-value" v-model="formData.value" placeholder="value" />
        </div>
        <div class="form-group">
          <label for="env-scope">Scope</label>
          <InputText id="env-scope" v-model="formData.scope" placeholder="global (default)" />
          <p class="text-muted-foreground text-sm">Empty = global</p>
        </div>
        <div class="form-group">
          <label for="env-desc">Description</label>
          <InputText id="env-desc" v-model="formData.description" placeholder="Optional" />
        </div>
      </div>
      <div class="sheet-footer">
        <Button outlined @click="isSheetOpen = false">Cancel</Button>
        <Button @click="saveEnvVar" :disabled="loading">{{ loading ? "Saving..." : isEditMode ? "Update" : "Add" }}</Button>
      </div>
    </Sidebar>

    <Dialog v-model:visible="isDeleteDialogOpen" modal header="Delete Variable" :style="{ width: '425px' }">
      <p>Delete "{{ envVarToDelete?.key }}"? This cannot be undone.</p>
      <template #footer>
        <Button outlined @click="isDeleteDialogOpen = false">Cancel</Button>
        <Button severity="danger" @click="confirmDeleteEnvVar" :disabled="loading">{{ loading ? "Deleting..." : "Delete" }}</Button>
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.page-root { display: flex; flex-direction: column; gap: 1.5rem; }
.page-header h1 { font-size: 1.5rem; font-weight: 700; margin-bottom: 0.25rem; }
.page-header p { font-size: 0.875rem; }
.page-header-row { display: flex; align-items: center; justify-content: space-between; }
.btn-icon-text { margin-left: 0.5rem; }

.toolbar { display: flex; align-items: center; gap: 0.5rem; }
.toolbar label { font-size: 0.875rem; font-weight: 500; }
.scope-select { width: 150px; }

.loading-state { text-align: center; padding: 2rem; color: var(--p-text-muted-color); }
.table-wrap { background: var(--p-surface-card); border-radius: var(--p-border-radius); overflow: hidden; }
.action-btns { display: flex; gap: 0.5rem; }
.font-mono { font-family: ui-monospace, monospace; }
.text-sm { font-size: 0.75rem; }
.pagination-bar { display: flex; justify-content: space-between; align-items: center; padding: 1rem; border-top: 1px solid var(--p-surface-border); font-size: 0.875rem; }
.pagination-btns { display: flex; gap: 0.5rem; }
.pagination-active { background: var(--p-primary-color) !important; color: var(--p-primary-contrast-color) !important; border-color: var(--p-primary-color) !important; }

.empty-state { padding: 3rem; text-align: center; background: var(--p-surface-card); border-radius: var(--p-border-radius); }
.empty-state i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
.empty-state .p-button { margin-top: 1rem; }

.sheet-header { padding: 0 1rem 1rem; }
.sheet-header h2 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.sheet-body { flex: 1; overflow-y: auto; padding: 0 1rem; display: flex; flex-direction: column; gap: 1rem; }
.sheet-footer { display: flex; justify-content: flex-end; gap: 0.5rem; padding: 1rem; border-top: 1px solid var(--p-surface-border); }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; }
</style>
