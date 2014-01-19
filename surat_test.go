package surat

import (
	"fmt"
	"testing"
)

func Test_Something(t *testing.T) {
	config := Config{Backends: []Backend{Backend{"127.0.0.1", "8000"}}}
	fmt.Printf("%v", config)
}
