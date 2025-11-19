<script setup lang="ts">
import { ref, onMounted, computed, nextTick, watch } from "vue";
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
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerDescription,
} from "@/components/ui/drawer";
import {
  Stepper,
  StepperItem,
  StepperTrigger,
  StepperTitle,
  StepperDescription,
  StepperSeparator,
} from "@/components/ui/stepper";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import {
  serviceApi,
  type DeployServiceRequest,
  type SaveServiceRequest,
  type Service,
} from "@/api/service";
import { applicationApi, type Application } from "@/api/application";
import { nodeApi, type Node } from "@/api/node";
import { showToast } from "@/lib/toast";
import { ApiResponseHelper } from "@/api/base";

const services = ref<Service[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
const totalPages = computed(() => Math.ceil(total.value / pageSize.value));

// Deploy dialog state
const isDeployDialogOpen = ref(false);
const currentStep = ref(1);
const previousClickStep = ref(-1);
const selectedApplication = ref<Application | null>(null);
const selectedNode = ref<Node | null>(null);
const qaValues = ref<Record<string, string>>({});
const applications = ref<Application[]>([]);
const nodes = ref<Node[]>([]);
const deployLoading = ref(false);
const applicationsLoading = ref(false);
const nodesLoading = ref(false);
const deployOutput = ref<string[]>([]);
const deployOutputRef = ref<HTMLDivElement | null>(null);
const isDeployOutputDrawerOpen = ref(false);

// Fetch services with pagination
const fetchServices = async (page: number = currentPage.value) => {
  loading.value = true;
  try {
    const response = await serviceApi.page(page, pageSize.value);
    if (ApiResponseHelper.isSuccess(response)) {
      services.value = response.data?.services || [];
      total.value = response.data?.total || 0;
      currentPage.value = response.data?.page || page;
      pageSize.value = response.data?.pageSize || 10;
    }
  } catch (error) {
    console.error("Failed to fetch services:", error);
  } finally {
    loading.value = false;
  }
};

// Fetch applications for deploy dialog
const fetchApplications = async () => {
  if (applicationsLoading.value) return; // Prevent duplicate requests
  applicationsLoading.value = true;
  try {
    const response = await applicationApi.list();
    if (ApiResponseHelper.isSuccess(response)) {
      applications.value = response.data || [];
    }
  } catch (error) {
    console.error("Failed to fetch applications:", error);
  } finally {
    applicationsLoading.value = false;
  }
};

// Fetch nodes for deploy dialog
const fetchNodes = async () => {
  if (nodesLoading.value) return; // Prevent duplicate requests
  nodesLoading.value = true;
  try {
    const response = await nodeApi.list();
    if (ApiResponseHelper.isSuccess(response)) {
      nodes.value = response.data || [];
    }
  } catch (error) {
    console.error("Failed to fetch nodes:", error);
  } finally {
    nodesLoading.value = false;
  }
};

// Open deploy sheet
const openDeployDialog = async () => {
  currentStep.value = 1;
  previousClickStep.value = -1;
  selectedApplication.value = null;
  selectedNode.value = null;
  qaValues.value = {};
  // Clear previous data
  applications.value = [];
  nodes.value = [];
  isDeployDialogOpen.value = true;

  // Initialize first step data
  await goToStep(1);
};

// Close deploy sheet
const closeDeployDialog = () => {
  isDeployDialogOpen.value = false;
  currentStep.value = 1;
  previousClickStep.value = -1;
  selectedApplication.value = null;
  selectedNode.value = null;
  qaValues.value = {};
  deployOutput.value = [];
};

// Handle step navigation (for Next button)
const handleStepChange = async (step: number) => {
  // Prevent navigation if prerequisites are not met
  if (step === 2 && !selectedApplication.value) return;
  // Only update if step is different to prevent unnecessary updates
  if (previousClickStep.value === step) return;

  const previousStep = previousClickStep.value;
  previousClickStep.value = step;

  // Fetch data when navigating to a step (only if step actually changed)
  if (previousStep !== step) {
    if (
      step === 1 &&
      applications.value.length === 0 &&
      !applicationsLoading.value
    ) {
      await fetchApplications();
    } else if (step === 3 && nodes.value.length === 0 && !nodesLoading.value) {
      await fetchNodes();
    }
  }
};

const goToStep = async (step: number) => {
  currentStep.value = step;
  await handleStepChange(step);
};

// Initialize QA values when application is selected
const handleApplicationSelect = (appId: string) => {
  const app = applications.value.find((a) => a.id === Number(appId));
  if (app) {
    selectedApplication.value = app;
    // Initialize QA values with default values
    if (app.qa && app.qa.length > 0) {
      const initialValues: Record<string, string> = {};
      app.qa.forEach((item) => {
        initialValues[item.name] = item.default_value || "";
      });
      qaValues.value = initialValues;
    } else {
      qaValues.value = {};
    }
  }
};

// Handle node select
const handleNodeSelect = (nodeId: string) => {
  const node = nodes.value.find((n) => n.id === Number(nodeId));
  if (node) {
    selectedNode.value = node;
  }
};

// Validate current step
const canProceedToNextStep = computed(() => {
  if (currentStep.value === 1) {
    return selectedApplication.value !== null;
  } else if (currentStep.value === 2) {
    if (!selectedApplication.value?.qa) return true;
    // Check if all required QA items are filled
    return selectedApplication.value.qa.every((item) => {
      if (!item.required) return true;
      const value = qaValues.value[item.name];
      return value !== undefined && value !== null && value !== "";
    });
  } else if (currentStep.value === 3) {
    return selectedNode.value !== null;
  }
  return false;
});

// Handle QA value change
const handleQAValueChange = (name: string, value: string) => {
  qaValues.value[name] = value;
};

// Scroll to bottom of output
const scrollToBottom = () => {
  nextTick(() => {
    if (deployOutputRef.value) {
      deployOutputRef.value.scrollTop = deployOutputRef.value.scrollHeight;
    }
  });
};

// Watch deployOutput to auto-scroll and open drawer
watch(
  deployOutput,
  (newOutput) => {
    scrollToBottom();
    // 当有输出时自动打开 drawer
    if (newOutput.length > 0 && !isDeployOutputDrawerOpen.value) {
      isDeployOutputDrawerOpen.value = true;
    }
  },
  { deep: true }
);

// Deploy service
const deployService = async ({
  onlySave = false,
}: { onlySave?: boolean } = {}) => {
  if (!selectedApplication.value || !selectedNode.value) {
    return;
  }
  const saveReq: SaveServiceRequest = {
    id: 0,
    name: selectedApplication.value.name,
    app_id: selectedApplication.value.id,
    node_id: selectedNode.value.id,
    qa_values: qaValues.value,
  };
  const response = await serviceApi.save(saveReq);
  if (!ApiResponseHelper.isSuccess(response)) {
    showToast(response.message ?? "Failed to save service", "error");
    return;
  }
  if (onlySave) {
    return;
  }

  deployLoading.value = true;
  deployOutput.value = [];
  isDeployOutputDrawerOpen.value = true;
  const req: DeployServiceRequest = {
    app_id: selectedApplication.value.id,
    node_id: selectedNode.value.id,
    qa_values: qaValues.value,
  };
  await serviceApi.deployStream(
    req,
    (data) => {
      deployOutput.value.push(data);
    },
    (error) => {
      deployOutput.value.push(`[错误] ${error.message}`);
      console.error("Failed to deploy service:", error);
      deployLoading.value = false;
    },
    () => {
      // 部署完成
      deployLoading.value = false;
      console.log("Deployment completed");
      deployOutput.value.push("Deployment completed!!!");
    }
  );
};

// Get QA value for an item
const getQAValue = (itemName: string) => {
  return qaValues.value[itemName] || "";
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    running: "bg-green-500/10 text-green-500 border-green-500/20",
    stopped: "bg-red-500/10 text-red-500 border-red-500/20",
    pending: "bg-yellow-500/10 text-yellow-500 border-yellow-500/20",
  };
  return colors[status] || colors.stopped;
};

