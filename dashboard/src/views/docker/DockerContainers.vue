<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
const isLogDialogOpen = ref(false);
const logContent = ref<string[]>([]);
const logController = ref<AbortController | null>(null);
const logContainerName = ref<string>("");
const logContentRef = ref<HTMLDivElement | null>(null);

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

const getStatusColor = (status: string) => {
  return status === "running"
    ? "bg-green-500/10 text-green-500 border-green-500/20"
    : "bg-red-500/10 text-red-500 border-red-500/20";
};

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
  action: "start" | "stop" | "restart" | "remove" | "log" | "terminal",
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
  <Card>
    <CardHeader>
      <CardTitle>Containers</CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="loading" class="text-center py-8 text-muted-foreground">
        Loading...
      </div>
      <div v-else-if="containers.length > 0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Image</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Ports</TableHead>
              <TableHead>Created</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="container in containers" :key="container.id">
              <TableCell class="font-mono text-xs">
                {{ DockerUtils.getContainerShortId(container.id) }}
              </TableCell>
              <TableCell class="font-medium">{{ container.name }}</TableCell>
              <TableCell>{{ container.image }}</TableCell>
              <TableCell>
                <span
                  :class="[
                    'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium',
                    getStatusColor(container.status),
                  ]"
                >
                  {{ container.status }}
                </span>
              </TableCell>
              <TableCell class="font-mono text-xs">{{
                container.ports
              }}</TableCell>
              <TableCell class="text-muted-foreground">{{
                container.created
              }}</TableCell>
              <TableCell>
                <div class="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    class="h-8 px-2 text-green-600 hover:text-green-700 hover:bg-green-50 dark:hover:bg-green-950"
                    :disabled="container.status === 'running'"
                    @click="handleContainerAction('start', container)"
                  >
                    <Icon icon="lucide:play" class="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    class="h-8 px-2 text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
                    :disabled="container.status !== 'running'"
                    @click="handleContainerAction('stop', container)"
                  >
                    <Icon icon="lucide:square" class="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    class="h-8 px-2"
                    :disabled="container.status !== 'running'"
                    @click="handleContainerAction('terminal', container)"
                  >
                    <Icon icon="lucide:terminal" class="h-4 w-4" />
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button variant="ghost" size="sm" class="h-8 px-2">
                        <Icon icon="lucide:more-horizontal" class="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                      <DropdownMenuItem
                        :disabled="container.status === 'running'"
                        @click="handleContainerAction('start', container)"
                      >
                        <Icon icon="lucide:play" class="h-4 w-4 mr-2" />
                        Start
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        :disabled="container.status !== 'running'"
                        @click="handleContainerAction('stop', container)"
                      >
                        <Icon icon="lucide:square" class="h-4 w-4 mr-2" />
                        Stop
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        :disabled="container.status !== 'running'"
                        @click="handleContainerAction('restart', container)"
                      >
                        <Icon icon="lucide:rotate-cw" class="h-4 w-4 mr-2" />
                        Restart
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        @click="handleContainerAction('log', container)"
                      >
                        <Icon icon="lucide:file-text" class="h-4 w-4 mr-2" />
                        Log
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        variant="destructive"
                        @click="handleContainerAction('remove', container)"
                      >
                        <Icon icon="lucide:trash-2" class="h-4 w-4 mr-2" />
                        Remove
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
      <div v-else class="text-center py-8 text-muted-foreground">
        <Icon icon="lucide:box" class="h-12 w-12 mx-auto mb-4 opacity-50" />
        <p>No containers found</p>
      </div>
    </CardContent>
  </Card>

  <!-- Log Dialog -->
  <Dialog
    v-model:open="isLogDialogOpen"
    @update:open="(open) => !open && closeLogDialog()"
  >
    <DialogContent
      class="!max-w-none !w-screen !h-screen !max-h-screen flex flex-col p-0 gap-0 !rounded-none !translate-x-0 !translate-y-0 !top-0 !left-0 !right-0 !bottom-0"
    >
      <DialogHeader class="px-6 pt-6 pb-4 border-b shrink-0">
        <DialogTitle>Container Logs - {{ logContainerName }}</DialogTitle>
        <DialogDescription>
          Real-time logs from the container
        </DialogDescription>
      </DialogHeader>
      <div
        ref="logContentRef"
        class="flex-1 overflow-y-auto bg-muted/30 p-4 font-mono text-xs"
      >
        <div
          v-if="logContent.length === 0"
          class="text-center text-muted-foreground py-8"
        >
          Waiting for logs...
        </div>
        <div
          v-for="(line, index) in logContent"
          :key="index"
          class="mb-0.5 whitespace-pre-wrap break-words leading-relaxed"
        >
          {{ line }}
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
