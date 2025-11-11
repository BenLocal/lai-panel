import { post, type ApiResponse } from "./base";

export interface Node {
    id: number;
    is_local: boolean;
    name: string;
    display_name?: string | null;
    address: string;
    status?: string;
}

export interface CreateNodeRequest {
    name: string;
    address: string;
    ssh_port: number;
    agent_port: number;
    ssh_user: string;
    ssh_password: string;
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

export interface NodePageResponse {
    total: number;
    page: number;
    pageSize: number;
    nodes: Node[];
}

export const nodeApi = {
    async list(): Promise<ApiResponse<Node[]>> {
        return post<Node[]>("/api/node/list", null);
    },

    async get(id: number): Promise<ApiResponse<Node>> {
        return post<Node>("/api/node/get", { id });
    },

    async create(node: CreateNodeRequest): Promise<ApiResponse<Node>> {
        return post<Node>("/api/node/add", node);
    },

    async update(node: UpdateNodeRequest): Promise<ApiResponse<Node>> {
        return post<Node>("/api/node/update", node);
    },

    async delete(id: number): Promise<ApiResponse<void>> {
        return post<void>("/api/node/delete", { id });
    },

    async page(page: number, pageSize: number): Promise<ApiResponse<NodePageResponse>> {
        return post<NodePageResponse>("/api/node/page", { page, pageSize });
    },
};
