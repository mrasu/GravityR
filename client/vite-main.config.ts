import { defineConfig } from "vite";
import { config } from "./vite-common.config";

config["build"]["rollupOptions"]["input"] = ["src/main.ts"];

export default defineConfig(config);
