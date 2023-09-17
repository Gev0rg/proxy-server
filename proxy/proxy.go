package proxy

import (
	"io"
	"net/http"
)

type Proxy struct {

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

	// TODO: save request
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	// TODO: save response

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
	
}

func copyHeader(dst, src http.Header) {
	for key, value := range src {
		for _, v := range value {
			dst.Add(key, v)
		}
	}
}
