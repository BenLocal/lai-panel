<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import Button from 'primevue/button'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import XTermTerminal from "@/components/application/XTermTerminal.vue";
import { showToast } from "@/lib/toast";
import { HubConnectionBuilder, HubConnectionState } from "@microsoft/signalr";
import type { HubConnection } from "@microsoft/signalr";

const route = useRoute();
const router = useRouter();

// 从路由参数获取 nodeId 和 containerId
const nodeId = ref<number>(Number(route.query.nodeId) || 0);
const containerId = ref<string>((route.query.containerId as string) || "");
const containerName = ref<string>(
  (route.query.containerName as string) || "Unknown"
);

const terminalRef = ref<InstanceType<typeof XTermTerminal> | null>(null);
const connection = ref<HubConnection | null>(null);
const isConnected = ref(false);
const isConnecting = ref(false);
const connectionStatus = ref<
  "disconnected" | "connecting" | "connected" | "reconnecting"
>("disconnected");
const selectedShell = ref<"sh" | "bash" | "custom">("sh");
const customShell = ref<string>("");
const shellOptions = [
  { label: 'sh', value: 'sh' },
  { label: 'bash', value: 'bash' },
  { label: 'Custom', value: 'custom' }
];

// 初始化 SignalR 连接
const initConnection = () => {
  if (connection.value) {
    return;
  }

  connection.value = new HubConnectionBuilder()
    .withUrl("/api/signalr")
    .withAutomaticReconnect()
    .build();

  connection.value.serverTimeoutInMilliseconds = 6000;
  connection.value.keepAliveIntervalInMilliseconds = 2000;

  // 监听重连事件
  connection.value.onreconnecting(() => {
    connectionStatus.value = "reconnecting";
    showToast("Reconnecting...", "info");
  });

  // 监听重连成功事件
  connection.value.onreconnected(() => {
    connectionStatus.value = "connected";
    isConnected.value = true;
    showToast("Reconnected", "success");
  });

  // 监听关闭事件
  connection.value.onclose(() => {
    connectionStatus.value = "disconnected";
    isConnected.value = false;
  });

  // 监听容器数据
  connection.value.on("dockerExecData", (payload: string) => {
    if (terminalRef.value) {
      terminalRef.value.write(payload);
    }
  });

  // 监听容器关闭事件
  connection.value.on("dockerExecClosed", () => {
    if (terminalRef.value) {
      terminalRef.value.writeln("\r\n[Container session closed]");
    }
    isConnected.value = false;
    connectionStatus.value = "disconnected";
  });
};

// 连接到容器
const connect = async () => {
  if (!nodeId.value || !containerId.value) {
    showToast("Invalid node ID or container ID", "error");
    return;
  }

  if (isConnecting.value || isConnected.value) {
    return;
  }

  try {
    isConnecting.value = true;
    connectionStatus.value = "connecting";

    // 初始化连接（如果还没有）
    if (!connection.value) {
      initConnection();
    }

    // 启动连接
    if (connection.value?.state !== HubConnectionState.Connected) {
      console.log("Starting SignalR connection...");
      await connection.value?.start();
      console.log("Connection started successfully");
    }

    const size = terminalRef.value?.getSize();
    const cols = size?.cols || 80;
    const rows = size?.rows || 24;
    const shellToUse =
      selectedShell.value === "custom"
        ? customShell.value
        : selectedShell.value;
    if (!shellToUse || shellToUse.trim() === "") {
      showToast("Please enter a valid shell command", "error");
      return;
    }
    await connection.value?.invoke(
      "StartDockerExec",
      nodeId.value,
      containerId.value,
      cols,
      rows,
      shellToUse
    );

    isConnected.value = true;
    connectionStatus.value = "connected";
    showToast("Connected to container", "success");
  } catch (error) {
    showToast(
      `Failed to connect: ${
        error instanceof Error ? error.message : "Unknown error"
      }`,
      "error"
    );
    connectionStatus.value = "disconnected";
    isConnected.value = false;

    // 如果连接失败，清理连接
    if (connection.value) {
      try {
        await connection.value.stop();
      } catch (stopError) {
        // ignore
      }
      connection.value = null;
    }
  } finally {
    isConnecting.value = false;
  }
};

// 断开连接
const disconnect = async () => {
  try {
    if (connection.value?.state === HubConnectionState.Connected) {
      await connection.value.invoke("StopDockerExec");
    }
    await connection.value?.stop();
  } catch (error) {
    console.warn("Disconnect error:", error);
  } finally {
    isConnected.value = false;
    connectionStatus.value = "disconnected";
    if (terminalRef.value) {
      terminalRef.value.writeln("\r\n[Disconnected]");
    }
  }
};

// 处理终端数据输入
const handleTerminalData = async (data: string) => {
  if (connection.value?.state === HubConnectionState.Connected) {
    try {
      await connection.value.invoke("SendDockerExecInput", data);
    } catch (error) {
      console.error("Send data error:", error);
      showToast("Failed to send data", "error");
    }
  }
};

