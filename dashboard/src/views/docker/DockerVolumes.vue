<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import TooltipWithCopy from "@/components/application/TooltipWithCopy.vue";
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
const searchQuery = ref("");
const menuRefs = ref<Record<string, any>>({});

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

// 过滤卷
const filteredVolumes = computed(() => {
  if (!searchQuery.value.trim()) {
    return volumes.value;
  }
  const query = searchQuery.value.toLowerCase();
  return volumes.value.filter(
    (volume) =>
      volume.name.toLowerCase().includes(query) ||
      volume.driver.toLowerCase().includes(query) ||
      volume.mountpoint.toLowerCase().includes(query)
  );
});

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
  <div class="docker-panel">
    <div class="docker-toolbar">
      <h3>Volumes</h3>
      <InputText v-model="searchQuery" placeholder="Search..." class="search-input" />
    </div>
    <div v-if="loading" class="docker-loading text-muted-foreground">Loading...</div>
    <div v-else-if="filteredVolumes.length > 0" class="table-wrap">
      <DataTable :value="filteredVolumes" size="small" striped-rows>
        <Column header="Name">
          <template #body="{ data }"><TooltipWithCopy :text="data.name" max-width="200px" /></template>
        </Column>
        <Column header="Driver">
          <template #body="{ data }"><Tag :value="data.driver" severity="info" /></template>
        </Column>
        <Column header="Mountpoint">
          <template #body="{ data }"><TooltipWithCopy :text="data.mountpoint" max-width="300px" /></template>
        </Column>
        <Column field="size" header="Size" />
        <Column field="created" header="Created">
          <template #body="{ data }"><span class="text-muted-foreground">{{ data.created }}</span></template>
        </Column>
        <Column header="Actions">
          <template #body="{ data }">
            <div class="action-btns">
              <Button text rounded size="small" severity="danger" @click="handleVolumeAction('remove', data)" v-tooltip.top="'Remove'">
                <i class="pi pi-trash"></i>
              </Button>
              <Menu :ref="(el: any) => { if (el) menuRefs[data.name] = el }" :model="[
                { label: 'Inspect', icon: 'pi pi-info', command: () => handleVolumeAction('inspect', data) },
                { separator: true },
                { label: 'Remove', icon: 'pi pi-trash', command: () => handleVolumeAction('remove', data) }
              ]" popup />
              <Button text rounded size="small" @click="(e) => menuRefs[data.name]?.toggle(e)">
                <i class="pi pi-ellipsis-h"></i>
              </Button>
            </div>
          </template>
        </Column>
      </DataTable>
    </div>
    <div v-else-if="!volumes.length" class="docker-empty text-muted-foreground">
      <i class="pi pi-database"></i>
      <p>No volumes found</p>
    </div>
    <div v-else class="docker-empty text-muted-foreground">
      <i class="pi pi-search"></i>
      <p>No volumes match your search</p>
    </div>
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
.docker-empty { text-align: center; padding: 2rem; }
.docker-empty i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
</style>
