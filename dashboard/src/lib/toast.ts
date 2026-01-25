import { useToast } from 'primevue/usetoast'

// Create a global toast instance
// Note: This needs to be called within a component context
// We'll initialize it in App.vue
let toastInstance: ReturnType<typeof useToast> | null = null

export const setToastService = (toast: ReturnType<typeof useToast>) => {
  toastInstance = toast
}

export const showToast = (
  message: string,
  type: "success" | "error" | "warning" | "info" = "info",
  options?: {
    description?: string;
    duration?: number;
  }
) => {
  if (!toastInstance) {
    console.warn('Toast service not initialized. Make sure App.vue is mounted.')
    return
  }
  
  const severity = type === "error" ? "error" : type === "warning" ? "warn" : type === "success" ? "success" : "info"
  
  toastInstance.add({
    severity,
    summary: message,
    detail: options?.description,
    life: options?.duration ?? 3000,
  })
}
