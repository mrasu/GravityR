import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vitejs.dev/config/
export const config = {
  plugins: [svelte()],
  build: {
    emptyOutDir: false,
    rollupOptions: {
      input: ["src/dummy.ts"],
      output: {
        assetFileNames: "assets/[name].[ext]",
        entryFileNames: "assets/[name].js",
      },
    },
  },
  resolve: {
    alias: [{ find: "@", replacement: "/src" }],
  },
};
