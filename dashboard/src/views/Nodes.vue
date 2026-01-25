<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Checkbox from 'primevue/checkbox'
import Sidebar from 'primevue/sidebar'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import { nodeApi } from "@/api/node";
import { ApiResponseHelper } from "@/api/base";
import { showToast } from "@/lib/toast";

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
  ssh_user: string;
  ssh_password: string;
  is_local: boolean;
  display_name?: string;
}

const router = useRouter();
const nodes = ref<Node[]>([]);
const loading = ref(false);
const isSheetOpen = ref(false);
const isEditMode = ref(false);
const editingNode = ref<Node | null>(null);
const isDeleteDialogOpen = ref(false);
const nodeToDelete = ref<Node | null>(null);
const totalRecords = ref(0);
const lazyParams = ref({ first: 0, rows: 10 });

const formData = ref<CreateNodeRequest>({
  name: "", address: "", ssh_port: 22, ssh_user: "root", ssh_password: "",
  is_local: false, display_name: "",
});

const fetchNodes = async () => {
  loading.value = true;
  const page = Math.floor(lazyParams.value.first / lazyParams.value.rows) + 1;
  const res = await nodeApi.page(page, lazyParams.value.rows);
  loading.value = false;
  if (!ApiResponseHelper.isSuccess(res)) return;
  const d = res.data!;
  nodes.value = (d.nodes ?? []).map((n) => ({
    id: n.id,
    name: n.name,
    address: n.address,
    is_local: n.is_local,
    display_name: n.display_name ?? null,
    ssh_port: n.ssh_port,
    agent_port: n.agent_port,
    ssh_user: n.ssh_user,
    ssh_password: "",
    status: n.status || "offline",
  }));
  totalRecords.value = d.total ?? 0;
};

const onPage = (e: { first: number; rows: number }) => {
  lazyParams.value = { first: e.first, rows: e.rows };
  fetchNodes();
};

const openAddDialog = () => {
  isEditMode.value = false;
  editingNode.value = null;
  formData.value = { name: "", address: "", ssh_port: 22, ssh_user: "root", ssh_password: "", is_local: false, display_name: "" };
  isSheetOpen.value = true;
};

const openEditDialog = (node: Node) => {
  isEditMode.value = true;
  editingNode.value = node;
  formData.value = {
    name: node.name,
    address: node.address,
    ssh_port: node.ssh_port,
    ssh_user: node.ssh_user,
    ssh_password: node.ssh_password,
    is_local: node.is_local,
    display_name: node.display_name ?? "",
  };
  isSheetOpen.value = true;
};

const saveNode = async () => {
  if (!formData.value.name || !formData.value.address) {
    showToast("Name and address are required", "error");
    return;
  }
  loading.value = true;
  try {
    if (isEditMode.value && editingNode.value) {
      const res = await nodeApi.update({
        id: editingNode.value.id,
        ...formData.value,
        display_name: formData.value.display_name || "",
      });
      if (ApiResponseHelper.isSuccess(res)) {
        isSheetOpen.value = false;
        fetchNodes();
      } else showToast(res.message ?? "Failed to update node", "error");
    } else {
      const res = await nodeApi.create(formData.value);
      if (ApiResponseHelper.isSuccess(res)) {
        isSheetOpen.value = false;
        fetchNodes();
      } else showToast(res.message ?? "Failed to add node", "error");
    }
  } finally {
    loading.value = false;
  }
};

const openDeleteDialog = (node: Node) => {
  nodeToDelete.value = node;
  isDeleteDialogOpen.value = true;
};

const confirmDeleteNode = async () => {
  if (!nodeToDelete.value) return;
  loading.value = true;
  const res = await nodeApi.delete(nodeToDelete.value.id);
  loading.value = false;
  if (ApiResponseHelper.isSuccess(res)) fetchNodes();
  else showToast(res.message ?? "Failed to delete node", "error");
  isDeleteDialogOpen.value = false;
  nodeToDelete.value = null;
};

const openTerminal = (node: Node) => {
  router.push({ name: "NodeTerminal", query: { nodeId: String(node.id), nodeName: node.display_name || node.name } });
};

onMounted(fetchNodes);
</script>

