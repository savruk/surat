package surat

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/savruk/cacher"
)

var path string = "./conf.json"
var conf = LoadConfig(path)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	instance := cacher.Cacher{cacher.NewRedisEngine()}
	key := r.URL.Path
	item, _ := instance.Get(key)
	if string(item.Value) != "" {
		headers := GetHeadersFromCache(instance, key)
		SetHttpHeaders(w, headers)
		fileType := http.DetectContentType(item.Value)

		if strings.HasSuffix(key, ".js") || strings.HasSuffix(key, ".css") || strings.Contains(fileType, "image") {
			w.WriteHeader(http.StatusNotModified)
		} else {
			w.Write(item.Value)
		}
	} else {
		url := strings.Join([]string{fmt.Sprintf("%s://%s:%d", conf.Backend.Protocol, conf.Backend.Host, conf.Backend.Port), key}, "")
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			http.Error(w, resp.Status, resp.StatusCode)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			memerr := instance.Set(key, []byte(body))
			if memerr != nil {
				log.Println(memerr)
			}
			headers := ConvertHeaders(resp.Header)
			SetHeadersToCache(instance, key, headers)
			SetHttpHeaders(w, headers)
			w.Write(body)
		}
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	log.Panicln("Post is not supported")
}
