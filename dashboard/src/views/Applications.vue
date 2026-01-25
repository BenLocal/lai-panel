<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Sidebar from 'primevue/sidebar'
import Dialog from 'primevue/dialog'
import Menu from 'primevue/menu'
import {
  applicationApi,
  type Application,
  type ApplicationQAItem,
} from "@/api/application";
import ApplicationQAEditor from "@/components/application/ApplicationQAEditor.vue";
import ApplicationWorkspaceManager from "@/components/application/ApplicationWorkspaceManager.vue";
import YamlEditor from "@/components/application/YamlEditor.vue";
import { ApiResponseHelper } from "@/api/base";
import { showToast } from "@/lib/toast";

interface ApplicationForm {
  name: string;
  version?: string;
  display?: string;
  description?: string;
  icon?: string;
  qa: ApplicationQAItem[];
  dockerCompose: string;
  static_path?: string;
}

const applications = ref<Application[]>([]);
const currentPage = ref(1);
const pageSize = ref(6);
const totalPages = ref(1);
const isSheetOpen = ref(false);
const isEditMode = ref(false);
const loading = ref(false);
const editingApplicationId = ref<number | null>(null);
const isComposeEditorOpen = ref(false);
const composeDraft = ref("");
const isWorkspaceDialogOpen = ref(false);
const workspaceDialogAppName = ref("");
const workspaceDialogDisplayName = ref("");
const menuRefs = ref<Record<number, any>>({});

const createDefaultForm = (): ApplicationForm => ({
  display: "", name: "", description: "", version: "", icon: "pi-th-large",
  qa: [], dockerCompose: "", static_path: "",
});

/** Only use icon if it's a valid PrimeIcon (pi-xxx), else fallback. */
const validAppIcon = (icon?: string | null) => {
  const s = (icon ?? "").trim();
  return s.startsWith("pi-") ? s : "pi-th-large";
};

const formData = reactive<ApplicationForm>(createDefaultForm());

const goToPage = (p: number) => {
  if (p >= 1 && p <= totalPages.value) currentPage.value = p;
};

const fetchApplications = async () => {
  const res = await applicationApi.page(currentPage.value, pageSize.value);
  if (!ApiResponseHelper.isSuccess(res)) return;
  const d = res.data!;
  applications.value = d.apps ?? [];
  totalPages.value = Math.ceil((d.total ?? 0) / pageSize.value);
  currentPage.value = d.currentPage ?? 1;
  pageSize.value = d.pageSize ?? 6;
};

const dockerComposePreview = computed(() => {
  const c = formData.dockerCompose?.trim();
  if (!c) return "";
  const lines = c.split("\n");
  const s = lines.slice(0, 6).join("\n");
  return lines.length > 6 ? `${s}\n...` : s;
});

const namePattern = /^[A-Za-z]*$/;
const isNameValid = computed(() => namePattern.test(formData.name));

const handleNameInput = (e: Event) => {
  const t = e.target as HTMLInputElement;
  const s = (t.value.match(/[A-Za-z]/g) ?? []).join("");
  if (s !== t.value) { t.value = s; }
  formData.name = s;
};

const isSaveDisabled = computed(() =>
  !formData.name.trim() || !isNameValid.value || loading.value);

const resetForm = () => {
  Object.assign(formData, createDefaultForm());
  editingApplicationId.value = null;
  isComposeEditorOpen.value = false;
  composeDraft.value = "";
};

const openAddApplicationDialog = () => {
  isEditMode.value = false;
  resetForm();
  isSheetOpen.value = true;
};

const openEditApplicationDialog = (app: Application) => {
  isEditMode.value = true;
  editingApplicationId.value = app.id;
  Object.assign(formData, {
    name: app.name ?? "",
    description: app.description ?? "",
    version: app.version ?? "",
    icon: app.icon ?? "pi-th-large",
    display: app.display ?? "",
    qa: app.qa ? app.qa.map((i) => ({ ...i, options: i.options ? [...i.options] : undefined })) : [],
    dockerCompose: app.docker_compose ?? "",
    static_path: app.static_path ?? "",
  });
  isComposeEditorOpen.value = false;
  composeDraft.value = app.docker_compose ?? "";
  isSheetOpen.value = true;
};

