<script setup lang="ts">
import { ref, watch, nextTick } from "vue";
import { useRoute } from "vue-router";
import { Icon } from "@iconify/vue";
import { PanelLeft } from "lucide-vue-next";
import { useDark, useToggle } from "@vueuse/core";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar/index";
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from "@/components/ui/resizable";
import { Breadcrumb } from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";

const route = useRoute();
const isDark = useDark();
const toggleDark = useToggle(isDark);

const panelSizes = ref<number[]>([20, 80]);
const defaultSidebarSize = 20;
const isSidebarCollapsed = ref(false);

const toggleSidebar = () => {
  const newCollapsed = !isSidebarCollapsed.value;
  isSidebarCollapsed.value = newCollapsed;

  // Force update panel sizes immediately
  if (newCollapsed) {
    // Hide sidebar - set first panel to 0
    panelSizes.value = [0, 100];
  } else {
    // Show sidebar - restore default sizes
    panelSizes.value = [defaultSidebarSize, 100 - defaultSidebarSize];
  }

  // Force reactivity update
  nextTick(() => {
    const sizes = [...panelSizes.value];
    panelSizes.value = sizes;
  });
};

const handleResize = (sizes: number[]) => {
  // Update panel sizes when user manually resizes
  panelSizes.value = sizes;
  // Update collapsed state based on first panel size
  const firstPanelSize = sizes[0] ?? 0;
  isSidebarCollapsed.value = firstPanelSize < 1;
};

const navigation = [
  {
    title: "Dashboard",
    icon: "lucide:layout-dashboard",
    path: "/dashboard",
  },
  {
    title: "Applications",
    icon: "lucide:layers",
    path: "/applications",
  },
  {
    title: "Nodes",
    icon: "lucide:server",
    path: "/nodes",
  },
  {
    title: "Docker",
    icon: "lucide:container",
    path: "/docker",
  },
  {
    title: "Services",
    icon: "lucide:rocket",
    path: "/services",
  },
  {
    title: "Environment",
    icon: "lucide:key",
    path: "/environment-variables",
  },
  {
    title: "Settings",
    icon: "lucide:settings",
    path: "/settings",
  },
];

const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + "/");
};
</script>

<template>
  <SidebarProvider>
    <ResizablePanelGroup :model-value="panelSizes" @update:model-value="handleResize" direction="horizontal"
      class="h-screen w-full overflow-hidden">
      <ResizablePanel :default-size="defaultSidebarSize" :min-size="0" :max-size="40"
        :class="{ collapsed: isSidebarCollapsed }">
        <Sidebar v-show="!isSidebarCollapsed" collapsible="none" class="h-screen w-full">
          <SidebarHeader>
            <div class="flex items-center gap-2 px-2 py-1.5">
              <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <Icon icon="lucide:layers" class="h-4 w-4" />
              </div>
              <div class="flex flex-col">
                <span class="font-semibold text-sm">Panel Manager</span>
                <span class="text-xs text-muted-foreground">Enterprise</span>
              </div>
              <Icon icon="lucide:chevron-up" class="ml-auto h-4 w-4 text-muted-foreground" />
            </div>
          </SidebarHeader>

          <SidebarContent class="overflow-hidden">
            <SidebarGroup>
              <SidebarGroupContent>
                <SidebarMenu>
                  <SidebarMenuItem v-for="item in navigation" :key="item.path">
                    <SidebarMenuButton :as-child="true" :data-active="isActive(item.path)">
                      <router-link :to="item.path" class="flex w-full items-center gap-2">
                        <Icon :icon="item.icon" class="h-4 w-4 shrink-0" />
                        <span class="flex-1">{{ item.title }}</span>
                        <Icon icon="lucide:chevron-right" class="ml-auto h-4 w-4 shrink-0" />
                      </router-link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </SidebarContent>

          <SidebarFooter>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton class="w-full">
                  <div class="flex items-center gap-2">
                    <div
                      class="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground">
                      <Icon icon="lucide:user" class="h-4 w-4" />
                    </div>
                    <div class="flex flex-col flex-1 min-w-0">
                      <span class="font-semibold text-sm truncate">admin</span>
                      <span class="text-xs text-muted-foreground truncate">m@example.com</span>
                    </div>
                    <Icon icon="lucide:chevron-up" class="h-4 w-4 text-muted-foreground" />
                  </div>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarFooter>
        </Sidebar>
      </ResizablePanel>

      <ResizableHandle v-if="!isSidebarCollapsed" with-handle />

      <ResizablePanel :default-size="80" :min-size="60">
        <SidebarInset class="flex flex-col h-screen overflow-hidden">
          <header class="flex h-16 shrink-0 items-center justify-between gap-2 border-b px-4">
            <div class="flex items-center gap-2">
              <Button variant="ghost" size="icon" class="h-7 w-7 -ml-1" @click="toggleSidebar">
                <PanelLeft class="h-4 w-4" />
                <span class="sr-only">Toggle Sidebar</span>
              </Button>
              <Breadcrumb />
            </div>
            <div class="flex items-center gap-2">
              <Button variant="ghost" size="sm" @click="toggleDark()" class="h-9 w-9 p-0"
                :title="isDark ? 'Switch to light mode' : 'Switch to dark mode'">
                <Icon :icon="isDark ? 'lucide:sun' : 'lucide:moon'" class="h-4 w-4" />
              </Button>
              <Button variant="ghost" size="sm" as-child class="h-9 w-9 p-0" title="Gitee">
                <a href="https://gitee.com" target="_blank" rel="noopener noreferrer">
                  <Icon icon="simple-icons:gitee" class="h-4 w-4" />
                </a>
              </Button>
              <Button variant="ghost" size="sm" as-child class="h-9 w-9 p-0" title="GitHub">
                <a href="https://github.com" target="_blank" rel="noopener noreferrer">
                  <Icon icon="lucide:github" class="h-4 w-4" />
                </a>
              </Button>
            </div>
          </header>
          <main class="flex flex-1 flex-col gap-4 p-4 overflow-y-auto">
            <router-view />
          </main>
        </SidebarInset>
      </ResizablePanel>
    </ResizablePanelGroup>
  </SidebarProvider>
</template>

<style scoped>
[data-active="true"] {
  background-color: var(--color-sidebar-accent);
  color: var(--color-sidebar-accent-foreground);
}

:deep([data-sidebar="content"]) {
  overflow: hidden !important;
}

/* Force sidebar panel to collapse when isSidebarCollapsed is true */
:deep([data-slot="resizable-panel"]:first-child) {
  transition: width 0.2s ease;
}

:deep([data-slot="resizable-panel"]:first-child.collapsed) {
  width: 0 !important;
  min-width: 0 !important;
  max-width: 0 !important;
  overflow: hidden;
}
</style>
