import { post, type ApiResponse } from "./base";

export interface WorkspaceEntry {
  name: string;
  path: string;
  is_dir: boolean;
  size: number;
  mod_time: string;
}

export interface WorkspaceListResponse {
  currentPath: string;
  entries: WorkspaceEntry[];
}

export const workspaceApi = {
  async list(appName: string, path = "") {
    return post<WorkspaceListResponse>("/api/workspace/list", {
      app_name: appName,
      path,
    });
  },

  async read(appName: string, path: string) {
    return post<{ content: string }>("/api/workspace/read", {
      app_name: appName,
      path,
    });
  },

  async save(appName: string, path: string, content: string) {
    return post<void>("/api/workspace/save", {
      app_name: appName,
      path,
      content,
    });
  },

  async remove(appName: string, path: string) {
    return post<void>("/api/workspace/delete", {
      app_name: appName,
      path,
    });
  },

  async mkdir(appName: string, path: string) {
    return post<void>("/api/workspace/mkdir", {
      app_name: appName,
      path,
    });
  },
};

export async function uploadWorkspaceFile(
  appName: string,
  directory: string,
  file: File
): Promise<ApiResponse<any>> {
  const formData = new FormData();
  formData.append("app_name", appName);
  formData.append("path", directory || "");
  formData.append("file", file);

  try {
    const response = await fetch("/api/workspace/upload", {
      method: "POST",
      body: formData,
    });

    const data = await response.json();

    if (!response.ok) {
      return {
        code: -1,
        message: data?.message || "Upload failed",
        data,
      };
    }

    return data;
  } catch (error) {
    return {
      code: -1,
      message: "Network error",
      error: error as Error,
    };
  }
}
