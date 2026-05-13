import unittest

from nexus_ai import get_health_status
from nexus_ai.__main__ import HealthHandler
from http.server import ThreadingHTTPServer
from threading import Thread
from urllib.request import urlopen


class HealthStatusTest(unittest.TestCase):
    def test_health_status_is_ok(self) -> None:
        status = get_health_status()

        self.assertEqual(status.service, "ai")
        self.assertEqual(status.status, "ok")
        self.assertIn("time", status.to_dict())

    def test_health_endpoint_responds(self) -> None:
        server = ThreadingHTTPServer(("127.0.0.1", 0), HealthHandler)
        thread = Thread(target=server.serve_forever, daemon=True)
        thread.start()

        try:
            with urlopen(f"http://127.0.0.1:{server.server_port}/healthz", timeout=2) as response:
                self.assertEqual(response.status, 200)
                self.assertIn(b'"service": "ai"', response.read())
        finally:
            server.shutdown()
            server.server_close()


if __name__ == "__main__":
    unittest.main()
