from dataclasses import dataclass
from datetime import UTC, datetime


@dataclass(frozen=True)
class HealthStatus:
    service: str
    status: str
    time: str

    def to_dict(self) -> dict[str, str]:
        return {
            "service": self.service,
            "status": self.status,
            "time": self.time,
        }


def get_health_status() -> HealthStatus:
    return HealthStatus(
        service="ai",
        status="ok",
        time=datetime.now(UTC).isoformat(),
    )

