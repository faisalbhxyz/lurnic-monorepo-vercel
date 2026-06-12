import type { NextConfig } from "next";
import path from "path";
import { fileURLToPath } from "url";

const frontendDir = path.dirname(fileURLToPath(import.meta.url));

const nextConfig: NextConfig = {
  output: "standalone",
  // Monorepo: trace from repo root, but resolve React only from frontend/node_modules.
  outputFileTracingRoot: path.join(frontendDir, ".."),
  webpack: (config) => {
    const react = path.join(frontendDir, "node_modules/react");
    const reactDom = path.join(frontendDir, "node_modules/react-dom");
    config.resolve.alias = {
      ...config.resolve.alias,
      react,
      "react-dom": reactDom,
      "react/jsx-runtime": path.join(react, "jsx-runtime.js"),
      "react/jsx-dev-runtime": path.join(react, "jsx-dev-runtime.js"),
    };
    return config;
  },
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    remotePatterns: [
      {
        protocol: "http",
        hostname: "*", // Allow any hostname
        port: "", // Optional: allow any port
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "*", // Allow any hostname
        port: "", // Optional: allow any port
        pathname: "/**",
      },
      {
        protocol: "http",
        hostname: "localhost",
        port: "5002",
      },
      {
        protocol: "http",
        hostname: "backend",
        port: "5000",
        pathname: "/uploads/**",
      },
    ],
  },
};

export default nextConfig;
