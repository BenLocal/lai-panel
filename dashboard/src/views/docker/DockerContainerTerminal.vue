<script setup lang="ts">
import { ref, onUnmounted, watch, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import XTermTerminal from "@/components/application/XTermTerminal.vue";
import { showToast } from "@/lib/toast";
import { HubConnectionBuilder, HubConnectionState } from "@microsoft/signalr";
import type { HubConnection } from "@microsoft/signalr";

const route = useRoute();
const router = useRouter();

// 从路由参数获取 nodeId 和 containerId
const nodeId = ref<number>(Number(route.query.nodeId) || 0);
const containerId = ref<string>((route.query.containerId as string) || "");
const containerName = ref<string>((route.query.containerName as string) || "Unknown");

const terminalRef = ref<InstanceType<typeof XTermTerminal> | null>(null);
const connection = ref<HubConnection | null>(null);
const isConnected = ref(false);
const isConnecting = ref(false);
const connectionStatus = ref<"disconnected" | "connecting" | "connected" | "reconnecting">("disconnected");

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
    showToast("Connection closed", "info");
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
    showToast("Container session closed", "info");
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

    // 获取终端尺寸
    const size = terminalRef.value?.getSize();
    const cols = size?.cols || 80;
    const rows = size?.rows || 24;

    // 启动 docker exec 会话
    // 注意：后端需要实现以下 SignalR Hub 方法：
    // - StartDockerExec(nodeId int64, containerId string, cols int, rows int)
    // - SendDockerExecInput(data string)
    // - ResizeDockerExec(cols int, rows int)
    // - StopDockerExec()
    // 以及事件：
    // - dockerExecData (payload string)
    // - dockerExecClosed ()
    await connection.value?.invoke("StartDockerExec", nodeId.value, containerId.value, cols, rows);

    isConnected.value = true;
    connectionStatus.value = "connected";
    showToast("Connected to container", "success");
  } catch (error) {
    console.error("Connection error:", error);
    console.error("Error details:", {
      message: error instanceof Error ? error.message : String(error),
      stack: error instanceof Error ? error.stack : undefined,
      connectionState: connection.value?.state,
    });
    showToast(`Failed to connect: ${error instanceof Error ? error.message : "Unknown error"}`, "error");
    connectionStatus.value = "disconnected";
    isConnected.value = false;

    // 如果连接失败，清理连接
    if (connection.value) {
      try {
        await connection.value.stop();
      } catch (stopError) {
        console.error("Error stopping connection:", stopError);
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
const handleTerminalReady = () => {
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
  () => [route.query.nodeId, route.query.containerId, route.query.containerName],
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
  <div class="h-screen flex flex-col">
    <!-- Header -->
    <Card class="rounded-none border-b shrink-0">
      <CardHeader class="pb-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <Button variant="ghost" size="sm" @click="back">
              <Icon icon="lucide:arrow-left" class="h-4 w-4 mr-2" />
              Back
            </Button>
            <div>
              <CardTitle>Container Terminal - {{ containerName }}</CardTitle>
              <p class="text-sm text-muted-foreground mt-1">
                Node ID: {{ nodeId }} | Container: {{ containerId.substring(0, 12) }}
              </p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span :class="[
              'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium',
              connectionStatus === 'connected'
                ? 'bg-green-500/10 text-green-500 border-green-500/20'
                : connectionStatus === 'connecting' || connectionStatus === 'reconnecting'
                  ? 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20'
                  : 'bg-red-500/10 text-red-500 border-red-500/20',
            ]">
              <span :class="[
                'h-1.5 w-1.5 rounded-full',
                connectionStatus === 'connected'
                  ? 'bg-green-500'
                  : connectionStatus === 'connecting' || connectionStatus === 'reconnecting'
                    ? 'bg-yellow-500 animate-pulse'
                    : 'bg-red-500',
              ]"></span>
              {{ connectionStatus === 'connected' ? 'Connected' : connectionStatus === 'connecting' ? 'Connecting...' :
                connectionStatus === 'reconnecting' ? 'Reconnecting...' : 'Disconnected' }}
            </span>
            <Button v-if="!isConnected && !isConnecting" @click="connect" :disabled="!nodeId || !containerId">
              <Icon icon="lucide:plug" class="h-4 w-4 mr-2" />
              Connect
            </Button>
            <Button v-else-if="isConnected" variant="destructive" @click="disconnect">
              <Icon icon="lucide:plug-zap" class="h-4 w-4 mr-2" />
              Disconnect
            </Button>
            <Button v-else variant="outline" @click="disconnect" :disabled="true">
              <Icon icon="lucide:loader-2" class="h-4 w-4 mr-2 animate-spin" />
              Connecting...
            </Button>
            <Button v-if="isConnected" variant="outline" @click="reconnect">
              <Icon icon="lucide:refresh-cw" class="h-4 w-4 mr-2" />
              Reconnect
            </Button>
          </div>
        </div>
      </CardHeader>
    </Card>

    <!-- Terminal -->
    <CardContent class="flex-1 p-0 overflow-hidden">
      <div class="h-full w-full">
        <XTermTerminal ref="terminalRef" :auto-fit="true" :font-size="13" :readonly="false" @data="handleTerminalData"
          @ready="handleTerminalReady" @resize="handleTerminalResize" />
      </div>
    </CardContent>
  </div>
</template>

<style scoped>
:deep(.xterm-terminal-container) {
  height: 100%;
}
</style>
