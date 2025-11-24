<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
  <div class="flex h-full flex-col gap-4 rounded-lg border p-4">
    <div
      class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <p class="text-sm font-medium">Workspace Files</p>
        <p class="text-xs text-muted-foreground">
          Manage files stored under the workspace directory for this
          application.
        </p>
      </div>
      <div class="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          :disabled="!hasWorkspace || isUploading"
          @click="triggerUpload"
        >
          <Icon icon="lucide:upload" class="mr-2 h-4 w-4" />
          Upload
        </Button>
        <Button variant="outline" size="sm" @click="openCreateFileDialog">
          <Icon icon="lucide:file-plus-2" class="mr-2 h-4 w-4" />
          New File
        </Button>
        <Button variant="outline" size="sm" @click="openCreateDirDialog">
          <Icon icon="lucide:folder-plus" class="mr-2 h-4 w-4" />
          New Folder
        </Button>
      </div>
    </div>

    <input
      ref="fileInputRef"
      type="file"
      class="hidden"
      @change="handleFileInputChange"
    />

    <div
      v-if="!hasWorkspace"
      class="flex flex-1 items-center justify-center rounded-md border border-dashed p-4 text-sm text-muted-foreground"
    >
      Set an application name to access its workspace folder.
    </div>

    <div v-else class="grid flex-1 gap-4 lg:grid-cols-2">
      <div class="flex min-h-0 flex-col space-y-3 rounded-md border p-3">
        <div
          class="flex flex-wrap items-center gap-2 text-xs text-muted-foreground"
        >
          <button
            v-for="crumb in breadcrumbs"
            :key="crumb.path || 'root'"
            class="inline-flex items-center gap-1 rounded px-2 py-1 transition hover:bg-muted"
            @click="handleBreadcrumbClick(crumb.path)"
          >
            <Icon icon="lucide:folder" class="h-3 w-3" />
            {{ crumb.label || "/" }}
          </button>
        </div>
        <div class="flex items-center justify-between text-xs">
          <span class="text-muted-foreground">
            {{ currentPath ? `/${currentPath}` : "/" }}
          </span>
          <div class="flex items-center gap-2">
            <Button variant="ghost" size="sm" @click="refreshEntries">
              Refresh
            </Button>
            <Button
              variant="ghost"
              size="sm"
              :disabled="!selectedEntry || isDeleting"
              @click="requestDeleteSelection"
            >
              Delete
            </Button>
          </div>
        </div>
        <div class="flex flex-1 flex-col rounded-md border bg-muted/20">
          <div v-if="loading" class="p-4 text-sm text-muted-foreground">
            Loading...
          </div>
          <div
            v-else-if="entries.length === 0"
            class="p-4 text-sm text-muted-foreground"
          >
            Folder is empty.
          </div>
          <ul v-else class="flex-1 divide-y overflow-auto">
            <li
              v-for="entry in entries"
              :key="entry.path || entry.name"
              class="flex items-center justify-between px-3 py-2 transition hover:bg-muted"
              :class="selectedEntry?.path === entry.path ? 'bg-muted' : ''"
              @click="selectEntry(entry)"
              @dblclick.prevent="enterEntry(entry)"
            >
              <div class="flex items-center gap-3">
                <div
                  class="flex h-8 w-8 items-center justify-center rounded bg-primary/10"
                >
                  <Icon
                    :icon="entry.is_dir ? 'lucide:folder' : 'lucide:file'"
                    class="h-4 w-4 text-primary"
                  />
                </div>
                <div>
                  <p class="text-sm font-medium">{{ entry.name }}</p>
                  <p class="text-xs text-muted-foreground">
                    {{ entry.is_dir ? "Folder" : formatBytes(entry.size) }}
                  </p>
                </div>
              </div>
              <div
                class="flex items-center gap-2 text-xs text-muted-foreground"
              >
                <span>{{ new Date(entry.mod_time).toLocaleString() }}</span>
                <Button
                  v-if="entry.is_dir"
                  variant="ghost"
                  size="icon"
                  class="h-6 w-6 text-muted-foreground hover:text-foreground"
                  @click.stop="enterEntry(entry)"
                >
                  <Icon icon="lucide:corner-down-right" class="h-3.5 w-3.5" />
                </Button>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <div class="flex min-h-0 flex-col space-y-3 rounded-md border p-3">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium">Editor</p>
            <p class="text-xs text-muted-foreground">
              {{
                editorPath ? editorPath : "Select a file to preview and edit"
              }}
            </p>
          </div>
          <Button
            size="sm"
            :disabled="!isFileSelected || isSaving"
            @click="saveFile"
          >
            Save
          </Button>
        </div>
        <div class="flex flex-1 flex-col rounded-md border bg-background">
          <div
            v-if="!isFileSelected"
            class="flex flex-1 items-center justify-center px-4 py-6 text-center text-sm text-muted-foreground"
          >
            Select a file from the list to start editing.
          </div>
          <MonacoEditor
            v-else
            class="h-full min-h-[280px]"
            v-model:value="editorContent"
            theme="vs-dark"
            :language="editorLanguage"
            :options="editorOptions"
          />
        </div>
      </div>
    </div>
  </div>
  <Dialog v-model:open="isCreateFileDialogOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Create New File</DialogTitle>
        <DialogDescription>
          Enter a file name to create it inside the current workspace directory.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-2">
        <label class="text-sm font-medium" for="workspace-new-file">
          File Name
        </label>
        <Input
          id="workspace-new-file"
          v-model="newFileName"
          placeholder="e.g. config.yaml"
          :disabled="isCreatingFile"
        />
      </div>
      <DialogFooter>
        <Button
          variant="outline"
          @click="isCreateFileDialogOpen = false"
          :disabled="isCreatingFile"
        >
          Cancel
        </Button>
        <Button @click="confirmCreateFile" :disabled="isCreatingFile">
          {{ isCreatingFile ? "Creating..." : "Create" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
  <Dialog v-model:open="isCreateDirDialogOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Create New Folder</DialogTitle>
        <DialogDescription>
          Enter a folder name to create it inside the current workspace
          directory.
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-2">
        <label class="text-sm font-medium" for="workspace-new-dir">
          Folder Name
        </label>
        <Input
          id="workspace-new-dir"
          v-model="newDirName"
          placeholder="e.g. configs"
          :disabled="isCreatingDir"
        />
      </div>
      <DialogFooter>
        <Button
          variant="outline"
          @click="isCreateDirDialogOpen = false"
          :disabled="isCreatingDir"
        >
          Cancel
        </Button>
        <Button @click="confirmCreateDir" :disabled="isCreatingDir">
          {{ isCreatingDir ? "Creating..." : "Create" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
  <Dialog v-model:open="isDeleteDialogOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Delete Entry</DialogTitle>
        <DialogDescription>
          This action will permanently remove
          <span class="font-medium">{{ selectedEntry?.name }}</span>
          from the workspace.
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <Button
          variant="outline"
          @click="isDeleteDialogOpen = false"
          :disabled="isDeleting"
        >
          Cancel
        </Button>
        <Button
          variant="destructive"
          @click="deleteSelection"
          :disabled="isDeleting"
        >
          {{ isDeleting ? "Deleting..." : "Delete" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
