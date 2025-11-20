import { toast } from "vue-sonner";

export const showToast = (
  message: string,
  type: "success" | "error" | "warning" | "info" = "info",
  options?: {
    description?: string;
    duration?: number;
  }
) => {
  const baseOptions = {
    duration: options?.duration ?? 3000,
    closeButton: false,
    ...(options?.description && { description: options.description }),
  };

  switch (type) {
    case "success":
      toast.success(message, baseOptions);
      break;
    case "error":
      toast.error(message, {
        ...baseOptions,
        style: {
          backgroundColor: "#ef4444",
          color: "white",
          border: "1px solid #dc2626",
        },
      });
      break;
    case "warning":
      toast.warning(message, {
        ...baseOptions,
        style: {
          backgroundColor: "#fbbf24",
          color: "#78350f",
          border: "1px solid #f59e0b",
        },
      });
      break;
    case "info":
      toast.info(message, baseOptions);
      break;
  }
};
