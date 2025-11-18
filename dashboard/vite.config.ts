import tailwindcss from "@tailwindcss/vite";
import vue from "@vitejs/plugin-vue";
import { defineConfig } from "vite";
import path from "node:path";
import tsconfigPaths from "vite-tsconfig-paths";
import vueDevTools from "vite-plugin-vue-devtools";

// https://vite.dev/config/
const isDebug = process.env.NODE_ENV === "development" || process.env.DEBUG === "true";

export default defineConfig({
    plugins: [
        vue(),
        tailwindcss(),
        tsconfigPaths(),
        ...(isDebug ? [vueDevTools()] : []),
    ],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
    server: {
        proxy: {
            "/api": {
                target: "http://localhost:8080",
                changeOrigin: true,
            },
        },
    },
});
