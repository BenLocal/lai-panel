<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from "vue";
import { Terminal } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import "@xterm/xterm/css/xterm.css";

interface Props {
  /** 是否自动调整大小 */
  autoFit?: boolean;
  /** 字体大小 */
  fontSize?: number;
  /** 主题配置 */
  theme?: {
    background?: string;
    foreground?: string;
    cursor?: string;
    selection?: string;
  };
  /** 是否启用光标闪烁 */
  cursorBlink?: boolean;
  /** 是否转换行尾符 */
  convertEol?: boolean;
  /** 初始数据 */
  initialData?: string;
  /** 是否只读 */
  readonly?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  autoFit: true,
  fontSize: 13,
  cursorBlink: true,
  convertEol: true,
  readonly: false,
  theme: () => ({
    background: "#000000",
    foreground: "#ffffff",
  }),
});

const emit = defineEmits<{
  data: [data: string];
  ready: [];
  resize: [cols: number, rows: number];
}>();

const terminalRef = ref<HTMLDivElement | null>(null);
const terminal = ref<Terminal | null>(null);
const fitAddon = ref<FitAddon | null>(null);
const isTerminalInitialized = ref(false);

// 初始化终端
onMounted(async () => {
  if (!terminalRef.value) return;

  // 创建终端实例
  terminal.value = new Terminal({
    convertEol: props.convertEol,
    cursorBlink: props.cursorBlink,
    fontSize: props.fontSize,
    theme: props.theme,
    disableStdin: props.readonly,
  });

  // 加载 fit addon
  if (props.autoFit) {
    fitAddon.value = new FitAddon();
    terminal.value.loadAddon(fitAddon.value);
  }

  // 打开终端
  terminal.value.open(terminalRef.value);
  isTerminalInitialized.value = true;

  // 自动调整大小
  if (props.autoFit && fitAddon.value) {
    await nextTick();
    fitAddon.value.fit();
  }

  // 监听数据输入
  if (!props.readonly) {
    terminal.value.onData((data) => {
      emit("data", data);
    });
  }

  // 监听窗口大小变化
  if (props.autoFit && fitAddon.value) {
    const handleResize = () => {
      if (fitAddon.value && terminal.value) {
        fitAddon.value.fit();
        emit("resize", terminal.value.cols, terminal.value.rows);
      }
    };
    window.addEventListener("resize", handleResize);

    // 清理函数
    onUnmounted(() => {
      window.removeEventListener("resize", handleResize);
    });
  }

  // 写入初始数据
  if (props.initialData) {
    terminal.value.write(props.initialData);
  }

  emit("ready");
});

// 清理资源
onUnmounted(() => {
  try {
    // 只有在终端已初始化时才进行清理
    if (terminal.value && isTerminalInitialized.value) {
      // terminal.dispose() 会自动清理所有已加载的 addon
      // 不需要手动清理 fitAddon
      terminal.value.dispose();
    }
  } catch (error) {
    // 忽略清理错误，避免影响组件卸载
    console.warn("Error disposing terminal:", error);
  } finally {
    terminal.value = null;
    fitAddon.value = null;
    isTerminalInitialized.value = false;
  }
});

// 写入数据到终端
const write = (data: string) => {
  if (terminal.value) {
    terminal.value.write(data);
  }
};

// 写入一行数据
const writeln = (data: string = "") => {
  if (terminal.value) {
    terminal.value.writeln(data);
  }
};

// 清空终端
const clear = () => {
  if (terminal.value) {
    terminal.value.clear();
  }
};

// 获取终端列数和行数
const getSize = () => {
  if (terminal.value) {
    return {
      cols: terminal.value.cols,
      rows: terminal.value.rows,
    };
  }
  return { cols: 0, rows: 0 };
};

// 调整大小
const fit = () => {
  if (fitAddon.value) {
    fitAddon.value.fit();
    if (terminal.value) {
      emit("resize", terminal.value.cols, terminal.value.rows);
    }
  }
};

// 重置终端
const reset = () => {
  if (terminal.value) {
    terminal.value.reset();
  }
};

// 暴露方法给父组件
defineExpose({
  write,
  writeln,
  clear,
  getSize,
  fit,
  reset,
  terminal: terminal.value,
});
</script>

<template>
  <div ref="terminalRef" class="xterm-terminal-container"></div>
</template>

<style scoped>
.xterm-terminal-container {
  width: 100%;
  height: 100%;
  padding: 0.5rem;
}

.xterm-terminal-container :deep(.xterm) {
  height: 100%;
}

.xterm-terminal-container :deep(.xterm-viewport) {
  background-color: transparent !important;
}

.xterm-terminal-container :deep(.xterm-screen) {
  background-color: transparent !important;
}
</style>
