import { post, stream, type ApiResponse } from "./base";

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

export interface DeployServiceRequest {
  app_id: number;
  node_id: number;
  qa_values: Record<string, string>;
}

export interface SaveServiceRequest {
  id: number;
  name: string;
  app_id: number;
  node_id: number;
  qa_values: Record<string, string>;
}

export const serviceApi = {
  async page(
    page: number,
    pageSize: number
  ): Promise<ApiResponse<ServicePageResponse>> {
    return post<ServicePageResponse>("/api/service/page", { page, pageSize });
  },

  async save(req: SaveServiceRequest): Promise<ApiResponse<void>> {
    return post<void>("/api/service/save", req);
  },

  async deployStream(
    req: DeployServiceRequest,
    onMessage?: (data: string) => void,
    onError?: (error: Error) => void,
    onEnd?: () => void
  ): Promise<AbortController> {
    return stream("/api/docker/compose/deploy", req, onMessage, onError, onEnd);
  },
};
