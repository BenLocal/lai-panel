<script setup lang="ts">
import { ref, watch } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Select,
    SelectTrigger,
    SelectValue,
    SelectContent,
    SelectItem,
} from "@/components/ui/select";
import { Checkbox } from "@/components/ui/checkbox";
import type { ApplicationQAItem } from "@/api/application";

type QAType = ApplicationQAItem["type"];
type QAEditorEmits = (
    event: "update:modelValue",
    value: ApplicationQAItem[]
) => void;

const props = defineProps<{
    modelValue: ApplicationQAItem[];
}>();

const emit = defineEmits<QAEditorEmits>();

const qaItems = ref<ApplicationQAItem[]>([]);

const normalizeQAItem = (item: ApplicationQAItem): ApplicationQAItem => {
    const normalized: ApplicationQAItem = {
        name: item.name ?? "",
        value: item.value ?? "",
        type: item.type ?? "text",
        default_value: item.default_value ?? "",
        options: item.options,
        required: item.required ?? false,
        description: item.description ?? "",
    };

    if (normalized.type === "select") {
        const options =
            normalized.options && normalized.options.length > 0
                ? normalized.options
                : ["Option 1", "Option 2"];
        normalized.options = options;
        if (!options.includes(normalized.value)) {
            normalized.value = options[0] ?? "";
        }
        if (normalized.default_value && !options.includes(normalized.default_value)) {
            normalized.default_value = options[0] ?? "";
        }
    } else if (normalized.type === "boolean") {
        normalized.value = normalized.value === "true" ? "true" : "false";
        normalized.default_value =
            normalized.default_value === "true" ? "true" : "false";
        normalized.options = undefined;
    } else {
        normalized.options = undefined;
    }

    return normalized;
};

const syncFromProps = (items: ApplicationQAItem[] | undefined) => {
    qaItems.value = (items ?? []).map((item) =>
        normalizeQAItem({
            ...item,
        })
    );
};

watch(
    () => props.modelValue,
    (value) => {
        syncFromProps(value);
    },
    { immediate: true, deep: true }
);

const emitChange = () => {
    emit(
        "update:modelValue",
        qaItems.value.map((item) => normalizeQAItem({ ...item }))
    );
};

const createEmptyItem = (): ApplicationQAItem => ({
    name: "",
    value: "",
    type: "text",
    default_value: "",
    required: false,
    description: "",
});

const addQAItem = () => {
    qaItems.value = [...qaItems.value, createEmptyItem()];
    emitChange();
};

const removeQAItem = (index: number) => {
    const next = [...qaItems.value];
    next.splice(index, 1);
    qaItems.value = next;
    emitChange();
};

const updateItemField = <K extends keyof ApplicationQAItem>(
    index: number,
    key: K,
    value: ApplicationQAItem[K]
) => {
    const current = qaItems.value[index];
    if (!current) {
        return;
    }

    const nextItems = [...qaItems.value];
    const nextItem = { ...current, [key]: value } as ApplicationQAItem;

    if (key === "type") {
        const newType = value as QAType;
        nextItem.type = newType;
        if (newType === "select") {
            const options =
                current.options && current.options.length > 0
                    ? current.options
                    : ["Option 1", "Option 2"];
            nextItem.options = options;
            nextItem.value = options[0] ?? "";
            nextItem.default_value = options[0] ?? "";
        } else if (newType === "boolean") {
            nextItem.options = undefined;
            nextItem.value = current.value === "true" ? "true" : "false";
            nextItem.default_value =
                current.default_value === "true" ? "true" : nextItem.value;
        } else {
            nextItem.options = undefined;
        }
    }

    if (key === "options") {
        const options = (value as string[]).filter(
            (itemOption) => itemOption.trim().length
        );
        nextItem.options = options;
        if (!options.includes(nextItem.value)) {
            nextItem.value = options[0] ?? "";
        }
        if (
            nextItem.default_value &&
            !options.includes(nextItem.default_value ?? "")
        ) {
            nextItem.default_value = options[0] ?? "";
        }
    }

    nextItems[index] = normalizeQAItem(nextItem);
    qaItems.value = nextItems;
    emitChange();
};

