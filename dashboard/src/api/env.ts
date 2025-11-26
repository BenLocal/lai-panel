import { post, type ApiResponse } from "./base";

export interface Env {
  id: number;
  key: string;
  value: string;
  scope: string;
  description: string;
  metadata: string;
  created_at: string;
  updated_at: string;
}

export interface GetEnvPageRequest {
  scope: string;
  page: number;
  page_size: number;
}

export interface GetEnvPageResponse {
  total: number;
  current_page: number;
  page_size: number;
  list: Env[];
}

export interface GetEnvScopesResponse {
  scopes: string[];
}

export interface AddOrUpdateEnvRequest {
  id: number | null;
  key: string;
  value: string;
  scope: string;
}

export const envApi = {
  async page(request: GetEnvPageRequest): Promise<ApiResponse<GetEnvPageResponse>> {
    return post<GetEnvPageResponse>("/api/env/page", request);
  },

  async scopes(): Promise<ApiResponse<string[]>> {
    return post<string[]>("/api/env/scopes", null);
  },

  async addOrUpdate(request: AddOrUpdateEnvRequest): Promise<ApiResponse<void>> {
    return post<void>("/api/env/addOrUpdate", request);
  },

  async delete(id: number): Promise<ApiResponse<void>> {
    return post<void>("/api/env/delete", { id });
  },
}
