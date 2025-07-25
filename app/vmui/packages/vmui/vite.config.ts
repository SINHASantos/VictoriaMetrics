import * as path from "path";

import { defineConfig, ProxyOptions } from "vite";
import preact from "@preact/preset-vite";
import dynamicIndexHtmlPlugin from "./config/plugins/dynamicIndexHtml";

const getProxy = (): Record<string, ProxyOptions> | undefined => {
  const playground = process.env.PLAYGROUND;

  switch (playground) {
    case "METRICS": {
      return {
        "^/vmalert/.*": {
          target: "https://play.victoriametrics.com",
          changeOrigin: true,
          configure: (proxy) => {
            proxy.on("error", (err) => {
              console.error("[proxy error]", err.message);
            });
          }
        },
        "^/api/.*": {
          target: "https://play.victoriametrics.com/select/0/prometheus/",
          changeOrigin: true,
          configure: (proxy) => {
            proxy.on("error", (err) => {
              console.error("[proxy error]", err.message);
            });
          }
        },
        "^/prometheus/api.*": {
          target: "https://play.victoriametrics.com/select/0/",
          changeOrigin: true,
          configure: (proxy) => {
            proxy.on("error", (err) => {
              console.error("[proxy error]", err.message);
            });
          }
        }
      };
    }
    default: {
      return undefined;
    }
  }
};

export default defineConfig(({ mode }) => {
  return {
    base: "",
    plugins: [
      preact(),
      dynamicIndexHtmlPlugin({ mode })
    ],
    assetsInclude: ["**/*.md"],
    server: {
      open: true,
      port: 3000,
      proxy: getProxy(),
    },
    resolve: {
      alias: {
        "src": path.resolve(__dirname, "src"),
      },
    },
    build: {
      outDir: "./build",
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes("node_modules")) {
              return "vendor";
            }
          }
        }
      }
    },
  };
});



