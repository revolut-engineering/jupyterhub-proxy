package main

import (
	"net/http"
	"strings"
)

func correct_headers(header *http.Header, prefix string) {
	location := header.Get("Location")
	if location != "" {
		new_location := strings.TrimSuffix(prefix, "/") + location
		header.Set("Location", new_location)
	}

	referer := header.Get("Referer")
	if referer != "" {
		new_referer := strings.Replace(referer, "/hue/", "/", -1)
		header.Set("Referer", new_referer)
	}
}
