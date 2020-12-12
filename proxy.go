package main

import (
	"bytes"
	"fmt"
	"github.com/andybalholm/brotli"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

//将request转发给 http://localhost:2003
func helloHandler(w http.ResponseWriter, r *http.Request) {

	trueServer := "https://www.notion.so"

	url, err := url.Parse(trueServer)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(r.Host)

	proxy := NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)

}

type A struct {
}

func (a A) RoundTrip(request *http.Request) (*http.Response, error) {
	fmt.Println(request.URL.Scheme, request.URL.Host, request.URL.Path)

	resp, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("response body read", err)
	}
	//fmt.Println(string(body))

	isBr := true
	r := brotli.NewReader(bytes.NewReader(body))
	unbr, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("un br read", err)
		unbr = body
		isBr = false
	}

	nb := strings.Replace(string(unbr), "//www.notion.so", "//localhost:8888", -1)

	//fuckJs := `didMount(){const{device:e}=this.environment;"external"===Object(m.l)({url:window.location.href,isMobile:e.isMobile,baseUrl:g.default.baseURL,protocol:g.default.protocol,currentUrl:window.location.href}).name&&Ve.showDialog({message:"Mismatch between origin and baseUrl (dev).",showCancel:!1,keepFocus:!1,items:[{label:n.createElement(o.FormattedMessage,{defaultMessage:"Okay",id:"notionAppContainer.dialog.mismatchedOriginURL.okayButton.label"}),onAccept:()=>{const e=xi.j({relativeUrl:xi.f(window.location.href),baseUrl:g.default.baseURL}),t=xi.e(e);t.protocol=xi.e(window.location.href).protocol,window.location.href=xi.b(t)}}]})}`
	//fuckJs = strings.Replace(fuckJs, " ", "", -1)
	//fuckJs = strings.Replace(fuckJs, "\n", "", -1)
	//fmt.Println(fuckJs)
	//nb = strings.Replace(nb, fuckJs, "", -1)
	nb = strings.Replace(nb, "//notion.so", "//localhost:8888", -1)
	nb = strings.Replace(nb, "/www.notion.so", "/localhost:8888", -1)
	//nb = strings.Replace(nb, "window.location.href", `"http://localhost:8888/xxxxxx"`, -1)
	//nb = strings.Replace(nb, `"http://localhost:8888/xxxxxx"=`, `"http://localhost:8888/xxxxxx"`, -1)

	//fmt.Println(nb)

	buff := &bytes.Buffer{}
	brWriter := brotli.NewWriter(buff)
	n, err := brWriter.Write([]byte(nb))
	brWriter.Flush()
	brWriter.Close()
	if err != nil {
		fmt.Println("brWriter", err, n)
	}
	fmt.Println("write", len(nb), n, len(buff.Bytes()))

	if isBr {
		resp.Body = ioutil.NopCloser(bytes.NewReader(buff.Bytes()))
	} else {
		resp.Body = ioutil.NopCloser(strings.NewReader(nb))
	}
	resp.Header.Del("X-Frame-Options")
	return resp, err
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		//if req.URL.Scheme == "https" {
		//	req.URL.Scheme = "http"
		//}
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		orgin := strings.Replace(req.Header.Get("Origin"), "localhost:8888", "notion.so", -1)
		req.Header.Set("Origin", orgin)
		referer := strings.Replace(req.Header.Get("Referer"), "localhost:8888", "notion.so", -1)
		req.Header.Set("Referer", referer)

		fmt.Println(req.URL.Host, req.RequestURI)
	}

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: &A{},
	}
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
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
