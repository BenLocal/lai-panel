<script setup lang="ts">
import { computed, ref, watch, onMounted, onUnmounted } from "vue";
import MonacoEditor from "@guolao/vue-monaco-editor";

const props = withDefaults(
  defineProps<{
    modelValue?: string;
    height?: number | string;
    readOnly?: boolean;
  }>(),
  {
    modelValue: "",
    height: 620,
    readOnly: false,
  }
);

const emit = defineEmits<{
  (event: "update:modelValue", value: string): void;
  (event: "valid-state-change", value: boolean): void;
}>();

const yamlContent = ref(props.modelValue ?? "");

watch(yamlContent, (value) => {
  emit("update:modelValue", value);
});

watch(
  () => props.modelValue,
  (value) => {
    if (value !== undefined && value !== yamlContent.value) {
      yamlContent.value = value;
    }
  }
);

const editorOptions = computed(() => ({
  automaticLayout: true,
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  fontSize: 14,
  tabSize: 2,
  insertSpaces: true,
  detectIndentation: false,
  wordWrap: "on",
  readOnly: props.readOnly,
}));

const containerRef = ref<HTMLElement | null>(null);
const editorHeight = ref<number | string>(props.height);

// Calculate height when using percentage
const updateHeight = () => {
  if (containerRef.value && typeof props.height === "string" && props.height.includes("%")) {
    const parent = containerRef.value.parentElement;
    if (parent) {
      const parentHeight = parent.clientHeight;
      const percentValue = parseFloat(props.height) / 100;
      editorHeight.value = parentHeight * percentValue;
    }
  } else {
    editorHeight.value = props.height;
  }
};

onMounted(() => {
  updateHeight();
  window.addEventListener("resize", updateHeight);
});

onUnmounted(() => {
  window.removeEventListener("resize", updateHeight);
});

watch(() => props.height, updateHeight);
</script>

<template>
  <div ref="containerRef" class="yaml-root" :class="{ readonly: readOnly }">
    <div class="yaml-wrap">
      <MonacoEditor
        class="yaml-editor-container"
        v-model:value="yamlContent"
        theme="vs-dark"
        language="yaml"
        :options="editorOptions"
        :style="{
          height: typeof editorHeight === 'number' ? `${editorHeight}px` : editorHeight,
          minHeight: typeof props.height === 'number' ? `${props.height}px` : props.height,
        }"
      />
    </div>
  </div>
</template>

<style scoped>
.yaml-root { height: 100%; display: flex; flex-direction: column; min-height: 0; }
.yaml-root.readonly { opacity: 0.85; }
.yaml-wrap { flex: 1; min-height: 0; border-radius: var(--p-border-radius); border: 1px solid var(--p-surface-border); overflow: hidden; }
.yaml-editor-container { height: 100%; }
.yaml-editor-container :deep(.monaco-editor) { border-radius: var(--p-border-radius); height: 100% !important; }
.yaml-editor-container :deep(.monaco-editor .monaco-editor-background) { height: 100%; }
.yaml-editor-container :deep(.monaco-editor .overflow-guard) { height: 100%; }
</style>
