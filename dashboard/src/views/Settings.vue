<script setup lang="ts">
import { ref, nextTick } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";

const isNavOpen = ref(false);
const activeSection = ref("general");

// Form data
const theme = ref("light");
const language = ref("en");
const logLevel = ref("info");
const enable2FA = ref(false);
const emailNotifications = ref(true);
const smsNotifications = ref(false);
const pushNotifications = ref(true);
const enableDebug = ref(false);

const sections = [
  { id: "general", title: "General Settings", icon: "lucide:settings" },
  { id: "security", title: "Security", icon: "lucide:shield" },
  { id: "notifications", title: "Notifications", icon: "lucide:bell" },
  { id: "appearance", title: "Appearance", icon: "lucide:palette" },
  { id: "advanced", title: "Advanced", icon: "lucide:sliders" },
];

const scrollToSection = async (sectionId: string) => {
  activeSection.value = sectionId;
  isNavOpen.value = false;

  // Wait for sheet to close and DOM to update
  await nextTick();

  // Add a delay to ensure sheet is fully closed and DOM is ready
  setTimeout(() => {
    const element = document.getElementById(sectionId);
    if (element) {
      // Use scrollIntoView which automatically finds the correct scroll container
      // The scroll-mt-20 class on the element provides offset for fixed headers
      element.scrollIntoView({
        behavior: "smooth",
        block: "start",
        inline: "nearest",
      });
    }
  }, 300);
};
</script>

<template>
  <div class="space-y-6 relative">
    <div>
      <h1 class="text-3xl font-bold">Settings</h1>
      <p class="text-muted-foreground mt-1">Manage your application settings</p>
    </div>

    <!-- Floating Navigation Button -->
    <Button class="fixed bottom-6 right-6 z-50 rounded-full h-12 w-12 shadow-lg" @click="isNavOpen = true">
      <Icon icon="lucide:menu" class="h-5 w-5" />
    </Button>

    <!-- Navigation Sheet -->
    <Sheet v-model:open="isNavOpen">
      <SheetContent side="right" class="w-[300px]">
        <SheetHeader>
          <SheetTitle>Settings Navigation</SheetTitle>
        </SheetHeader>
        <div class="mt-6 space-y-2">
          <Button v-for="section in sections" :key="section.id" variant="ghost" class="w-full justify-start"
            :class="activeSection === section.id ? 'bg-accent' : ''" @click="scrollToSection(section.id)">
            <Icon :icon="section.icon" class="h-4 w-4 mr-2" />
            {{ section.title }}
          </Button>
        </div>
      </SheetContent>
    </Sheet>

    <!-- General Settings -->
    <div id="general" class="scroll-mt-20">
      <Card>
        <CardHeader>
          <CardTitle>General Settings</CardTitle>
          <CardDescription>
            Configure general application settings
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="app-name" class="text-sm font-medium">Application Name</Label>
            <Input id="app-name" placeholder="Enter application name" value="Lai Panel" />
          </div>
          <div class="space-y-2">
            <Label for="app-description" class="text-sm font-medium">Description</Label>
            <Input id="app-description" placeholder="Enter description" value="Server management panel" />
          </div>
          <div class="space-y-2">
            <Label for="app-version" class="text-sm font-medium">Version</Label>
            <Input id="app-version" placeholder="1.0.0" value="1.0.0" />
          </div>
        </CardContent>
      </Card>
    </div>

    <Separator />

    <!-- Security Settings -->
    <div id="security" class="scroll-mt-20">
      <Card>
        <CardHeader>
          <CardTitle>Security</CardTitle>
          <CardDescription>
            Manage security and authentication settings
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="session-timeout" class="text-sm font-medium">Session Timeout (minutes)</Label>
            <Input id="session-timeout" type="number" placeholder="30" value="30" />
          </div>
          <div class="space-y-2">
            <Label for="max-login-attempts" class="text-sm font-medium">Max Login Attempts</Label>
            <Input id="max-login-attempts" type="number" placeholder="5" value="5" />
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox id="enable-2fa" v-model:checked="enable2FA" />
            <Label for="enable-2fa" class="text-sm font-medium cursor-pointer">Enable Two-Factor Authentication</Label>
          </div>
        </CardContent>
      </Card>
    </div>

    <Separator />

    <!-- Notifications Settings -->
    <div id="notifications" class="scroll-mt-20">
      <Card>
        <CardHeader>
          <CardTitle>Notifications</CardTitle>
          <CardDescription>
            Configure notification preferences
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="flex items-center space-x-2">
            <Checkbox id="email-notifications" v-model:checked="emailNotifications" />
            <Label for="email-notifications" class="text-sm font-medium cursor-pointer">Email Notifications</Label>
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox id="sms-notifications" v-model:checked="smsNotifications" />
            <Label for="sms-notifications" class="text-sm font-medium cursor-pointer">SMS Notifications</Label>
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox id="push-notifications" v-model:checked="pushNotifications" />
            <Label for="push-notifications" class="text-sm font-medium cursor-pointer">Push Notifications</Label>
          </div>
        </CardContent>
      </Card>
    </div>

    <Separator />

    <!-- Appearance Settings -->
    <div id="appearance" class="scroll-mt-20">
      <Card>
        <CardHeader>
          <CardTitle>Appearance</CardTitle>
          <CardDescription>
            Customize the appearance of the application
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="theme" class="text-sm font-medium">Theme</Label>
            <Select v-model="theme">
              <SelectTrigger id="theme">
                <SelectValue placeholder="Select theme" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="light">Light</SelectItem>
                <SelectItem value="dark">Dark</SelectItem>
                <SelectItem value="system">System</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-2">
            <Label for="language" class="text-sm font-medium">Language</Label>
            <Select v-model="language">
              <SelectTrigger id="language">
                <SelectValue placeholder="Select language" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="en">English</SelectItem>
                <SelectItem value="zh">中文</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>
    </div>

    <Separator />

    <!-- Advanced Settings -->
    <div id="advanced" class="scroll-mt-20">
      <Card>
        <CardHeader>
          <CardTitle>Advanced</CardTitle>
          <CardDescription> Advanced configuration options </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="api-timeout" class="text-sm font-medium">API Timeout (seconds)</Label>
            <Input id="api-timeout" type="number" placeholder="30" value="30" />
          </div>
          <div class="space-y-2">
            <Label for="log-level" class="text-sm font-medium">Log Level</Label>
            <Select v-model="logLevel">
              <SelectTrigger id="log-level">
                <SelectValue placeholder="Select log level" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="debug">Debug</SelectItem>
                <SelectItem value="info">Info</SelectItem>
                <SelectItem value="warn">Warning</SelectItem>
                <SelectItem value="error">Error</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox id="enable-debug" v-model:checked="enableDebug" />
            <Label for="enable-debug" class="text-sm font-medium cursor-pointer">Enable Debug Mode</Label>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
