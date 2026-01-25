import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import AppLayout from "@/layouts/AppLayout.vue";
import { auth } from "@/auth";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: {
      title: "Login",
      requiresAuth: false,
    },
  },
  {
    path: "/",
    component: AppLayout,
    children: [
      {
        path: "",
        redirect: "/dashboard",
      },
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/Dashboard.vue"),
        meta: {
          title: "Dashboard",
          icon: "pi-chart-line",
          requiresAuth: true,
        },
      },
      {
        path: "applications",
        name: "Applications",
        component: () => import("@/views/Applications.vue"),
        meta: {
          title: "Applications",
          icon: "pi-th-large",
          requiresAuth: true,
        },
      },
      {
        path: "nodes",
        name: "Nodes",
        component: () => import("@/views/Nodes.vue"),
        meta: {
          title: "Nodes",
          icon: "pi-server",
          requiresAuth: true,
        },
      },
      {
        path: "docker",
        name: "Docker",
        component: () => import("@/views/Docker.vue"),
        meta: {
          title: "Docker",
          icon: "pi-box",
          requiresAuth: true,
        },
      },
      {
        path: "docker/container/terminal",
        name: "DockerContainerTerminal",
        component: () => import("@/views/docker/DockerContainerTerminal.vue"),
        meta: {
          title: "Container Terminal",
          icon: "pi-terminal",
          requiresAuth: true,
        },
      },
      {
        path: "nodes/terminal",
        name: "NodeTerminal",
        component: () => import("@/views/NodeTerminal.vue"),
        meta: {
          title: "Node Terminal",
          icon: "pi-terminal",
          requiresAuth: true,
        },
      },
      {
        path: "services",
        name: "Services",
        component: () => import("@/views/Services.vue"),
        meta: {
          title: "Services",
          icon: "pi-send",
          requiresAuth: true,
        },
      },
      {
        path: "environment-variables",
        name: "EnvironmentVariables",
        component: () => import("@/views/EnvironmentVariables.vue"),
        meta: {
          title: "Environment Variables",
          icon: "pi-key",
          requiresAuth: true,
        },
      },
      {
        path: "settings",
        name: "Settings",
        component: () => import("@/views/Settings.vue"),
        meta: {
          title: "Settings",
          icon: "pi-cog",
          requiresAuth: true,
        },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation guard to check authentication
router.beforeEach((to, from, next) => {
  const isAuthenticated = auth.isAuthenticated();
  const requiresAuth = to.meta.requiresAuth !== false; // Default to true if not specified

  // If route requires authentication and user is not authenticated
  if (requiresAuth && !isAuthenticated) {
    next({ name: "Login", query: { redirect: to.fullPath } });
    return;
  }

  // If user is authenticated and trying to access login page, redirect to dashboard
  if (to.name === "Login" && isAuthenticated) {
    next({ name: "Dashboard" });
    return;
  }

  next();
});

export default router;
