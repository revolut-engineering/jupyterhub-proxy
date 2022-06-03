package main

import (
	"net/http"
	"testing"
)

func TestRequestHeader(t *testing.T) {
	header := make(map[string][]string)
	header["Location"] = []string{"/hue"}
	header["Referer"] = []string{"https://jove.com/hue/user/user@comp/hue"}

	prefix := "/user/user@comp/"
	req := http.Request{
		Proto:            "HTTP/1.0",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Header:           header,
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Trailer:          nil,
		TLS:              nil,
	}

	correct_headers(&req.Header, prefix)

	location := req.Header.Get("Location")
	want_location := "/user/user@comp/hue"
	if location != want_location {
		t.Fatalf("correct_headers() = %s want %s", location, want_location)
	}

	referer := req.Header.Get("Referer")
	want_referer := "https://jove.com/user/user@comp/hue"
	if referer != want_referer {
		t.Fatalf("correct_headers() = %s want %s", referer, want_referer)
	}
}

func TestResponseHeader(t *testing.T) {
	header := make(map[string][]string)
	header["Location"] = []string{"/hue"}
	header["Referer"] = []string{"https://jove.com/hue/user/user@comp/hue"}

	prefix := "/user/user@comp/"
	resp := http.Response{
		Status:           "200 OK",
		StatusCode:       200,
		Proto:            "HTTP/1.0",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Header:           header,
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}

	correct_headers(&resp.Header, prefix)

	location := resp.Header.Get("Location")
	want_location := "/user/user@comp/hue"
	if location != want_location {
		t.Fatalf("correct_headers() = %s want %s", location, want_location)
	}

	referer := resp.Header.Get("Referer")
	want_referer := "https://jove.com/user/user@comp/hue"
	if referer != want_referer {
		t.Fatalf("correct_headers() = %s want %s", referer, want_referer)
	}
}
