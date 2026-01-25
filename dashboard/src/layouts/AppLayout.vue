<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import Menu from 'primevue/menu'
import { authApi } from "@/api/auth";
import { showToast } from "@/lib/toast";
import { auth, userStorage } from "@/auth";
import type { UserInfo } from "@/api/auth";

const route = useRoute();
const router = useRouter();

const userInfo = ref<UserInfo | null>(null);
const isChangePasswordDialogOpen = ref(false);
const oldPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const isChangingPassword = ref(false);

const userMenu = ref();
const userMenuItems = ref([
  { label: 'Change Password', icon: 'pi pi-key', command: () => openChangePasswordDialog() },
  { separator: true },
  { label: 'Logout', icon: 'pi pi-sign-out', command: () => handleLogout() }
]);

const loadUserInfo = () => { userInfo.value = userStorage.get(); };
const displayName = computed(() => userInfo.value?.name || userInfo.value?.username || "User");
const displaySecondary = computed(() => userInfo.value?.username || "");

onMounted(() => loadUserInfo());

const openChangePasswordDialog = () => {
  isChangePasswordDialogOpen.value = true;
  oldPassword.value = "";
  newPassword.value = "";
  confirmPassword.value = "";
};

const closeChangePasswordDialog = () => {
  isChangePasswordDialogOpen.value = false;
  oldPassword.value = "";
  newPassword.value = "";
  confirmPassword.value = "";
};

const handleChangePassword = async () => {
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
    const response = await authApi.changePassword({
      old_password: btoa(oldPassword.value),
      new_password: btoa(newPassword.value),
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

const handleLogout = async () => {
  try { await authApi.logout(); } catch (e) { console.error("Logout error:", e); }
  auth.clear();
  userInfo.value = null;
  showToast("Logged out successfully", "success");
  router.push("/login");
};

const isSidebarCollapsed = ref(false);
const toggleSidebar = () => { isSidebarCollapsed.value = !isSidebarCollapsed.value; };

const navigation = [
  { title: "Dashboard", icon: "pi-chart-line", path: "/dashboard" },
  { title: "Applications", icon: "pi-th-large", path: "/applications" },
  { title: "Nodes", icon: "pi-server", path: "/nodes" },
  { title: "Docker", icon: "pi-box", path: "/docker" },
  { title: "Services", icon: "pi-send", path: "/services" },
  { title: "Environment", icon: "pi-key", path: "/environment-variables" },
  { title: "Settings", icon: "pi-cog", path: "/settings" },
];

const isActive = (path: string) => route.path === path || route.path.startsWith(path + "/");
</script>

<template>
  <div class="app-layout">
    <aside v-show="!isSidebarCollapsed" class="app-sidebar">
      <div class="app-sidebar-header">
        <div class="app-sidebar-logo"><i class="pi pi-th-large"></i></div>
      </div>
      <nav class="app-sidebar-nav">
        <router-link
          v-for="item in navigation"
          :key="item.path"
          :to="item.path"
          v-tooltip.right="item.title"
          class="app-nav-link"
          :class="{ 'app-nav-link-active': isActive(item.path) }"
        >
          <i :class="'pi ' + item.icon"></i>
        </router-link>
      </nav>
      <div class="app-sidebar-footer">
        <Button text rounded v-tooltip.right="displayName" @click="(e) => userMenu.toggle(e)" class="app-user-btn">
          <i class="pi pi-user"></i>
        </Button>
        <Menu ref="userMenu" :model="userMenuItems" popup />
      </div>
    </aside>

    <div class="app-main">
      <header class="app-header">
        <Button text rounded @click="toggleSidebar" aria-label="Toggle sidebar">
          <i :class="isSidebarCollapsed ? 'pi pi-angle-right' : 'pi pi-angle-left'"></i>
        </Button>
      </header>
      <main class="app-content">
        <router-view />
      </main>
    </div>

    <Dialog
      v-model:visible="isChangePasswordDialogOpen"
      modal
      header="Change Password"
      :style="{ width: '425px' }"
      :closable="!isChangingPassword"
    >
      <div class="app-dialog-form">
        <div class="app-field">
          <label for="old-password">Old Password</label>
          <InputText id="old-password" v-model="oldPassword" type="password" placeholder="Enter old password"
            :disabled="isChangingPassword" @keypress.enter="handleChangePassword" />
        </div>
        <div class="app-field">
          <label for="new-password">New Password</label>
          <InputText id="new-password" v-model="newPassword" type="password" placeholder="Enter new password"
            :disabled="isChangingPassword" @keypress.enter="handleChangePassword" />
        </div>
        <div class="app-field">
          <label for="confirm-password">Confirm New Password</label>
          <InputText id="confirm-password" v-model="confirmPassword" type="password" placeholder="Confirm new password"
            :disabled="isChangingPassword" @keypress.enter="handleChangePassword" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="closeChangePasswordDialog" :disabled="isChangingPassword" />
        <Button label="Change Password" @click="handleChangePassword" :disabled="isChangingPassword"
          :loading="isChangingPassword" />
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  width: 100%;
  overflow: hidden;
}

.app-sidebar {
  width: 64px;
  height: 100vh;
  background: var(--p-surface-0);
  border-right: 1px solid var(--p-surface-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.app-sidebar-header {
  min-height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid var(--p-surface-border);
}

.app-sidebar-logo {
  width: 2rem;
  height: 2rem;
  border-radius: var(--p-border-radius);
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.app-nav-link {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 0.75rem;
  border-radius: var(--p-border-radius);
  color: var(--p-text-color);
  text-decoration: none;
  transition: background-color 0.2s, color 0.2s;
}

.app-nav-link:hover {
  background: var(--p-surface-hover);
}

.app-nav-link-active {
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
}

.app-nav-link-active:hover {
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
}

.app-nav-link i {
  font-size: 1.25rem;
}

.app-sidebar-footer {
  border-top: 1px solid var(--p-surface-border);
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-user-btn {
  width: 2.5rem;
  height: 2.5rem;
  padding: 0;
}

.app-user-btn i {
  font-size: 1.25rem;
}

.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-header {
  height: 4rem;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  padding: 0 1rem;
  border-bottom: 1px solid var(--p-surface-border);
  background: var(--p-surface-0);
}

.app-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
  overflow-y: auto;
  background: var(--p-surface-ground);
}

.app-dialog-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 0.5rem 0;
}

.app-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.app-field label {
  font-size: 0.875rem;
  font-weight: 500;
}
</style>
