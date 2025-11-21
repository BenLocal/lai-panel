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
import TooltipWithCopy from "@/components/application/TooltipWithCopy.vue";
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

interface Volume {
  name: string;
  driver: string;
  mountpoint: string;
  size: string;
  created: string;
}

interface Props {
  nodeId?: string;
}

const props = defineProps<Props>();

const volumes = ref<Volume[]>([]);
const loading = ref(false);

const fetchVolumes = async () => {
  loading.value = true;

  const response = await dockerApi.volumes(Number(props.nodeId));
  volumes.value =
    response.data?.map((volume) => ({
      name: volume.Name,
      driver: volume.Driver ?? "",
      mountpoint: volume.Mountpoint ?? "",
      size: volume.Size ? DockerUtils.formatDisplaySize(volume.Size) : "",
      created: volume.CreatedAt ?? "",
    })) ?? [];
  loading.value = false;
};

onMounted(() => {
  fetchVolumes();
});

// Watch for nodeId changes to refetch data
watch(
  () => props.nodeId,
  () => {
    if (props.nodeId) {
      fetchVolumes();
    }
  }
);

const handleVolumeAction = async (
  action: "remove" | "inspect",
  volume: Volume
) => {
  const nodeId = Number(props.nodeId);
  if (Number.isNaN(nodeId)) {
    showToast("Invalid node ID", "error");
    return;
  }

  try {
    switch (action) {
      case "remove":
        // TODO: Implement volume remove API
        showToast("Volume remove not implemented yet", "info");
        // const response = await dockerApi.volumeRemove(nodeId, volume.name);
        // if (ApiResponseHelper.isSuccess(response)) {
        //   showToast("Volume removed successfully", "success");
        //   await fetchVolumes();
        // } else {
        //   showToast(response.message ?? "Failed to remove volume", "error");
        // }
        break;
      case "inspect":
        // TODO: Implement volume inspect
        showToast("Volume inspect not implemented yet", "info");
        break;
      default:
        return;
    }
  } catch (error) {
    showToast(`Failed to ${action} volume`, "error");
    console.error(`Failed to ${action} volume:`, error);
  }
};
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Volumes</CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="loading" class="text-center py-8 text-muted-foreground">
        Loading...
      </div>
      <div v-else-if="volumes.length > 0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Driver</TableHead>
              <TableHead>Mountpoint</TableHead>
              <TableHead>Size</TableHead>
              <TableHead>Created</TableHead>
              <TableHead class="sticky right-0 z-10 bg-background border-l">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="volume in volumes" :key="volume.name">
              <TableCell class="font-medium">
                <TooltipWithCopy :text="volume.name" max-width="200px" />
              </TableCell>
              <TableCell>
                <span
                  class="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium bg-gray-500/10 text-gray-500 border-gray-500/20">
                  {{ volume.driver }}
                </span>
              </TableCell>
              <TableCell class="font-mono text-xs text-muted-foreground">
                <TooltipWithCopy :text="volume.mountpoint" max-width="300px" />
              </TableCell>
              <TableCell>{{ volume.size }}</TableCell>
              <TableCell class="text-muted-foreground">{{
                volume.created
              }}</TableCell>
              <TableCell class="sticky right-0 z-10 bg-background border-l">
                <div class="flex items-center gap-2">
                  <Button variant="ghost" size="sm"
                    class="h-8 px-2 text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
                    @click="handleVolumeAction('remove', volume)">
                    <Icon icon="lucide:trash-2" class="h-4 w-4" />
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button variant="ghost" size="sm" class="h-8 px-2">
                        <Icon icon="lucide:more-horizontal" class="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                      <DropdownMenuItem @click="handleVolumeAction('inspect', volume)">
                        <Icon icon="lucide:info" class="h-4 w-4 mr-2" />
                        Inspect
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem variant="destructive" @click="handleVolumeAction('remove', volume)">
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
        <Icon icon="lucide:database" class="h-12 w-12 mx-auto mb-4 opacity-50" />
        <p>No volumes found</p>
      </div>
    </CardContent>
  </Card>
</template>
