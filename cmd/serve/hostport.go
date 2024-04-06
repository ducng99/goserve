package serve

import (
	"net"
	"strings"

	"r.tomng.dev/goserve/internal/logger"
)

func parseHostPort(hostport string) (string, string) {
	// Cannot split without a colon
	// Add a colon to split then use default port
	if !strings.Contains(hostport, ":") {
		hostport = hostport + ":"
	}

	_host, _port, err := net.SplitHostPort(hostport)
	if err != nil {
		logger.Fatalf("Invalid address (%v)\n", err)
	}

	host := DefaultListenHost
	port := DefaultListenPort

	if _host != "" {
		host = _host
	}

	if _port != "" {
		port = _port
	}

	return host, port
}
