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



const fetchImages = async () => {
    loading.value = true;
    const response = await dockerApi.images(Number(props.nodeId));
    images.value = response.data?.map((image) => {
        // image.repoTags is usually an array like ["repo:tag"]
        const id = DockerUtils.getShortImageId(image.Id);
        const { repository, tag } = DockerUtils.getImageRepository(image.RepoDigests, image.RepoTags);
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

onMounted(() => {
    fetchImages();
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
                            <TableHead>Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        <TableRow v-for="image in images" :key="image.id">
                            <TableCell class="font-medium">{{ image.repository }}</TableCell>
                            <TableCell>
                                <span
                                    class="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium bg-blue-500/10 text-blue-500 border-blue-500/20">
                                    {{ image.tag }}
                                </span>
                            </TableCell>
                            <TableCell class="font-mono text-xs">
                                {{ image.id.substring(0, 12) }}
                            </TableCell>
                            <TableCell>{{ image.size }}</TableCell>
                            <TableCell class="text-muted-foreground">{{ image.created }}</TableCell>
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
                <Icon icon="lucide:layers" class="h-12 w-12 mx-auto mb-4 opacity-50" />
                <p>No images found</p>
            </div>
        </CardContent>
    </Card>
</template>
