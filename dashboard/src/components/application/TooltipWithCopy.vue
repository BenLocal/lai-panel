<script setup lang="ts">
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { showToast } from "@/lib/toast";

interface Props {
  /** 要显示的文本内容 */
  text: string;
  /** 最大宽度（用于截断） */
  maxWidth?: string;
  /** 是否显示复制按钮 */
  showCopy?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  maxWidth: "200px",
  showCopy: true,
});

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    showToast("Copied to clipboard", "success");
  } catch (error) {
    console.error("Failed to copy:", error);
    showToast("Failed to copy to clipboard", "error");
  }
};
</script>

<template>
  <TooltipProvider>
    <Tooltip>
      <TooltipTrigger as-child>
        <slot>
          <div class="truncate" :style="{ maxWidth: maxWidth }">
            {{ text }}
          </div>
        </slot>
      </TooltipTrigger>
      <TooltipContent class="flex items-center gap-2">
        <p class="max-w-xs break-all">{{ text }}</p>
        <Button
          v-if="showCopy"
          variant="ghost"
          size="icon"
          class="h-6 w-6 shrink-0"
          @click.stop="copyToClipboard(text)"
        >
          <Icon icon="lucide:copy" class="h-3 w-3" />
        </Button>
      </TooltipContent>
    </Tooltip>
  </TooltipProvider>
</template>

