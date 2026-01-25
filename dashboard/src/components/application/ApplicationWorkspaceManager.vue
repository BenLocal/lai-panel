<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import { ApiResponseHelper } from "@/api/base";
import {
  workspaceApi,
  uploadWorkspaceFile,
  type WorkspaceEntry,
} from "@/api/workspace";
import { showToast } from "@/lib/toast";
import MonacoEditor from "@guolao/vue-monaco-editor";

interface Props {
  appName?: string;
}

const props = defineProps<Props>();

const entries = ref<WorkspaceEntry[]>([]);
const loading = ref(false);
const currentPath = ref("");
const selectedEntry = ref<WorkspaceEntry | null>(null);
const editorContent = ref("");
const editorPath = ref("");
const isSaving = ref(false);
const isDeleting = ref(false);
const isUploading = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);
const isCreateFileDialogOpen = ref(false);
const newFileName = ref("");
const isCreatingFile = ref(false);
const isCreateDirDialogOpen = ref(false);
const newDirName = ref("");
const isCreatingDir = ref(false);
const isDeleteDialogOpen = ref(false);
const editorOptions = {
  automaticLayout: true,
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  fontSize: 14,
};

const hasWorkspace = computed(() => Boolean(props.appName?.trim()));

const breadcrumbs = computed(() => {
  const segments = currentPath.value
    ? currentPath.value.split("/").filter(Boolean)
    : [];

  const crumbs = [
    {
      label: props.appName ? `/${props.appName}` : "/workspace",
      path: "",
    },
  ];

  let cumulative = "";
  for (const segment of segments) {
    cumulative = cumulative ? `${cumulative}/${segment}` : segment;
    crumbs.push({
      label: segment,
      path: cumulative,
    });
  }

  return crumbs;
});

const loadEntries = async (path = "") => {
  if (!hasWorkspace.value || !props.appName) {
    entries.value = [];
    currentPath.value = "";
    selectedEntry.value = null;
    editorContent.value = "";
    editorPath.value = "";
    return;
  }

  loading.value = true;
  const response = await workspaceApi.list(props.appName, path);
  loading.value = false;

  if (!ApiResponseHelper.isSuccess(response)) {
    showToast(response.message ?? "Failed to load workspace", "error");
    return;
  }

  currentPath.value = response.data?.currentPath ?? "";
  entries.value = response.data?.entries ?? [];
  selectedEntry.value = null;
  editorContent.value = "";
  editorPath.value = "";
};

const refreshEntries = () => loadEntries(currentPath.value);

const selectEntry = (entry: WorkspaceEntry) => {
  selectedEntry.value = entry;

  if (entry.is_dir) {
    editorPath.value = "";
    editorContent.value = "";
    return;
  }

  loadFile(entry);
};

const enterEntry = (entry: WorkspaceEntry) => {
  if (entry.is_dir) {
    loadEntries(entry.path);
  } else {
    loadFile(entry);
  }
};

const loadFile = async (entry: WorkspaceEntry) => {
  if (!props.appName) return;

  const response = await workspaceApi.read(props.appName, entry.path);
  if (!ApiResponseHelper.isSuccess(response)) {
    showToast(response.message ?? "Failed to open file", "error");
    return;
  }

  selectedEntry.value = entry;
  editorPath.value = entry.path;
  editorContent.value = response.data?.content ?? "";
};

const saveFile = async () => {
  if (!props.appName || !editorPath.value) return;

  isSaving.value = true;
  const response = await workspaceApi.save(
    props.appName,
    editorPath.value,
    editorContent.value
  );
  isSaving.value = false;

  if (!ApiResponseHelper.isSuccess(response)) {
    showToast(response.message ?? "Failed to save file", "error");
    return;
  }

  showToast("File saved", "success");
  refreshEntries();
};

const joinPath = (base: string, name: string) => {
  if (!base) return name;
  return `${base.replace(/\/$/, "")}/${name}`;
};

const openCreateFileDialog = () => {
  newFileName.value = "";
  isCreateFileDialogOpen.value = true;
};

const confirmCreateFile = async () => {
  if (!props.appName) {
    showToast("Set an application name first", "error");
    return;
  }
  const filename = newFileName.value.trim();
  if (!filename) {
    showToast("File name is required", "error");
    return;
  }

  const targetPath = joinPath(currentPath.value, filename);
  isCreatingFile.value = true;
  const response = await workspaceApi.save(props.appName, targetPath, "");
  isCreatingFile.value = false;

  if (!ApiResponseHelper.isSuccess(response)) {
    showToast(response.message ?? "Failed to create file", "error");
    return;
  }

  showToast("File created", "success");
  isCreateFileDialogOpen.value = false;
  await loadEntries(currentPath.value);
  const entry = entries.value.find((item) => item.path === targetPath);
  if (entry) {
    loadFile(entry);
  }
};

