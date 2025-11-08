import { post, type ApiResponse } from "./base";

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

export const nodeApi = {
  async list(): Promise<ApiResponse<Node[]>> {
    return post<Node[]>("/node/list", null);
  },

  async get(id: number): Promise<ApiResponse<Node>> {
    return post<Node>("/node/get", { id });
  },

  async create(node: CreateNodeRequest): Promise<ApiResponse<Node>> {
    return post<Node>("/node/add", node);
  },

  async update(node: UpdateNodeRequest): Promise<ApiResponse<Node>> {
    return post<Node>("/node/update", node);
  },

  async delete(id: number): Promise<ApiResponse<void>> {
    return post<void>("/node/delete", { id });
  },
};
