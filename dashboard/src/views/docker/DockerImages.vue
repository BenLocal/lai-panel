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

// Mock data
const mockImages: Image[] = [
  {
    id: "sha256:abc123def456",
    repository: "nginx",
    tag: "latest",
    size: "150 MB",
    created: "2024-01-10 08:00:00",
  },
  {
    id: "sha256:def456ghi789",
    repository: "postgres",
    tag: "15",
    size: "380 MB",
    created: "2024-01-08 12:30:00",
  },
  {
    id: "sha256:ghi789jkl012",
    repository: "redis",
    tag: "7",
    size: "120 MB",
    created: "2024-01-05 14:20:00",
  },
  {
    id: "sha256:jkl012mno345",
    repository: "node",
    tag: "18",
    size: "950 MB",
    created: "2024-01-03 09:15:00",
  },
  {
    id: "sha256:mno345pqr678",
    repository: "ubuntu",
    tag: "22.04",
    size: "77 MB",
    created: "2024-01-01 10:00:00",
  },
];

const fetchImages = () => {
  loading.value = true;
  // Simulate API call
  setTimeout(() => {
    images.value = [...mockImages];
    loading.value = false;
  }, 300);
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
                <span class="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium bg-blue-500/10 text-blue-500 border-blue-500/20">
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

