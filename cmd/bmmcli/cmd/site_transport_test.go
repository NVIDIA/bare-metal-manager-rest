// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type captureRoundTripper struct {
	path    string
	rawPath string
}

func (c *captureRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	c.path = req.URL.Path
	c.rawPath = req.URL.RawPath
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func TestRewriteAPINamePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		apiName string
		want    string
	}{
		{
			name:    "rewrites org scoped v2 path",
			path:    "/v2/org/test-org/forge/site",
			apiName: "carbide",
			want:    "/v2/org/test-org/carbide/site",
		},
		{
			name:    "leaves non org path unchanged",
			path:    "/healthz",
			apiName: "carbide",
			want:    "/healthz",
		},
		{
			name:    "leaves path unchanged when api name empty",
			path:    "/v2/org/test-org/forge/site",
			apiName: "",
			want:    "/v2/org/test-org/forge/site",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rewriteAPINamePath(tt.path, tt.apiName)
			if got != tt.want {
				t.Fatalf("rewriteAPINamePath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPINameRewriteTransportRoundTrip(t *testing.T) {
	capture := &captureRoundTripper{}
	transport := &apiNameRewriteTransport{
		apiName: "carbide",
		next:    capture,
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8388/v2/org/test-org/forge/site", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	_, err = transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip: %v", err)
	}

	if capture.path != "/v2/org/test-org/carbide/site" {
		t.Fatalf("rewritten path = %q, want %q", capture.path, "/v2/org/test-org/carbide/site")
	}
}
