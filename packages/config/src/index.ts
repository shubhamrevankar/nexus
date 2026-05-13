export interface PublicConfig {
  appName: string;
}

export function getPublicConfig(): PublicConfig {
  return {
    appName: process.env.NEXT_PUBLIC_APP_NAME ?? "Nexus"
  };
}

