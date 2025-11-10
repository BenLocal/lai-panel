<script setup lang="ts">
import { computed, ref, watch } from "vue";
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
    height: 320,
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

const parsedState = computed(() => {
  const source = yamlContent.value;

  if (!source || !source.trim()) {
    return {
      error: null as string | null,
      content: "",
    };
  }

  try {
    const parsed = load(source);
    let stringified = "";

    if (parsed !== undefined) {
      stringified =
        typeof parsed === "string" ? parsed : JSON.stringify(parsed, null, 2);
    }

    return {
      error: null as string | null,
      content: stringified,
    };
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Unable to parse YAML content";
    return {
      error: message,
      content: "",
    };
  }
});

const previewContent = computed(() => parsedState.value.content);
const parseError = computed(() => parsedState.value.error);

watch(
  () => parseError.value,
  (error) => {
    emit("valid-state-change", !error);
  },
  { immediate: true }
);
</script>

<template>
  <div class="space-y-3">
    <div
      class="grid gap-4 lg:grid-cols-[minmax(0,1fr),minmax(0,1fr)]"
      :class="readOnly ? 'opacity-80' : ''"
    >
      <div class="rounded-lg border">
        <div
          class="border-b px-3 py-2 text-xs font-medium uppercase text-muted-foreground"
        >
          Editor
        </div>
        <MonacoEditor
          class="yaml-editor-container"
          v-model:value="yamlContent"
          theme="vs-dark"
          language="yaml"
          :options="editorOptions"
          :style="{
            minHeight: typeof height === 'number' ? `${height}px` : height,
          }"
        />
      </div>
      <div class="rounded-lg border">
        <div
          class="border-b px-3 py-2 flex items-center justify-between text-xs font-medium uppercase text-muted-foreground"
        >
          <span>Preview (JSON)</span>
          <span
            v-if="!parseError"
            class="text-[10px] font-normal text-muted-foreground"
          >
            Parsed successfully
          </span>
          <span v-else class="text-[10px] font-normal text-destructive">
            Invalid YAML
          </span>
        </div>
        <div
          class="preview-container whitespace-pre-wrap break-words bg-muted/40 px-3 py-2 text-xs font-mono text-muted-foreground"
          :class="[parseError ? 'border-destructive/50' : 'border-transparent']"
        >
          <template v-if="parseError">
            {{ parseError }}
          </template>
          <template v-else-if="previewContent">
            {{ previewContent }}
          </template>
          <template v-else>
            <span class="italic text-muted-foreground/70">
              The parsed JSON view will appear here.
            </span>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.yaml-editor-container :deep(.monaco-editor) {
  border-radius: 0 0 0.5rem 0.5rem;
}

.preview-container {
  border-top: 1px solid transparent;
  border-radius: 0 0 0.5rem 0.5rem;
  min-height: 220px;
}
</style>
