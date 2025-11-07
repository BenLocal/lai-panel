<script setup lang="ts">
import { ref, onMounted } from "vue";
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
import {
  useVueTable,
  getCoreRowModel,
  getPaginationRowModel,
  type ColumnDef,
} from "@tanstack/vue-table";

interface Service {
  id: number;
  name: string;
  version: string;
  status: string;
  node: string;
  port: number;
  created_at: string;
  updated_at: string;
}

// Mock data
const mockServices: Service[] = [
  {
    id: 1,
    name: "web-service",
    version: "1.0.0",
    status: "running",
    node: "Node 01",
    port: 8080,
    created_at: "2024-01-15 10:00:00",
    updated_at: "2024-01-20 14:30:00",
  },
  {
    id: 2,
    name: "api-service",
    version: "2.1.0",
    status: "running",
    node: "Node 02",
    port: 3000,
    created_at: "2024-01-14 09:15:00",
    updated_at: "2024-01-19 16:45:00",
  },
  {
    id: 3,
    name: "database-service",
    version: "1.5.2",
    status: "stopped",
    node: "Node 03",
    port: 5432,
    created_at: "2024-01-13 11:20:00",
    updated_at: "2024-01-18 10:00:00",
  },
  {
    id: 4,
    name: "cache-service",
    version: "1.2.3",
    status: "running",
    node: "Local Node",
    port: 6379,
    created_at: "2024-01-12 08:30:00",
    updated_at: "2024-01-17 13:20:00",
  },
  {
    id: 5,
    name: "auth-service",
    version: "3.0.1",
    status: "running",
    node: "Node 01",
    port: 4000,
    created_at: "2024-01-11 15:45:00",
    updated_at: "2024-01-16 09:10:00",
  },
  {
    id: 6,
    name: "notification-service",
    version: "1.8.5",
    status: "running",
    node: "Node 02",
    port: 5000,
    created_at: "2024-01-10 12:00:00",
    updated_at: "2024-01-15 11:30:00",
  },
  {
    id: 7,
    name: "analytics-service",
    version: "2.5.0",
    status: "stopped",
    node: "Node 05",
    port: 6000,
    created_at: "2024-01-09 14:20:00",
    updated_at: "2024-01-14 08:15:00",
  },
  {
    id: 8,
    name: "search-service",
    version: "1.4.1",
    status: "running",
    node: "Node 06",
    port: 7000,
    created_at: "2024-01-08 10:30:00",
    updated_at: "2024-01-13 15:40:00",
  },
  {
    id: 9,
    name: "payment-service",
    version: "2.3.0",
    status: "running",
    node: "Node 01",
    port: 8000,
    created_at: "2024-01-07 09:00:00",
    updated_at: "2024-01-12 12:25:00",
  },
  {
    id: 10,
    name: "logging-service",
    version: "1.2.8",
    status: "stopped",
    node: "Node 04",
    port: 9000,
    created_at: "2024-01-06 11:15:00",
    updated_at: "2024-01-11 14:50:00",
  },
  {
    id: 11,
    name: "monitoring-service",
    version: "1.9.0",
    status: "running",
    node: "Local Node",
    port: 9100,
    created_at: "2024-01-05 13:45:00",
    updated_at: "2024-01-10 10:20:00",
  },
  {
    id: 12,
    name: "storage-service",
    version: "1.1.0",
    status: "running",
    node: "Node 07",
    port: 9200,
    created_at: "2024-01-04 08:20:00",
    updated_at: "2024-01-09 16:10:00",
  },
];

const services = ref<Service[]>([]);
const loading = ref(false);

// Fetch services
const fetchServices = () => {
  loading.value = true;
  // Simulate API delay
  setTimeout(() => {
    services.value = [...mockServices];
    loading.value = false;
  }, 300);
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    running: "bg-green-500/10 text-green-500 border-green-500/20",
    stopped: "bg-red-500/10 text-red-500 border-red-500/20",
    pending: "bg-yellow-500/10 text-yellow-500 border-yellow-500/20",
  };
  return colors[status] || colors.stopped;
};

