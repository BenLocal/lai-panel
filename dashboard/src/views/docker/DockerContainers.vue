<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import Menu from 'primevue/menu'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import TooltipWithCopy from "@/components/application/TooltipWithCopy.vue";
import { dockerApi, DockerUtils } from "@/api/docker";
import { showToast } from "@/lib/toast";
import { ApiResponseHelper } from "@/api/base";
import { nextTick } from "vue";

interface Container {
  id: string;
  name: string;
  image: string;
  status: string;
  ports: string;
  created: string;
}

interface Props {
  nodeId?: string;
}

const props = defineProps<Props>();
const router = useRouter();

const containers = ref<Container[]>([]);
const loading = ref(false);
const searchQuery = ref("");
const isLogDialogOpen = ref(false);
const logContent = ref<string[]>([]);
const logController = ref<AbortController | null>(null);
const logContainerName = ref<string>("");
const logContentRef = ref<HTMLDivElement | null>(null);
const isInspectDialogOpen = ref(false);
const inspectData = ref<any>(null);
const inspectContainerName = ref<string>("");
const inspectContentRef = ref<HTMLDivElement | null>(null);
const menuRefs = ref<Record<string, any>>({});

const fetchContainers = async () => {
  loading.value = true;

  const nodeId = Number(props.nodeId);
  if (Number.isNaN(nodeId)) {
    return;
  }
  const response = await dockerApi.containers(Number(props.nodeId));
  containers.value =
    response.data?.map((container) => ({
      id: container.Id,
      name: DockerUtils.getContainerName(container.Names),
      image: container.Image,
      status: container.State,
      ports: DockerUtils.getContainerDisplayPorts(container.Ports),
      created: new Date(container.Created * 1000).toLocaleString(),
    })) ?? [];
  loading.value = false;
};

const filteredContainers = computed(() => {
  if (!searchQuery.value.trim()) {
    return containers.value;
  }
  const query = searchQuery.value.toLowerCase();
  return containers.value.filter(
    (container) =>
      container.name.toLowerCase().includes(query) ||
      container.image.toLowerCase().includes(query) ||
      container.status.toLowerCase().includes(query) ||
      container.id.toLowerCase().includes(query) ||
      container.ports.toLowerCase().includes(query)
  );
});


const scrollLogToBottom = () => {
  if (logContentRef.value) {
    logContentRef.value.scrollTop = logContentRef.value.scrollHeight;
  }
};

const closeLogDialog = () => {
  // Abort the controller when closing dialog
  if (logController.value) {
    logController.value.abort();
    logController.value = null;
  }
  isLogDialogOpen.value = false;
  logContent.value = [];
  logContainerName.value = "";
};

const handleContainerAction = async (
  action:
    | "start"
    | "stop"
    | "restart"
    | "remove"
    | "log"
    | "terminal"
    | "inspect",
  container: Container
) => {
  const nodeId = Number(props.nodeId);
  if (Number.isNaN(nodeId)) {
    showToast("Invalid node ID", "error");
    return;
  }

  try {
    let response;
    switch (action) {
      case "start":
        response = await dockerApi.containerStart(nodeId, container.id);
        break;
      case "stop":
        response = await dockerApi.containerStop(nodeId, container.id);
        break;
      case "restart":
        response = await dockerApi.containerRestart(nodeId, container.id);
        break;
      case "remove":
        response = await dockerApi.containerRemove(nodeId, container.id);
        break;
      case "inspect":
        response = await dockerApi.containerInspect(nodeId, container.id);
        if (ApiResponseHelper.isSuccess(response)) {
          inspectData.value = response.data;
          inspectContainerName.value = container.name;
          isInspectDialogOpen.value = true;
        } else {
          showToast(response.message ?? "Failed to inspect container", "error");
        }
        return;
      case "log":
        // Open dialog and start log stream
        logContent.value = [];
        logContainerName.value = container.name;
        isLogDialogOpen.value = true;

        logController.value = await dockerApi.containerLogStream(
          nodeId,
          container.id,
          (data) => {
            logContent.value.push(data);
            // Auto scroll to bottom
            nextTick(() => {
              scrollLogToBottom();
            });
          },
          (error) => {
            showToast(`Log stream error: ${error.message}`, "error");
            console.error("Log stream error:", error);
          },
          () => {
            showToast("Log stream ended", "info");
          }
        );
        return;
      case "terminal":
        // 跳转到终端页面
        router.push({
          name: "DockerContainerTerminal",
          query: {
            nodeId: nodeId.toString(),
            containerId: container.id,
            containerName: container.name,
          },
        });
        return;
      default:
        return;
    }

    if (ApiResponseHelper.isSuccess(response)) {
      showToast(`Container ${action}ed successfully`, "success");
      await fetchContainers();
    } else {
      showToast(response.message ?? `Failed to ${action} container`, "error");
    }
  } catch (error) {
    showToast(`Failed to ${action} container`, "error");
    console.error(`Failed to ${action} container:`, error);
  }
};

