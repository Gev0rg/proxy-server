package proxy

import (
	"github.com/Gev0rg/proxy-server/storage"
	"io"
	"net"
	"net/http"
)

type Proxy struct {
	Store *storage.Storage
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.HandleHttps(w, r)
	} else {
		p.HandleHttp(w, r)
	}
}

func (p *Proxy) HandleHttp(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	r.RequestURI = ""

	p.Store.SaveRequest(r)

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	p.Store.SaveResponse(resp)

	w.WriteHeader(resp.StatusCode)
	copyHeader(w.Header(), resp.Header)

	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func (p *Proxy) HandleHttps(w http.ResponseWriter, r *http.Request) {
	p.Store.SaveRequest(r)

	dest, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	client, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(dest, client)
	go transfer(client, dest)
}

func transfer(dest io.WriteCloser, src io.ReadCloser) {
	defer dest.Close()
	defer src.Close()

	io.Copy(dest, src)
}

func copyHeader(dst, src http.Header) {
	for key, value := range src {
		for _, v := range value {
			dst.Add(key, v)
		}
	}
}
