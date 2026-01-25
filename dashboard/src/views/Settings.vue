<script setup lang="ts">
import { ref, nextTick } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Checkbox from 'primevue/checkbox'
import Select from 'primevue/select'
import Sidebar from 'primevue/sidebar'

const isNavOpen = ref(false);
const activeSection = ref("general");
const theme = ref("light");
const language = ref("en");
const logLevel = ref("info");
const enable2FA = ref(false);
const emailNotifications = ref(true);
const smsNotifications = ref(false);
const pushNotifications = ref(true);
const enableDebug = ref(false);
const appName = ref("Lai Panel");
const appDesc = ref("Server management panel");
const appVer = ref("1.0.0");
const sessionTimeout = ref(30);
const maxAttempts = ref(5);
const apiTimeout = ref(30);

const themeOptions = [
  { label: 'Light', value: 'light' },
  { label: 'Dark', value: 'dark' },
  { label: 'System', value: 'system' }
];

const languageOptions = [
  { label: 'English', value: 'en' },
  { label: '中文', value: 'zh' }
];

const logLevelOptions = [
  { label: 'Debug', value: 'debug' },
  { label: 'Info', value: 'info' },
  { label: 'Warning', value: 'warn' },
  { label: 'Error', value: 'error' }
];

const sections = [
  { id: "general", title: "General", icon: "pi-cog" },
  { id: "security", title: "Security", icon: "pi-shield" },
  { id: "notifications", title: "Notifications", icon: "pi-bell" },
  { id: "appearance", title: "Appearance", icon: "pi-palette" },
  { id: "advanced", title: "Advanced", icon: "pi-sliders-h" },
];

const scrollToSection = async (id: string) => {
  activeSection.value = id;
  isNavOpen.value = false;
  await nextTick();
  setTimeout(() => {
    document.getElementById(id)?.scrollIntoView({ behavior: "smooth", block: "start" });
  }, 300);
};
</script>

