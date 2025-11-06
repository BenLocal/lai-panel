import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import AppLayout from "@/layouts/AppLayout.vue";

const routes: RouteRecordRaw[] = [
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
          icon: "lucide:layout-dashboard",
        },
      },
      {
        path: "applications",
        name: "Applications",
        component: () => import("@/views/Applications.vue"),
        meta: {
          title: "Applications",
          icon: "lucide:layers",
        },
      },
      {
        path: "nodes",
        name: "Nodes",
        component: () => import("@/views/Nodes.vue"),
        meta: {
          title: "Nodes",
          icon: "lucide:server",
        },
      },
      {
        path: "docker",
        name: "Docker",
        component: () => import("@/views/Docker.vue"),
        meta: {
          title: "Docker",
          icon: "lucide:container",
        },
      },
      {
        path: "services",
        name: "Services",
        component: () => import("@/views/Services.vue"),
        meta: {
          title: "Services",
          icon: "lucide:rocket",
        },
      },
      {
        path: "environment-variables",
        name: "EnvironmentVariables",
        component: () => import("@/views/EnvironmentVariables.vue"),
        meta: {
          title: "Environment Variables",
          icon: "lucide:key",
        },
      },
      {
        path: "settings",
        name: "Settings",
        component: () => import("@/views/Settings.vue"),
        meta: {
          title: "Settings",
          icon: "lucide:settings",
        },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
