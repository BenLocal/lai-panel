declare module "js-yaml" {
  export function load(str: string, opts?: unknown): any;
}

declare module "@guolao/vue-monaco-editor" {
  import type { DefineComponent } from "vue";
  const MonacoEditor: DefineComponent<
    Record<string, unknown>,
    Record<string, unknown>,
    any
  >;
  export default MonacoEditor;
}
