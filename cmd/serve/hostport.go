package serve

import (
	"net"

	"r.tomng.dev/goserve/internal/logger"
)

func parseHostPort(hostport string) (string, string) {
	// Cannot split without a colon for port
	// Add a colon to split then use default port
	if !hasPort(hostport) {
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

func hasPort(hostport string) bool {
	for i := len(hostport) - 1; i >= 0; i-- {
		if hostport[i] == ':' {
			return true
		} else if hostport[i] < '0' || hostport[i] > '9' {
			return false
		}
	}

	return false
}
