<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import { authApi } from "@/api/auth";
import { showToast } from "@/lib/toast";
import { auth } from "@/auth";

const router = useRouter();
const route = useRoute();

const username = ref("");
const password = ref("");
const isLoading = ref(false);

onMounted(() => {
  if (auth.isAuthenticated()) {
    router.push((route.query.redirect as string) || "/dashboard");
  }
});

const handleLogin = async () => {
  if (!username.value || !password.value) {
    showToast("Please enter username and password", "error");
    return;
  }
  isLoading.value = true;
  try {
    const response = await authApi.login({
      username: username.value,
      password: btoa(password.value),
    });
    if (response.code === 0 && response.data?.token) {
      showToast("Login successful", "success");
      auth.setAuth(response.data.token, {
        user_id: response.data.user_id!,
        username: response.data.username!,
        role: response.data.role!,
        name: response.data.name!,
      });
      router.push((route.query.redirect as string) || "/dashboard");
    } else {
      showToast(response.message || "Login failed", "error");
    }
  } catch (e) {
    showToast("Login failed, please try again later", "error");
    console.error("Login error:", e);
  } finally {
    isLoading.value = false;
  }
};

const onKeyPress = (e: KeyboardEvent) => {
  if (e.key === "Enter") handleLogin();
};
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="login-logo"><i class="pi pi-th-large"></i></div>
        <h2>Welcome</h2>
        <p class="text-muted-foreground">Enter your username and password to continue</p>
      </div>
      <div class="login-form">
        <div class="login-field">
          <label for="username">Username</label>
          <InputText id="username" v-model="username" placeholder="Enter username"
            :disabled="isLoading" @keypress="onKeyPress" class="w-full" />
        </div>
        <div class="login-field">
          <label for="password">Password</label>
          <InputText id="password" v-model="password" type="password" placeholder="Enter password"
            :disabled="isLoading" @keypress="onKeyPress" class="w-full" />
        </div>
        <Button label="Login" :loading="isLoading" :disabled="isLoading" @click="handleLogin" class="w-full" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  background: var(--p-surface-ground);
}

.login-card {
  width: 100%;
  max-width: 24rem;
  padding: 1.5rem;
  background: var(--p-surface-card);
  border-radius: var(--p-border-radius);
  box-shadow: var(--p-overlay-shadow);
}

.login-header {
  text-align: center;
  margin-bottom: 1.5rem;
}

.login-logo {
  width: 4rem;
  height: 4rem;
  margin: 0 auto 1rem;
  border-radius: var(--p-border-radius);
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-logo i {
  font-size: 1.5rem;
}

.login-header h2 {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.login-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.login-field label {
  font-size: 0.875rem;
  font-weight: 500;
}

.w-full { width: 100%; }
</style>