const openComposeEditor = () => {
  composeDraft.value = formData.dockerCompose ?? "";
  isComposeEditorOpen.value = true;
};

const cancelComposeEdit = () => {
  isComposeEditorOpen.value = false;
  composeDraft.value = formData.dockerCompose ?? "";
};

const confirmComposeEdit = () => {
  formData.dockerCompose = composeDraft.value ?? "";
  isComposeEditorOpen.value = false;
};

const openWorkspace = (app: Application) => {
  if (!app?.name?.trim()) { showToast("Workspace path unavailable", "error"); return; }
  workspaceDialogAppName.value = app.name;
  workspaceDialogDisplayName.value = app.display || app.name;
  isWorkspaceDialogOpen.value = true;
};

const handleCancel = () => {
  isSheetOpen.value = false;
  isEditMode.value = false;
  resetForm();
};

const saveApplication = async () => {
  if (isSaveDisabled.value) return;
  loading.value = true;
  const payload: Application = {
    id: editingApplicationId.value ?? 0,
    name: formData.name.trim(),
    display: formData.display?.trim() ?? "",
    description: formData.description?.trim() ?? "",
    version: formData.version?.trim() ?? "",
    icon: formData.icon?.trim() ?? "",
    qa: formData.qa,
    docker_compose: formData.dockerCompose?.trim() ?? "",
    static_path: formData.static_path?.trim() ?? "",
  };
  try {
    if (isEditMode.value && editingApplicationId.value != null) {
      await applicationApi.update(payload);
    } else {
      await applicationApi.add(payload);
    }
    await fetchApplications();
    isSheetOpen.value = false;
    resetForm();
  } catch (e) {
    console.error(e);
    showToast("Failed to save application", "error");
  } finally {
    loading.value = false;
  }
};

onMounted(fetchApplications);
</script>