onMounted(() => {
  fetchContainers();
});

// Watch for nodeId changes to refetch data
watch(
  () => props.nodeId,
  () => {
    if (props.nodeId) {
      fetchContainers();
    }
  }
);
</script>

<template>
  <div class="docker-panel">
    <div class="docker-toolbar">
      <h3>Containers</h3>
      <InputText v-model="searchQuery" placeholder="Search..." class="search-input" />
    </div>
    <div v-if="loading" class="docker-loading text-muted-foreground">Loading...</div>
    <div v-else-if="filteredContainers.length > 0" class="table-wrap">
      <DataTable :value="filteredContainers" size="small" striped-rows>
        <Column header="ID">
          <template #body="{ data }">
            <span class="font-mono">{{ DockerUtils.getContainerShortId(data.id) }}</span>
          </template>
        </Column>
        <Column header="Name">
          <template #body="{ data }"><TooltipWithCopy :text="data.name" max-width="200px" /></template>
        </Column>
        <Column header="Image">
          <template #body="{ data }"><TooltipWithCopy :text="data.image" max-width="300px" /></template>
        </Column>
        <Column header="Status">
          <template #body="{ data }">
            <Tag :value="data.status" :severity="data.status === 'running' ? 'success' : 'danger'" />
          </template>
        </Column>
        <Column header="Ports">
          <template #body="{ data }"><TooltipWithCopy :text="data.ports" max-width="200px" /></template>
        </Column>
        <Column header="Created">
          <template #body="{ data }"><span class="text-muted-foreground">{{ data.created }}</span></template>
        </Column>
        <Column header="Actions">
          <template #body="{ data }">
            <div class="action-btns">
              <Button text rounded size="small" :disabled="data.status === 'running'" @click="handleContainerAction('start', data)" v-tooltip.top="'Start'">
                <i class="pi pi-play"></i>
              </Button>
              <Button text rounded size="small" :disabled="data.status !== 'running'" @click="handleContainerAction('stop', data)" v-tooltip.top="'Stop'">
                <i class="pi pi-stop"></i>
              </Button>
              <Button text rounded size="small" :disabled="data.status !== 'running'" @click="handleContainerAction('terminal', data)" v-tooltip.top="'Terminal'">
                <i class="pi pi-terminal"></i>
              </Button>
              <Menu :ref="(el: any) => { if (el) menuRefs[data.id] = el }" :model="[
                { label: 'Start', icon: 'pi pi-play', command: () => handleContainerAction('start', data), disabled: data.status === 'running' },
                { label: 'Stop', icon: 'pi pi-stop', command: () => handleContainerAction('stop', data), disabled: data.status !== 'running' },
                { label: 'Restart', icon: 'pi pi-refresh', command: () => handleContainerAction('restart', data), disabled: data.status !== 'running' },
                { separator: true },
                { label: 'Inspect', icon: 'pi pi-info', command: () => handleContainerAction('inspect', data) },
                { label: 'Log', icon: 'pi pi-file', command: () => handleContainerAction('log', data) },
                { separator: true },
                { label: 'Remove', icon: 'pi pi-trash', command: () => handleContainerAction('remove', data) }
              ]" popup />
              <Button text rounded size="small" @click="(e) => menuRefs[data.id]?.toggle(e)">
                <i class="pi pi-ellipsis-h"></i>
              </Button>
            </div>
          </template>
        </Column>
      </DataTable>
    </div>
    <div v-else-if="!containers.length" class="docker-empty text-muted-foreground">
      <i class="pi pi-box"></i>
      <p>No containers found</p>
    </div>
    <div v-else class="docker-empty text-muted-foreground">
      <i class="pi pi-search"></i>
      <p>No containers match your search</p>
    </div>

    <Dialog v-model:visible="isLogDialogOpen" modal :style="{ width: '100vw', height: '100vh', maxWidth: 'none' }" :contentStyle="{ display: 'flex', flexDirection: 'column', padding: 0 }" @update:visible="(v: boolean) => !v && closeLogDialog()">
      <template #header>
        <div class="dialog-custom-header">
          <h3>Container Logs – {{ logContainerName }}</h3>
          <p class="text-muted-foreground">Real-time logs from the container</p>
        </div>
      </template>
      <div ref="logContentRef" class="docker-log-body">
        <div v-if="!logContent.length" class="text-muted-foreground">Waiting for logs...</div>
        <div v-for="(line, i) in logContent" :key="i" class="docker-log-line">{{ line }}</div>
      </div>
    </Dialog>

    <Dialog v-model:visible="isInspectDialogOpen" modal :style="{ width: '90vw', height: '90vh', maxWidth: '56rem' }" :contentStyle="{ display: 'flex', flexDirection: 'column', padding: 0 }">
      <template #header>
        <div class="dialog-custom-header">
          <h3>Container Inspect – {{ inspectContainerName }}</h3>
          <p class="text-muted-foreground">JSON representation of the container configuration</p>
        </div>
      </template>
      <div ref="inspectContentRef" class="docker-log-body">
        <pre v-if="inspectData" class="docker-inspect-pre">{{ JSON.stringify(inspectData, null, 2) }}</pre>
        <div v-else class="text-muted-foreground">Loading...</div>
      </div>
    </Dialog>
  </div>
