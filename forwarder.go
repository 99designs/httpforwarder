package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// AsyncHttpForwarder forwards HTTP requests
// asynchronously, returning an empty 200 response
type AsyncHttpForwarder struct {
	// The transport used to perform proxy requests.
	// If nil, http.DefaultTransport is used.
	Transport http.RoundTripper
}

func NewAsyncHttpForwarder() *AsyncHttpForwarder {
	return &AsyncHttpForwarder{}
}

func (p *AsyncHttpForwarder) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// forward the request asyncronously
	go p.forward(p.copyRequest(req))
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// Hop-by-hop headers. These are removed when the request is forwarded
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func (p *AsyncHttpForwarder) copyRequest(req *http.Request) *http.Request {
	outreq := new(http.Request)
	*outreq = *req // includes shallow copies of maps, but okay

	target, err := url.Parse(req.RequestURI)
	if err != nil {
		log.Printf("http: url parse error: %v", err)
		return nil
	}

	outreq.URL = target
	outreq.Proto = "HTTP/1.1"
	outreq.ProtoMajor = 1
	outreq.ProtoMinor = 1
	outreq.Close = false

	// Remove hop-by-hop headers.  Especially
	// important is "Connection" because we want a persistent
	// connection, regardless of what the client sent to us.  This
	// is modifying the same underlying map from req (shallow
	// copied above) so we only copy it if necessary.
	copiedHeaders := false
	for _, h := range hopHeaders {
		if outreq.Header.Get(h) != "" {
			if !copiedHeaders {
				outreq.Header = make(http.Header)
				copyHeader(outreq.Header, req.Header)
				copiedHeaders = true
			}
			outreq.Header.Del(h)
		}
	}

	b := bytes.NewBuffer([]byte{})
	b.ReadFrom(req.Body)
	outreq.Body = ioutil.NopCloser(b)

	return outreq
}

func (p *AsyncHttpForwarder) forward(outreq *http.Request) {
	transport := p.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	res, err := transport.RoundTrip(outreq)

	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	log.Println(res.StatusCode, outreq.Method, outreq.URL.String())
	defer res.Body.Close()
}