const handleOptionsInput = (index: number, raw: string) => {
    const normalized = raw
        .split("\n")
        .map((option) => option.trim())
        .filter((option) => option.length > 0);
    updateItemField(index, "options", normalized);
};

const typeOptions: Array<{ label: string; value: QAType }> = [
    { label: "Text", value: "text" },
    { label: "Number", value: "number" },
    { label: "Boolean", value: "boolean" },
    { label: "Select", value: "select" },
    { label: "Textarea", value: "textarea" },
];

type QAStringField = "name" | "default_value" | "value" | "description";

const handleStringFieldInput = (
    index: number,
    key: QAStringField,
    value: string | number
) => {
    updateItemField(
        index,
        key,
        String(value) as ApplicationQAItem[QAStringField]
    );
};

const handleTypeChange = (index: number, value: unknown) => {
    if (typeof value !== "string") {
        return;
    }
    updateItemField(index, "type", value as QAType);
};

const handleSelectFieldChange = (
    index: number,
    key: "default_value" | "value",
    value: unknown
) => {
    if (typeof value !== "string") {
        return;
    }
    updateItemField(index, key, value);
};
</script>

<template>
    <div class="space-y-4">
        <div v-if="qaItems.length === 0"
            class="rounded-lg border border-dashed bg-muted/30 p-6 text-center text-sm text-muted-foreground">
            No QA configuration yet. Click the button below to add one.
        </div>

        <div v-for="(qa, index) in qaItems" :key="index" class="space-y-4 rounded-lg border p-4">
            <div class="flex items-center justify-between">
                <div class="text-sm font-medium">QA Item {{ index + 1 }}</div>
                <Button variant="ghost" size="icon" class="h-8 w-8 text-muted-foreground hover:text-destructive"
                    @click="removeQAItem(index)">
                    <Icon icon="lucide:trash" class="h-4 w-4" />
                </Button>
            </div>

            <div class="grid gap-4 md:grid-cols-2">
                <div class="space-y-2">
                    <label :for="`qa-name-${index}`" class="text-sm font-medium">
                        Field Name
                    </label>
                    <Input :id="`qa-name-${index}`" :model-value="qa.name" placeholder="Unique key" @update:model-value="(value) =>
                        handleStringFieldInput(index, 'name', value ?? '')
                    " />
                </div>

                <div class="space-y-2">
                    <div class="text-sm font-medium">Type</div>
                    <Select :model-value="qa.type" @update:model-value="(value) => handleTypeChange(index, value)">
                        <SelectTrigger>
                            <SelectValue placeholder="Select a type" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem v-for="option in typeOptions" :key="option.value" :value="option.value">
                                {{ option.label }}
                            </SelectItem>
                        </SelectContent>
                    </Select>
                </div>
            </div>

            <div class="grid gap-4 md:grid-cols-2">
                <div class="space-y-2">
                    <div class="text-sm font-medium">Default Value</div>
                    <Input v-if="qa.type === 'text' || qa.type === 'number'"
                        :type="qa.type === 'number' ? 'number' : 'text'" :model-value="qa.default_value ?? ''"
                        placeholder="Default value" @update:model-value="(value) =>
                            handleStringFieldInput(index, 'default_value', value ?? '')
                        " />
                    <Select v-else-if="qa.type === 'boolean'" :model-value="qa.default_value ?? 'false'"
                        @update:model-value="
                            (value) => handleSelectFieldChange(index, 'default_value', value)
                        ">
                        <SelectTrigger>
                            <SelectValue placeholder="Select a default value" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="true">True</SelectItem>
                            <SelectItem value="false">False</SelectItem>
                        </SelectContent>
                    </Select>
                    <Select v-else-if="qa.type === 'select'" :model-value="qa.default_value ?? ''" @update:model-value="
                        (value) => handleSelectFieldChange(index, 'default_value', value)
                    ">
                        <SelectTrigger>
                            <SelectValue placeholder="Select a default value" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem v-for="option in qa.options ?? []" :key="option" :value="option">
                                {{ option }}
                            </SelectItem>
                        </SelectContent>
                    </Select>
                    <textarea v-else
                        class="border-input placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 min-h-[38px] w-full rounded-md border bg-transparent px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"
                        :value="qa.default_value ?? ''" placeholder="Default value" @input="(event: Event) =>
                            handleStringFieldInput(
                                index,
                                'default_value',
                                (event.target as HTMLTextAreaElement).value
                            )"></textarea>
                </div>

                <div class="space-y-2">
                    <div class="text-sm font-medium">Required</div>
                    <div class="flex h-10 items-center space-x-3 rounded-md border px-3">
                        <Checkbox :id="`qa-required-${index}`" :model-value="qa.required ?? false" @update:model-value="(checked) =>
                            updateItemField(index, 'required', checked as ApplicationQAItem['required'])" />
                        <label :for="`qa-required-${index}`" class="text-sm text-muted-foreground">
                            Must be provided
                        </label>
                    </div>
                </div>
            </div>

            <div class="space-y-2">
                <div class="text-sm font-medium">Description</div>
                <textarea
                    class="border-input placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 min-h-[80px] w-full rounded-md border bg-transparent px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"
                    :value="qa.description ?? ''" placeholder="Optional description" @input="(event: Event) =>
                        handleStringFieldInput(
                            index,
                            'description',
                            (event.target as HTMLTextAreaElement).value
                        )"></textarea>
            </div>

            <div class="space-y-2">
                <div class="text-sm font-medium">Current Value</div>
                <Input v-if="qa.type === 'text' || qa.type === 'number'"
                    :type="qa.type === 'number' ? 'number' : 'text'" :model-value="qa.value ?? ''"
                    placeholder="User-provided value" @update:model-value="(value) =>
                        handleStringFieldInput(index, 'value', value ?? '')
                    " />
                <Select v-else-if="qa.type === 'boolean'" :model-value="qa.value" @update:model-value="
                    (value) => handleSelectFieldChange(index, 'value', value)
                ">
                    <SelectTrigger>
                        <SelectValue placeholder="Select a value" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="true">True</SelectItem>
                        <SelectItem value="false">False</SelectItem>
                    </SelectContent>
                </Select>
                <Select v-else-if="qa.type === 'select'" :model-value="qa.value" @update:model-value="
                    (value) => handleSelectFieldChange(index, 'value', value)
                ">
                    <SelectTrigger>
                        <SelectValue placeholder="Select a value" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem v-for="option in qa.options ?? []" :key="option" :value="option">
                            {{ option }}
                        </SelectItem>
                    </SelectContent>
                </Select>
                <textarea v-else
                    class="border-input placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 min-h-[80px] w-full rounded-md border bg-transparent px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"
                    :value="qa.value" placeholder="User-provided value" @input="(event) =>
                        updateItemField(
                            index,
                            'value',
                            (event.target as HTMLTextAreaElement).value
                        )"></textarea>
            </div>

            <div v-if="qa.type === 'select'" class="space-y-2">
                <div class="text-sm font-medium">Options (one per line)</div>
                <textarea
                    class="border-input placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 min-h-[80px] w-full rounded-md border bg-transparent px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"
                    :value="(qa.options ?? []).join('\n')" placeholder="Example: Option 1" @input="(event: Event) =>
                        handleOptionsInput(
                            index,
                            (event.target as HTMLTextAreaElement).value
                        )"></textarea>
            </div>
        </div>

        <div class="pt-2">
            <Button variant="outline" @click="addQAItem">
                <Icon icon="lucide:plus" class="mr-2 h-4 w-4" />
                Add QA Item
            </Button>
        </div>
    </div>
</template>
