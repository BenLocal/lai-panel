<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
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
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import TooltipWithCopy from "@/components/application/TooltipWithCopy.vue";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { dockerApi, DockerUtils } from "@/api/docker";
import { nodeApi, type Node } from "@/api/node";
import { showToast } from "@/lib/toast";
import { ApiResponseHelper } from "@/api/base";

interface Image {
  id: string;
  repository: string;
  tag: string;
  size: string;
  created: string;
}

interface Props {
  nodeId?: string;
}

const props = defineProps<Props>();

const images = ref<Image[]>([]);
const loading = ref(false);
const isPushDialogOpen = ref(false);
const selectedImageForPush = ref<Image | null>(null);
const targetNodeId = ref<string>("");
const nodes = ref<Node[]>([]);
const pushLoading = ref(false);

// 计算可用的目标节点（排除当前节点）
const availableNodes = computed(() => {
  const currentNodeId = props.nodeId ? Number(props.nodeId) : null;
  if (currentNodeId === null) {
    return nodes.value;
  }
  return nodes.value.filter((n) => n.id !== currentNodeId);
});

const fetchImages = async () => {
  loading.value = true;
  const response = await dockerApi.images(Number(props.nodeId));
  images.value =
    response.data?.map((image) => {
      // image.repoTags is usually an array like ["repo:tag"]
      const id = image.Id;
      const { repository, tag } = DockerUtils.getImageRepository(
        image.RepoDigests,
        image.RepoTags
      );
      const size = DockerUtils.formatDisplaySize(image.Size);
      return {
        id: id,
        repository: repository,
        tag: tag,
        // Format size from bytes to human-readable string
        size: size,
        created: new Date(image.Created * 1000).toLocaleString(),
      };
    }) ?? [];
  loading.value = false;
};

const fetchNodes = async () => {
  try {
    const response = await nodeApi.list();
    if (ApiResponseHelper.isSuccess(response)) {
      nodes.value = response.data ?? [];
    } else {
      console.error("Failed to fetch nodes:", response.message);
      showToast("Failed to fetch nodes", "error");
    }
  } catch (error) {
    console.error("Error fetching nodes:", error);
    showToast("Failed to fetch nodes", "error");
  }
};

onMounted(() => {
  fetchImages();
  fetchNodes();
});

// Watch for nodeId changes to refetch data
watch(
  () => props.nodeId,
  () => {
    if (props.nodeId) {
      fetchImages();
    }
  }
);

const handleImageAction = async (
  action: "remove" | "inspect" | "push",
  image: Image
) => {
  const nodeId = Number(props.nodeId);
  if (Number.isNaN(nodeId)) {
    showToast("Invalid node ID", "error");
    return;
  }

  try {
    switch (action) {
      case "remove":
        // TODO: Implement image remove API
        showToast("Image remove not implemented yet", "info");
        // const response = await dockerApi.imageRemove(nodeId, image.id);
        // if (ApiResponseHelper.isSuccess(response)) {
        //   showToast("Image removed successfully", "success");
        //   await fetchImages();
        // } else {
        //   showToast(response.message ?? "Failed to remove image", "error");
        // }
        break;
      case "inspect":
        // TODO: Implement image inspect
        showToast("Image inspect not implemented yet", "info");
        break;
      case "push":
        // Open push dialog
        selectedImageForPush.value = image;
        targetNodeId.value = "";
        isPushDialogOpen.value = true;
        break;
      default:
        return;
    }
  } catch (error) {
    showToast(`Failed to ${action} image`, "error");
    console.error(`Failed to ${action} image:`, error);
  }
};

const handlePushImage = async () => {
  if (!selectedImageForPush.value || !targetNodeId.value) {
    showToast("Please select a target node", "error");
    return;
  }

  const sourceNodeId = Number(props.nodeId);
  const targetNodeIdNum = Number(targetNodeId.value);

  if (Number.isNaN(sourceNodeId) || Number.isNaN(targetNodeIdNum)) {
    showToast("Invalid node ID", "error");
    return;
  }

  if (sourceNodeId === targetNodeIdNum) {
    showToast("Source and target nodes cannot be the same", "error");
    return;
  }

  try {
    pushLoading.value = true;
    // TODO: Implement image push API
    showToast("Image push not implemented yet", "info");
    // const response = await dockerApi.imagePush(sourceNodeId, targetNodeIdNum, selectedImageForPush.value.id);
    // if (ApiResponseHelper.isSuccess(response)) {
    //   showToast("Image pushed successfully", "success");
    //   isPushDialogOpen.value = false;
    //   selectedImageForPush.value = null;
    //   targetNodeId.value = "";
    // } else {
    //   showToast(response.message ?? "Failed to push image", "error");
    // }
  } catch (error) {
    showToast("Failed to push image", "error");
    console.error("Failed to push image:", error);
  } finally {
    pushLoading.value = false;
  }
};

