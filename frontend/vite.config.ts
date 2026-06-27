import path from 'node:path';
import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';

export default defineConfig({
    plugins: [solidPlugin()],
    build: {
        target: 'esnext',
    },
    resolve: {
        alias: {
            src: path.resolve(__dirname, "src"),
            wailsjs: path.resolve(__dirname, "wailsjs"),
        },
    },
    esbuild: {
        sourcemap: false,
    }
});
