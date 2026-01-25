<script setup lang="ts">
import Button from 'primevue/button'
import { showToast } from "@/lib/toast";

const props = withDefaults(
  defineProps<{ text: string; maxWidth?: string; showCopy?: boolean }>(),
  { maxWidth: "200px", showCopy: true }
);

const copy = async () => {
  try {
    await navigator.clipboard.writeText(props.text);
    showToast("Copied to clipboard", "success");
  } catch (e) {
    console.error(e);
    showToast("Failed to copy", "error");
  }
};
</script>

<template>
  <div class="tooltip-copy" v-tooltip.top="text">
    <span class="tooltip-copy-text" :style="{ maxWidth: maxWidth }">
      <slot>{{ text }}</slot>
    </span>
    <Button
      v-if="showCopy"
      text
      rounded
      size="small"
      class="tooltip-copy-btn"
      @click.stop="copy"
    >
      <i class="pi pi-copy"></i>
    </Button>
  </div>
</template>

<style scoped>
.tooltip-copy {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  position: relative;
}
.tooltip-copy-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: default;
}
.tooltip-copy-btn { flex-shrink: 0; }
</style>
