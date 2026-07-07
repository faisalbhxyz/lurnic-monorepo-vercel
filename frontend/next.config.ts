import type { NextConfig } from "next";
import path from "path";
import { fileURLToPath } from "url";

const frontendDir = path.dirname(fileURLToPath(import.meta.url));

const nextConfig: NextConfig = {
  output: "standalone",
  // Monorepo: trace deps from repo root so standalone layout is standalone/frontend/ in Docker.
  outputFileTracingRoot: path.join(frontendDir, ".."),
  // Keep webpack for raw CSS imports; Next 16 defaults to Turbopack.
  turbopack: {},
  typescript: {
    ignoreBuildErrors: true,
  },
  webpack: (config) => {
    config.module.rules.forEach((rule) => {
      if (typeof rule === "object" && rule !== null && "oneOf" in rule && Array.isArray(rule.oneOf)) {
        rule.oneOf.unshift({
          test: /\.css$/i,
          resourceQuery: /raw/,
          type: "asset/source",
        });
      }
    });
    return config;
  },
  images: {
    remotePatterns: [
      {
        protocol: "http",
        hostname: "*",
        port: "",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "*",
        port: "",
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