</template>

<style scoped>
.docker-panel { padding: 1rem; border: 1px solid var(--p-surface-border); border-radius: var(--p-border-radius); background: var(--p-surface-card); }
.docker-toolbar { display: flex; align-items: center; justify-content: space-between; gap: 1rem; margin-bottom: 1rem; }
.docker-toolbar h3 { font-size: 1.125rem; font-weight: 600; margin: 0; }
.search-input { width: 12rem; }
.docker-loading { text-align: center; padding: 2rem; }
.table-wrap { border-radius: var(--p-border-radius); overflow: hidden; }
.action-btns { display: flex; align-items: center; gap: 0.5rem; }
.font-mono { font-family: ui-monospace, monospace; font-size: 0.75rem; }
.docker-empty { text-align: center; padding: 2rem; }
.docker-empty i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
.dialog-custom-header { padding: 1rem 1.5rem; border-bottom: 1px solid var(--p-surface-border); }
.dialog-custom-header h3 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.dialog-custom-header p { font-size: 0.875rem; }
.docker-log-body { flex: 1; overflow-y: auto; padding: 1rem; background: var(--p-surface-50); font-family: ui-monospace, monospace; font-size: 0.75rem; }
.docker-log-line { white-space: pre-wrap; word-break: break-word; margin-bottom: 0.125rem; }
.docker-inspect-pre { white-space: pre-wrap; word-break: break-word; font-size: 0.875rem; margin: 0; }
</style>
