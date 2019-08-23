package utils

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type ProxyHandle struct {
	proxy_host       string
	proxy_port       string
	proxy_fixed_path string
}

var (
	RealDataProxy  *ProxyHandle
	RealNewsProxy  *ProxyHandle
	RealUserProxy  *ProxyHandle
	RealGroupProxy *ProxyHandle
	RealMallProxy  *ProxyHandle
)

func (p *ProxyHandle) GetProxyAddress() string {
	return fmt.Sprintf("http://%s:%s%s", p.proxy_host, p.proxy_port, p.proxy_fixed_path)
}

func (p *ProxyHandle) GetProxyHostUrl() string {
	return fmt.Sprintf("http://%s:%s", p.proxy_host, p.proxy_port)
}

var onExitFlushLoop func()

type ReverseProxy struct {
	Director      func(*http.Request)
	Transport     http.RoundTripper
	FlushInterval time.Duration
	ErrorLog      *log.Logger
	Outreq        *http.Request  //代理请求
	Outres        *http.Response //代理应答
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func NewSingleHostReverseProxy(target *url.URL) *ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.RequestURI = target.String()
		//		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &ReverseProxy{Director: director}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// Hop-by-hop headers. These are removed when sent to the backend.
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

func (p *ReverseProxy) HandleRequest(req *http.Request) {
	if p.Transport == nil {
		p.Transport = http.DefaultTransport
	}

	p.Outreq = new(http.Request)
	*(p.Outreq) = *req // includes shallow copies of maps, but okay

	p.Director(p.Outreq)
	p.Outreq.Proto = "HTTP/1.1"
	p.Outreq.ProtoMajor = 1
	p.Outreq.ProtoMinor = 1
	p.Outreq.Close = false

	// Remove hop-by-hop headers to the backend.  Especially
	// important is "Connection" because we want a persistent
	// connection, regardless of what the client sent to us.  This
	// is modifying the same underlying map from req (shallow
	// copied above) so we only copy it if necessary.
	copiedHeaders := false
	for _, h := range hopHeaders {
		if p.Outreq.Header.Get(h) != "" {
			if !copiedHeaders {
				p.Outreq.Header = make(http.Header)
				copyHeader(p.Outreq.Header, req.Header)
				copiedHeaders = true
			}
			p.Outreq.Header.Del(h)
		}
	}

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {

		if prior, ok := p.Outreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		p.Outreq.Header.Set("X-Forwarded-For", clientIP)
	}
}

func (p *ReverseProxy) Request() error {
	res, err := p.Transport.RoundTrip(p.Outreq)
	if err != nil {
		p.logf("http: proxy error: %v", err)
		return err
	}
	p.Outres = res
	return nil
}

func (p *ReverseProxy) HandleResponse(rw http.ResponseWriter, bytebody []byte) {
	for _, h := range hopHeaders {
		p.Outres.Header.Del(h)
	}

	copyHeader(rw.Header(), p.Outres.Header)

	rw.WriteHeader(p.Outres.StatusCode)
	rw.Write(bytebody)
	//	p.CopyResponse(rw, p.Outres.Body)
	p.Outres.Body.Close()
}

func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	transport := p.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	outreq := new(http.Request)
	*outreq = *req // includes shallow copies of maps, but okay

	p.Director(outreq)
	outreq.Proto = "HTTP/1.1"
	outreq.ProtoMajor = 1
	outreq.ProtoMinor = 1
	outreq.Close = false

	// Remove hop-by-hop headers to the backend.  Especially
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

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := outreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outreq.Header.Set("X-Forwarded-For", clientIP)
	}

	res, err := transport.RoundTrip(outreq)
	if err != nil {
		p.logf("http: proxy error: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	for _, h := range hopHeaders {
		res.Header.Del(h)
	}

	copyHeader(rw.Header(), res.Header)

	rw.WriteHeader(res.StatusCode)
	p.CopyResponse(rw, res.Body)
}

func (p *ReverseProxy) CopyResponse(dst io.Writer, src io.Reader) {
	if p.FlushInterval != 0 {
		if wf, ok := dst.(writeFlusher); ok {
			mlw := &maxLatencyWriter{
				dst:     wf,
				latency: p.FlushInterval,
				done:    make(chan bool),
			}
			go mlw.flushLoop()
			defer mlw.stop()
			dst = mlw
		}
	}

	io.Copy(dst, src)
}

func (p *ReverseProxy) logf(format string, args ...interface{}) {
	if p.ErrorLog != nil {
		p.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

type writeFlusher interface {
	io.Writer
	http.Flusher
}

type maxLatencyWriter struct {
	dst     writeFlusher
	latency time.Duration

	lk   sync.Mutex // protects Write + Flush
	done chan bool
}

func (m *maxLatencyWriter) Write(p []byte) (int, error) {
	m.lk.Lock()
	defer m.lk.Unlock()
	return m.dst.Write(p)
}

func (m *maxLatencyWriter) flushLoop() {
	t := time.NewTicker(m.latency)
	defer t.Stop()
	for {
		select {
		case <-m.done:
			if onExitFlushLoop != nil {
				onExitFlushLoop()
			}
			return
		case <-t.C:
			m.lk.Lock()
			m.dst.Flush()
			m.lk.Unlock()
		}
	}
}

func (m *maxLatencyWriter) stop() { m.done <- true }

func init() {
	RealDataProxy = &ProxyHandle{
		//proxy_host:       "10.10.178.19",
		proxy_host:       "123.59.84.71",
		proxy_port:       "8080",
		proxy_fixed_path: "/data/v1/",
	}
	RealNewsProxy = &ProxyHandle{
		//proxy_host:       "10.10.167.183",
		proxy_host:       "123.59.84.71",
		proxy_port:       "8080",
		proxy_fixed_path: "/news/v1/",
	}
	RealUserProxy = &ProxyHandle{
		//proxy_host:       "10.10.152.93",
		proxy_host:       "123.59.84.71",
		proxy_port:       "8080",
		proxy_fixed_path: "/user/v1/",
	}
	RealGroupProxy = &ProxyHandle{
		//proxy_host:       "10.10.152.93",
		proxy_host:       "123.59.84.71",
		proxy_port:       "8082",
		proxy_fixed_path: "/group/v1/",
	}
	RealMallProxy = &ProxyHandle{
		proxy_host:       "localhost",
		proxy_port:       "8099",
		proxy_fixed_path: "/mall/v1/",
	}
}
