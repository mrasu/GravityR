package lib

import (
	"log"
	"net/http"
	"time"
)

type HttpTransport struct {
	Transport http.RoundTripper
}

func NewHttpTransport() *HttpTransport {
	return &HttpTransport{Transport: http.DefaultTransport}
}

func (ht *HttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("[http] --> %s %s", req.Method, req.URL)
	startTime := time.Now()

	resp, err := ht.Transport.RoundTrip(req)

	duration := time.Since(startTime)
	duration /= time.Millisecond

	if err != nil {
		log.Printf("[http] <-- ERROR method=%s host=%s path=%s status=error durationMs=%d error=%q", req.Method, req.Host, req.URL.Path, duration, err.Error())
		return nil, err
	}

	log.Printf("[http] <-- method=%s host=%s path=%s status=%d durationMs=%d", req.Method, req.Host, req.URL.Path, resp.StatusCode, duration)
	return resp, nil
}
