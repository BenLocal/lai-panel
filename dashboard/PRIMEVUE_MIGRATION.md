# PrimeVue 迁移指南

## 已完成的迁移

1. ✅ 安装 PrimeVue 和 PrimeIcons
2. ✅ 配置 main.ts 和 App.vue
3. ✅ 替换 Toast 组件（vue-sonner -> PrimeVue Toast）
4. ✅ 替换 Login.vue
5. ✅ 替换 Dashboard.vue
6. ✅ 替换 AppLayout.vue
7. ✅ 替换 Docker.vue
8. ✅ 替换 NodeTerminal.vue
9. ✅ 替换 Settings.vue
10. ✅ 替换 Applications.vue
11. ✅ 删除 src/components/ui 目录
12. ✅ 删除 components.json 文件
13. ✅ 从 package.json 中移除 shadcn-vue 相关依赖
14. ✅ 运行 npm install 清理依赖

## 组件映射表

### 基础组件
- `Button` → `primevue/button` (Button)
- `Input` → `primevue/inputtext` (InputText)
- `Label` → 使用原生 `<label>` 或 PrimeVue 的 Fieldset
- `Separator` → 使用 `<hr>` 或 PrimeVue 的 Divider
- `Checkbox` → `primevue/checkbox` (Checkbox)

### 复杂组件
- `Card` → 使用自定义 div + Tailwind 类，或 PrimeVue 的 Panel
- `Dialog` → `primevue/dialog` (Dialog)
- `Select` → `primevue/select` (Select)
- `Table` → `primevue/datatable` (DataTable) + Column
- `DropdownMenu` → `primevue/menu` (Menu) 或 `primevue/tieredmenu` (TieredMenu)

### 布局组件
- `Sidebar` → 自定义实现（已替换）
- `Sheet` → `primevue/sidebar` (Sidebar) 或 Dialog
- `Drawer` → `primevue/sidebar` (Sidebar)
- `ResizablePanel` → 自定义实现（已替换）
- `Breadcrumb` → `primevue/breadcrumb` (Breadcrumb)
- `Tabs` → `primevue/tabview` (TabView) + TabPanel

### 高级组件
- `AlertDialog` → `primevue/confirmdialog` (ConfirmDialog) 或 Dialog
- `Stepper` → 自定义实现或使用 PrimeVue 的 Steps
- `Tooltip` → `primevue/tooltip` (Tooltip)
- `Skeleton` → `primevue/skeleton` (Skeleton)
- `Calendar` → `primevue/calendar` (Calendar)

## 需要替换的文件

### 视图文件
- [x] Settings.vue ✅
- [x] Applications.vue ✅
- [ ] EnvironmentVariables.vue
- [ ] Nodes.vue
- [ ] Services.vue
- [ ] docker/DockerContainers.vue
- [ ] docker/DockerImages.vue
- [ ] docker/DockerVolumes.vue
- [ ] docker/DockerNetworks.vue
- [ ] docker/DockerContainerTerminal.vue

### 组件文件
- [ ] application/TooltipWithCopy.vue
- [ ] application/MetadataEditor.vue
- [ ] application/ApplicationWorkspaceManager.vue
- [ ] application/ApplicationQAEditor.vue

## 替换步骤示例

### 示例 1: 替换 Button
```vue
<!-- 之前 -->
import { Button } from "@/components/ui/button";
<Button variant="ghost" size="sm">Click</Button>

<!-- 之后 -->
import Button from 'primevue/button'
<Button text rounded size="small">Click</Button>
```

### 示例 2: 替换 Input
```vue
<!-- 之前 -->
import { Input } from "@/components/ui/input";
<Input v-model="value" placeholder="Enter text" />

<!-- 之后 -->
import InputText from 'primevue/inputtext'
<InputText v-model="value" placeholder="Enter text" />
```

### 示例 3: 替换 Select
```vue
<!-- 之前 -->
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
<Select v-model="value">
  <SelectTrigger>
    <SelectValue placeholder="Select" />
  </SelectTrigger>
  <SelectContent>
    <SelectItem value="1">Option 1</SelectItem>
  </SelectContent>
</Select>

<!-- 之后 -->
import Select from 'primevue/select'
<Select v-model="value" :options="options" optionLabel="label" optionValue="value" placeholder="Select" />
```

### 示例 4: 替换 Table
```vue
<!-- 之前 -->
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Name</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    <TableRow>
      <TableCell>Value</TableCell>
    </TableRow>
  </TableBody>
</Table>

<!-- 之后 -->
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
<DataTable :value="data">
  <Column field="name" header="Name"></Column>
</DataTable>
```

### 示例 5: 替换 Dialog
```vue
<!-- 之前 -->
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
<Dialog v-model:open="visible">
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Title</DialogTitle>
    </DialogHeader>
    Content
    <DialogFooter>
      <Button>OK</Button>
    </DialogFooter>
  </DialogContent>
</Dialog>

<!-- 之后 -->
import Dialog from 'primevue/dialog'
<Dialog v-model:visible="visible" modal header="Title">
  Content
  <template #footer>
    <Button label="OK" />
  </template>
</Dialog>
```

### 示例 6: 替换 Card
```vue
<!-- 之前 -->
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
<Card>
  <CardHeader>
    <CardTitle>Title</CardTitle>
  </CardHeader>
  <CardContent>Content</CardContent>
</Card>

<!-- 之后 -->
<!-- 使用自定义 div 或 PrimeVue Panel -->
<div class="p-6 bg-card rounded-lg border">
  <div class="mb-4">
    <h3 class="text-lg font-semibold">Title</h3>
  </div>
  Content
</div>
```

## 注意事项

1. PrimeVue 的 Select 组件使用 `v-model` 绑定整个对象，而不是值
2. PrimeVue 的 Dialog 使用 `v-model:visible` 而不是 `v-model:open`
3. PrimeVue 的 Button 使用 `severity` 属性而不是 `variant`
4. PrimeVue 的 DataTable 需要单独导入 Column 组件
5. 某些组件（如 Sidebar, ResizablePanel）需要自定义实现

## 清理步骤

完成所有替换后：
1. ✅ 删除 `src/components/ui` 目录
2. ✅ 删除 `components.json` 文件
3. ✅ 从 package.json 中移除 shadcn-vue 相关依赖：
   - reka-ui ✅
   - vaul-vue ✅
   - vue-sonner ✅
   - class-variance-authority ✅
   - clsx ✅
   - tailwind-merge ✅
4. ✅ 运行 `npm install` 清理依赖

## 剩余工作

以下文件仍需要按照组件映射表继续替换：
- EnvironmentVariables.vue
- Nodes.vue
- Services.vue
- docker/DockerContainers.vue
- docker/DockerImages.vue
- docker/DockerVolumes.vue
- docker/DockerNetworks.vue
- docker/DockerContainerTerminal.vue
- application/TooltipWithCopy.vue
- application/MetadataEditor.vue
- application/ApplicationWorkspaceManager.vue
- application/ApplicationQAEditor.vue

替换方法：按照上面的组件映射表和替换步骤示例，逐个文件替换即可。
