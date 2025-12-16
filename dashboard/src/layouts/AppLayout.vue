<script setup lang="ts">
import { ref, nextTick, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
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
} from "@/components/ui/sidebar/index";
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from "@/components/ui/resizable";
import { Breadcrumb } from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { authApi } from "@/api/auth";
import { showToast } from "@/lib/toast";
import { auth, userStorage } from "@/auth";
import type { UserInfo } from "@/api/auth";

const route = useRoute();
const router = useRouter();
const isDark = useDark();
const toggleDark = useToggle(isDark);

// User information
const userInfo = ref<UserInfo | null>(null);

// Change password dialog state
const isChangePasswordDialogOpen = ref(false);
const oldPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const isChangingPassword = ref(false);

// Load user information from localStorage
const loadUserInfo = () => {
  userInfo.value = userStorage.get();
};

// Display name (use name if available, otherwise username)
const displayName = computed(() => {
  if (userInfo.value?.name) {
    return userInfo.value.name;
  }
  return userInfo.value?.username || "User";
});

// Display email/username (use username as secondary info)
const displaySecondary = computed(() => {
  return userInfo.value?.username || "";
});

onMounted(() => {
  loadUserInfo();
});

// Open change password dialog
const openChangePasswordDialog = () => {
  isChangePasswordDialogOpen.value = true;
  oldPassword.value = "";
  newPassword.value = "";
  confirmPassword.value = "";
};

// Close change password dialog
const closeChangePasswordDialog = () => {
  isChangePasswordDialogOpen.value = false;
  oldPassword.value = "";
  newPassword.value = "";
  confirmPassword.value = "";
};

// Handle change password
const handleChangePassword = async () => {
  // Validation
  if (!oldPassword.value || !newPassword.value || !confirmPassword.value) {
    showToast("Please fill in all fields", "error");
    return;
  }

  if (newPassword.value !== confirmPassword.value) {
    showToast("New password and confirm password do not match", "error");
    return;
  }

  if (newPassword.value.length < 6) {
    showToast("New password must be at least 6 characters", "error");
    return;
  }

  isChangingPassword.value = true;

  try {
    const base64OldPassword = btoa(oldPassword.value);
    const base64NewPassword = btoa(newPassword.value);

    const response = await authApi.changePassword({
      old_password: base64OldPassword,
      new_password: base64NewPassword,
    });

    if (response.code === 0) {
      showToast("Password changed successfully", "success");
      closeChangePasswordDialog();
    } else {
      showToast(response.message || "Failed to change password", "error");
    }
  } catch (error) {
    showToast("Failed to change password, please try again later", "error");
    console.error("Change password error:", error);
  } finally {
    isChangingPassword.value = false;
  }
};

// Logout function
const handleLogout = async () => {
  try {
    // Call logout API if needed
    await authApi.logout();
  } catch (error) {
    console.error("Logout error:", error);
  } finally {
    // Clear all authentication data
    auth.clear();
    // Clear user info
    userInfo.value = null;
    // Show toast
    showToast("Logged out successfully", "success");
    // Redirect to login page
    router.push("/login");
  }
};

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
                  <div class="flex items-center gap-2 w-full">
                    <div
                      class="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground">
                      <Icon icon="lucide:user" class="h-4 w-4" />
                    </div>
                    <div class="flex flex-col flex-1 min-w-0">
                      <span class="font-semibold text-sm truncate">{{ displayName }}</span>
                      <span class="text-xs text-muted-foreground truncate">{{ displaySecondary }}</span>
                    </div>
                    <DropdownMenu>
                      <DropdownMenuTrigger as-child>
                        <button
                          class="flex h-6 w-6 items-center justify-center rounded-md hover:bg-accent focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 ml-auto"
                          type="button">
                          <Icon icon="lucide:chevron-down" class="h-4 w-4 text-muted-foreground" />
                        </button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end" class="w-56">
                        <DropdownMenuLabel>
                          <div class="flex flex-col">
                            <span class="font-semibold">{{ displayName }}</span>
                            <span class="text-xs text-muted-foreground font-normal">{{ displaySecondary }}</span>
                          </div>
                        </DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem @click="openChangePasswordDialog" class="cursor-pointer">
                          <Icon icon="lucide:key" class="mr-2 h-4 w-4" />
                          <span>Change Password</span>
                        </DropdownMenuItem>
                        <DropdownMenuItem @click="handleLogout" class="cursor-pointer">
                          <Icon icon="lucide:log-out" class="mr-2 h-4 w-4" />
                          <span>Logout</span>
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
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
            </div>
          </header>
          <main class="flex flex-1 flex-col gap-4 p-4 overflow-y-auto">
            <router-view />
          </main>
        </SidebarInset>
      </ResizablePanel>
    </ResizablePanelGroup>
  </SidebarProvider>

  <!-- Change Password Dialog -->
  <Dialog v-model:open="isChangePasswordDialogOpen">
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>Change Password</DialogTitle>
        <DialogDescription>
          Enter your old password and choose a new password.
        </DialogDescription>
      </DialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid gap-2">
          <Label for="old-password">Old Password</Label>
          <Input id="old-password" v-model="oldPassword" type="password" placeholder="Enter old password"
            :disabled="isChangingPassword" @keypress.enter="handleChangePassword" />
        </div>
        <div class="grid gap-2">
          <Label for="new-password">New Password</Label>
          <Input id="new-password" v-model="newPassword" type="password" placeholder="Enter new password"
            :disabled="isChangingPassword" @keypress.enter="handleChangePassword" />
        </div>
        <div class="grid gap-2">
          <Label for="confirm-password">Confirm New Password</Label>
          <Input id="confirm-password" v-model="confirmPassword" type="password" placeholder="Confirm new password"
            :disabled="isChangingPassword" @keypress.enter="handleChangePassword" />
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" @click="closeChangePasswordDialog" :disabled="isChangingPassword">
          Cancel
        </Button>
        <Button @click="handleChangePassword" :disabled="isChangingPassword">
          <span v-if="!isChangingPassword">Change Password</span>
          <span v-else class="flex items-center gap-2">
            <Icon icon="lucide:loader-2" class="h-4 w-4 animate-spin" />
            Changing...
          </span>
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
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
