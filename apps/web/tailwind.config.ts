import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        nexus: {
          bg: "#09090b",
          panel: "#111827",
          accent: "#8b5cf6"
        }
      }
    }
  },
  plugins: []
};

export default config;

