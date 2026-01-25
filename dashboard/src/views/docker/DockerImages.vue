<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import Menu from 'primevue/menu'
import Select from 'primevue/select'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import TooltipWithCopy from "@/components/application/TooltipWithCopy.vue";
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
const searchQuery = ref("");
const isPushDialogOpen = ref(false);
const selectedImageForPush = ref<Image | null>(null);
const targetNodeId = ref<string>("");
const targetNode = ref<Node | null>(null);
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

// 过滤镜像
const filteredImages = computed(() => {
  if (!searchQuery.value.trim()) {
    return images.value;
  }
  const query = searchQuery.value.toLowerCase();
  return images.value.filter(
    (image) =>
      image.repository.toLowerCase().includes(query) ||
      image.tag.toLowerCase().includes(query) ||
      image.id.toLowerCase().includes(query)
  );
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
        targetNode.value = null;
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
  if (!selectedImageForPush.value || !targetNode.value) {
    showToast("Please select a target node", "error");
    return;
  }

  const sourceNodeId = Number(props.nodeId);
  const targetNodeIdNum = targetNode.value.id;

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
    //   targetNode.value = null;
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
  targetNode.value = null;
};
const menuRefs = ref<Record<string, any>>({});
</script>

<template>
  <div class="docker-panel">
    <div class="docker-toolbar">
      <h3>Images</h3>
      <InputText v-model="searchQuery" placeholder="Search..." class="search-input" />
    </div>
    <div v-if="loading" class="docker-loading text-muted-foreground">Loading...</div>
    <div v-else-if="filteredImages.length > 0" class="table-wrap">
      <DataTable :value="filteredImages" size="small" striped-rows>
        <Column header="Repository">
          <template #body="{ data }"><TooltipWithCopy :text="data.repository" max-width="300px" /></template>
        </Column>
        <Column header="Tag">
          <template #body="{ data }">
            <TooltipWithCopy :text="data.tag" max-width="200px">
              <Tag :value="data.tag" severity="info" />
            </TooltipWithCopy>
          </template>
        </Column>
        <Column header="Image ID">
          <template #body="{ data }"><span class="font-mono">{{ DockerUtils.getShortImageId(data.id) }}</span></template>
        </Column>
        <Column field="size" header="Size" />
        <Column field="created" header="Created">
          <template #body="{ data }"><span class="text-muted-foreground">{{ data.created }}</span></template>
        </Column>
        <Column header="Actions">
          <template #body="{ data }">
            <div class="action-btns">
              <Button text rounded size="small" severity="danger" @click="handleImageAction('remove', data)" v-tooltip.top="'Remove'">
                <i class="pi pi-trash"></i>
              </Button>
              <Menu :ref="(el: any) => { if (el) menuRefs[data.id] = el }" :model="[
                { label: 'Inspect', icon: 'pi pi-info', command: () => handleImageAction('inspect', data) },
                { separator: true },
                { label: 'Push to...', icon: 'pi pi-upload', command: () => handleImageAction('push', data) },
                { separator: true },
                { label: 'Remove', icon: 'pi pi-trash', command: () => handleImageAction('remove', data) }
              ]" popup />
              <Button text rounded size="small" @click="(e) => menuRefs[data.id]?.toggle(e)">
                <i class="pi pi-ellipsis-h"></i>
              </Button>
            </div>
          </template>
        </Column>
      </DataTable>
    </div>
    <div v-else-if="!images.length" class="docker-empty text-muted-foreground">
      <i class="pi pi-th-large"></i>
      <p>No images found</p>
    </div>
    <div v-else class="docker-empty text-muted-foreground">
      <i class="pi pi-search"></i>
      <p>No images match your search</p>
    </div>

    <Dialog v-model:visible="isPushDialogOpen" modal header="Push Image to Node" :style="{ width: '425px' }" @update:visible="(v: boolean) => !v && closePushDialog()">
      <p class="dialog-desc">Select a target node to push the image to.</p>
      <div class="form-group">
        <label>Image</label>
        <p class="text-muted-foreground">{{ selectedImageForPush?.repository }}:{{ selectedImageForPush?.tag }}</p>
      </div>
      <div class="form-group">
        <label for="target-node">Target Node</label>
        <Select id="target-node" v-model="targetNode" :options="availableNodes" option-label="display_name" placeholder="Select target node">
          <template #value="slot">
            <span v-if="slot.value">{{ slot.value.display_name || slot.value.name }} ({{ slot.value.address }})</span>
            <span v-else>{{ slot.placeholder }}</span>
          </template>
          <template #option="slot">
            {{ slot.option.display_name || slot.option.name }} ({{ slot.option.address }})
          </template>
          <template #empty>No available nodes</template>
        </Select>
      </div>
      <template #footer>
        <Button outlined @click="closePushDialog" :disabled="pushLoading">Cancel</Button>
        <Button @click="handlePushImage" :disabled="pushLoading || !targetNode" :loading="pushLoading">
          <i class="pi pi-upload"></i><span class="btn-icon-text">{{ pushLoading ? "Pushing..." : "Push" }}</span>
        </Button>
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.docker-panel { padding: 1rem; border: 1px solid var(--p-surface-border); border-radius: var(--p-border-radius); background: var(--p-surface-card); }
.docker-toolbar { display: flex; align-items: center; justify-content: space-between; gap: 1rem; margin-bottom: 1rem; }
.docker-toolbar h3 { font-size: 1.125rem; font-weight: 600; margin: 0; }
.search-input { width: 12rem; }
.docker-loading { text-align: center; padding: 2rem; }
.table-wrap { border-radius: var(--p-border-radius); overflow: hidden; }
.action-btns { display: flex; align-items: center; gap: 0.5rem; }
.font-mono { font-family: ui-monospace, monospace; font-size: 0.75rem; }
.docker-empty { text-align: center; padding: 2rem; }
.docker-empty i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
.dialog-desc { margin-bottom: 1rem; }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; margin-bottom: 1rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; }
.btn-icon-text { margin-left: 0.5rem; }
</style>