<template>
  <div class="page-root">
    <div class="page-header page-header-row">
      <div>
        <h1>Applications</h1>
        <p class="text-muted-foreground">Manage and monitor your applications</p>
      </div>
      <Button label="New Application" icon="pi pi-plus" @click="openAddApplicationDialog" />
    </div>

    <div v-if="applications.length > 0">
      <div class="app-grid">
        <div
          v-for="app in applications"
          :key="app.id"
          class="app-card"
          @click="openEditApplicationDialog(app)"
        >
          <div class="app-card-head">
            <div class="app-card-title">
              <div class="app-card-icon">
                <i :class="'pi ' + (validAppIcon(app.icon))"></i>
              </div>
              <div>
                <h3>{{ app.name }}</h3>
                <p class="text-muted-foreground">{{ app.version }}</p>
              </div>
            </div>
          </div>
          <p class="app-card-desc">{{ app.description }}</p>
          <div class="app-card-actions">
            <Menu
              :ref="(el: any) => { if (el) menuRefs[app.id] = el }"
              :model="[{ label: 'Open Workspace', icon: 'pi pi-folder-open', command: () => openWorkspace(app) }]"
              popup
            />
            <Button text rounded size="small" @click.stop="(e) => menuRefs[app.id]?.toggle(e)">
              <i class="pi pi-ellipsis-h"></i>
            </Button>
          </div>
        </div>
      </div>

      <div v-if="totalPages > 1" class="pagination-bar">
        <span class="text-muted-foreground">
          Showing {{ (currentPage - 1) * pageSize + 1 }}–{{ Math.min(currentPage * pageSize, applications.length) }} of {{ applications.length }}
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
      <i class="pi pi-th-large"></i>
      <p>No applications found</p>
      <Button label="Add First Application" icon="pi pi-plus" @click="openAddApplicationDialog" />
    </div>

    <Sidebar v-model:visible="isSheetOpen" position="right" :style="{ width: '90vw', maxWidth: '1200px' }" class="app-sheet">
      <div class="sheet-header">
        <h2>{{ isEditMode ? "Edit Application" : "Add Application" }}</h2>
        <p class="text-muted-foreground">
          {{ isEditMode ? "Update application information" : "Fill in the application details" }}
        </p>
      </div>
      <div class="sheet-body">
        <div class="form-group">
          <label for="app-name">Name *</label>
          <InputText id="app-name" v-model="formData.name" placeholder="Application name, English letters only" @input="handleNameInput" />
          <p v-if="formData.name && !isNameValid" class="text-destructive">Only English letters (A–Z) are allowed.</p>
        </div>
        <div class="form-group">
          <label for="app-display">Display Name</label>
          <InputText id="app-display" v-model="formData.display" placeholder="Display name (optional)" />
        </div>
        <div class="form-group">
          <label for="app-version">Version</label>
          <InputText id="app-version" v-model="formData.version" placeholder="v1.0.0" />
        </div>
        <div class="form-group">
          <label for="app-icon">Icon</label>
          <InputText id="app-icon" v-model="formData.icon" placeholder="pi-th-large" />
        </div>
        <div class="form-group">
          <label for="app-description">Description</label>
          <Textarea id="app-description" v-model="formData.description" placeholder="Describe the application" rows="4" />
        </div>
        <div class="form-group">
          <label for="app-static-path">Static Path</label>
          <InputText id="app-static-path" v-model="formData.static_path" placeholder="URL or path to installer" />
          <p class="text-muted-foreground text-sm">Path to installer file. Tar.gz will be extracted.</p>
        </div>
        <hr class="form-divider" />
        <div class="form-group">
          <div class="form-group-head">
            <span class="font-medium">Docker Compose (YAML)</span>
            <span class="text-muted-foreground text-sm">Edit and preview</span>
          </div>
          <div class="compose-preview-block">
            <pre v-if="dockerComposePreview" class="compose-preview">{{ dockerComposePreview }}</pre>
            <p v-else class="text-muted-foreground text-sm">No Docker Compose definition yet.</p>
            <div class="compose-actions">
              <Button size="small" outlined @click="openComposeEditor">
                <i class="pi pi-pencil"></i>
                <span class="btn-icon-text">Open Editor</span>
              </Button>
              <Button v-if="formData.dockerCompose" text rounded size="small" :disabled="loading" @click="formData.dockerCompose = ''">
                <i class="pi pi-trash"></i>
                <span class="btn-icon-text">Clear</span>
              </Button>
            </div>
          </div>
        </div>
        <hr class="form-divider" />
        <div class="form-group">
          <span class="font-medium">QA Configuration</span>
          <ApplicationQAEditor v-model="formData.qa" />
        </div>
        <hr class="form-divider" />
        <ApplicationWorkspaceManager :app-name="formData.name" />
      </div>
      <div class="sheet-footer">
        <Button outlined @click="handleCancel" :disabled="loading">Cancel</Button>
        <Button @click="saveApplication" :disabled="isSaveDisabled">
          {{ loading ? "Saving..." : isEditMode ? "Update Application" : "Add Application" }}
        </Button>
      </div>
    </Sidebar>

    <Dialog v-model:visible="isComposeEditorOpen" modal :style="{ width: '95vw', height: '100vh', maxWidth: 'none' }" :contentStyle="{ display: 'flex', flexDirection: 'column', padding: 0 }">
      <template #header>
        <div class="dialog-custom-header">
          <h3>Docker Compose Editor</h3>
          <p class="text-muted-foreground">Edit and preview the Docker Compose YAML.</p>
        </div>
      </template>
      <div class="compose-editor-wrap">
        <YamlEditor v-model="composeDraft" :height="'100%'" />
      </div>
      <template #footer>
        <div class="dialog-custom-footer">
          <Button outlined @click="cancelComposeEdit">Cancel</Button>
          <Button @click="confirmComposeEdit">Confirm</Button>
        </div>
      </template>
    </Dialog>

    <Dialog v-model:visible="isWorkspaceDialogOpen" modal :style="{ width: '95vw', height: '100vh', maxWidth: 'none' }" :contentStyle="{ display: 'flex', flexDirection: 'column', padding: 0 }">
      <template #header>
        <div class="dialog-custom-header">
          <h3>Workspace · {{ workspaceDialogDisplayName }}</h3>
          <p class="text-muted-foreground">Manage files under the workspace for this application.</p>
        </div>
      </template>
      <div class="workspace-dialog-body">
        <ApplicationWorkspaceManager :app-name="workspaceDialogAppName" />
      </div>
      <template #footer>
        <Button outlined @click="isWorkspaceDialogOpen = false">Close</Button>
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

