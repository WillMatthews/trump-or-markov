import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/v1": {
        target: "http://localhost:1776",
        changeOrigin: true,
        secure: false, // This is optional, but might be necessary if you're using HTTPS.
        rewrite: (path) => path.replace(/^\/v1/, "/v1"), // Optional: Adjust the path if needed
      },
    },
  },
});
