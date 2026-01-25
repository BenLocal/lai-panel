<script setup lang="ts">
import { ref, watch } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
interface MetadataProperty {
  key: string;
  value: string;
}

interface EditableMetadata {
  name: string;
  properties: MetadataProperty[];
}

export type { EditableMetadata, MetadataProperty };

type MetadataEditorEmits = (
  event: "update:modelValue",
  value: EditableMetadata[]
) => void;

const props = defineProps<{
  modelValue: EditableMetadata[];
}>();

const emit = defineEmits<MetadataEditorEmits>();

const editableMetadata = ref<EditableMetadata[]>([]);

const cloneMetadata = (items: EditableMetadata[] = []): EditableMetadata[] =>
  items.map((item) => ({
    name: item.name ?? "",
    properties: item.properties?.map((property) => ({
      key: property.key ?? "",
      value: property.value ?? "",
    })) ?? [
      {
        key: "",
        value: "",
      },
    ],
  }));

watch(
  () => props.modelValue,
  (value) => {
    editableMetadata.value = cloneMetadata(value);
  },
  { immediate: true, deep: true }
);

const emitChange = () => {
  emit("update:modelValue", cloneMetadata(editableMetadata.value));
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
  if (!current) return;
  const properties = [...current.properties];
  const property = properties[propertyIndex];
  if (!property) return;
  properties[propertyIndex] = { ...property, value: String(value) };
  const next = [...editableMetadata.value];
  next[metadataIndex] = { ...current, properties };
  editableMetadata.value = next;
  emitChange();
};
</script>

<style scoped>
.metadata-editor { display: flex; flex-direction: column; gap: 1rem; }
.metadata-empty {
  padding: 1.5rem;
  text-align: center;
  border: 1px dashed var(--p-surface-border);
  border-radius: var(--p-border-radius);
  background: var(--p-surface-50);
  font-size: 0.875rem;
}
.metadata-item {
  padding: 1rem;
  border: 1px solid var(--p-surface-border);
  border-radius: var(--p-border-radius);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.metadata-head { display: flex; justify-content: space-between; align-items: center; }
.font-medium { font-weight: 500; font-size: 0.875rem; }
.metadata-props { display: flex; flex-direction: column; gap: 0.75rem; }
.text-uppercase { font-size: 0.75rem; font-weight: 500; }
.prop-row { display: grid; grid-template-columns: 1fr 1fr auto; gap: 1rem; align-items: end; }
@media (max-width: 768px) { .prop-row { grid-template-columns: 1fr; } }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.75rem; font-weight: 500; }
.btn-icon-text { margin-left: 0.5rem; }
</style>

<template>
  <div class="metadata-editor">
    <div v-if="!editableMetadata.length" class="metadata-empty text-muted-foreground">
      No metadata defined yet. Click the button below to add one.
    </div>

    <div v-for="(metadata, mi) in editableMetadata" :key="mi" class="metadata-item">
      <div class="metadata-head">
        <span class="font-medium">Metadata {{ mi + 1 }}</span>
        <Button text rounded size="small" severity="danger" @click="removeMetadataItem(mi)" v-tooltip.top="'Remove'">
          <i class="pi pi-trash"></i>
        </Button>
      </div>

      <div class="form-group">
        <label :for="`metadata-name-${mi}`">Name</label>
        <InputText
          :id="`metadata-name-${mi}`"
          :model-value="metadata.name"
          placeholder="Example: database"
          @update:model-value="(v) => updateMetadataName(mi, v ?? '')"
        />
      </div>

      <div class="metadata-props">
        <div class="text-muted-foreground text-uppercase">Properties</div>
        <div v-for="(prop, pi) in metadata.properties" :key="`${mi}-${pi}`" class="prop-row">
          <div class="form-group">
            <label class="text-muted-foreground">Key</label>
            <InputText
              :model-value="prop.key"
              placeholder="Example: host"
              @update:model-value="(v) => updatePropertyKey(mi, pi, v ?? '')"
            />
          </div>
          <div class="form-group">
            <label class="text-muted-foreground">Value</label>
            <InputText
              :model-value="prop.value"
              placeholder="Example: 127.0.0.1"
              @update:model-value="(v) => updatePropertyValue(mi, pi, v ?? '')"
            />
          </div>
          <Button text rounded size="small" severity="danger" @click="removeProperty(mi, pi)" v-tooltip.top="'Remove'">
            <i class="pi pi-minus"></i>
          </Button>
        </div>
        <Button outlined size="small" @click="addProperty(mi)">
          <i class="pi pi-plus"></i>
          <span class="btn-icon-text">Add Property</span>
        </Button>
      </div>
    </div>

    <Button outlined @click="addMetadataItem">
      <i class="pi pi-plus"></i>
      <span class="btn-icon-text">Add Metadata</span>
    </Button>
  </div>
</template>
