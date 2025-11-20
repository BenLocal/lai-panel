<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
} from "@/components/ui/sheet";
import {
    applicationApi,
    type Application,
    type ApplicationQAItem,
} from "@/api/application";
import ApplicationQAEditor from "@/components/application/ApplicationQAEditor.vue";
import { Separator } from "@/components/ui/separator";
import YamlEditor from "@/components/application/YamlEditor.vue";
import { load as parseYaml } from "js-yaml";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { ApiResponseHelper } from "@/api/base";

interface ApplicationForm {
    name: string;
    version?: string;
    display?: string;
    description?: string;
    icon?: string;
    qa: ApplicationQAItem[];
    dockerCompose: string;
}

const applications = ref<Application[]>([]);
const currentPage = ref(1);
const pageSize = ref(6);
const totalPages = ref(1);
const isSheetOpen = ref(false);
const isEditMode = ref(false);
const loading = ref(false);
const editingApplicationId = ref<number | null>(null);
const isDockerComposeValid = ref(true);
const isComposeEditorOpen = ref(false);

const createDefaultForm = (): ApplicationForm => ({
    display: "",
    name: "",
    description: "",
    version: "",
    icon: "lucide:app-window",
    qa: [],
    dockerCompose: "",
});

const formData = reactive<ApplicationForm>(createDefaultForm());

const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page;
    }
};

const fetchApplications = async () => {
    const response = await applicationApi.page(currentPage.value, pageSize.value);
    if (!ApiResponseHelper.isSuccess(response)) {
        return;
    }
    const data = response.data;
    applications.value = data?.apps || [];
    totalPages.value = Math.ceil((data?.total || 0) / pageSize.value);
    currentPage.value = data?.currentPage || 1;
    pageSize.value = data?.pageSize || 6;
};

const handleDockerComposeValidityChange = (isValid: boolean) => {
    isDockerComposeValid.value = isValid;
};

watch(
    () => formData.dockerCompose,
    (value) => {
        if (!value || !value.trim()) {
            isDockerComposeValid.value = true;
            return;
        }
        try {
            parseYaml(value);
            isDockerComposeValid.value = true;
        } catch {
            isDockerComposeValid.value = false;
        }
    },
    { immediate: true }
);

const dockerComposePreview = computed(() => {
    const content = formData.dockerCompose?.trim();
    if (!content) {
        return "";
    }

    const lines = content.split("\n");
    const snippet = lines.slice(0, 6).join("\n");
    return lines.length > 6 ? `${snippet}\n...` : snippet;
});

const namePattern = /^[A-Za-z]*$/;

const isNameValid = computed(() => namePattern.test(formData.name));

const handleNameInput = (event: Event) => {
    const target = event.target as HTMLInputElement;
    const sanitized = (target.value.match(/[A-Za-z]/g) ?? []).join("");
    if (sanitized !== target.value) {
        target.value = sanitized;
    }
    formData.name = sanitized;
};

const isSaveDisabled = computed(() => {
    return (
        !formData.name.trim() ||
        !isNameValid.value ||
        !isDockerComposeValid.value ||
        loading.value
    );
});

const resetForm = () => {
    Object.assign(formData, createDefaultForm());
    editingApplicationId.value = null;
    isDockerComposeValid.value = true;
    isComposeEditorOpen.value = false;
};

const openAddApplicationDialog = () => {
    isEditMode.value = false;
    resetForm();
    isSheetOpen.value = true;
};

const openEditApplicationDialog = (application: Application) => {
    isEditMode.value = true;
    editingApplicationId.value = application.id;
    Object.assign(formData, {
        name: application.name ?? "",
        description: application.description ?? "",
        version: application.version ?? "",
        icon: application.icon ?? "lucide:app-window",
        display: application.display ?? "",
        qa: application.qa
            ? application.qa.map((item) => ({
                ...item,
                options: item.options ? [...item.options] : undefined,
            }))
            : [],
        metadata: application.metadata
            ? application.metadata.map((item) => ({
                name: item.name ?? "",
                properties: Object.entries(item.properties ?? {}).map(
                    ([key, value]) => ({
                        key,
                        value,
                    })
                ),
            }))
            : [],
        dockerCompose: application.docker_compose ?? "",
    });
    isComposeEditorOpen.value = false;
    isSheetOpen.value = true;
};

const handleCancel = () => {
    isSheetOpen.value = false;
    isEditMode.value = false;
    resetForm();
};

