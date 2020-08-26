package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

func handleRequest(w http.ResponseWriter, domain string, path string) {
	log.Printf("REQ: %s %s", domain, path)
	item, err := GetItem(domain)
	if err != nil {
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<h1>Site Not Found</h1>")
	}

	req, err := http.NewRequest(
		"GET",
		path,
		http.NoBody)

	if err != nil {
		log.Fatalf("Invalid URL")
	}

	proxy := PagesReverseProxy(item)
	proxy.ServeHTTP(w, req)
}

func guessContentType(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return "text/html"
	}
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

func PagesReverseProxy(item *Item) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		path := req.URL.Path
		if strings.HasSuffix(path, "/") {
			path += "index.html"
		}
		path = path[1:] // strip leading /

		req.URL.Scheme = "https"
		req.URL.Host = item.D1

		req.URL.Path = "/view_archive.php"
		req.URL.RawPath = req.URL.Path

		q := url.Values{}
		q.Set("archive", item.Dir+"/"+item.Root)
		q.Set("file", path)
		req.URL.RawQuery = q.Encode()

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	modifyResponse := func(res *http.Response) error {
		// delete the "Content-Disposition" header as that interfers with opening pages in browser
		delete(res.Header, "Content-Disposition")

		// Set the content-type using the extension of the requested resource
		contentType := guessContentType(res.Request.URL.Query().Get("file"))
		res.Header.Set("Content-Type", contentType)
		return nil
	}

	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
	}
}

func getItem(domain string) (int, error) {
	// return 0, errors.New("item not found")
	return 0, nil
}

func getSubdomain(host string) string {
	host = strings.Split(host, ":")[0] // strip port
	return strings.Split(host, ".")[0]
}

// Serves all the archive pages
func Serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, getSubdomain(r.Host), r.URL.Path)
	})

	http.ListenAndServe(":8080", nil)
}
