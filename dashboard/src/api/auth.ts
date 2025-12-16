import { post, type ApiResponse } from "./base";

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token?: string;
  user_id?: number;
  username?: string;
  role?: string;
  name?: string;
}

export interface UserInfo {
  user_id: number;
  username: string;
  role: string;
  name: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export const authApi = {
  async login(req: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return post<LoginResponse>("/api/auth/login", req);
  },

  async logout(): Promise<ApiResponse<void>> {
    return post<void>("/api/auth/logout", {});
  },

  async changePassword(req: ChangePasswordRequest): Promise<ApiResponse<void>> {
    return post<void>("/api/auth/change/password", req);
  },
};