.app-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

.app-card {
  padding: 1.25rem;
  background: var(--p-surface-card);
  border-radius: var(--p-border-radius);
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.app-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.08); }

.app-card-head { margin-bottom: 1rem; }
.app-card-title { display: flex; align-items: flex-start; gap: 0.75rem; }
.app-card-icon {
  width: 2.5rem; height: 2.5rem;
  border-radius: var(--p-border-radius);
  background: var(--p-primary-color);
  opacity: 0.15;
  display: flex; align-items: center; justify-content: center;
}
.app-card-icon i { color: var(--p-primary-color); opacity: 1; }
.app-card-title h3 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.app-card-title p { font-size: 0.75rem; }
.app-card-desc {
  font-size: 0.875rem;
  color: var(--p-text-muted-color);
  margin-bottom: 1rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.app-card-actions { display: flex; justify-content: flex-end; }

.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 1.5rem;
  margin-top: 1.5rem;
  border-top: 1px solid var(--p-surface-border);
  font-size: 0.875rem;
}
.pagination-btns { display: flex; align-items: center; gap: 0.5rem; }
.pagination-active { background: var(--p-primary-color) !important; color: var(--p-primary-contrast-color) !important; border-color: var(--p-primary-color) !important; }

.empty-state {
  padding: 3rem;
  text-align: center;
  background: var(--p-surface-card);
  border-radius: var(--p-border-radius);
}
.empty-state i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
.empty-state .p-button { margin-top: 1rem; }

.app-sheet { display: flex; flex-direction: column; }
.sheet-header { padding: 0 1rem 1rem; }
.sheet-header h2 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.sheet-header p { font-size: 0.875rem; }
.sheet-body { flex: 1; overflow-y: auto; padding: 0 1rem; display: flex; flex-direction: column; gap: 1rem; }
.sheet-footer { display: flex; justify-content: flex-end; gap: 0.5rem; padding: 1rem; border-top: 1px solid var(--p-surface-border); }

.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; }
.form-group-head { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 0.5rem; }
.form-divider { border: none; border-top: 1px solid var(--p-surface-border); margin: 0.5rem 0; }
.text-sm { font-size: 0.75rem; }
.font-medium { font-weight: 500; }

.compose-preview-block {
  padding: 1rem;
  border-radius: var(--p-border-radius);
  border: 1px dashed var(--p-surface-border);
  background: var(--p-surface-50);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.compose-preview {
  max-height: 12rem;
  overflow: auto;
  font-size: 0.75rem;
  padding: 0.75rem;
  background: var(--p-surface-0);
  border-radius: var(--p-border-radius);
  white-space: pre-wrap;
  word-break: break-word;
}
.compose-actions { display: flex; gap: 0.5rem; flex-wrap: wrap; }

.dialog-custom-header { padding: 1rem 1.5rem; }
.dialog-custom-header h3 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.dialog-custom-header p { font-size: 0.875rem; }
.compose-editor-wrap { flex: 1; min-height: 0; padding: 0 1.5rem; }
.dialog-custom-footer { display: flex; gap: 0.5rem; padding: 1rem 1.5rem; border-top: 1px solid var(--p-surface-border); }
.workspace-dialog-body { flex: 1; overflow: auto; padding: 1rem 1.5rem; }
</style>
