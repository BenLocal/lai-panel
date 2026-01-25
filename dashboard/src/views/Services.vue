<script setup lang="ts">
import { ref, onMounted, computed, nextTick, watch } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Checkbox from 'primevue/checkbox'
import Select from 'primevue/select'
import Sidebar from 'primevue/sidebar'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Stepper from 'primevue/stepper'
import StepList from 'primevue/steplist'
import Step from 'primevue/step'
import StepPanels from 'primevue/steppanels'
import StepPanel from 'primevue/steppanel'
import {
  serviceApi,
  type DeployServiceRequest,
  type SaveServiceRequest,
  type Service,
} from "@/api/service";
import { applicationApi, type Application } from "@/api/application";
import { nodeApi, type Node } from "@/api/node";
import { showToast } from "@/lib/toast";
import { ApiResponseHelper } from "@/api/base";

const services = ref<Service[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
const totalPages = computed(() => Math.ceil(total.value / pageSize.value));

const isDeployDialogOpen = ref(false);
const currentStep = ref(1);
const prevClickStep = ref(-1);
const selectedApplication = ref<Application | null>(null);
const selectedNode = ref<Node | null>(null);
const qaValues = ref<Record<string, string>>({});
const serverName = ref("");
const applications = ref<Application[]>([]);
const nodes = ref<Node[]>([]);
const deployLoading = ref(false);
const applicationsLoading = ref(false);
const nodesLoading = ref(false);
const deployOutput = ref<string[]>([]);
const deployOutputRef = ref<HTMLDivElement | null>(null);
const isDeployOutputDrawerOpen = ref(false);
const editingServiceId = ref<number | null>(null);
const isDeleteDialogOpen = ref(false);
const serviceToDelete = ref<Service | null>(null);
const forceDelete = ref(false);

const fetchServices = async (page: number = currentPage.value) => {
  loading.value = true;
  try {
    const res = await serviceApi.page(page, pageSize.value);
    if (!ApiResponseHelper.isSuccess(res)) return;
    const d = res.data!;
    services.value = d.services ?? [];
    total.value = d.total ?? 0;
    currentPage.value = d.page ?? page;
    pageSize.value = d.pageSize ?? 10;
  } catch {
    showToast("Failed to fetch services", "error");
  } finally {
    loading.value = false;
  }
};

const fetchApplications = async () => {
  if (applicationsLoading.value) return;
  applicationsLoading.value = true;
  try {
    const res = await applicationApi.list();
    if (ApiResponseHelper.isSuccess(res)) applications.value = res.data ?? [];
    else showToast("Failed to fetch applications", "error");
  } finally {
    applicationsLoading.value = false;
  }
};

const fetchNodesForDeploy = async () => {
  if (nodesLoading.value) return;
  nodesLoading.value = true;
  try {
    const res = await nodeApi.list();
    if (ApiResponseHelper.isSuccess(res)) nodes.value = res.data ?? [];
    else showToast("Failed to fetch nodes", "error");
  } finally {
    nodesLoading.value = false;
  }
};

const openDeployDialog = async () => {
  currentStep.value = 1;
  prevClickStep.value = -1;
  selectedApplication.value = null;
  selectedNode.value = null;
  qaValues.value = {};
  serverName.value = "";
  editingServiceId.value = null;
  applications.value = [];
  nodes.value = [];
  isDeployDialogOpen.value = true;
  await goToStep(1);
};

const closeDeployDialog = () => {
  isDeployDialogOpen.value = false;
  currentStep.value = 1;
  prevClickStep.value = -1;
  selectedApplication.value = null;
  selectedNode.value = null;
  qaValues.value = {};
  serverName.value = "";
  editingServiceId.value = null;
  deployOutput.value = [];
};

const openEditDialog = async (s: Service) => {
  editingServiceId.value = s.id;
  serverName.value = s.name;
  currentStep.value = 1;
  prevClickStep.value = -1;
  qaValues.value = {};
  deployOutput.value = [];
  applications.value = [];
  nodes.value = [];
  isDeployDialogOpen.value = true;
  await fetchApplications();
  await fetchNodesForDeploy();
  const app = applications.value.find((a) => a.id === s.app_id);
  const node = nodes.value.find((n) => n.id === s.node_id);
  if (app) {
    selectedApplication.value = app;
    if (app.qa?.length) {
      const init: Record<string, string> = {};
      app.qa.forEach((item) => {
        init[item.name] = s.qa_values?.[item.name] ?? item.default_value ?? "";
      });
      qaValues.value = init;
    }
  }
  if (node) selectedNode.value = node;
  if (app && node) await goToStep(3);
  else await goToStep(1);
};

const handleStepChange = async (step: number) => {
  if (step === 2 && !selectedApplication.value) return;
  if (prevClickStep.value === step) return;
  const prev = prevClickStep.value;
  prevClickStep.value = step;
  if (prev !== step) {
    if (step === 1 && !applications.value.length && !applicationsLoading.value) await fetchApplications();
    else if (step === 3 && !nodes.value.length && !nodesLoading.value) await fetchNodesForDeploy();
  }
};

const goToStep = async (step: number) => {
  currentStep.value = step;
  await handleStepChange(step);
};

const onStepperChange = async (step: number) => {
  if (step === 2 && !selectedApplication.value) return;
  await goToStep(step);
};

const handleApplicationSelect = (app: Application | null) => {
  if (!app) {
    qaValues.value = {};
    return;
  }
  selectedApplication.value = app;
  if (!app.qa?.length) return;
  const init: Record<string, string> = {};
  let existing: Record<string, string> | undefined;
  if (editingServiceId.value) {
    const s = services.value.find((x) => x.id === editingServiceId.value);
    if (s?.qa_values) existing = s.qa_values;
  }
  app.qa.forEach((item) => {
    init[item.name] = existing?.[item.name] ?? item.default_value ?? "";
  });
  qaValues.value = init;
};

const handleNodeSelect = (node: Node | null) => {
  if (node) selectedNode.value = node;
};

const canProceedToNextStep = computed(() => {
  if (currentStep.value === 1) return selectedApplication.value != null;
  if (currentStep.value === 2) {
    const qa = selectedApplication.value?.qa;
    if (!qa?.length) return true;
    return qa.every((item) => !item.required || (qaValues.value[item.name] ?? "").trim() !== "");
  }
  if (currentStep.value === 3) return selectedNode.value != null;
  return false;
});

const getQAValue = (name: string) => qaValues.value[name] ?? "";
const handleQAValueChange = (name: string, value: string) => {
  qaValues.value[name] = value;
};

const scrollToBottom = () => {
  nextTick(() => {
    if (deployOutputRef.value) deployOutputRef.value.scrollTop = deployOutputRef.value.scrollHeight;
  });
};

watch(deployOutput, () => {
  scrollToBottom();
  if (deployOutput.value.length > 0 && !isDeployOutputDrawerOpen.value) isDeployOutputDrawerOpen.value = true;
}, { deep: true });

const deployService = async (opts?: { onlySave?: boolean }) => {
  const onlySave = opts?.onlySave ?? false;
  if (!selectedApplication.value || !selectedNode.value) return;
  const saveReq: SaveServiceRequest = {
    id: editingServiceId.value ?? 0,
    name: serverName.value || selectedApplication.value.name,
    app_id: selectedApplication.value.id,
    node_id: selectedNode.value.id,
    qa_values: qaValues.value,
  };
  const res = await serviceApi.save(saveReq);
  if (!ApiResponseHelper.isSuccess(res)) {
    showToast(res.message ?? "Failed to save service", "error");
    return;
  }
  const serviceId = res.data?.id ?? editingServiceId.value ?? 0;
  await fetchServices(currentPage.value);
  if (onlySave) {
    showToast(editingServiceId.value ? "Service updated" : "Service saved", "success");
    closeDeployDialog();
    return;
  }
  if (!serviceId) {
    showToast("Failed to get service ID", "error");
    closeDeployDialog();
    return;
  }
  editingServiceId.value = serviceId;
  deployLoading.value = true;
  deployOutput.value = [];
  isDeployOutputDrawerOpen.value = true;
  const req: DeployServiceRequest = {
    service_id: serviceId,
    app_id: selectedApplication.value.id,
    node_id: selectedNode.value.id,
    qa_values: qaValues.value,
  };
  await serviceApi.deployStream(
    req,
    (data) => deployOutput.value.push(data),
    (err) => {
      deployOutput.value.push(`[Error] ${err.message}`);
      showToast("Deploy failed", "error");
      deployLoading.value = false;
    },
    () => {
      deployLoading.value = false;
      deployOutput.value.push("Deployment completed!!!");
      fetchServices(currentPage.value);
    }
  );
  closeDeployDialog();
};

const openDeleteDialog = (s: Service) => {
  serviceToDelete.value = s;
  forceDelete.value = false;
  isDeleteDialogOpen.value = true;
};

const confirmDeleteService = async () => {
  if (!serviceToDelete.value) return;
  loading.value = true;
  try {
    const res = await serviceApi.delete(serviceToDelete.value.id, forceDelete.value);
    if (ApiResponseHelper.isSuccess(res)) {
      showToast("Service deleted", "success");
      await fetchServices(currentPage.value);
      if (services.value.length === 0 && currentPage.value > 1) await fetchServices(currentPage.value - 1);
    } else showToast(res.message ?? "Failed to delete", "error");
  } catch {
    showToast("Failed to delete service", "error");
  } finally {
    loading.value = false;
    isDeleteDialogOpen.value = false;
    serviceToDelete.value = null;
    forceDelete.value = false;
  }
};

const goToPage = (p: number) => {
  if (p < 1 || p > totalPages.value) return;
  currentPage.value = p;
  fetchServices(p);
};

const copyDeployOutput = async () => {
  if (!deployOutput.value.length) {
    showToast("No output to copy", "info");
    return;
  }
  try {
    await navigator.clipboard.writeText(deployOutput.value.join("\n"));
    showToast("Output copied", "success");
  } catch {
    showToast("Failed to copy", "error");
  }
};

const statusSeverity = (s?: string) => {
  if (s === "running") return "success";
  if (s === "stopped") return "danger";
  if (s === "pending") return "warn";
  return "secondary";
};

onMounted(fetchServices);
</script>

<template>
  <div class="page-root">
    <div class="page-header page-header-row">
      <div>
        <h1>Services</h1>
        <p class="text-muted-foreground">Manage and monitor your deployed services</p>
      </div>
      <Button label="Deploy Service" icon="pi pi-plus" @click="openDeployDialog" />
    </div>

    <div v-if="loading && !services.length" class="loading-state">Loading...</div>

    <div v-else-if="services.length > 0" class="table-wrap">
      <DataTable :value="services" size="small" striped-rows>
        <Column field="id" header="ID" />
        <Column field="name" header="Name" />
        <Column field="status" header="Status">
          <template #body="{ data }">
            <Tag :value="data.status || 'Unknown'" :severity="statusSeverity(data.status)" />
          </template>
        </Column>
        <Column field="app_name" header="App Name" />
        <Column field="node_name" header="Node Name" />
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
      <i class="pi pi-send"></i>
      <p>No services found</p>
      <Button label="Deploy First Service" icon="pi pi-plus" @click="openDeployDialog" />
    </div>

    <Sidebar v-model:visible="isDeployOutputDrawerOpen" position="bottom" :style="{ height: '70vh' }">
      <div class="deploy-output-header">
        <h2>Deploy Output</h2>
        <Button text rounded :disabled="!deployOutput.length" @click="copyDeployOutput">
          <i class="pi pi-copy"></i>
        </Button>
      </div>
      <div ref="deployOutputRef" class="deploy-output-body">
        <div v-if="!deployOutput.length" class="text-muted-foreground">No output yet</div>
        <div v-for="(line, i) in deployOutput" :key="i" class="deploy-output-line">{{ line }}</div>
      </div>
    </Sidebar>

    <Sidebar v-model:visible="isDeployDialogOpen" position="right" :style="{ width: '90vw', maxWidth: '1200px' }" class="sheet">
      <template #header>
        <div class="sheet-header">
          <h2>{{ editingServiceId ? "Edit Service" : "Deploy Service" }}</h2>
          <p class="text-muted-foreground">
            {{ editingServiceId ? "Update configuration" : "Follow steps to deploy" }}
          </p>
        </div>
      </template>

      <div class="sheet-body">
        <Stepper :value="currentStep" @update:value="onStepperChange" class="deploy-stepper">
          <StepList>
          <Step :value="1">Application</Step>
          <Step :value="2" :disabled="!selectedApplication">QA</Step>
          <Step :value="3">Node</Step>
        </StepList>
        <StepPanels>
          <StepPanel :value="1">
            <div class="step-content">
              <div v-if="!editingServiceId" class="form-group">
                <label>Server Name</label>
                <InputText v-model="serverName" placeholder="Optional" />
              </div>
              <div v-else class="form-group">
                <label>Server Name</label>
                <p>{{ serverName }}</p>
              </div>
              <div class="form-group">
                <label>Application</label>
                <div v-if="applicationsLoading" class="loading-inline">
                  <i class="pi pi-spin pi-spinner"></i> Loading...
                </div>
                <Select
                  v-else
                  v-model="selectedApplication"
                  :options="applications"
                  option-label="display"
                  placeholder="Select application"
                  @update:model-value="handleApplicationSelect"
                >
                  <template #value="slot">
                    <span v-if="slot.value">{{ slot.value.display || slot.value.name }}</span>
                    <span v-else>{{ slot.placeholder }}</span>
                  </template>
                  <template #option="slot">
                    {{ slot.option.display || slot.option.name }}
                    <span v-if="slot.option.version" class="text-muted-foreground"> ({{ slot.option.version }})</span>
                  </template>
                </Select>
                <p v-if="selectedApplication?.description" class="text-muted-foreground text-sm">{{ selectedApplication.description }}</p>
              </div>
            </div>
          </StepPanel>
          <StepPanel :value="2">
            <div class="step-content">
              <div v-if="!selectedApplication?.qa?.length" class="qa-empty text-muted-foreground">
                No QA configuration required.
              </div>
              <div v-else class="qa-list">
                <div v-for="item in selectedApplication?.qa ?? []" :key="item.name" class="form-group">
                  <label :for="`qa-${item.name}`">
                    {{ item.name }}
                    <span v-if="item.required" class="text-destructive">*</span>
                  </label>
                  <p v-if="item.description" class="text-muted-foreground text-sm">{{ item.description }}</p>
                  <InputText
                    v-if="item.type === 'text'"
                    :id="`qa-${item.name}`"
                    :model-value="getQAValue(item.name)"
                    :placeholder="item.description || `Enter ${item.name}`"
                    @update:model-value="(v: string | undefined) => handleQAValueChange(item.name, v ?? '')"
                  />
                  <Textarea
                    v-else-if="item.type === 'textarea'"
                    :id="`qa-${item.name}`"
                    :model-value="getQAValue(item.name)"
                    :placeholder="item.description || `Enter ${item.name}`"
                    rows="3"
                    @update:model-value="(v: string | undefined) => handleQAValueChange(item.name, v ?? '')"
                  />
                  <InputText
                    v-else-if="item.type === 'number'"
                    :id="`qa-${item.name}`"
                    type="number"
                    :model-value="getQAValue(item.name)"
                    :placeholder="item.description || `Enter ${item.name}`"
                    @update:model-value="(v: string | undefined) => handleQAValueChange(item.name, v ?? '')"
                  />
                  <div v-else-if="item.type === 'boolean'" class="qa-boolean">
                    <Checkbox
                      :id="`qa-${item.name}`"
                      :model-value="getQAValue(item.name) === 'true'"
                      binary
                      @update:model-value="(c: boolean) => handleQAValueChange(item.name, c ? 'true' : 'false')"
                    />
                    <label :for="`qa-${item.name}`">{{ item.description || item.name }}</label>
                  </div>
                  <Select
                    v-else-if="item.type === 'select'"
                    :model-value="getQAValue(item.name)"
                    :options="item.options ?? []"
                    :placeholder="item.description || `Select ${item.name}`"
                    @update:model-value="(v: unknown) => handleQAValueChange(item.name, v != null ? String(v) : '')"
                  />
                </div>
              </div>
            </div>
          </StepPanel>
          <StepPanel :value="3">
            <div class="step-content">
              <div class="form-group">
                <label>Node</label>
                <div v-if="nodesLoading" class="loading-inline">
                  <i class="pi pi-spin pi-spinner"></i> Loading...
                </div>
                <Select
                  v-else
                  v-model="selectedNode"
                  :options="nodes"
                  option-label="display_name"
                  placeholder="Select node"
                  @update:model-value="handleNodeSelect"
                >
                  <template #value="slot">
                    <span v-if="slot.value">{{ slot.value.display_name || slot.value.name }} ({{ slot.value.address }})</span>
                    <span v-else>{{ slot.placeholder }}</span>
                  </template>
                  <template #option="slot">
                    {{ slot.option.display_name || slot.option.name }} ({{ slot.option.address }})
                  </template>
                </Select>
                <p v-if="selectedNode" class="text-muted-foreground text-sm">
                  {{ selectedNode.display_name || selectedNode.name }} – {{ selectedNode.address }}
                </p>
              </div>
            </div>
          </StepPanel>
        </StepPanels>
        </Stepper>
      </div>

      <template #footer>
        <div class="sheet-footer">
          <Button outlined @click="closeDeployDialog" :disabled="deployLoading">Cancel</Button>
          <Button v-if="currentStep > 1" label="Previous" icon="pi pi-chevron-left" @click="goToStep(currentStep - 1)" :disabled="deployLoading" />
          <Button v-if="currentStep < 3" label="Next" icon="pi pi-chevron-right" icon-pos="right" @click="goToStep(currentStep + 1)" :disabled="!canProceedToNextStep" />
          <Button 
            v-if="currentStep === 3" 
            :label="editingServiceId ? 'Update' : 'Save'" 
            icon="pi pi-save" 
            @click="() => deployService({ onlySave: true })"
            :disabled="!canProceedToNextStep || deployLoading" 
          />
          <Button 
            v-if="currentStep === 3" 
            :label="deployLoading ? 'Deploying...' : (editingServiceId ? 'Update & Deploy' : 'Save & Deploy')" 
            icon="pi pi-upload" 
            @click="() => deployService()"
            :disabled="!canProceedToNextStep || deployLoading" 
            :loading="deployLoading" 
          />
        </div>
      </template>
    </Sidebar>

    <Dialog v-model:visible="isDeleteDialogOpen" modal header="Delete Service" :style="{ width: '425px' }">
      <p>Are you sure you want to delete "{{ serviceToDelete?.name }}"? This cannot be undone.</p>
      <div class="form-group flex-row">
        <Checkbox v-model="forceDelete" binary inputId="force-delete" />
        <label for="force-delete">Force delete (undeploy if deployed)</label>
      </div>
      <template #footer>
        <Button outlined @click="isDeleteDialogOpen = false">Cancel</Button>
        <Button severity="danger" @click="confirmDeleteService" :disabled="loading">
          {{ loading ? "Deleting..." : "Delete" }}
        </Button>
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
.action-btns { display: flex; gap: 0.5rem; }
.pagination-bar { display: flex; justify-content: space-between; align-items: center; padding: 1rem; border-top: 1px solid var(--p-surface-border); font-size: 0.875rem; }
.pagination-btns { display: flex; gap: 0.5rem; }
.pagination-active { background: var(--p-primary-color) !important; color: var(--p-primary-contrast-color) !important; border-color: var(--p-primary-color) !important; }

.empty-state { padding: 3rem; text-align: center; background: var(--p-surface-card); border-radius: var(--p-border-radius); }
.empty-state i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
.empty-state .p-button { margin-top: 1rem; }

.deploy-output-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
.deploy-output-header h2 { font-size: 1.125rem; font-weight: 600; }
.deploy-output-body { max-height: 60vh; overflow-y: auto; padding: 1rem; background: var(--p-surface-50); border-radius: var(--p-border-radius); font-family: ui-monospace, monospace; font-size: 0.75rem; }
.deploy-output-line { white-space: pre-wrap; word-break: break-word; margin-bottom: 0.25rem; }

.sheet-header h2 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.sheet-header p { font-size: 0.875rem; }
.sheet-body { flex: 1; overflow-y: auto; padding: 1rem; display: flex; flex-direction: column; gap: 1rem; }
.sheet-footer { display: flex; justify-content: flex-end; gap: 0.5rem; padding: 1rem; border-top: 1px solid var(--p-surface-border); }

.deploy-stepper { margin-bottom: 1rem; }
.step-content { display: flex; flex-direction: column; gap: 1rem; }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; }
.form-group.flex-row { flex-direction: row; align-items: center; gap: 0.5rem; }
.text-sm { font-size: 0.75rem; }
.qa-empty { padding: 1.5rem; text-align: center; border: 1px dashed var(--p-surface-border); border-radius: var(--p-border-radius); background: var(--p-surface-50); }
.qa-list { display: flex; flex-direction: column; gap: 1rem; }
.qa-boolean { display: flex; align-items: center; gap: 0.5rem; }
.loading-inline { display: flex; align-items: center; gap: 0.5rem; color: var(--p-text-muted-color); }
</style>
