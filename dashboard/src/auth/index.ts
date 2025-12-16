import type { UserInfo } from "@/api/auth";

const TOKEN_KEY = "panel-token";
const USER_KEY = "panel-user";

/**
 * Token management
 */
export const tokenStorage = {
  /**
   * Get the authentication token
   */
  get(): string | null {
    return localStorage.getItem(TOKEN_KEY);
  },

  /**
   * Set the authentication token
   */
  set(token: string): void {
    localStorage.setItem(TOKEN_KEY, token);
  },

  /**
   * Remove the authentication token
   */
  remove(): void {
    localStorage.removeItem(TOKEN_KEY);
  },
};

/**
 * User information management
 */
export const userStorage = {
  /**
   * Get the user information
   */
  get(): UserInfo | null {
    try {
      const userInfoStr = localStorage.getItem(USER_KEY);
      if (userInfoStr) {
        return JSON.parse(userInfoStr) as UserInfo;
      }
      return null;
    } catch (error) {
      console.error("Failed to parse user info:", error);
      return null;
    }
  },

  /**
   * Set the user information
   */
  set(userInfo: UserInfo): void {
    try {
      localStorage.setItem(USER_KEY, JSON.stringify(userInfo));
    } catch (error) {
      console.error("Failed to save user info:", error);
    }
  },

  /**
   * Remove the user information
   */
  remove(): void {
    localStorage.removeItem(USER_KEY);
  },
};

/**
 * Authentication state management
 */
export const auth = {
  /**
   * Check if user is authenticated
   */
  isAuthenticated(): boolean {
    return !!tokenStorage.get();
  },

  /**
   * Get the authentication token
   */
  getToken(): string | null {
    return tokenStorage.get();
  },

  /**
   * Get the user information
   */
  getUser(): UserInfo | null {
    return userStorage.get();
  },

  /**
   * Set authentication data (token and user info)
   */
  setAuth(token: string, userInfo: UserInfo): void {
    tokenStorage.set(token);
    userStorage.set(userInfo);
  },

  /**
   * Clear all authentication data
   */
  clear(): void {
    tokenStorage.remove();
    userStorage.remove();
  },
};

