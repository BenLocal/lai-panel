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

interface Network {
  id: string;
  name: string;
  driver: string;
  scope: string;
  subnet: string;
  gateway: string;
  containers: number;
}

interface Props {
  nodeId?: string;
}

const props = defineProps<Props>();

const networks = ref<Network[]>([]);
const loading = ref(false);

// Mock data
const mockNetworks: Network[] = [
  {
    id: "abc123def456",
    name: "bridge",
    driver: "bridge",
    scope: "local",
    subnet: "172.17.0.0/16",
    gateway: "172.17.0.1",
    containers: 3,
  },
  {
    id: "def456ghi789",
    name: "app_network",
    driver: "bridge",
    scope: "local",
    subnet: "172.18.0.0/16",
    gateway: "172.18.0.1",
    containers: 2,
  },
  {
    id: "ghi789jkl012",
    name: "database_network",
    driver: "bridge",
    scope: "local",
    subnet: "172.19.0.0/16",
    gateway: "172.19.0.1",
    containers: 1,
  },
  {
    id: "jkl012mno345",
    name: "host",
    driver: "host",
    scope: "local",
    subnet: "N/A",
    gateway: "N/A",
    containers: 0,
  },
];

const fetchNetworks = () => {
  loading.value = true;
  // Simulate API call
  setTimeout(() => {
    networks.value = [...mockNetworks];
    loading.value = false;
  }, 300);
};

onMounted(() => {
  fetchNetworks();
});

// Watch for nodeId changes to refetch data
watch(
  () => props.nodeId,
  () => {
    if (props.nodeId) {
      fetchNetworks();
    }
  }
);
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Networks</CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="loading" class="text-center py-8 text-muted-foreground">
        Loading...
      </div>
      <div v-else-if="networks.length > 0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Driver</TableHead>
              <TableHead>Scope</TableHead>
              <TableHead>Subnet</TableHead>
              <TableHead>Gateway</TableHead>
              <TableHead>Containers</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="network in networks" :key="network.id">
              <TableCell class="font-medium">{{ network.name }}</TableCell>
              <TableCell>
                <span
                  class="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium bg-blue-500/10 text-blue-500 border-blue-500/20"
                >
                  {{ network.driver }}
                </span>
              </TableCell>
              <TableCell>
                <span
                  class="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium bg-gray-500/10 text-gray-500 border-gray-500/20"
                >
                  {{ network.scope }}
                </span>
              </TableCell>
              <TableCell class="font-mono text-xs">{{
                network.subnet
              }}</TableCell>
              <TableCell class="font-mono text-xs">{{
                network.gateway
              }}</TableCell>
              <TableCell>
                <span class="font-medium">{{ network.containers }}</span>
              </TableCell>
              <TableCell>
                <div class="flex items-center gap-2">
                  <Button variant="ghost" size="sm" class="h-8 px-2">
                    <Icon icon="lucide:trash-2" class="h-4 w-4" />
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
        <Icon icon="lucide:network" class="h-12 w-12 mx-auto mb-4 opacity-50" />
        <p>No networks found</p>
      </div>
    </CardContent>
  </Card>
</template>