<template>
  <div class="page-root">
    <div class="page-header">
      <h1>Settings</h1>
      <p class="text-muted-foreground">Manage your application settings</p>
    </div>

    <Button class="nav-fab" rounded @click="isNavOpen = true" v-tooltip.left="'Navigation'">
      <i class="pi pi-bars"></i>
    </Button>

    <Sidebar v-model:visible="isNavOpen" position="right" :style="{ width: '300px' }">
      <h2 class="sidebar-title">Settings</h2>
      <div class="nav-list">
        <Button
          v-for="s in sections"
          :key="s.id"
          text
          class="nav-item"
          :class="{ 'nav-item-active': activeSection === s.id }"
          @click="scrollToSection(s.id)"
        >
          <i :class="'pi ' + s.icon"></i>
          <span>{{ s.title }}</span>
        </Button>
      </div>
    </Sidebar>

    <section id="general" class="settings-section">
      <div class="panel-card">
        <div class="panel-header">
          <h3>General</h3>
          <p class="text-muted-foreground">Configure general application settings</p>
        </div>
        <div class="form-list">
          <div class="form-group">
            <label for="app-name">Application Name</label>
            <InputText id="app-name" v-model="appName" placeholder="Enter name" />
          </div>
          <div class="form-group">
            <label for="app-desc">Description</label>
            <InputText id="app-desc" v-model="appDesc" placeholder="Enter description" />
          </div>
          <div class="form-group">
            <label for="app-ver">Version</label>
            <InputText id="app-ver" v-model="appVer" placeholder="1.0.0" />
          </div>
        </div>
      </div>
    </section>

    <section id="security" class="settings-section">
      <div class="panel-card">
        <div class="panel-header">
          <h3>Security</h3>
          <p class="text-muted-foreground">Authentication and security</p>
        </div>
        <div class="form-list">
          <div class="form-group">
            <label for="session-timeout">Session Timeout (min)</label>
            <InputNumber id="session-timeout" v-model="sessionTimeout" placeholder="30" />
          </div>
          <div class="form-group">
            <label for="max-attempts">Max Login Attempts</label>
            <InputNumber id="max-attempts" v-model="maxAttempts" placeholder="5" />
          </div>
          <div class="form-group flex-row">
            <Checkbox v-model="enable2FA" binary inputId="enable-2fa" />
            <label for="enable-2fa">Two-Factor Authentication</label>
          </div>
        </div>
      </div>
    </section>

    <section id="notifications" class="settings-section">
      <div class="panel-card">
        <div class="panel-header">
          <h3>Notifications</h3>
          <p class="text-muted-foreground">Notification preferences</p>
        </div>
        <div class="form-list">
          <div class="form-group flex-row">
            <Checkbox v-model="emailNotifications" binary inputId="email-notify" />
            <label for="email-notify">Email</label>
          </div>
          <div class="form-group flex-row">
            <Checkbox v-model="smsNotifications" binary inputId="sms-notify" />
            <label for="sms-notify">SMS</label>
          </div>
          <div class="form-group flex-row">
            <Checkbox v-model="pushNotifications" binary inputId="push-notify" />
            <label for="push-notify">Push</label>
          </div>
        </div>
      </div>
    </section>

    <section id="appearance" class="settings-section">
      <div class="panel-card">
        <div class="panel-header">
          <h3>Appearance</h3>
          <p class="text-muted-foreground">Theme and language</p>
        </div>
        <div class="form-list">
          <div class="form-group">
            <label for="theme">Theme</label>
            <Select v-model="theme" :options="themeOptions" option-label="label" option-value="value" placeholder="Select" />
          </div>
          <div class="form-group">
            <label for="lang">Language</label>
            <Select v-model="language" :options="languageOptions" option-label="label" option-value="value" placeholder="Select" />
          </div>
        </div>
      </div>
    </section>

    <section id="advanced" class="settings-section">
      <div class="panel-card">
        <div class="panel-header">
          <h3>Advanced</h3>
          <p class="text-muted-foreground">Advanced configuration</p>
        </div>
        <div class="form-list">
          <div class="form-group">
            <label for="api-timeout">API Timeout (s)</label>
            <InputNumber id="api-timeout" v-model="apiTimeout" placeholder="30" />
          </div>
          <div class="form-group">
            <label for="log-level">Log Level</label>
            <Select v-model="logLevel" :options="logLevelOptions" option-label="label" option-value="value" placeholder="Select" />
          </div>
          <div class="form-group flex-row">
            <Checkbox v-model="enableDebug" binary inputId="enable-debug" />
            <label for="enable-debug">Debug Mode</label>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.page-root { display: flex; flex-direction: column; gap: 1.5rem; position: relative; }
.page-header h1 { font-size: 1.5rem; font-weight: 700; margin-bottom: 0.25rem; }
.page-header p { font-size: 0.875rem; }

.nav-fab { position: fixed; bottom: 1.5rem; right: 1.5rem; z-index: 50; width: 3rem; height: 3rem; box-shadow: var(--p-overlay-shadow); }

.sidebar-title { font-size: 1.125rem; font-weight: 600; margin-bottom: 1.5rem; }
.nav-list { display: flex; flex-direction: column; gap: 0.5rem; }
.nav-item { justify-content: flex-start; width: 100%; }
.nav-item i { margin-right: 0.5rem; }
.nav-item-active { background: var(--p-surface-hover); }

.settings-section { scroll-margin-top: 2rem; }
.panel-card { padding: 1.5rem; background: var(--p-surface-card); border-radius: var(--p-border-radius); }
.panel-header { margin-bottom: 1rem; }
.panel-header h3 { font-size: 1.125rem; font-weight: 600; margin-bottom: 0.25rem; }
.panel-header p { font-size: 0.875rem; }
.form-list { display: flex; flex-direction: column; gap: 1rem; }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; }
.form-group.flex-row { flex-direction: row; align-items: center; gap: 0.5rem; }
</style>
