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
import { dockerApi, DockerUtils } from "@/api/docker";

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
    volumes.value = response.data?.map((volume) => ({
        name: volume.Name,
        driver: volume.Driver ?? '',
        mountpoint: volume.Mountpoint ?? '',
        size: volume.Size ? DockerUtils.formatDisplaySize(volume.Size) : '',
        created: volume.CreatedAt ?? '',
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
                            <TableHead>Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        <TableRow v-for="volume in volumes" :key="volume.name">
                            <TableCell class="font-medium">{{ volume.name }}</TableCell>
                            <TableCell>
                                <span
                                    class="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium bg-gray-500/10 text-gray-500 border-gray-500/20">
                                    {{ volume.driver }}
                                </span>
                            </TableCell>
                            <TableCell class="font-mono text-xs text-muted-foreground">
                                {{ volume.mountpoint }}
                            </TableCell>
                            <TableCell>{{ volume.size }}</TableCell>
                            <TableCell class="text-muted-foreground">{{ volume.created }}</TableCell>
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
                <Icon icon="lucide:database" class="h-12 w-12 mx-auto mb-4 opacity-50" />
                <p>No volumes found</p>
            </div>
        </CardContent>
    </Card>
</template>
