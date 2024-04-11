package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"r.tomng.dev/goserve/internal/logger"
)

// Creates a new reverse proxy handler to the target URL.
//
// incHeaders: a boolean that determines whether to include X-Forwarded-For and X-Forwarded-Proto headers in the proxy request.
//
// ignoreRedirect: a boolean that determines whether to ignore redirects from the target server.
// Basically strips out Location header and return.
func New(targetURL string, incHeaders bool, ignoreRedirect bool) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			if incHeaders {
				req.Header.Set("X-Forwarded-For", req.RemoteAddr)
				req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
			}

			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = strings.TrimSuffix(target.Path, "/") + "/" + strings.TrimPrefix(req.URL.Path, "/")

			if req.URL.RawQuery == "" {
				req.URL.RawQuery = target.RawQuery
			} else {
				req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
			}
		},
		ModifyResponse: func(res *http.Response) error {
			location := res.Header.Get("Location")

			if ignoreRedirect {
				res.Header.Del("Location")
				res.Header.Set("X-Original-Location", location)
			} else {
				if location != "" && location != res.Request.URL.String() {
					newLocation := strings.Replace(location, targetURL, "", 1)
					res.Header.Set("Location", newLocation)
				}
			}

			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Printf(logger.LogError, "%v\n", err)
			http.Error(w, "Error proxying request", http.StatusBadGateway)
		},
	}

	return reverseProxy, nil
}
