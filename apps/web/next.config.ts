import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  poweredByHeader: false,
  reactStrictMode: true,
  transpilePackages: ["@nexus/config", "@nexus/logger", "@nexus/types"]
};

export default nextConfig;
