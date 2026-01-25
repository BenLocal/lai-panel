<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
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
const searchQuery = ref("");
const menuRefs = ref<Record<string, any>>({});

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

// 过滤网络
const filteredNetworks = computed(() => {
  if (!searchQuery.value.trim()) {
    return networks.value;
  }
  const query = searchQuery.value.toLowerCase();
  return networks.value.filter(
    (network) =>
      network.name.toLowerCase().includes(query) ||
      network.driver.toLowerCase().includes(query) ||
      network.scope.toLowerCase().includes(query) ||
      network.subnet.toLowerCase().includes(query) ||
      network.gateway.toLowerCase().includes(query) ||
      network.id.toLowerCase().includes(query)
  );
});

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
  <div class="docker-panel">
    <div class="docker-toolbar">
      <h3>Networks</h3>
      <InputText v-model="searchQuery" placeholder="Search..." class="search-input" />
    </div>
    <div v-if="loading" class="docker-loading text-muted-foreground">Loading...</div>
    <div v-else-if="filteredNetworks.length > 0" class="table-wrap">
      <DataTable :value="filteredNetworks" size="small" striped-rows>
        <Column field="name" header="Name" />
        <Column header="Driver">
          <template #body="{ data }"><Tag :value="data.driver" severity="info" /></template>
        </Column>
        <Column header="Scope">
          <template #body="{ data }"><Tag :value="data.scope" severity="contrast" /></template>
        </Column>
        <Column header="Subnet">
          <template #body="{ data }"><span class="font-mono">{{ data.subnet }}</span></template>
        </Column>
        <Column header="Gateway">
          <template #body="{ data }"><span class="font-mono">{{ data.gateway }}</span></template>
        </Column>
        <Column field="containers" header="Containers" />
        <Column header="Actions">
          <template #body="{ data }">
            <div class="action-btns">
              <Button text rounded size="small" severity="danger" @click="handleNetworkAction('remove', data)" v-tooltip.top="'Remove'">
                <i class="pi pi-trash"></i>
              </Button>
              <Menu :ref="(el: any) => { if (el) menuRefs[data.id] = el }" :model="[
                { label: 'Inspect', icon: 'pi pi-info', command: () => handleNetworkAction('inspect', data) },
                { separator: true },
                { label: 'Remove', icon: 'pi pi-trash', command: () => handleNetworkAction('remove', data) }
              ]" popup />
              <Button text rounded size="small" @click="(e) => menuRefs[data.id]?.toggle(e)">
                <i class="pi pi-ellipsis-h"></i>
              </Button>
            </div>
          </template>
        </Column>
      </DataTable>
    </div>
    <div v-else-if="!networks.length" class="docker-empty text-muted-foreground">
      <i class="pi pi-sitemap"></i>
      <p>No networks found</p>
    </div>
    <div v-else class="docker-empty text-muted-foreground">
      <i class="pi pi-search"></i>
      <p>No networks match your search</p>
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
.font-mono { font-family: ui-monospace, monospace; font-size: 0.75rem; }
.docker-empty { text-align: center; padding: 2rem; }
.docker-empty i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
</style>
