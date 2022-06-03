package main

import (
	"net/http"
	"strings"
)

func correct_headers(resp *http.Response, prefix string) {
	location := resp.Header.Get("Location")
	if location != "" {
		new_location := strings.TrimSuffix(prefix, "/") + location
		resp.Header.Set("Location", new_location)
	}

	referer := resp.Header.Get("Referer")
	if referer != "" {
		new_referer := strings.Replace(referer, "/hue/", "/", -1)
		resp.Header.Set("Referer", new_referer)
	}
}