const saveApplication = async () => {
    if (isSaveDisabled.value) {
        return;
    }

    loading.value = true;

    const payload: Application = {
        id: editingApplicationId.value ?? 0,
        name: formData.name.trim(),
        display: formData.display?.trim() ?? "",
        description: formData.description?.trim() ?? "",
        version: formData.version?.trim() ?? "",
        icon: formData.icon?.trim() ?? "",
        qa: formData.qa,
        docker_compose: formData.dockerCompose?.trim() ?? "",
    };

    try {
        if (isEditMode.value && editingApplicationId.value !== null) {
            await applicationApi.update(payload);
        } else {
            await applicationApi.add(payload);
        }
        await fetchApplications();
        isSheetOpen.value = false;
        isEditMode.value = false;
        resetForm();
    } catch (error) {
        console.error("Failed to save application", error);
    } finally {
        loading.value = false;
    }
};

onMounted(() => {
    fetchApplications();
});
</script>

<template>
    <div class="space-y-6">
        <div class="flex items-center justify-between">
            <div>
                <h1 class="text-3xl font-bold">Applications</h1>
                <p class="text-muted-foreground mt-1">
                    Manage and monitor your applications
                </p>
            </div>
            <Button @click="openAddApplicationDialog">
                <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
                New Application
            </Button>
        </div>

        <!-- Applications Cards -->
        <div v-if="applications.length > 0">
            <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                <div v-for="app in applications" :key="app.id"
                    class="group rounded-lg border bg-card p-6 hover:shadow-md transition-shadow cursor-pointer"
                    @click="openEditApplicationDialog(app)">
                    <div class="flex items-start justify-between mb-4">
                        <div class="flex items-center gap-3">
                            <div
                                class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                                <Icon :icon="app.icon || 'lucide:app-window'" class="h-5 w-5" />
                            </div>
                            <div>
                                <h3 class="font-semibold text-lg">{{ app.name }}</h3>
                                <p class="text-xs text-muted-foreground">{{ app.version }}</p>
                            </div>
                        </div>
                    </div>

                    <p class="text-sm text-muted-foreground mb-4 line-clamp-2">
                        {{ app.description }}
                    </p>

                    <div class="flex items-center justify-end">
                        <Button variant="ghost" size="sm" @click.stop>
                            <Icon icon="lucide:more-horizontal" class="h-4 w-4" />
                        </Button>
                    </div>
                </div>
            </div>

            <!-- Pagination Controls -->
            <div v-if="totalPages > 1" class="flex items-center justify-between border-t pt-6 mt-6">
                <div class="text-sm text-muted-foreground">
                    Showing {{ (currentPage - 1) * pageSize + 1 }} -
                    {{ Math.min(currentPage * pageSize, applications.length) }} of
                    {{ applications.length }} applications
                </div>
                <div class="flex items-center gap-2">
                    <Button variant="outline" size="sm" :disabled="currentPage === 1"
                        @click="goToPage(currentPage - 1)">
                        <Icon icon="lucide:chevron-left" class="h-4 w-4" />
                    </Button>
                    <div class="flex items-center gap-1">
                        <Button v-for="page in totalPages" :key="page" variant="outline" size="sm" :class="[
                            'min-w-[40px]',
                            currentPage === page
                                ? 'bg-primary text-primary-foreground'
                                : '',
                        ]" @click="goToPage(page)">
                            {{ page }}
                        </Button>
                    </div>
                    <Button variant="outline" size="sm" :disabled="currentPage === totalPages"
                        @click="goToPage(currentPage + 1)">
                        <Icon icon="lucide:chevron-right" class="h-4 w-4" />
                    </Button>
                </div>
            </div>
        </div>

        <!-- Empty State -->
        <div v-else class="rounded-lg border bg-card p-12 text-center">
            <Icon icon="lucide:layers" class="h-12 w-12 mx-auto text-muted-foreground mb-4" />
            <p class="text-muted-foreground">No applications found</p>
            <Button class="mt-4" @click="openAddApplicationDialog">
                <Icon icon="lucide:plus" class="h-4 w-4 mr-2" />
                Add First Application
            </Button>
        </div>

        <Sheet v-model:open="isSheetOpen">
            <SheetContent class="flex h-full w-full max-w-[90vw] sm:max-w-none lg:max-w-[1200px] flex-col">
                <SheetHeader class="px-3 sm:px-5">
                    <SheetTitle>{{
                        isEditMode ? "Edit Application" : "Add Application"
                        }}</SheetTitle>
                    <SheetDescription>
                        {{
                            isEditMode
                                ? "Update application information"
                                : "Fill in the application details to create a new application"
                        }}
                    </SheetDescription>
                </SheetHeader>

                <div class="overflow-y-auto">
                    <div class="space-y-4 px-3 sm:px-5">
                        <div class="space-y-2">
                            <label for="app-name" class="text-sm font-medium">Name *</label>
                            <Input id="app-name" v-model="formData.name"
                                placeholder="Application name, English letters only" @input="handleNameInput" />
                            <p v-if="formData.name && !isNameValid" class="text-xs text-destructive">
                                Only English letters (A-Z) are allowed.
                            </p>
                        </div>
                        <div class="space-y-2">
                            <label for="app-display" class="text-sm font-medium">Display Name</label>
                            <Input id="app-display" v-model="formData.display"
                                placeholder="Application display name, if not set, use name" />
                        </div>

                        <div class="space-y-2">
                            <label for="app-version" class="text-sm font-medium">Version</label>
                            <Input id="app-version" v-model="formData.version" placeholder="v1.0.0" />
                        </div>

                        <div class="space-y-2">
                            <label for="app-icon" class="text-sm font-medium">Icon</label>
                            <Input id="app-icon" v-model="formData.icon" placeholder="lucide:app-window" />
                        </div>

                        <div class="space-y-2">
                            <label for="app-description" class="text-sm font-medium">Description</label>
                            <textarea id="app-description" v-model="formData.description"
                                class="border-input placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 min-h-[120px] w-full rounded-md border bg-transparent px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"
                                placeholder="Describe the application"></textarea>
                        </div>
                        <div class="space-y-2">
                            <Separator class="my-4" />
                        </div>
                        <div class="space-y-3">
                            <div class="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
                                <span class="text-sm font-medium">Docker Compose (YAML)</span>
                                <span class="text-xs text-muted-foreground">
                                    Edit and preview the compose definition
                                </span>
                            </div>
                            <div
                                class="space-y-3 rounded-md border border-dashed border-muted-foreground/40 bg-muted/20 p-4">
                                <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                                    <div class="flex-1">
                                        <pre v-if="dockerComposePreview"
                                            class="max-h-48 overflow-auto whitespace-pre-wrap break-words rounded-md bg-background/80 px-3 py-2 text-xs text-muted-foreground shadow-inner">{{ dockerComposePreview }}</pre>
                                        <p v-else class="text-xs italic text-muted-foreground">
                                            No Docker Compose definition yet.
                                        </p>
                                    </div>
                                    <div class="flex shrink-0 gap-2 sm:justify-end">
                                        <Button size="sm" variant="outline" @click="isComposeEditorOpen = true">
                                            <Icon icon="lucide:pen-square" class="mr-2 h-4 w-4" />
                                            Open Editor
                                        </Button>
                                        <Button v-if="formData.dockerCompose" size="sm" variant="ghost"
                                            :disabled="loading" @click="formData.dockerCompose = ''">
                                            <Icon icon="lucide:trash-2" class="mr-2 h-4 w-4" />
                                            Clear
                                        </Button>
                                    </div>
                                </div>
                                <p v-if="!isDockerComposeValid" class="text-xs text-destructive">
                                    YAML contains syntax errors. Please fix them in the editor
                                    before saving.
                                </p>
                            </div>
                        </div>
                        <div class="space-y-2">
                            <Separator class="my-4" />
                        </div>
                        <div class="space-y-2">
                            <div class="text-sm font-medium">QA Configuration</div>
                            <ApplicationQAEditor v-model="formData.qa" />
                        </div>
                    </div>
                </div>

                <SheetFooter class="flex flex-row items-center justify-end gap-2 px-3 sm:px-5">
                    <Button variant="outline" @click="handleCancel" :disabled="loading">
                        Cancel
                    </Button>
                    <Button @click="saveApplication" :disabled="isSaveDisabled">
                        {{
                            loading
                                ? "Saving..."
                                : isEditMode
                                    ? "Update Application"
                                    : "Add Application"
                        }}
                    </Button>
                </SheetFooter>
            </SheetContent>
        </Sheet>
        <Dialog v-model:open="isComposeEditorOpen">
            <DialogContent
                class="flex h-screen w-screen max-w-none flex-col bg-background p-0 sm:rounded-none sm:max-w-none lg:w-[95vw]">
                <DialogHeader class="border-b px-6 py-4">
                    <DialogTitle>Docker Compose Editor</DialogTitle>
                    <DialogDescription>
                        Edit and preview the Docker Compose YAML in a dedicated workspace.
                    </DialogDescription>
                </DialogHeader>
                <div class="flex-1 overflow-hidden">
                    <YamlEditor v-model="formData.dockerCompose" :height="'100%'"
                        @valid-state-change="handleDockerComposeValidityChange" />
                </div>
                <DialogFooter class="border-t px-6 py-4">
                    <div class="flex w-full flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                        <span class="text-xs" :class="isDockerComposeValid
                            ? 'text-muted-foreground'
                            : 'text-destructive'
                            ">
                            {{
                                isDockerComposeValid
                                    ? "YAML parsed successfully."
                                    : "YAML contains syntax errors."
                            }}
                        </span>
                        <div class="flex gap-2">
                            <Button variant="outline" @click="isComposeEditorOpen = false">
                                Close
                            </Button>
                        </div>
                    </div>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    </div>
</template>
