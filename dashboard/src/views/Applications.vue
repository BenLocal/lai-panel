<script setup lang="ts">
import { ref, computed } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";

interface Application {
  id: number;
  name: string;
  description: string;
  status: string;
  version: string;
  icon: string;
}

// Mock data - expanded for pagination
const allApplications: Application[] = [
  {
    id: 1,
    name: "Web Application",
    description: "A modern web application built with Vue.js",
    status: "running",
    version: "1.2.3",
    icon: "lucide:globe",
  },
  {
    id: 2,
    name: "API Service",
    description: "RESTful API service for backend operations",
    status: "running",
    version: "2.0.1",
    icon: "lucide:server",
  },
  {
    id: 3,
    name: "Database Service",
    description: "PostgreSQL database management service",
    status: "stopped",
    version: "1.0.0",
    icon: "lucide:database",
  },
  {
    id: 4,
    name: "Cache Service",
    description: "Redis cache service for improved performance",
    status: "running",
    version: "1.5.2",
    icon: "lucide:zap",
  },
  {
    id: 5,
    name: "Message Queue",
    description: "RabbitMQ message queue service",
    status: "running",
    version: "1.3.4",
    icon: "lucide:message-square",
  },
  {
    id: 6,
    name: "File Storage",
    description: "S3-compatible file storage service",
    status: "stopped",
    version: "1.1.0",
    icon: "lucide:folder",
  },
  {
    id: 7,
    name: "Authentication Service",
    description: "OAuth2 and JWT authentication service",
    status: "running",
    version: "2.1.0",
    icon: "lucide:shield",
  },
  {
    id: 8,
    name: "Notification Service",
    description: "Email and SMS notification service",
    status: "running",
    version: "1.8.5",
    icon: "lucide:bell",
  },
  {
    id: 9,
    name: "Analytics Service",
    description: "Real-time analytics and reporting service",
    status: "stopped",
    version: "3.0.2",
    icon: "lucide:bar-chart",
  },
  {
    id: 10,
    name: "Search Service",
    description: "Elasticsearch-based search service",
    status: "running",
    version: "1.4.1",
    icon: "lucide:search",
  },
  {
    id: 11,
    name: "Payment Gateway",
    description: "Payment processing and gateway service",
    status: "running",
    version: "2.5.0",
    icon: "lucide:credit-card",
  },
  {
    id: 12,
    name: "Logging Service",
    description: "Centralized logging and monitoring service",
    status: "stopped",
    version: "1.2.8",
    icon: "lucide:file-text",
  },
];

const applications = ref<Application[]>(allApplications);

// Pagination
const currentPage = ref(1);
const pageSize = ref(6); // 6 cards per page for 3 columns layout

const paginatedApplications = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return applications.value.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(applications.value.length / pageSize.value);
});

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
  }
};

const getStatusColor = (status: string) => {
  return status === "running" ? "text-green-500" : "text-red-500";
};

const getStatusBadgeColor = (status: string) => {
  return status === "running"
    ? "bg-green-500/10 text-green-500 border-green-500/20"
    : "bg-red-500/10 text-red-500 border-red-500/20";
};
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold">Applications</h1>
        <p class="text-muted-foreground mt-1">
          Manage and monitor your applications
        </p>
      </div>
      <Button>
        <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
        New Application
      </Button>
    </div>

    <!-- Applications Cards -->
    <div v-if="applications.length > 0">
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="app in paginatedApplications"
          :key="app.id"
          class="group rounded-lg border bg-card p-6 hover:shadow-md transition-shadow cursor-pointer"
        >
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center gap-3">
              <div
                class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary"
              >
                <Icon :icon="app.icon" class="h-5 w-5" />
              </div>
              <div>
                <h3 class="font-semibold text-lg">{{ app.name }}</h3>
                <p class="text-xs text-muted-foreground">v{{ app.version }}</p>
              </div>
            </div>
          </div>

          <p class="text-sm text-muted-foreground mb-4 line-clamp-2">
            {{ app.description }}
          </p>

          <div class="flex items-center justify-between">
            <span
              :class="[
                'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium',
                getStatusBadgeColor(app.status),
              ]"
            >
              <span
                :class="[
                  'mr-1.5 h-1.5 w-1.5 rounded-full',
                  getStatusColor(app.status),
                ]"
              ></span>
              {{ app.status }}
            </span>
            <Button variant="ghost" size="sm">
              <Icon icon="lucide:more-horizontal" class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      <!-- Pagination Controls -->
      <div
        v-if="totalPages > 1"
        class="flex items-center justify-between border-t pt-6 mt-6"
      >
        <div class="text-sm text-muted-foreground">
          Showing {{ (currentPage - 1) * pageSize + 1 }} -
          {{ Math.min(currentPage * pageSize, applications.length) }} of
          {{ applications.length }} applications
        </div>
        <div class="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            :disabled="currentPage === 1"
            @click="goToPage(currentPage - 1)"
          >
            <Icon icon="lucide:chevron-left" class="h-4 w-4" />
          </Button>
          <div class="flex items-center gap-1">
            <Button
              v-for="page in totalPages"
              :key="page"
              variant="outline"
              size="sm"
              :class="[
                'min-w-[40px]',
                currentPage === page
                  ? 'bg-primary text-primary-foreground'
                  : '',
              ]"
              @click="goToPage(page)"
            >
              {{ page }}
            </Button>
          </div>
          <Button
            variant="outline"
            size="sm"
            :disabled="currentPage === totalPages"
            @click="goToPage(currentPage + 1)"
          >
            <Icon icon="lucide:chevron-right" class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="rounded-lg border bg-card p-12 text-center">
      <Icon
        icon="lucide:layers"
        class="h-12 w-12 mx-auto text-muted-foreground mb-4"
      />
      <p class="text-muted-foreground">No applications found</p>
      <Button class="mt-4">
        <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
        Add First Application
      </Button>
    </div>
  </div>
</template>
