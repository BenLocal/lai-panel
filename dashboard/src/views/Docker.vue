<script setup lang="ts">
import { ref, onMounted } from "vue";
import Select from 'primevue/select'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import Tag from 'primevue/tag'
import DockerContainers from "./docker/DockerContainers.vue";
import DockerImages from "./docker/DockerImages.vue";
import DockerVolumes from "./docker/DockerVolumes.vue";
import DockerNetworks from "./docker/DockerNetworks.vue";
import { nodeApi } from "@/api/node";

interface Node {
  id: number;
  name: string;
  display_name?: string | null;
  address: string;
  status?: string;
}

const nodes = ref<Node[]>([]);
const selectedNodeId = ref<string>("");
const selectedNode = ref<Node | null>(null);

const onNodeChange = (value: Node | null | undefined) => {
  if (value == null) {
    selectedNodeId.value = "";
    selectedNode.value = null;
    return;
  }
  selectedNodeId.value = value.id.toString();
  selectedNode.value = value;
};

const fetchNodes = async () => {
  const res = await nodeApi.list();
  nodes.value = res.data ?? [];
  if (nodes.value.length > 0) onNodeChange(nodes.value[0]);
};

onMounted(fetchNodes);
</script>

<template>
  <div class="page-root">
    <div class="page-header">
      <h1>Docker</h1>
      <p class="text-muted-foreground">Manage Docker containers and images</p>
    </div>

    <div class="docker-toolbar">
      <label for="docker-node-select">Select Node:</label>
      <Select
        v-model="selectedNode"
        :options="nodes"
        option-label="name"
        placeholder="Select a node"
        class="node-select"
        input-id="docker-node-select"
        @update:model-value="onNodeChange"
      >
        <template #value="slot">
          <span v-if="slot.value">{{ slot.value.display_name || slot.value.name }} ({{ slot.value.address }})</span>
          <span v-else>{{ slot.placeholder }}</span>
        </template>
      </Select>
      <Tag v-if="selectedNode" :value="selectedNode.status || 'offline'"
        :severity="selectedNode.status === 'online' ? 'success' : 'danger'" />
    </div>

    <TabView v-if="selectedNode">
      <TabPanel value="containers">
        <template #header>
          <span class="tab-header-inner"><i class="pi pi-box tab-icon"></i>Containers</span>
        </template>
        <DockerContainers :node-id="selectedNodeId" />
      </TabPanel>
      <TabPanel value="images">
        <template #header>
          <span class="tab-header-inner"><i class="pi pi-th-large tab-icon"></i>Images</span>
        </template>
        <DockerImages :node-id="selectedNodeId" />
      </TabPanel>
      <TabPanel value="volumes">
        <template #header>
          <span class="tab-header-inner"><i class="pi pi-database tab-icon"></i>Volumes</span>
        </template>
        <DockerVolumes :node-id="selectedNodeId" />
      </TabPanel>
      <TabPanel value="networks">
        <template #header>
          <span class="tab-header-inner"><i class="pi pi-sitemap tab-icon"></i>Networks</span>
        </template>
        <DockerNetworks :node-id="selectedNodeId" />
      </TabPanel>
    </TabView>

    <div v-else class="empty-state">
      <i class="pi pi-server"></i>
      <p>Please select a node to view Docker information</p>
    </div>
  </div>
</template>

<style scoped>
.page-root { display: flex; flex-direction: column; gap: 1.5rem; }
.page-header h1 { font-size: 1.5rem; font-weight: 700; margin-bottom: 0.25rem; }
.page-header p { font-size: 0.875rem; }

.docker-toolbar {
  display: flex;
  align-items: center;
  gap: 1rem;
}
.docker-toolbar label { font-size: 0.875rem; font-weight: 500; }
.node-select { width: 250px; }
.tab-header-inner { display: inline-flex; align-items: center; }
.tab-icon { margin-right: 0.5rem; flex-shrink: 0; }

.empty-state {
  padding: 3rem;
  text-align: center;
  background: var(--p-surface-card);
  border-radius: var(--p-border-radius);
}
.empty-state i { font-size: 3rem; opacity: 0.5; display: block; margin-bottom: 1rem; }
</style>
