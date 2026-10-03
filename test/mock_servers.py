import json
import signal
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


SERVERS = {
    3000: "app",
    4000: "app-test",
    5000: "api",
}


class MockHandler(BaseHTTPRequestHandler):
    def do_GET(self) -> None:
        body = json.dumps(
            {
                "server": self.server.server_name,
                "port": self.server.server_port,
                "method": self.command,
                "host": self.headers.get("Host"),
                "path": self.path,
                "headers": dict(self.headers),
            },
            indent=2,
        ).encode()

        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, format: str, *args: object) -> None:
        print(f"[{self.server.server_name}:{self.server.server_port}] {format % args}")


def main() -> None:
    servers = []

    for port, name in SERVERS.items():
        server = ThreadingHTTPServer(("127.0.0.1", port), MockHandler)
        server.server_name = name
        thread = threading.Thread(target=server.serve_forever, daemon=True)
        thread.start()
        servers.append(server)
        print(f"Mock {name} server listening on http://localhost:{port}")

    stop = threading.Event()

    def request_shutdown(_signum: int, _frame: object) -> None:
        stop.set()

    signal.signal(signal.SIGINT, request_shutdown)
    signal.signal(signal.SIGTERM, request_shutdown)
    stop.wait()

    for server in servers:
        server.shutdown()
        server.server_close()


if __name__ == "__main__":
    main()
