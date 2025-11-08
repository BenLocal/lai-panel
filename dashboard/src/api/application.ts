import { post, type ApiResponse, type Metadata } from "./base";

export interface Application {
  id: number;
  // english name
  name: string;
  // display name, if not set, use name
  display?: string;
  description?: string;
  version?: string;
  icon?: string;
  qa?: ApplicationQAItem[];
  metadata?: Metadata[] | null;
}

export interface ApplicationQAItem {
  name: string;
  value: string;
  type: "text" | "number" | "boolean" | "select" | "textarea";
  defaultValue?: string;
  options?: string[];
  required?: boolean;
  description?: string;
}

export const applicationApi = {
  async list(): Promise<ApiResponse<Application[]>> {
    return post<Application[]>("/api/application/list", null);
  },

  async add(application: Application): Promise<ApiResponse<Application>> {
    return post<Application>("/api/application/add", application);
  },

  async update(application: Application): Promise<ApiResponse<Application>> {
    return post<Application>("/api/application/update", application);
  },

  async delete(id: number): Promise<ApiResponse<void>> {
    return post<void>("/api/application/delete", { id });
  },

  async get(id: number): Promise<ApiResponse<Application>> {
    return post<Application>("/api/application/get", { id });
  },
};