const openCreateDirDialog = () => {
  newDirName.value = "";
  isCreateDirDialogOpen.value = true;
};

const confirmCreateDir = async () => {
  if (!props.appName) {
    showToast("Set an application name first", "error");
    return;
  }
  const dirname = newDirName.value.trim();
  if (!dirname) {
    showToast("Directory name is required", "error");
    return;
  }

  const targetPath = joinPath(currentPath.value, dirname);
  isCreatingDir.value = true;
  const response = await workspaceApi.mkdir(props.appName, targetPath);
  isCreatingDir.value = false;

  if (!ApiResponseHelper.isSuccess(response)) {
    showToast(response.message ?? "Failed to create directory", "error");
    return;
  }

  showToast("Directory created", "success");
  isCreateDirDialogOpen.value = false;
  loadEntries(currentPath.value);
};

const requestDeleteSelection = () => {
  if (!selectedEntry.value) {
    showToast("Select an entry to delete", "error");
    return;
  }
  isDeleteDialogOpen.value = true;
};

const deleteSelection = async () => {
  if (!props.appName || !selectedEntry.value) return;

  isDeleting.value = true;
  const response = await workspaceApi.remove(
    props.appName,
    selectedEntry.value.path
  );
  isDeleting.value = false;

  if (!ApiResponseHelper.isSuccess(response)) {
    showToast(response.message ?? "Failed to delete entry", "error");
    return;
  }

  showToast("Entry deleted", "success");
  selectedEntry.value = null;
  editorContent.value = "";
  editorPath.value = "";
  isDeleteDialogOpen.value = false;
  loadEntries(currentPath.value);
};

const triggerUpload = () => {
  if (!hasWorkspace.value) return;
  fileInputRef.value?.click();
};

const handleFileInputChange = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  target.value = "";
  if (!file || !props.appName) return;

  isUploading.value = true;
  const response = await uploadWorkspaceFile(
    props.appName,
    currentPath.value,
    file
  );
  isUploading.value = false;

  if (!ApiResponseHelper.isSuccess(response)) {
    showToast(response.message ?? "Failed to upload file", "error");
    return;
  }

  showToast("File uploaded", "success");
  loadEntries(currentPath.value);
};

const handleBreadcrumbClick = (path: string) => {
  loadEntries(path);
};

const formatBytes = (bytes: number) => {
  if (!bytes) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const index = Math.min(
    Math.floor(Math.log(bytes) / Math.log(1024)),
    units.length - 1
  );
  const value = bytes / Math.pow(1024, index);
  return `${value.toFixed(value >= 10 || index === 0 ? 0 : 1)} ${units[index]}`;
};

const isFileSelected = computed(() =>
  Boolean(selectedEntry.value && !selectedEntry.value.is_dir)
);

const detectEditorLanguage = (path: string): string => {
  const lower = path.toLowerCase();
  if (lower.endsWith(".json")) return "json";
  if (lower.endsWith(".yaml") || lower.endsWith(".yml")) return "yaml";
  if (lower.endsWith(".ts")) return "typescript";
  if (lower.endsWith(".js")) return "javascript";
  if (lower.endsWith(".sh")) return "shell";
  if (lower.endsWith(".md")) return "markdown";
  if (lower.endsWith(".vue")) return "vue";
  if (lower.endsWith(".css")) return "css";
  if (lower.endsWith(".html")) return "html";
  return "plaintext";
};

const editorLanguage = computed(() =>
  editorPath.value ? detectEditorLanguage(editorPath.value) : "plaintext"
);

const resetState = () => {
  entries.value = [];
  currentPath.value = "";
  selectedEntry.value = null;
  editorContent.value = "";
  editorPath.value = "";
};

watch(
  () => props.appName,
  () => {
    resetState();
    if (hasWorkspace.value) {
      loadEntries();
    }
  },
  { immediate: true }
);

watch(isCreateFileDialogOpen, (open) => {
  if (!open) {
    newFileName.value = "";
    isCreatingFile.value = false;
  }
});

watch(isCreateDirDialogOpen, (open) => {
  if (!open) {
    newDirName.value = "";
    isCreatingDir.value = false;
  }
});

onMounted(() => {
  if (hasWorkspace.value) {
    loadEntries();
  }
});
</script>

