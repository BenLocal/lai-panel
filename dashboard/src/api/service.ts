import { post, type ApiResponse } from "./base";

export interface Service {
    id: number;
    name: string;
    app_id: number;
    node_id: number;
    status: string;
    metadata?: Record<string, string>;
}

export interface ServicePageResponse {
    total: number;
    page: number;
    pageSize: number;
    services: Service[];
}

export const serviceApi = {
    async page(page: number, pageSize: number): Promise<ApiResponse<ServicePageResponse>> {
        return post<ServicePageResponse>("/api/service/page", { page, pageSize });
    },
}