<template>
  <div class="page-root">
    <div class="page-header page-header-row">
      <div>
        <h1>Nodes</h1>
        <p class="text-muted-foreground">Manage and monitor your server nodes</p>
      </div>
      <Button label="Add Node" icon="pi pi-plus" @click="openAddDialog" />
    </div>

    <div v-if="loading && nodes.length === 0" class="loading-state">
      Loading...
    </div>

    <div v-else-if="nodes.length > 0" class="table-wrap">
      <DataTable
        :value="nodes"
        :lazy="true"
        :paginator="true"
        :first="lazyParams.first"
        :rows="lazyParams.rows"
        :totalRecords="totalRecords"
        :loading="loading"
        data-key="id"
        @page="onPage"
        size="small"
      >
        <Column field="id" header="ID" />
        <Column field="name" header="Name">
          <template #body="{ data }">{{ data.display_name || data.name }}</template>
        </Column>
        <Column field="address" header="Address" />
        <Column field="status" header="Status">
          <template #body="{ data }">
            <Tag :value="data.status || 'offline'" :severity="data.status === 'online' ? 'success' : 'danger'" />
          </template>
        </Column>
        <Column field="is_local" header="Type">
          <template #body="{ data }">
            <Tag :value="data.is_local ? 'Local' : 'Remote'" :severity="data.is_local ? 'info' : 'contrast'" />
          </template>
        </Column>
        <Column field="ssh_port" header="SSH Port" />
        <Column field="agent_port" header="Agent Port" />
        <Column field="ssh_user" header="SSH User" />
        <Column header="Actions">
          <template #body="{ data }">
            <div class="action-btns">
              <Button text rounded size="small" :disabled="data.status !== 'online'" @click="openTerminal(data)" v-tooltip.top="'Terminal'">
                <i class="pi pi-terminal"></i>
              </Button>
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
    </div>

    <div v-else class="empty-state">
      <i class="pi pi-server"></i>
      <p>No nodes found</p>
      <Button label="Add First Node" icon="pi pi-plus" @click="openAddDialog" />
    </div>

    <Sidebar v-model:visible="isSheetOpen" position="right" :style="{ width: '90vw', maxWidth: '1200px' }" class="sheet">
      <div class="sheet-header">
        <h2>{{ isEditMode ? "Edit Node" : "Add Node" }}</h2>
        <p class="text-muted-foreground">{{ isEditMode ? "Update node information" : "Fill in the node information" }}</p>
      </div>
      <div class="sheet-body">
        <div class="form-group">
          <label for="node-name">Name *</label>
          <InputText id="node-name" v-model="formData.name" placeholder="Node name" />
        </div>
        <div class="form-group">
          <label for="node-display-name">Display Name</label>
          <InputText id="node-display-name" v-model="formData.display_name" placeholder="Display name (optional)" />
        </div>
        <div class="form-group">
          <label for="node-address">Address *</label>
          <InputText id="node-address" v-model="formData.address" placeholder="192.168.1.1" />
        </div>
        <div class="form-group">
          <label for="node-ssh-port">SSH Port</label>
          <InputNumber id="node-ssh-port" v-model="formData.ssh_port" placeholder="22" />
        </div>
        <div class="form-group">
          <label for="node-ssh-user">SSH User</label>
          <InputText id="node-ssh-user" v-model="formData.ssh_user" placeholder="root" />
        </div>
        <div class="form-group">
          <label for="node-ssh-password">SSH Password</label>
          <InputText id="node-ssh-password" v-model="formData.ssh_password" type="password" placeholder="SSH password" />
        </div>
        <div class="form-group flex-row">
          <Checkbox v-model="formData.is_local" binary inputId="is_local" />
          <label for="is_local">Local Node</label>
        </div>
      </div>
      <div class="sheet-footer">
        <Button outlined @click="isSheetOpen = false">Cancel</Button>
        <Button @click="saveNode" :disabled="loading">{{ loading ? "Saving..." : isEditMode ? "Update" : "Add" }}</Button>
      </div>
    </Sidebar>

    <Dialog v-model:visible="isDeleteDialogOpen" modal header="Confirm Delete" :style="{ width: '425px' }">
      <p>Are you sure you want to delete node "{{ nodeToDelete?.name || nodeToDelete?.display_name }}"? This action cannot be undone.</p>
      <template #footer>
        <Button outlined @click="isDeleteDialogOpen = false">Cancel</Button>
        <Button severity="danger" @click="confirmDeleteNode" :disabled="loading">{{ loading ? "Deleting..." : "Delete" }}</Button>
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

.loading-state { text-align: center; padding: 2rem; color: var(--p-text-muted-color); }
.table-wrap { background: var(--p-surface-card); border-radius: var(--p-border-radius); overflow: hidden; }
.action-btns { display: flex; align-items: center; gap: 0.5rem; }

.empty-state { padding: 3rem; text-align: center; background: var(--p-surface-card); border-radius: var(--p-border-radius); }
.empty-state i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
.empty-state .p-button { margin-top: 1rem; }

.sheet { display: flex; flex-direction: column; }
.sheet-header { padding: 0 1rem 1rem; }
.sheet-header h2 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.sheet-body { flex: 1; overflow-y: auto; padding: 0 1rem; display: flex; flex-direction: column; gap: 1rem; }
.sheet-footer { display: flex; justify-content: flex-end; gap: 0.5rem; padding: 1rem; border-top: 1px solid var(--p-surface-border); }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; }
.form-group.flex-row { flex-direction: row; align-items: center; gap: 0.5rem; }
</style>