<template>
  <div class="workspace-root">
    <div class="workspace-header">
      <div>
        <p class="workspace-title">Workspace Files</p>
        <p class="text-muted-foreground">Manage files under the workspace for this application.</p>
      </div>
      <div class="workspace-actions">
        <Button outlined size="small" :disabled="!hasWorkspace || isUploading" @click="triggerUpload">
          <i class="pi pi-upload"></i><span class="btn-icon-text">Upload</span>
        </Button>
        <Button outlined size="small" @click="openCreateFileDialog">
          <i class="pi pi-file-plus"></i><span class="btn-icon-text">New File</span>
        </Button>
        <Button outlined size="small" @click="openCreateDirDialog">
          <i class="pi pi-folder-plus"></i><span class="btn-icon-text">New Folder</span>
        </Button>
      </div>
    </div>

    <input ref="fileInputRef" type="file" class="hidden" @change="handleFileInputChange" />

    <div v-if="!hasWorkspace" class="workspace-empty text-muted-foreground">
      Set an application name to access its workspace folder.
    </div>

    <div v-else class="workspace-grid">
      <div class="workspace-panel">
        <div class="breadcrumbs">
          <button
            v-for="crumb in breadcrumbs"
            :key="crumb.path || 'root'"
            class="breadcrumb-btn"
            @click="handleBreadcrumbClick(crumb.path)"
          >
            <i class="pi pi-folder"></i> {{ crumb.label || "/" }}
          </button>
        </div>
        <div class="workspace-toolbar">
          <span class="text-muted-foreground">{{ currentPath ? `/${currentPath}` : "/" }}</span>
          <div class="workspace-toolbar-actions">
            <Button text size="small" @click="refreshEntries">Refresh</Button>
            <Button text size="small" :disabled="!selectedEntry || isDeleting" @click="requestDeleteSelection">Delete</Button>
          </div>
        </div>
        <div class="workspace-list-wrap">
          <div v-if="loading" class="workspace-list-msg text-muted-foreground">Loading...</div>
          <div v-else-if="!entries.length" class="workspace-list-msg text-muted-foreground">Folder is empty.</div>
          <ul v-else class="workspace-list">
            <li
              v-for="entry in entries"
              :key="entry.path || entry.name"
              class="workspace-item"
              :class="{ selected: selectedEntry?.path === entry.path }"
              @click="selectEntry(entry)"
              @dblclick.prevent="enterEntry(entry)"
            >
              <div class="workspace-item-main">
                <div class="workspace-item-icon">
                  <i :class="entry.is_dir ? 'pi pi-folder' : 'pi pi-file'"></i>
                </div>
                <div>
                  <p class="workspace-item-name">{{ entry.name }}</p>
                  <p class="text-muted-foreground">{{ entry.is_dir ? "Folder" : formatBytes(entry.size) }}</p>
                </div>
              </div>
              <div class="workspace-item-meta">
                <span>{{ new Date(entry.mod_time).toLocaleString() }}</span>
                <Button v-if="entry.is_dir" text rounded size="small" @click.stop="enterEntry(entry)">
                  <i class="pi pi-arrow-down-right"></i>
                </Button>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <div class="workspace-panel">
        <div class="editor-header">
          <div>
            <p class="workspace-title">Editor</p>
            <p class="text-muted-foreground">{{ editorPath || "Select a file to preview and edit" }}</p>
          </div>
          <Button size="small" :disabled="!isFileSelected || isSaving" @click="saveFile">Save</Button>
        </div>
        <div class="workspace-editor-wrap">
          <div v-if="!isFileSelected" class="workspace-editor-placeholder text-muted-foreground">
            Select a file from the list to start editing.
          </div>
          <MonacoEditor
            v-else
            class="workspace-monaco"
            v-model:value="editorContent"
            theme="vs-dark"
            :language="editorLanguage"
            :options="editorOptions"
          />
        </div>
      </div>
    </div>
  </div>
  <Dialog v-model:visible="isCreateFileDialogOpen" modal header="Create New File" :style="{ width: '425px' }">
    <p class="dialog-desc">Enter a file name to create it inside the current workspace directory.</p>
    <div class="form-group">
      <label for="workspace-new-file">File Name</label>
      <InputText id="workspace-new-file" v-model="newFileName" placeholder="e.g. config.yaml" :disabled="isCreatingFile" />
    </div>
    <template #footer>
      <Button outlined @click="isCreateFileDialogOpen = false" :disabled="isCreatingFile">Cancel</Button>
      <Button @click="confirmCreateFile" :disabled="isCreatingFile">{{ isCreatingFile ? "Creating..." : "Create" }}</Button>
    </template>
  </Dialog>
  <Dialog v-model:visible="isCreateDirDialogOpen" modal header="Create New Folder" :style="{ width: '425px' }">
    <p class="dialog-desc">Enter a folder name to create it inside the current workspace directory.</p>
    <div class="form-group">
      <label for="workspace-new-dir">Folder Name</label>
      <InputText id="workspace-new-dir" v-model="newDirName" placeholder="e.g. configs" :disabled="isCreatingDir" />
    </div>
    <template #footer>
      <Button outlined @click="isCreateDirDialogOpen = false" :disabled="isCreatingDir">Cancel</Button>
      <Button @click="confirmCreateDir" :disabled="isCreatingDir">{{ isCreatingDir ? "Creating..." : "Create" }}</Button>
    </template>
  </Dialog>
  <Dialog v-model:visible="isDeleteDialogOpen" modal header="Delete Entry" :style="{ width: '425px' }">
    <p>Permanently remove <span class="font-medium">{{ selectedEntry?.name }}</span> from the workspace?</p>
    <template #footer>
      <Button outlined @click="isDeleteDialogOpen = false" :disabled="isDeleting">Cancel</Button>
      <Button severity="danger" @click="deleteSelection" :disabled="isDeleting">{{ isDeleting ? "Deleting..." : "Delete" }}</Button>
    </template>
  </Dialog>
