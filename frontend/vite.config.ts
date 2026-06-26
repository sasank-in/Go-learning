import { defineConfig } from "vite";

export default defineConfig({
  server: {
    port: 5173,
    // Proxy API calls to the Go backend during development so the frontend
    // can use same-origin "/api/..." paths.
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
