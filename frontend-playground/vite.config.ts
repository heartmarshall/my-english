import path from "path"
import tailwindcss from "@tailwindcss/vite"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      // @graphql-typed-document-node/core is types-only, provide a runtime shim
      "@graphql-typed-document-node/core": path.resolve(__dirname, "./src/gql/typed-document-node-shim.ts"),
    },
  },
  optimizeDeps: {
    include: ['@apollo/client', '@apollo/client/react', 'graphql'],
  },
  server: {
    proxy: {
      '/query': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