</template>

<style scoped>
.workspace-root { display: flex; flex-direction: column; gap: 1rem; height: 100%; padding: 1rem; border: 1px solid var(--p-surface-border); border-radius: var(--p-border-radius); }
.workspace-header { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 0.5rem; }
.workspace-title { font-size: 0.875rem; font-weight: 500; margin-bottom: 0.25rem; }
.workspace-actions { display: flex; gap: 0.5rem; }
.btn-icon-text { margin-left: 0.5rem; }
.workspace-empty { flex: 1; display: flex; align-items: center; justify-content: center; padding: 1rem; border: 1px dashed var(--p-surface-border); border-radius: var(--p-border-radius); font-size: 0.875rem; }
.workspace-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; flex: 1; min-height: 0; }
@media (max-width: 1024px) { .workspace-grid { grid-template-columns: 1fr; } }
.workspace-panel { display: flex; flex-direction: column; gap: 0.75rem; min-height: 0; padding: 0.75rem; border: 1px solid var(--p-surface-border); border-radius: var(--p-border-radius); }
.breadcrumbs { display: flex; flex-wrap: wrap; gap: 0.5rem; font-size: 0.75rem; }
.breadcrumb-btn { display: inline-flex; align-items: center; gap: 0.25rem; padding: 0.25rem 0.5rem; border-radius: var(--p-border-radius); background: transparent; border: none; cursor: pointer; color: var(--p-text-muted-color); }
.breadcrumb-btn:hover { background: var(--p-surface-hover); color: var(--p-text-color); }
.workspace-toolbar { display: flex; align-items: center; justify-content: space-between; font-size: 0.75rem; }
.workspace-toolbar-actions { display: flex; gap: 0.5rem; }
.workspace-list-wrap { flex: 1; display: flex; flex-direction: column; min-height: 0; border: 1px solid var(--p-surface-border); border-radius: var(--p-border-radius); background: var(--p-surface-50); }
.workspace-list-msg { padding: 1rem; font-size: 0.875rem; }
.workspace-list { flex: 1; overflow: auto; list-style: none; margin: 0; padding: 0; }
.workspace-item { display: flex; align-items: center; justify-content: space-between; padding: 0.5rem 0.75rem; cursor: pointer; transition: background 0.2s; }
.workspace-item:hover { background: var(--p-surface-hover); }
.workspace-item.selected { background: var(--p-surface-100); }
.workspace-item-main { display: flex; align-items: center; gap: 0.75rem; }
.workspace-item-icon { width: 2rem; height: 2rem; display: flex; align-items: center; justify-content: center; border-radius: var(--p-border-radius); background: var(--p-primary-color); opacity: 0.15; }
.workspace-item-icon i { color: var(--p-primary-color); opacity: 1; }
.workspace-item-name { font-size: 0.875rem; font-weight: 500; }
.workspace-item-meta { display: flex; align-items: center; gap: 0.5rem; font-size: 0.75rem; }
.editor-header { display: flex; align-items: center; justify-content: space-between; }
.workspace-editor-wrap { flex: 1; display: flex; flex-direction: column; min-height: 0; border: 1px solid var(--p-surface-border); border-radius: var(--p-border-radius); background: var(--p-surface-0); }
.workspace-editor-placeholder { flex: 1; display: flex; align-items: center; justify-content: center; padding: 1.5rem; text-align: center; font-size: 0.875rem; }
.workspace-monaco { height: 100%; min-height: 280px; }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; }
.dialog-desc { margin-bottom: 1rem; }
.font-medium { font-weight: 500; }
.hidden { position: absolute; width: 0; height: 0; opacity: 0; pointer-events: none; }
</style>
