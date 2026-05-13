export type ServiceStatus = "ok" | "degraded" | "unavailable";

export interface HealthStatus {
  service: string;
  status: ServiceStatus;
  version?: string;
}