// 处理终端大小变化
const handleTerminalResize = async (cols: number, rows: number) => {
  if (connection.value?.state === HubConnectionState.Connected) {
    try {
      await connection.value.invoke("ResizeDockerExec", cols, rows);
    } catch (error) {
      console.error("Resize error:", error);
    }
  }
};

// 处理终端就绪
const handleTerminalReady = async () => {
  // 终端就绪后，确保终端大小正确
  await nextTick();
  if (terminalRef.value) {
    terminalRef.value.fit();
    // 自动聚焦到终端
    setTimeout(() => {
      if (terminalRef.value) {
        terminalRef.value.focus();
      }
    }, 100);
  }

  // 终端就绪后可以自动连接
  if (nodeId.value && containerId.value) {
    connect();
  }
};

// 重新连接
const reconnect = async () => {
  await disconnect();
  await nextTick();
  await connect();
};

// 页面挂载后尝试聚焦终端
onMounted(async () => {
  await nextTick();
  // 等待终端初始化完成后再聚焦
  setTimeout(() => {
    if (terminalRef.value) {
      terminalRef.value.focus();
    }
  }, 200);
});

// 清理资源
onUnmounted(async () => {
  await disconnect();
  if (connection.value) {
    connection.value = null;
  }
});

const back = () => {
  router.back();
};

// 监听路由参数变化
watch(
  () => [
    route.query.nodeId,
    route.query.containerId,
    route.query.containerName,
  ],
  ([newNodeId, newContainerId, newContainerName]) => {
    nodeId.value = Number(newNodeId) || 0;
    containerId.value = (newContainerId as string) || "";
    containerName.value = (newContainerName as string) || "Unknown";
    // 如果参数变化，断开当前连接
    if (isConnected.value) {
      disconnect();
    }
  }
);
</script>

<template>
  <div class="terminal-page">
    <div class="terminal-header">
      <div class="terminal-header-main">
        <Button text rounded @click="back">
          <i class="pi pi-arrow-left"></i>
          <span class="btn-icon-text">Back</span>
        </Button>
        <div>
          <h2>Container Terminal – {{ containerName }}</h2>
          <p class="text-muted-foreground">Node ID: {{ nodeId }} | Container: {{ containerId.substring(0, 12) }}</p>
        </div>
        <div class="shell-row">
          <label for="shell-select" class="text-muted-foreground">Shell:</label>
          <Select
            id="shell-select"
            v-model="selectedShell"
            :options="shellOptions"
            option-label="label"
            option-value="value"
            :disabled="isConnected || isConnecting"
            class="shell-select"
          />
          <InputText
            v-if="selectedShell === 'custom'"
            v-model="customShell"
            placeholder="Shell command"
            class="shell-input"
            :disabled="isConnected || isConnecting"
          />
        </div>
      </div>
      <div class="terminal-header-actions">
        <Tag
          :value="connectionStatus === 'connected' ? 'Connected' : connectionStatus === 'connecting' ? 'Connecting...' : connectionStatus === 'reconnecting' ? 'Reconnecting...' : 'Disconnected'"
          :severity="connectionStatus === 'connected' ? 'success' : connectionStatus === 'connecting' || connectionStatus === 'reconnecting' ? 'warn' : 'danger'"
        />
        <Button v-if="!isConnected && !isConnecting" @click="connect" :disabled="!nodeId || !containerId">
          <i class="pi pi-link"></i>
          <span class="btn-icon-text">Connect</span>
        </Button>
        <Button v-else-if="isConnected" severity="danger" @click="disconnect">
          <i class="pi pi-bolt"></i>
          <span class="btn-icon-text">Disconnect</span>
        </Button>
        <Button v-else outlined :disabled="true" :loading="true">Connecting...</Button>
        <Button v-if="isConnected" outlined @click="reconnect">
          <i class="pi pi-refresh"></i>
          <span class="btn-icon-text">Reconnect</span>
        </Button>
      </div>
    </div>
    <div class="terminal-body">
      <XTermTerminal
        ref="terminalRef"
        :auto-fit="true"
        :font-size="13"
        :readonly="false"
        @data="handleTerminalData"
        @ready="handleTerminalReady"
        @resize="handleTerminalResize"
      />
    </div>
  </div>
</template>

<style scoped>
.terminal-page { height: 100vh; display: flex; flex-direction: column; }
.terminal-header { border-bottom: 1px solid var(--p-surface-border); padding: 1.5rem; flex-shrink: 0; display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 1rem; }
.terminal-header-main { display: flex; flex-wrap: wrap; align-items: center; gap: 1rem; }
.terminal-header-main h2 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.terminal-header-main p { font-size: 0.875rem; }
.shell-row { display: flex; align-items: center; gap: 0.5rem; }
.shell-row label { font-size: 0.875rem; }
.shell-select { width: 6rem; }
.shell-input { width: 9rem; }
.terminal-header-actions { display: flex; flex-wrap: wrap; align-items: center; gap: 0.5rem; }
.btn-icon-text { margin-left: 0.5rem; }
.terminal-body { flex: 1; overflow: hidden; background: #1e1e1e; }
:deep(.xterm-terminal-container) { height: 100%; min-height: 200px; }
</style>
