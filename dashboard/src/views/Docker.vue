<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Icon } from "@iconify/vue";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
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

// Watch for selected node changes
const onNodeChange = (value: string | number | bigint | Record<string, any> | null) => {
    if (value === null || value === undefined) {
        selectedNodeId.value = "";
        selectedNode.value = null;
        return;
    }
    const stringValue = typeof value === 'string' ? value : String(value);
    selectedNodeId.value = stringValue;
    selectedNode.value =
        nodes.value.find((n) => n.id.toString() === stringValue) || null;
};

// Fetch nodes
const fetchNodes = async () => {
    const response = await nodeApi.list();

    nodes.value = response.data ?? [];
    if (nodes.value.length > 0 && nodes.value[0]) {
        onNodeChange(nodes.value[0].id.toString());
    }
};

onMounted(() => {
    fetchNodes();
});
</script>

<template>
    <div class="space-y-6">
        <div class="flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-bold">Docker</h1>
                <p class="text-muted-foreground mt-1">
                    Manage Docker containers and images
                </p>
            </div>
        </div>

        <!-- Node Selector -->
        <div class="flex items-center gap-4">
            <label class="text-sm font-medium">Select Node:</label>
            <Select :model-value="selectedNodeId" @update:model-value="onNodeChange">
                <SelectTrigger class="w-[250px]">
                    <SelectValue placeholder="Select a node" />
                </SelectTrigger>
                <SelectContent>
                    <SelectItem v-for="node in nodes" :key="node.id" :value="node.id.toString()">
                        {{ node.display_name || node.name }} ({{ node.address }})
                    </SelectItem>
                </SelectContent>
            </Select>
            <div v-if="selectedNode" class="text-sm text-muted-foreground">
                <span :class="[
                    'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium',
                    selectedNode.status === 'online'
                        ? 'bg-green-500/10 text-green-500 border-green-500/20'
                        : 'bg-red-500/10 text-red-500 border-red-500/20',
                ]">
                    <span :class="[
                        'h-1.5 w-1.5 rounded-full',
                        selectedNode.status === 'online' ? 'bg-green-500' : 'bg-red-500',
                    ]"></span>
                    {{ selectedNode.status || "offline" }}
                </span>
            </div>
        </div>

        <!-- Tabs Content -->
        <div v-if="selectedNode">
            <Tabs default-value="containers" class="w-full">
                <TabsList class="grid w-full grid-cols-4">
                    <TabsTrigger value="containers">
                        <Icon icon="lucide:box" class="h-4 w-4 mr-2" />
                        Containers
                    </TabsTrigger>
                    <TabsTrigger value="images">
                        <Icon icon="lucide:layers" class="h-4 w-4 mr-2" />
                        Images
                    </TabsTrigger>
                    <TabsTrigger value="volumes">
                        <Icon icon="lucide:database" class="h-4 w-4 mr-2" />
                        Volumes
                    </TabsTrigger>
                    <TabsTrigger value="networks">
                        <Icon icon="lucide:network" class="h-4 w-4 mr-2" />
                        Networks
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="containers" class="mt-6">
                    <DockerContainers :node-id="selectedNodeId" />
                </TabsContent>

                <TabsContent value="images" class="mt-6">
                    <DockerImages :node-id="selectedNodeId" />
                </TabsContent>

                <TabsContent value="volumes" class="mt-6">
                    <DockerVolumes :node-id="selectedNodeId" />
                </TabsContent>

                <TabsContent value="networks" class="mt-6">
                    <DockerNetworks :node-id="selectedNodeId" />
                </TabsContent>
            </Tabs>
        </div>

        <!-- No Node Selected -->
        <div v-else class="rounded-lg border bg-card p-12 text-center">
            <Icon icon="lucide:server" class="h-12 w-12 mx-auto text-muted-foreground mb-4" />
            <p class="text-muted-foreground">Please select a node to view Docker information</p>
        </div>
    </div>
</template>
