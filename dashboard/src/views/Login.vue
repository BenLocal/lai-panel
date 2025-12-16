<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { authApi } from "@/api/auth";
import { showToast } from "@/lib/toast";
import { auth } from "@/auth";

const router = useRouter();
const route = useRoute();

const username = ref("");
const password = ref("");
const isLoading = ref(false);

// Check if user is already logged in
onMounted(() => {
  if (auth.isAuthenticated()) {
    // If already logged in, redirect to dashboard or the redirect path
    const redirect = (route.query.redirect as string) || "/dashboard";
    router.push(redirect);
  }
});

const handleLogin = async () => {
  if (!username.value || !password.value) {
    showToast("Please enter username and password", "error");
    return;
  }

  isLoading.value = true;

  try {
    const base64Password = btoa(password.value);
    const response = await authApi.login({
      username: username.value,
      password: base64Password,
    });

    if (response.code === 0 && response.data) {
      showToast("Login successful", "success");
      // Save token and user information
      if (response.data.token) {
        const userInfo = {
          user_id: response.data.user_id!,
          username: response.data.username!,
          role: response.data.role!,
          name: response.data.name!,
        };
        auth.setAuth(response.data.token, userInfo);
      }
      // Redirect to dashboard or the original target page
      const redirect = (route.query.redirect as string) || "/dashboard";
      router.push(redirect);
    } else {
      showToast(response.message || "Login failed, please check your username and password", "error");
    }
  } catch (error) {
    showToast("Login failed, please try again later", "error");
    console.error("Login error:", error);
  } finally {
    isLoading.value = false;
  }
};

const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === "Enter") {
    handleLogin();
  }
};
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gradient-to-br from-background via-background to-muted/20 p-4">
    <Card class="w-full max-w-md shadow-lg">
      <CardHeader class="space-y-1 text-center">
        <div class="flex justify-center mb-4">
          <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Icon icon="lucide:layers" class="h-6 w-6" />
          </div>
        </div>
        <CardTitle class="text-2xl font-bold">Welcome</CardTitle>
        <CardDescription>
          Please enter your username and password to continue
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-2">
          <Label for="username">Username</Label>
          <Input id="username" v-model="username" type="text" placeholder="Enter username" :disabled="isLoading"
            @keypress="handleKeyPress" class="h-11" />
        </div>
        <div class="space-y-2">
          <Label for="password">Password</Label>
          <Input id="password" v-model="password" type="password" placeholder="Enter password" :disabled="isLoading"
            @keypress="handleKeyPress" class="h-11" />
        </div>
        <Button class="w-full h-11" :disabled="isLoading" @click="handleLogin">
          <span v-if="!isLoading">Login</span>
          <span v-else class="flex items-center gap-2">
            <Icon icon="lucide:loader-2" class="h-4 w-4 animate-spin" />
            Logging in...
          </span>
        </Button>
      </CardContent>
    </Card>
  </div>
</template>
