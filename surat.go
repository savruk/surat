package surat

import (
	"net/http"

	"github.com/savruk/cacher"
)

func redishandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		GetHandler(w, r)
	case r.Method == "POST":
		PostHandler(w, r)
	}
}

func Flush() {
	instance := cacher.Cacher{cacher.NewRedisEngine()}
	instance.Flush()
}

func Run() {
	// http.HandleFunc("/", handler)
	http.HandleFunc("/", makeGzipHandler(redishandler))
	http.ListenAndServe(":8080", nil)
}
