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
    docker_compose?: string;
}

export interface ApplicationQAItem {
    name: string;
    type: "text" | "number" | "boolean" | "select" | "textarea";
    default_value?: string;
    options?: string[];
    required?: boolean;
    description?: string;
}

export interface ApplicationPageResponse {
    total: number;
    currentPage: number;
    pageSize: number;
    apps: Application[];
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

    async page(
        page: number,
        pageSize: number
    ): Promise<ApiResponse<ApplicationPageResponse>> {
        return post<ApplicationPageResponse>("/api/application/page", {
            page: page,
            page_size: pageSize,
        });
    },
};
