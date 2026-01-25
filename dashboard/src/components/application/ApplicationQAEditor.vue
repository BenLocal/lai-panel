<script setup lang="ts">
import { ref, watch } from "vue";
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import Checkbox from 'primevue/checkbox'
import type { ApplicationQAItem } from "@/api/application";

type QAType = ApplicationQAItem["type"];

const props = defineProps<{ modelValue: ApplicationQAItem[] }>();
const emit = defineEmits<{ "update:modelValue": [value: ApplicationQAItem[]] }>();

const qaItems = ref<ApplicationQAItem[]>([]);

const typeOptions: { label: string; value: QAType }[] = [
  { label: "Text", value: "text" },
  { label: "Number", value: "number" },
  { label: "Boolean", value: "boolean" },
  { label: "Select", value: "select" },
  { label: "Textarea", value: "textarea" },
];

const normalize = (item: ApplicationQAItem): ApplicationQAItem => {
  const n = {
    name: item.name ?? "",
    type: (item.type ?? "text") as QAType,
    default_value: item.default_value ?? "",
    options: item.options,
    required: item.required ?? false,
    description: item.description ?? "",
  };
  if (n.type === "select") {
    n.options = n.options?.length ? n.options : ["Option 1", "Option 2"];
    if (n.default_value && n.options && !n.options.includes(n.default_value)) n.default_value = n.options[0] ?? "";
  } else if (n.type === "boolean") {
    n.default_value = n.default_value === "true" ? "true" : "false";
    n.options = undefined;
  } else n.options = undefined;
  return n;
};

const syncFromProps = (items: ApplicationQAItem[] | undefined) => {
  qaItems.value = (items ?? []).map((i) => normalize({ ...i }));
};

watch(() => props.modelValue, syncFromProps, { immediate: true, deep: true });

const emitChange = () => {
  emit("update:modelValue", qaItems.value.map((i) => normalize({ ...i })));
};

const createEmpty = (): ApplicationQAItem => ({
  name: "",
  type: "text",
  default_value: "",
  required: false,
  description: "",
});

const addItem = () => {
  qaItems.value = [...qaItems.value, createEmpty()];
  emitChange();
};

const removeItem = (index: number) => {
  const next = [...qaItems.value];
  next.splice(index, 1);
  qaItems.value = next;
  emitChange();
};

const updateField = <K extends keyof ApplicationQAItem>(index: number, key: K, value: ApplicationQAItem[K]) => {
  const cur = qaItems.value[index];
  if (!cur) return;
  const next = [...qaItems.value];
  const item = { ...cur, [key]: value } as ApplicationQAItem;
  if (key === "type") {
    const t = value as QAType;
    item.type = t;
    if (t === "select") {
      item.options = cur.options?.length ? cur.options : ["Option 1", "Option 2"];
      item.default_value = item.options[0] ?? "";
    } else if (t === "boolean") {
      item.options = undefined;
      item.default_value = cur.default_value === "true" ? "true" : "false";
    } else item.options = undefined;
  }
  if (key === "options") {
    const opts = (value as string[]).filter((s) => s.trim().length);
    item.options = opts;
    if (item.default_value && !opts.includes(item.default_value)) item.default_value = opts[0] ?? "";
  }
  next[index] = normalize(item);
  qaItems.value = next;
  emitChange();
};

const handleOptionsInput = (index: number, raw: string) => {
  const opts = raw
    .split("\n")
    .map((s) => s.trim())
    .filter(Boolean);
  updateField(index, "options", opts);
};

type StrKey = "name" | "default_value" | "description";
const handleStr = (index: number, key: StrKey, v: string | number) => {
  updateField(index, key, String(v) as ApplicationQAItem[StrKey]);
};
</script>