// Define table columns
const columns: ColumnDef<Service>[] = [
  {
    accessorKey: "id",
    header: "ID",
  },
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "version",
    header: "Version",
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: (info) => {
      const status = info.getValue() as string;
      return { status };
    },
  },
  {
    accessorKey: "node",
    header: "Node",
  },
  {
    accessorKey: "port",
    header: "Port",
  },
  {
    accessorKey: "created_at",
    header: "Created At",
  },
  {
    accessorKey: "updated_at",
    header: "Updated At",
  },
];

// Create table instance
const table = useVueTable({
  get data() {
    return services.value;
  },
  columns,
  getCoreRowModel: getCoreRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
  initialState: {
    pagination: {
      pageSize: 10,
    },
  },
});

onMounted(() => {
  fetchServices();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold">Services</h1>
        <p class="text-muted-foreground mt-1">
          Manage and monitor your deployed services
        </p>
      </div>
      <Button>
        <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
        Deploy Service
      </Button>
    </div>

    <!-- Loading State -->
    <div
      v-if="loading && services.length === 0"
      class="text-center py-8 text-muted-foreground"
    >
      Loading...
    </div>

    <!-- Services Table -->
    <div v-else-if="services.length > 0" class="rounded-lg border bg-card">
      <Table>
        <TableHeader>
          <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
            <TableHead
              v-for="header in headerGroup.headers"
              :key="header.id"
              class="px-6"
            >
              <div v-if="!header.isPlaceholder">
                {{
                  typeof header.column.columnDef.header === "string"
                    ? header.column.columnDef.header
                    : header.column.id
                }}
              </div>
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="row in table.getRowModel().rows" :key="row.id">
            <TableCell
              v-for="cell in row.getVisibleCells()"
              :key="cell.id"
              class="px-6"
            >
              <template v-if="cell.column.id === 'status'">
                <span
                  :class="[
                    'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium',
                    getStatusColor((cell.getValue() as { status: string }).status),
                  ]"
                >
                  {{ (cell.getValue() as { status: string }).status }}
                </span>
              </template>
              <template v-else>
                {{ cell.getValue() }}
              </template>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <!-- Pagination Controls -->
      <div
        v-if="table.getPageCount() > 1"
        class="flex items-center justify-between border-t px-6 py-4"
      >
        <div class="text-sm text-muted-foreground">
          Showing
          {{
            table.getState().pagination.pageIndex *
              table.getState().pagination.pageSize +
            1
          }}
          -
          {{
            Math.min(
              (table.getState().pagination.pageIndex + 1) *
                table.getState().pagination.pageSize,
              services.length
            )
          }}
          of {{ services.length }} services
        </div>
        <div class="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            :disabled="!table.getCanPreviousPage()"
            @click="table.previousPage()"
          >
            <Icon icon="lucide:chevron-left" class="h-4 w-4" />
          </Button>
          <div class="flex items-center gap-1">
            <Button
              v-for="page in table.getPageCount()"
              :key="page"
              variant="outline"
              size="sm"
              :class="[
                'min-w-[40px]',
                table.getState().pagination.pageIndex + 1 === page
                  ? 'bg-primary text-primary-foreground'
                  : '',
              ]"
              @click="table.setPageIndex(page - 1)"
            >
              {{ page }}
            </Button>
          </div>
          <Button
            variant="outline"
            size="sm"
            :disabled="!table.getCanNextPage()"
            @click="table.nextPage()"
          >
            <Icon icon="lucide:chevron-right" class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="rounded-lg border bg-card p-12 text-center">
      <Icon
        icon="lucide:rocket"
        class="h-12 w-12 mx-auto text-muted-foreground mb-4"
      />
      <p class="text-muted-foreground">No services found</p>
      <Button class="mt-4">
        <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
        Deploy First Service
      </Button>
    </div>
  </div>
</template>

