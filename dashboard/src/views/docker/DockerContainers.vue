<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Container {
  id: string;
  name: string;
  image: string;
  status: string;
  ports: string;
  created: string;
}

interface Props {
  nodeId?: string;
}

const props = defineProps<Props>();

const containers = ref<Container[]>([]);
const loading = ref(false);

// Mock data
const mockContainers: Container[] = [
  {
    id: "a1b2c3d4e5f6",
    name: "web-server",
    image: "nginx:latest",
    status: "running",
    ports: "80:80, 443:443",
    created: "2024-01-15 10:30:00",
  },
  {
    id: "b2c3d4e5f6a1",
    name: "database",
    image: "postgres:15",
    status: "running",
    ports: "5432:5432",
    created: "2024-01-15 10:25:00",
  },
  {
    id: "c3d4e5f6a1b2",
    name: "redis-cache",
    image: "redis:7",
    status: "running",
    ports: "6379:6379",
    created: "2024-01-15 10:20:00",
  },
  {
    id: "d4e5f6a1b2c3",
    name: "api-service",
    image: "node:18",
    status: "stopped",
    ports: "3000:3000",
    created: "2024-01-14 15:45:00",
  },
];

const fetchContainers = () => {
  loading.value = true;
  // Simulate API call
  setTimeout(() => {
    containers.value = [...mockContainers];
    loading.value = false;
  }, 300);
};

const getStatusColor = (status: string) => {
  return status === "running"
    ? "bg-green-500/10 text-green-500 border-green-500/20"
    : "bg-red-500/10 text-red-500 border-red-500/20";
};

onMounted(() => {
  fetchContainers();
});

// Watch for nodeId changes to refetch data
watch(
  () => props.nodeId,
  () => {
    if (props.nodeId) {
      fetchContainers();
    }
  }
);
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Containers</CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="loading" class="text-center py-8 text-muted-foreground">
        Loading...
      </div>
      <div v-else-if="containers.length > 0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Image</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Ports</TableHead>
              <TableHead>Created</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="container in containers" :key="container.id">
              <TableCell class="font-mono text-xs">
                {{ container.id.substring(0, 12) }}
              </TableCell>
              <TableCell class="font-medium">{{ container.name }}</TableCell>
              <TableCell>{{ container.image }}</TableCell>
              <TableCell>
                <span
                  :class="[
                    'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium',
                    getStatusColor(container.status),
                  ]"
                >
                  {{ container.status }}
                </span>
              </TableCell>
              <TableCell class="font-mono text-xs">{{ container.ports }}</TableCell>
              <TableCell class="text-muted-foreground">{{ container.created }}</TableCell>
              <TableCell>
                <div class="flex items-center gap-2">
                  <Button variant="ghost" size="sm" class="h-8 px-2">
                    <Icon icon="lucide:play" class="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" class="h-8 px-2">
                    <Icon icon="lucide:square" class="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" class="h-8 px-2">
                    <Icon icon="lucide:more-horizontal" class="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
      <div v-else class="text-center py-8 text-muted-foreground">
        <Icon icon="lucide:box" class="h-12 w-12 mx-auto mb-4 opacity-50" />
        <p>No containers found</p>
      </div>
    </CardContent>
  </Card>
</template>