<template>
  <div class="qa-editor">
    <div v-if="!qaItems.length" class="qa-empty text-muted-foreground">
      No QA configuration yet. Click the button below to add one.
    </div>

    <div v-for="(qa, index) in qaItems" :key="index" class="qa-item">
      <div class="qa-item-head">
        <span class="font-medium">QA Item {{ index + 1 }}</span>
        <Button text rounded size="small" severity="danger" @click="removeItem(index)" v-tooltip.top="'Remove'">
          <i class="pi pi-trash"></i>
        </Button>
      </div>

      <div class="qa-grid">
        <div class="form-group">
          <label :for="`qa-name-${index}`">Field Name</label>
          <InputText
            :id="`qa-name-${index}`"
            :model-value="qa.name"
            placeholder="Unique key"
            @update:model-value="(v) => handleStr(index, 'name', v ?? '')"
          />
        </div>
        <div class="form-group">
          <label>Type</label>
          <Select
            :model-value="qa.type"
            :options="typeOptions"
            option-label="label"
            option-value="value"
            placeholder="Select type"
            @update:model-value="(v) => updateField(index, 'type', (v ?? 'text') as QAType)"
          />
        </div>
      </div>

      <div class="qa-grid">
        <div class="form-group">
          <label>Default Value</label>
          <InputText
            v-if="qa.type === 'text' || qa.type === 'number'"
            :model-value="qa.default_value ?? ''"
            :type="qa.type === 'number' ? 'number' : 'text'"
            placeholder="Default value"
            @update:model-value="(v: string | undefined) => handleStr(index, 'default_value', v ?? '')"
          />
          <Select
            v-else-if="qa.type === 'boolean'"
            :model-value="qa.default_value ?? 'false'"
            :options="['true', 'false']"
            placeholder="Default"
            @update:model-value="(v) => updateField(index, 'default_value', (v ?? 'false') as string)"
          />
          <Select
            v-else-if="qa.type === 'select'"
            :model-value="qa.default_value ?? ''"
            :options="qa.options ?? []"
            placeholder="Default"
            @update:model-value="(v: unknown) => updateField(index, 'default_value', (v ?? '') as string)"
          />
          <Textarea
            v-else
            :model-value="qa.default_value ?? ''"
            placeholder="Default value"
            rows="2"
            @update:model-value="(v) => handleStr(index, 'default_value', v ?? '')"
          />
        </div>
        <div class="form-group">
          <label>Required</label>
          <div class="qa-required-row">
            <Checkbox
              :id="`qa-required-${index}`"
              :model-value="qa.required ?? false"
              binary
              @update:model-value="(c) => updateField(index, 'required', !!c)"
            />
            <label :for="`qa-required-${index}`" class="text-muted-foreground">Must be provided</label>
          </div>
        </div>
      </div>

      <div class="form-group">
        <label>Description</label>
        <Textarea
          :model-value="qa.description ?? ''"
          placeholder="Optional description"
          rows="2"
          @update:model-value="(v: string | undefined) => handleStr(index, 'description', v ?? '')"
        />
      </div>

      <div v-if="qa.type === 'select'" class="form-group">
        <label>Options (one per line)</label>
        <Textarea
          :model-value="(qa.options ?? []).join('\n')"
          placeholder="Option 1"
          rows="3"
          @update:model-value="(v: string | undefined) => handleOptionsInput(index, v ?? '')"
        />
      </div>
    </div>

    <Button outlined @click="addItem">
      <i class="pi pi-plus"></i>
      <span class="btn-icon-text">Add QA Item</span>
    </Button>
  </div>
</template>

<style scoped>
.qa-editor { display: flex; flex-direction: column; gap: 1rem; }
.qa-empty {
  padding: 1.5rem;
  text-align: center;
  border: 1px dashed var(--p-surface-border);
  border-radius: var(--p-border-radius);
  background: var(--p-surface-50);
  font-size: 0.875rem;
}
.qa-item {
  padding: 1rem;
  border: 1px solid var(--p-surface-border);
  border-radius: var(--p-border-radius);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.qa-item-head { display: flex; justify-content: space-between; align-items: center; }
.font-medium { font-weight: 500; font-size: 0.875rem; }
.qa-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
@media (max-width: 640px) { .qa-grid { grid-template-columns: 1fr; } }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; }
.qa-required-row { display: flex; align-items: center; gap: 0.5rem; }
.btn-icon-text { margin-left: 0.5rem; }
</style>
