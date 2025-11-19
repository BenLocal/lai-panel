export interface ApiResponse<T> {
  code: number;
  message?: string;
  data?: T;
  error?: Error;
}

export const ApiResponseHelper = {
  isSuccess(response: ApiResponse<any>): boolean {
    return response.code === 0;
  },
};

export async function request<T>(
  url: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  try {
    const response = await fetch(`${url}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      return {
        code: -1,
        message: "网络错误",
        data: data,
        error: new Error("网络错误"),
      };
    }

    return data;
  } catch (error) {
    return {
      code: -1,
      message: "网络错误",
      error: error as Error,
    };
  }
}

export interface Metadata {
  name: string;
  properties: Record<string, string>;
}

export async function post<T>(
  url: string,
  data: any,
  headers: Record<string, string> = {}
): Promise<ApiResponse<T>> {
  const h = {
    "Content-Type": "application/json",
    ...headers,
  };

  return request<T>(`${url}`, {
    method: "POST",
    body: data ? JSON.stringify(data) : null,
    headers: h,
  });
}

export async function get<T>(url: string): Promise<ApiResponse<T>> {
  return request<T>(`${url}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });
}

export async function stream(
  url: string,
  data: any,
  onMessage?: (data: string) => void,
  onError?: (error: Error) => void,
  onEnd?: () => void,
  headers: Record<string, string> = {}
): Promise<AbortController> {
  let currentEvent = "";
  let currentData = "";
  let currentId = "";

  const processSSEEvent = () => {
    if (currentData) {
      console.log(currentData);
      onMessage?.(currentData);

      // 检测 done 事件，表示部署完成
      if (currentEvent === "done") {
        onEnd?.();
      }
    }

    // 重置当前事件数据
    currentEvent = "";
    currentData = "";
    currentId = "";
  };

  const processSSELine = (line: string) => {
    const trimmed = line.trim();

    // 空行表示一个事件结束
    if (!trimmed) {
      processSSEEvent();
      return;
    }

    // 跳过注释
    if (trimmed.startsWith(":")) return;

    if (trimmed.startsWith("event: ")) {
      currentEvent = trimmed.substring(7);
    } else if (trimmed.startsWith("data: ")) {
      const data = trimmed.substring(6);
      // 如果有多行 data，追加到当前数据
      if (currentData) {
        currentData += "\n" + data;
      } else {
        currentData = data;
      }
    } else if (trimmed.startsWith("id: ")) {
      currentId = trimmed.substring(4);
    }
  };

  const handleError = (error: unknown) => {
    if (error instanceof Error && error.name === "AbortError") {
      console.log("请求已取消");
    } else {
      const err = error instanceof Error ? error : new Error(String(error));
      console.error("SSE 请求错误:", err);
      onError?.(err);
    }
  };
  const controller = new AbortController();
  const handleStream = async () => {
    try {
      const h = {
        "Content-Type": "application/json",
        ...headers,
      };
      const response = await fetch(url, {
        method: "POST",
        headers: h,
        body: data ? JSON.stringify(data) : null,
        signal: controller.signal,
      });

      if (!response.ok || !response.body) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          // 流结束时，处理剩余的缓冲区数据
          if (buffer.trim()) {
            const lines = buffer.split("\n");
            for (const line of lines) {
              processSSELine(line);
            }
            // 处理最后一个事件（如果没有以换行符结尾）
            processSSEEvent();
          }
          // 调用结束回调
          onEnd?.();
          break;
        }

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          processSSELine(line);
        }
      }
    } catch (error) {
      handleError(error);
      // 发生错误时也调用结束回调
      onEnd?.();
    }
  };

  handleStream();
  return controller;
}