// Go to page
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
    fetchServices(page);
  }
};

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
      <Button @click="openDeployDialog">
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
          <TableRow>
            <TableHead class="px-6">ID</TableHead>
            <TableHead class="px-6">Name</TableHead>
            <TableHead class="px-6">Status</TableHead>
            <TableHead class="px-6">Node ID</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="service in services" :key="service.id">
            <TableCell class="px-6">{{ service.id }}</TableCell>
            <TableCell class="px-6">{{ service.name }}</TableCell>
            <TableCell class="px-6">
              <span
                :class="[
                  'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium',
                  getStatusColor(service.status),
                ]"
              >
                {{ service.status }}
              </span>
            </TableCell>
            <TableCell class="px-6">{{ service.node_id }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <!-- Pagination Controls -->
      <div
        v-if="totalPages > 1"
        class="flex items-center justify-between border-t px-6 py-4"
      >
        <div class="text-sm text-muted-foreground">
          Showing {{ (currentPage - 1) * pageSize + 1 }} -
          {{ Math.min(currentPage * pageSize, total) }} of {{ total }} services
        </div>
        <div class="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            :disabled="currentPage === 1"
            @click="goToPage(currentPage - 1)"
          >
            <Icon icon="lucide:chevron-left" class="h-4 w-4" />
          </Button>
          <div class="flex items-center gap-1">
            <Button
              v-for="page in totalPages"
              :key="page"
              variant="outline"
              size="sm"
              :class="[
                'min-w-[40px]',
                currentPage === page
                  ? 'bg-primary text-primary-foreground'
                  : '',
              ]"
              @click="goToPage(page)"
            >
              {{ page }}
            </Button>
          </div>
          <Button
            variant="outline"
            size="sm"
            :disabled="currentPage === totalPages"
            @click="goToPage(currentPage + 1)"
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
      <Button class="mt-4" @click="openDeployDialog">
        <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
        Deploy First Service
      </Button>
    </div>

    <!-- Deploy Output Drawer -->
    <Drawer v-model:open="isDeployOutputDrawerOpen">
      <DrawerContent>
        <DrawerHeader>
          <div class="flex items-center justify-between">
            <div>
              <DrawerTitle>Deploy Output</DrawerTitle>
            </div>
          </div>
        </DrawerHeader>
        <div class="px-4 pb-4">
          <div
            ref="deployOutputRef"
            class="max-h-[60vh] overflow-y-auto rounded-md bg-muted/30 p-3 font-mono text-xs"
          >
            <div
              v-if="deployOutput.length === 0"
              class="text-center text-muted-foreground py-8"
            >
              No output yet
            </div>
            <div
              v-for="(line, index) in deployOutput"
              :key="index"
              class="mb-1 whitespace-pre-wrap break-words"
            >
              {{ line }}
            </div>
          </div>
        </div>
      </DrawerContent>
    </Drawer>

    <!-- Deploy Sheet -->
    <Sheet v-model:open="isDeployDialogOpen">
      <SheetContent
        class="flex h-full w-full max-w-[90vw] sm:max-w-none lg:max-w-[1200px] flex-col overflow-y-auto"
      >
        <SheetHeader class="px-3 sm:px-5">
          <SheetTitle>Deploy Service</SheetTitle>
          <SheetDescription>
            Follow the steps to deploy a new service
          </SheetDescription>
        </SheetHeader>

        <div class="flex-1 overflow-y-auto">
          <div class="px-3 sm:px-5 py-6">
            <Stepper
              v-model="currentStep"
              class="flex w-full items-start gap-2"
            >
              <!-- Step 1: Select Application -->
              <StepperItem
                v-slot="{ state }"
                class="relative flex w-full flex-col items-center justify-center"
                :step="1"
              >
                <StepperSeparator
                  class="absolute left-[calc(50%+20px)] right-[calc(-50%+10px)] top-5 block h-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
                />
                <StepperTrigger as-child>
                  <Button
                    :variant="
                      state === 'completed' || state === 'active'
                        ? 'default'
                        : 'outline'
                    "
                    size="icon"
                    class="z-10 rounded-full shrink-0"
                    :class="[
                      state === 'active' &&
                        'ring-2 ring-ring ring-offset-2 ring-offset-background',
                    ]"
                    @click.stop="handleStepChange(1)"
                  >
                    <Icon
                      v-if="state === 'completed'"
                      icon="lucide:check"
                      class="h-5 w-5"
                    />
                    <Icon
                      v-else-if="state === 'active'"
                      icon="lucide:circle"
                      class="h-5 w-5"
                    />
                    <Icon v-else icon="lucide:dot" class="h-5 w-5" />
                  </Button>
                </StepperTrigger>
                <div class="mt-5 flex flex-col items-center text-center">
                  <StepperTitle
                    :class="[state === 'active' && 'text-primary']"
                    class="text-sm font-semibold transition lg:text-base"
                  >
                    Select Application
                  </StepperTitle>
                  <StepperDescription
                    :class="[state === 'active' && 'text-primary']"
                    class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm"
                  >
                    Choose an application to deploy
                  </StepperDescription>
                </div>
              </StepperItem>

              <!-- Step 2: Configure Settings -->
              <StepperItem
                v-slot="{ state }"
                class="relative flex w-full flex-col items-center justify-center"
                :step="2"
              >
                <StepperSeparator
                  class="absolute left-[calc(50%+20px)] right-[calc(-50%+10px)] top-5 block h-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
                />
                <StepperTrigger as-child>
                  <Button
                    :variant="
                      state === 'completed' || state === 'active'
                        ? 'default'
                        : 'outline'
                    "
                    size="icon"
                    class="z-10 rounded-full shrink-0"
                    :class="[
                      state === 'active' &&
                        'ring-2 ring-ring ring-offset-2 ring-offset-background',
                    ]"
                    @click.stop="handleStepChange(2)"
                  >
                    <Icon
                      v-if="state === 'completed'"
                      icon="lucide:check"
                      class="h-5 w-5"
                    />
                    <Icon
                      v-else-if="state === 'active'"
                      icon="lucide:circle"
                      class="h-5 w-5"
                    />
                    <Icon v-else icon="lucide:dot" class="h-5 w-5" />
                  </Button>
                </StepperTrigger>
                <div class="mt-5 flex flex-col items-center text-center">
                  <StepperTitle
                    :class="[state === 'active' && 'text-primary']"
                    class="text-sm font-semibold transition lg:text-base"
                  >
                    Configure Settings
                  </StepperTitle>
                  <StepperDescription
                    :class="[state === 'active' && 'text-primary']"
                    class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm"
                  >
                    Set application QA properties
                  </StepperDescription>
                </div>
              </StepperItem>

              <!-- Step 3: Select Node -->
              <StepperItem
                v-slot="{ state }"
                class="relative flex w-full flex-col items-center justify-center"
                :step="3"
              >
                <StepperTrigger as-child>
                  <Button
                    :variant="
                      state === 'completed' || state === 'active'
                        ? 'default'
                        : 'outline'
                    "
                    size="icon"
                    class="z-10 rounded-full shrink-0"
                    :class="[
                      state === 'active' &&
                        'ring-2 ring-ring ring-offset-2 ring-offset-background',
                    ]"
                    @click.stop="handleStepChange(3)"
                  >
                    <Icon
                      v-if="state === 'completed'"
                      icon="lucide:check"
                      class="h-5 w-5"
                    />
                    <Icon
                      v-else-if="state === 'active'"
                      icon="lucide:circle"
                      class="h-5 w-5"
                    />
                    <Icon v-else icon="lucide:dot" class="h-5 w-5" />
                  </Button>
                </StepperTrigger>
                <div class="mt-5 flex flex-col items-center text-center">
                  <StepperTitle
                    :class="[state === 'active' && 'text-primary']"
                    class="text-sm font-semibold transition lg:text-base"
                  >
                    Select Node
                  </StepperTitle>
                  <StepperDescription
                    :class="[state === 'active' && 'text-primary']"
                    class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm"
                  >
                    Choose a node to deploy on
                  </StepperDescription>
                </div>
              </StepperItem>
            </Stepper>

            <!-- Step Content -->
            <div class="mt-8 space-y-4">
              <!-- Step 1 Content: Application Selection -->
              <div v-if="currentStep === 1" class="space-y-4">
                <div class="space-y-2">
                  <Label for="deploy-app-select">Application</Label>
                  <div
                    v-if="applicationsLoading"
                    class="flex items-center justify-center py-8 text-muted-foreground"
                  >
                    <Icon
                      icon="lucide:loader-2"
                      class="h-5 w-5 mr-2 animate-spin"
                    />
                    Loading applications...
                  </div>
                  <Select
                    v-else
                    id="deploy-app-select"
                    :model-value="selectedApplication?.id?.toString()"
                    @update:model-value="
                      (val) => handleApplicationSelect(val?.toString() || '')
                    "
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select an application" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="app in applications"
                        :key="app.id"
                        :value="app.id.toString()"
                      >
                        {{ app.display || app.name }}
                        <span
                          v-if="app.version"
                          class="text-muted-foreground ml-2"
                        >
                          ({{ app.version }})
                        </span>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <p
                    v-if="selectedApplication?.description"
                    class="text-sm text-muted-foreground"
                  >
                    {{ selectedApplication.description }}
                  </p>
                </div>
              </div>

              <!-- Step 2 Content: QA Configuration -->
              <div v-if="currentStep === 2" class="space-y-4">
                <div
                  v-if="
                    !selectedApplication?.qa ||
                    selectedApplication.qa.length === 0
                  "
                  class="rounded-lg border border-dashed bg-muted/30 p-6 text-center text-sm text-muted-foreground"
                >
                  No QA configuration required for this application.
                </div>
                <div v-else class="space-y-4">
                  <div
                    v-for="item in selectedApplication?.qa || []"
                    :key="item.name"
                    class="space-y-2"
                  >
                    <Label :for="`qa-${item.name}`">
                      {{ item.name }}
                      <span v-if="item.required" class="text-destructive"
                        >*</span
                      >
                    </Label>
                    <p
                      v-if="item.description"
                      class="text-xs text-muted-foreground mb-2"
                    >
                      {{ item.description }}
                    </p>

                    <!-- Text Input -->
                    <Input
                      v-if="item.type === 'text'"
                      :id="`qa-${item.name}`"
                      :model-value="getQAValue(item.name)"
                      :placeholder="item.description || `Enter ${item.name}`"
                      @update:model-value="(val: string | number) => handleQAValueChange(item.name, String(val))"
                    />

                    <!-- Textarea -->
                    <textarea
                      v-else-if="item.type === 'textarea'"
                      :id="`qa-${item.name}`"
                      class="border-input placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 min-h-[80px] w-full rounded-md border bg-transparent px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"
                      :value="getQAValue(item.name)"
                      :placeholder="item.description || `Enter ${item.name}`"
                      @input="(e: Event) => handleQAValueChange(item.name, (e.target as HTMLTextAreaElement).value)"
                    />

                    <!-- Number Input -->
                    <Input
                      v-else-if="item.type === 'number'"
                      :id="`qa-${item.name}`"
                      type="number"
                      :model-value="getQAValue(item.name)"
                      :placeholder="item.description || `Enter ${item.name}`"
                      @update:model-value="(val: string | number) => handleQAValueChange(item.name, String(val))"
                    />

                    <!-- Boolean Checkbox -->
                    <div
                      v-else-if="item.type === 'boolean'"
                      class="flex items-center space-x-3 rounded-md border px-3 py-2"
                    >
                      <Checkbox
                        :id="`qa-${item.name}`"
                        :model-value="getQAValue(item.name) === 'true'"
                        @update:model-value="(checked: boolean | 'indeterminate') => handleQAValueChange(item.name, checked === true ? 'true' : 'false')"
                      />
                      <Label :for="`qa-${item.name}`" class="text-sm">
                        {{ item.description || item.name }}
                      </Label>
                    </div>

                    <!-- Select -->
                    <Select
                      v-else-if="item.type === 'select'"
                      :model-value="getQAValue(item.name)"
                      @update:model-value="
                        (val) =>
                          handleQAValueChange(item.name, val ? String(val) : '')
                      "
                    >
                      <SelectTrigger>
                        <SelectValue
                          :placeholder="
                            item.description || `Select ${item.name}`
                          "
                        />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem
                          v-for="option in item.options || []"
                          :key="option"
                          :value="option"
                        >
                          {{ option }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </div>

              <!-- Step 3 Content: Node Selection -->
              <div v-if="currentStep === 3" class="space-y-4">
                <div class="space-y-2">
                  <Label for="deploy-node-select">Node</Label>
                  <div
                    v-if="nodesLoading"
                    class="flex items-center justify-center py-8 text-muted-foreground"
                  >
                    <Icon
                      icon="lucide:loader-2"
                      class="h-5 w-5 mr-2 animate-spin"
                    />
                    Loading nodes...
                  </div>
                  <Select
                    v-else
                    id="deploy-node-select"
                    :model-value="selectedNode?.id?.toString()"
                    @update:model-value="
                      (val) => handleNodeSelect(val?.toString() || '')
                    "
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select a node" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="node in nodes"
                        :key="node.id"
                        :value="node.id.toString()"
                      >
                        {{ node.display_name || node.name }}
                        <span class="text-muted-foreground ml-2">
                          ({{ node.address }})
                        </span>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <p v-if="selectedNode" class="text-sm text-muted-foreground">
                    {{ selectedNode.display_name || selectedNode.name }} -
                    {{ selectedNode.address }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <SheetFooter
          class="flex flex-row items-center justify-end gap-2 px-3 sm:px-5"
        >
          <Button
            variant="outline"
            @click="closeDeployDialog"
            :disabled="deployLoading"
          >
            Cancel
          </Button>
          <Button
            v-if="currentStep > 1"
            @click="goToStep(currentStep - 1)"
            :disabled="currentStep == 1"
          >
            <Icon icon="lucide:chevron-left" class="h-4 w-4 mr-2" />
            Previous
          </Button>
          <Button
            v-if="currentStep < 3"
            @click="goToStep(currentStep + 1)"
            :disabled="!canProceedToNextStep"
          >
            Next
            <Icon icon="lucide:chevron-right" class="h-4 w-4 ml-2" />
          </Button>
          <Button
            v-if="currentStep === 3"
            @click="deployService({ onlySave: true })"
            :disabled="!canProceedToNextStep || deployLoading"
          >
            <Icon icon="lucide:save" class="h-4 w-4 mr-2" />
            Save
          </Button>
          <Button
            v-if="currentStep === 3"
            @click="deployService"
            :disabled="!canProceedToNextStep || deployLoading"
          >
            <Icon
              v-if="deployLoading"
              icon="lucide:loader-2"
              class="h-4 w-4 mr-2 animate-spin"
            />
            <Icon v-else icon="lucide:cloud-upload" class="h-4 w-4 mr-2" />
            {{ deployLoading ? "Deploying..." : "Save & Deploy" }}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  </div>
</template>
