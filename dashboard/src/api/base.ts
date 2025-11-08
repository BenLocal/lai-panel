export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

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
        success: false,
        error: data.error || "请求失败",
      };
    }

    return {
      success: true,
      data: data.data || data,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "网络错误",
    };
  }
}

export interface Metadata {
  name: string;
  properties: Record<string, string>;
}

export async function post<T>(url: string, data: any): Promise<ApiResponse<T>> {
  return request<T>(`${url}`, {
    method: "POST",
    body: data ? JSON.stringify(data) : null,
    headers: {
      "Content-Type": "application/json",
    },
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
