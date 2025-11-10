<script setup lang="ts">
import { ref, watch } from "vue";
import { Icon } from "@iconify/vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { Metadata } from "@/api/base";

interface MetadataProperty {
    key: string;
    value: string;
}

interface EditableMetadata {
    name: string;
    properties: MetadataProperty[];
}

type MetadataEditorEmits = (
    event: "update:modelValue",
    value: Metadata[]
) => void;

const props = defineProps<{
    modelValue: Metadata[];
}>();

const emit = defineEmits<MetadataEditorEmits>();

const editableMetadata = ref<EditableMetadata[]>([]);

const toEditable = (items: Metadata[] | undefined): EditableMetadata[] => {
    return (items ?? []).map((item) => {
        const propertyEntries = Object.entries(item.properties ?? {});
        const properties =
            propertyEntries.length > 0
                ? propertyEntries.map(([key, value]) => ({
                    key,
                    value,
                }))
                : [
                    {
                        key: "",
                        value: "",
                    },
                ];

        return {
            name: item.name ?? "",
            properties,
        };
    });
};

const toMetadata = (items: EditableMetadata[]): Metadata[] => {
    return items.map((item) => {
        const name = item.name.trim();
        const properties = item.properties.reduce<Record<string, string>>(
            (acc, property) => {
                const key = property.key.trim();
                if (key.length > 0) {
                    acc[key] = property.value;
                }
                return acc;
            },
            {}
        );
        return {
            name,
            properties,
        };
    });
};

watch(
    () => props.modelValue,
    (value) => {
        editableMetadata.value = toEditable(value);
    },
    { immediate: true, deep: true }
);

const emitChange = () => {
    emit("update:modelValue", toMetadata(editableMetadata.value));
};

const addMetadataItem = () => {
    editableMetadata.value = [
        ...editableMetadata.value,
        {
            name: "",
            properties: [
                {
                    key: "",
                    value: "",
                },
            ],
        },
    ];
    emitChange();
};

const removeMetadataItem = (index: number) => {
    const next = [...editableMetadata.value];
    next.splice(index, 1);
    editableMetadata.value = next;
    emitChange();
};

const updateMetadataName = (index: number, name: string | number) => {
    const current = editableMetadata.value[index];
    if (!current) {
        return;
    }
    const next = [...editableMetadata.value];
    next[index] = {
        ...current,
        name: String(name),
    };
    editableMetadata.value = next;
    emitChange();
};

const addProperty = (index: number) => {
    const current = editableMetadata.value[index];
    if (!current) {
        return;
    }
    const next = [...editableMetadata.value];
    next[index] = {
        ...current,
        properties: [
            ...current.properties,
            {
                key: "",
                value: "",
            },
        ],
    };
    editableMetadata.value = next;
    emitChange();
};

const removeProperty = (metadataIndex: number, propertyIndex: number) => {
    const current = editableMetadata.value[metadataIndex];
    if (!current) {
        return;
    }
    const next = [...editableMetadata.value];
    const properties = [...current.properties];
    properties.splice(propertyIndex, 1);
    next[metadataIndex] = {
        ...current,
        properties: properties.length
            ? properties
            : [
                {
                    key: "",
                    value: "",
                },
            ],
    };
    editableMetadata.value = next;
    emitChange();
};

const updatePropertyKey = (
    metadataIndex: number,
    propertyIndex: number,
    key: string | number
) => {
    const current = editableMetadata.value[metadataIndex];
    if (!current) {
        return;
    }
    const properties = [...current.properties];
    const property = properties[propertyIndex];
    if (!property) {
        return;
    }
    properties[propertyIndex] = {
        ...property,
        key: String(key),
    };
    const next = [...editableMetadata.value];
    next[metadataIndex] = {
        ...current,
        properties,
    };
    editableMetadata.value = next;
    emitChange();
};

const updatePropertyValue = (
    metadataIndex: number,
    propertyIndex: number,
    value: string | number
) => {
    const current = editableMetadata.value[metadataIndex];
    if (!current) {
        return;
    }
    const properties = [...current.properties];
    const property = properties[propertyIndex];
    if (!property) {
        return;
    }
    properties[propertyIndex] = {
        ...property,
        value: String(value),
    };
    const next = [...editableMetadata.value];
    next[metadataIndex] = {
        ...current,
        properties,
    };
    editableMetadata.value = next;
    emitChange();
};
</script>

<template>
    <div class="space-y-4">
        <div v-if="editableMetadata.length === 0"
            class="rounded-lg border border-dashed bg-muted/30 p-6 text-center text-sm text-muted-foreground">
            No metadata defined yet. Click the button below to add one.
        </div>

        <div v-for="(metadata, metadataIndex) in editableMetadata" :key="metadataIndex"
            class="space-y-4 rounded-lg border p-4">
            <div class="flex items-center justify-between">
                <div class="text-sm font-medium">Metadata {{ metadataIndex + 1 }}</div>
                <Button variant="ghost" size="icon" class="h-8 w-8 text-muted-foreground hover:text-destructive"
                    @click="removeMetadataItem(metadataIndex)">
                    <Icon icon="lucide:trash" class="h-4 w-4" />
                </Button>
            </div>

            <div class="space-y-2">
                <label :for="`metadata-name-${metadataIndex}`" class="text-sm font-medium">Name</label>
                <Input :id="`metadata-name-${metadataIndex}`" :model-value="metadata.name"
                    placeholder="Example: database" @update:model-value="(value) =>
                        updateMetadataName(metadataIndex, value ?? '')
                    " />
            </div>

            <div class="space-y-3">
                <div class="text-xs font-medium uppercase text-muted-foreground">
                    Properties
                </div>

                <div v-for="(property, propertyIndex) in metadata.properties" :key="`${metadataIndex}-${propertyIndex}`"
                    class="rounded-md border p-3">
                    <div class="grid gap-3 md:grid-cols-[1fr,1fr,auto] md:items-center">
                        <div class="space-y-2">
                            <div class="text-xs font-medium text-muted-foreground">Key</div>
                            <Input :model-value="property.key" placeholder="Example: host" @update:model-value="(value) =>
                                updatePropertyKey(metadataIndex, propertyIndex, value ?? '')
                            " />
                        </div>

                        <div class="space-y-2">
                            <div class="text-xs font-medium text-muted-foreground">Value</div>
                            <Input :model-value="property.value" placeholder="Example: 127.0.0.1" @update:model-value="(value) =>
                                updatePropertyValue(metadataIndex, propertyIndex, value ?? '')
                            " />
                        </div>

                        <div class="flex items-end justify-end">
                            <Button variant="ghost" size="icon"
                                class="h-8 w-8 text-muted-foreground hover:text-destructive"
                                @click="removeProperty(metadataIndex, propertyIndex)">
                                <Icon icon="lucide:minus" class="h-4 w-4" />
                            </Button>
                        </div>
                    </div>
                </div>

                <Button variant="outline" size="sm" @click="addProperty(metadataIndex)">
                    <Icon icon="lucide:plus" class="mr-2 h-4 w-4" />
                    Add Property
                </Button>
            </div>
        </div>

        <div class="pt-2">
            <Button variant="outline" @click="addMetadataItem">
                <Icon icon="lucide:plus" class="mr-2 h-4 w-4" />
                Add Metadata
            </Button>
        </div>
    </div>
</template>
