package surat

import (
	"log"
	"net/http"
	"strings"

	"github.com/savruk/cacher"
)

const (
	headerSeperator   = "____"
	keyValueSeperator = "::::"
)

func ConvertHeaders(hh http.Header) (headers string) {
	for key, val := range hh {
		headers = strings.Join([]string{strings.Join([]string{key, val[0]}, keyValueSeperator), headers}, headerSeperator)
	}
	return
}

func SetHeadersToCache(c cacher.Cacher, key string, headers string) {
	err := c.Set(strings.Join([]string{key, "headers"}, "@@"), []byte(headers))
	if err != nil {
		log.Panicln(err)
	}
}

func GetHeadersFromCache(c cacher.Cacher, key string) (headers string) {
	item, _ := c.Get(strings.Join([]string{key, "headers"}, "@@"))
	return string(item.Value)
}

func SetHttpHeaders(w http.ResponseWriter, headers string) {
	hh := strings.Split(headers, headerSeperator)
	for _, header := range hh {
		keyValuePair := strings.Split(header, keyValueSeperator)
		if len(keyValuePair) > 1 {
			w.Header().Set(keyValuePair[0], keyValuePair[1])
		}
	}
}
