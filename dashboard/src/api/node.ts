const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export interface Node {
  id: number;
  is_local: boolean;
  name: string;
  display_name?: string | null;
  address: string;
  ssh_port: number;
  agent_port: number;
  ssh_user: string;
  ssh_password: string;
  status?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CreateNodeRequest {
  name: string;
  address: string;
  ssh_port: number;
  agent_port: number;
  ssh_user: string;
  ssh_password: string;
  is_local: boolean;
  display_name?: string;
}

export interface UpdateNodeRequest {
  id: number;
  name: string;
  address: string;
  ssh_port: number;
  agent_port: number;
  ssh_user: string;
  ssh_password: string;
  is_local: boolean;
  display_name?: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

async function request<T>(
  url: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, {
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

export const nodeApi = {
  // 获取节点列表
  async list(): Promise<ApiResponse<Node[]>> {
    return request<Node[]>("/node/list", {
      method: "POST",
    });
  },

  // 获取单个节点
  async get(id: number): Promise<ApiResponse<Node>> {
    return request<Node>("/node/get", {
      method: "POST",
      body: JSON.stringify({ id }),
    });
  },

  // 添加节点
  async create(node: CreateNodeRequest): Promise<ApiResponse<Node>> {
    return request<Node>("/node/add", {
      method: "POST",
      body: JSON.stringify(node),
    });
  },

  // 更新节点
  async update(node: UpdateNodeRequest): Promise<ApiResponse<Node>> {
    return request<Node>("/node/update", {
      method: "POST",
      body: JSON.stringify(node),
    });
  },

  // 删除节点
  async delete(id: number): Promise<ApiResponse<void>> {
    return request<void>("/node/delete", {
      method: "POST",
      body: JSON.stringify({ id }),
    });
  },
};