const closePushDialog = () => {
  isPushDialogOpen.value = false;
  selectedImageForPush.value = null;
  targetNodeId.value = "";
};
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Images</CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="loading" class="text-center py-8 text-muted-foreground">
        Loading...
      </div>
      <div v-else-if="images.length > 0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Repository</TableHead>
              <TableHead>Tag</TableHead>
              <TableHead>Image ID</TableHead>
              <TableHead>Size</TableHead>
              <TableHead>Created</TableHead>
              <TableHead class="sticky right-0 z-10 bg-background border-l">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="image in images" :key="image.id">
              <TableCell class="font-medium">
                <TooltipWithCopy :text="image.repository" max-width="300px" />
              </TableCell>
              <TableCell>
                <TooltipWithCopy :text="image.tag" max-width="200px">
                <span
                    class="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium bg-blue-500/10 text-blue-500 border-blue-500/20 truncate max-w-full"
                >
                  {{ image.tag }}
                </span>
                </TooltipWithCopy>
              </TableCell>
              <TableCell class="font-mono text-xs">
                {{ DockerUtils.getShortImageId(image.id) }}
              </TableCell>
              <TableCell>{{ image.size }}</TableCell>
              <TableCell class="text-muted-foreground">{{
                image.created
              }}</TableCell>
              <TableCell class="sticky right-0 z-10 bg-background border-l">
                <div class="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    class="h-8 px-2 text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
                    @click="handleImageAction('remove', image)"
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
                        @click="handleImageAction('inspect', image)"
                      >
                        <Icon icon="lucide:info" class="h-4 w-4 mr-2" />
                        Inspect
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        @click="handleImageAction('push', image)"
                      >
                        <Icon icon="lucide:upload" class="h-4 w-4 mr-2" />
                        Push to...
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        variant="destructive"
                        @click="handleImageAction('remove', image)"
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
        <Icon icon="lucide:layers" class="h-12 w-12 mx-auto mb-4 opacity-50" />
        <p>No images found</p>
      </div>
    </CardContent>
  </Card>

  <!-- Push Image Dialog -->
  <Dialog
    v-model:open="isPushDialogOpen"
    @update:open="(open) => !open && closePushDialog()"
  >
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Push Image to Node</DialogTitle>
        <DialogDescription>
          Select a target node to push the image to
        </DialogDescription>
      </DialogHeader>
      <div class="space-y-4 py-4">
        <div class="space-y-2">
          <Label class="text-sm font-medium">Image</Label>
          <div class="text-sm text-muted-foreground">
            <span class="font-medium">{{
              selectedImageForPush?.repository
            }}</span>
            <span class="mx-1">:</span>
            <span class="font-medium">{{ selectedImageForPush?.tag }}</span>
          </div>
        </div>
        <div class="space-y-2">
          <Label for="target-node">Target Node</Label>
          <Select v-model="targetNodeId">
            <SelectTrigger id="target-node">
              <SelectValue placeholder="Select a target node" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="node in availableNodes"
                :key="node.id"
                :value="node.id.toString()"
              >
                {{ node.display_name || node.name }} ({{ node.address }})
              </SelectItem>
              <div
                v-if="availableNodes.length === 0"
                class="px-2 py-1.5 text-sm text-muted-foreground"
              >
                No available nodes
              </div>
            </SelectContent>
          </Select>
        </div>
      </div>
      <DialogFooter>
        <Button
          variant="outline"
          @click="closePushDialog"
          :disabled="pushLoading"
        >
          Cancel
        </Button>
        <Button
          @click="handlePushImage"
          :disabled="pushLoading || !targetNodeId"
        >
          <Icon
            v-if="pushLoading"
            icon="lucide:loader-2"
            class="h-4 w-4 mr-2 animate-spin"
          />
          <Icon v-else icon="lucide:upload" class="h-4 w-4 mr-2" />
          {{ pushLoading ? "Pushing..." : "Push" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
