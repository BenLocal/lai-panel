<script setup lang="ts">
import { ref, computed, type HTMLAttributes } from "vue"
import { useVModel } from "@vueuse/core"
import { cn } from "@/lib/utils"
import { Icon } from "@iconify/vue"

const props = defineProps<{
  defaultValue?: string | number
  modelValue?: string | number
  class?: HTMLAttributes["class"]
  type?: string
  placeholder?: string
  id?: string
  disabled?: boolean
}>()

const emits = defineEmits<{
  (e: "update:modelValue", payload: string | number): void
  (e: "keypress", event: KeyboardEvent): void
}>()

const modelValue = useVModel(props, "modelValue", emits, {
  passive: true,
  defaultValue: props.defaultValue,
})

const showPassword = ref(false)
const isPasswordType = computed(() => props.type === "password")
const inputType = computed(() => {
  if (isPasswordType.value) {
    return showPassword.value ? "text" : "password"
  }
  return props.type || "text"
})

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}
</script>

<template>
  <div class="relative">
    <input v-model="modelValue" :type="inputType" :placeholder="placeholder" :id="id" :disabled="disabled"
      data-slot="input" :class="cn(
        'file:text-foreground placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 border-input flex h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
        'focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]',
        'aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive',
        isPasswordType ? 'pr-10' : '',
        props.class,
      )" @keypress="emits('keypress', $event)" />
    <button v-if="isPasswordType" type="button"
      class="absolute right-2 top-1/2 -translate-y-1/2 w-6 h-6 flex items-center justify-center cursor-pointer text-muted-foreground hover:text-foreground transition-colors outline-none focus-visible:outline-none"
      @click="togglePasswordVisibility" tabindex="-1">
      <Icon :icon="showPassword ? 'lucide:eye-off' : 'lucide:eye'" class="h-4 w-4" />
    </button>
  </div>
</template>
