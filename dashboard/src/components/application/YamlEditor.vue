<script setup lang="ts">
import { computed, ref, watch, onMounted, onUnmounted } from "vue";
import MonacoEditor from "@guolao/vue-monaco-editor";
import { load } from "js-yaml";

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
  <div ref="containerRef" class="space-y-3 h-full flex flex-col">
    <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr),minmax(0,1fr)] flex-1 min-h-0"
      :class="readOnly ? 'opacity-80' : ''">
      <div class="rounded-lg border flex flex-col overflow-hidden">
        <MonacoEditor class="yaml-editor-container flex-1" v-model:value="yamlContent" theme="vs-dark" language="yaml"
          :options="editorOptions" :style="{
            height: typeof editorHeight === 'number' ? `${editorHeight}px` : editorHeight,
            minHeight: typeof props.height === 'number' ? `${props.height}px` : props.height,
          }" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.yaml-editor-container {
  height: 100%;
}

.yaml-editor-container :deep(.monaco-editor) {
  border-radius: 0 0 0.5rem 0.5rem;
  height: 100% !important;
}

.yaml-editor-container :deep(.monaco-editor .monaco-editor-background) {
  height: 100%;
}

.yaml-editor-container :deep(.monaco-editor .overflow-guard) {
  height: 100%;
}

.preview-container {
  border-top: 1px solid transparent;
  border-radius: 0 0 0.5rem 0.5rem;
  min-height: 220px;
}
</style>
