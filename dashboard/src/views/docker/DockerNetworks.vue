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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { dockerApi, DockerUtils } from "@/api/docker";
import { showToast } from "@/lib/toast";
import { ApiResponseHelper } from "@/api/base";

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

const fetchNetworks = async () => {
  loading.value = true;
  const response = await dockerApi.networks(Number(props.nodeId));
  networks.value =
    response.data?.map((network) => {
      const id = DockerUtils.getShortNetworkId(network.Id);
      const subnet = network.IPAM?.Config?.[0]?.Subnet ?? "";
      const gateway = network.IPAM?.Config?.[0]?.Gateway ?? "";
      return {
        id: id,
        name: network.Name,
        driver: network.Driver,
        scope: network.Scope,
        subnet: subnet,
        gateway: gateway,
        containers: network.Containers.length,
      };
    }) ?? [];
  loading.value = false;
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

const handleNetworkAction = async (
  action: "remove" | "inspect",
  network: Network
) => {
  const nodeId = Number(props.nodeId);
  if (Number.isNaN(nodeId)) {
    showToast("Invalid node ID", "error");
    return;
  }

  try {
    switch (action) {
      case "remove":
        // TODO: Implement network remove API
        showToast("Network remove not implemented yet", "info");
        // const response = await dockerApi.networkRemove(nodeId, network.id);
        // if (ApiResponseHelper.isSuccess(response)) {
        //   showToast("Network removed successfully", "success");
        //   await fetchNetworks();
        // } else {
        //   showToast(response.message ?? "Failed to remove network", "error");
        // }
        break;
      case "inspect":
        // TODO: Implement network inspect
        showToast("Network inspect not implemented yet", "info");
        break;
      default:
        return;
    }
  } catch (error) {
    showToast(`Failed to ${action} network`, "error");
    console.error(`Failed to ${action} network:`, error);
  }
};
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
                  <Button
                    variant="ghost"
                    size="sm"
                    class="h-8 px-2 text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
                    @click="handleNetworkAction('remove', network)"
                  >
                    <Icon icon="lucide:trash-2" class="h-4 w-4" />
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button variant="ghost" size="sm" class="h-8 px-2">
                        <Icon icon="lucide:more-horizontal" class="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                      <DropdownMenuItem
                        @click="handleNetworkAction('inspect', network)"
                      >
                        <Icon icon="lucide:info" class="h-4 w-4 mr-2" />
                        Inspect
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        variant="destructive"
                        @click="handleNetworkAction('remove', network)"
                      >
                        <Icon icon="lucide:trash-2" class="h-4 w-4 mr-2" />
                        Remove
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
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
