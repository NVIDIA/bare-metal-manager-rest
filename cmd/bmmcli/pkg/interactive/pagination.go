// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interactive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/nvidia/bare-metal-manager-rest/sdk/standard"
)

// paginationMeta is decoded from the X-Pagination response header.
type paginationMeta struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
}

func parsePaginationMeta(resp *http.Response) *paginationMeta {
	if resp == nil {
		return nil
	}
	raw := resp.Header.Get("X-Pagination")
	if raw == "" {
		return nil
	}
	var meta paginationMeta
	if err := json.Unmarshal([]byte(raw), &meta); err != nil {
		return nil
	}
	return &meta
}

// FetchPage is the signature expected by FetchAllPages.
type FetchPage[T any] func(pageNumber, pageSize int32) ([]T, *http.Response, error)

// FetchAllPages calls fn repeatedly until all pages have been retrieved.
//
// It recovers from the SDK's "undefined response type" error, which fires
// when the server returns a valid JSON body without a Content-Type header.
// In that case it falls back to JSON-decoding the raw body directly.
//
// Pass the session so that HTTP details (status, content-type, body) are
// logged when verbose mode is enabled.
func FetchAllPages[T any](s *Session, label string, fn FetchPage[T]) ([]T, *http.Response, error) {
	const pageSize = int32(100)
	var all []T
	var lastResp *http.Response
	page := int32(1)

	for {
		items, resp, err := fn(page, pageSize)
		if err != nil {
			if isUndefinedContentType(err) {
				logHTTPDetails(s, label, page, resp, err)

				// For HTTP error status codes, don't attempt JSON recovery —
				// return a clear human-readable error instead.
				if isHTTPError(resp) {
					return nil, resp, httpErrorFromResponse(label, page, resp, err)
				}

				decoded, decErr := recoverJSONBody[[]T](err, resp)
				if decErr == nil {
					items = decoded
					lastResp = resp
					all = append(all, items...)
					meta := parsePaginationMeta(resp)
					if meta != nil && meta.Total > 0 && len(all) >= meta.Total {
						break
					}
					if len(items) < int(pageSize) {
						break
					}
					page++
					continue
				}
				return nil, resp, enrichError(label, page, resp, err)
			}
			return nil, resp, enrichError(label, page, resp, err)
		}

		lastResp = resp
		all = append(all, items...)

		meta := parsePaginationMeta(resp)
		if meta != nil && meta.Total > 0 && len(all) >= meta.Total {
			break
		}
		if len(items) < int(pageSize) {
			break
		}
		page++
	}

	return all, lastResp, nil
}

// isUndefinedContentType returns true for the SDK's "undefined response type"
// error, which fires when the server returns a non-empty body without a
// recognisable Content-Type header.
func isUndefinedContentType(err error) bool {
	return err != nil && err.Error() == "undefined response type"
}

// isHTTPError returns true when the response status is a known HTTP error code
// that should not be treated as a recoverable decode failure.
func isHTTPError(resp *http.Response) bool {
	return resp != nil && resp.StatusCode >= 400
}

// httpErrorFromResponse returns a human-readable error for a non-2xx response,
// including the status code and a brief body excerpt.
func httpErrorFromResponse(label string, page int32, resp *http.Response, err error) error {
	body := ""
	if apiErr, ok := err.(standard.GenericOpenAPIError); ok {
		if b := apiErr.Body(); len(b) > 0 {
			preview := b
			if len(preview) > 200 {
				preview = preview[:200]
			}
			body = string(preview)
		}
	}

	switch resp.StatusCode {
	case 403:
		return fmt.Errorf("%s: permission denied (HTTP 403) — your org may not have the required role for this resource%s",
			label, bodyHint(body))
	case 401:
		return fmt.Errorf("%s: unauthorized (HTTP 401) — token may be expired; type 'login' to re-authenticate%s",
			label, bodyHint(body))
	case 400:
		return fmt.Errorf("%s: bad request (HTTP 400)%s", label, bodyHint(body))
	case 404:
		return fmt.Errorf("%s: not found (HTTP 404)%s", label, bodyHint(body))
	default:
		return fmt.Errorf("%s (page %d, HTTP %s): %w", label, page, resp.Status, err)
	}
}

func bodyHint(body string) string {
	if body == "" {
		return ""
	}
	return ": " + body
}

// recoverJSONBody tries to JSON-decode the raw body from a GenericOpenAPIError
// or from the buffered HTTP response body.
func recoverJSONBody[T any](err error, resp *http.Response) (T, error) {
	var zero T

	if apiErr, ok := err.(standard.GenericOpenAPIError); ok {
		if body := apiErr.Body(); len(body) > 0 {
			var result T
			if jsonErr := json.Unmarshal(body, &result); jsonErr == nil {
				return result, nil
			}
		}
	}

	if resp != nil && resp.Body != nil {
		body, readErr := io.ReadAll(resp.Body)
		if readErr == nil && len(body) > 0 {
			var result T
			if jsonErr := json.Unmarshal(body, &result); jsonErr == nil {
				return result, nil
			}
		}
	}

	return zero, fmt.Errorf("JSON recovery failed: %w", err)
}

// logHTTPDetails prints raw HTTP response info when session verbose mode is on.
func logHTTPDetails(s *Session, label string, page int32, resp *http.Response, err error) {
	if s == nil || !s.Verbose {
		return
	}

	fmt.Printf("\n[DEBUG] %s page=%d\n", label, page)

	if resp != nil {
		fmt.Printf("[DEBUG]   HTTP status:   %s\n", resp.Status)
		fmt.Printf("[DEBUG]   Content-Type:  %q\n", resp.Header.Get("Content-Type"))
		fmt.Printf("[DEBUG]   X-Pagination:  %q\n", resp.Header.Get("X-Pagination"))
	} else {
		fmt.Printf("[DEBUG]   response:      nil\n")
	}

	if apiErr, ok := err.(standard.GenericOpenAPIError); ok {
		body := apiErr.Body()
		if len(body) > 0 {
			preview := body
			if len(preview) > 512 {
				preview = preview[:512]
			}
			fmt.Printf("[DEBUG]   body (%d bytes): %s\n", len(body), string(preview))
		} else {
			fmt.Printf("[DEBUG]   body:          (empty)\n")
		}
	}
	fmt.Println()
}

// enrichError wraps an SDK error with context about which fetcher and page
// triggered it, plus Content-Type and status from the response.
func enrichError(label string, page int32, resp *http.Response, err error) error {
	if resp == nil {
		return fmt.Errorf("%s (page %d): %w", label, page, err)
	}

	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "(none)"
	}

	bodyPreview := ""
	if apiErr, ok := err.(standard.GenericOpenAPIError); ok {
		if b := apiErr.Body(); len(b) > 0 {
			preview := b
			if len(preview) > 256 {
				preview = preview[:256]
			}
			bodyPreview = fmt.Sprintf(", body: %s", string(preview))
		}
	}

	return fmt.Errorf("%s (page %d, HTTP %s, content-type: %s%s): %w",
		label, page, resp.Status, ct, bodyPreview, err)
}

// orgScopedPathPattern matches the path segment immediately after /v{n}/org/{org}/
// so it can be rewritten to the configured api.name value.
var orgScopedPathPattern = regexp.MustCompile(`(/v[0-9]+/org/[^/]+/)[^/]+`)

// APINameTransport rewrites the hardcoded path segment (e.g. "carbide") in SDK
// URLs to the value from api.name in the config (e.g. "forge").
// This is necessary because the generated SDK has the path segment baked in from
// the spec, but different environments may use a different value.
type APINameTransport struct {
	APIName string
	Wrapped http.RoundTripper
}

func (t *APINameTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	next := t.Wrapped
	if next == nil {
		next = http.DefaultTransport
	}

	name := strings.Trim(strings.TrimSpace(t.APIName), "/")
	if name == "" || req == nil || req.URL == nil {
		return next.RoundTrip(req)
	}

	rewritten := orgScopedPathPattern.ReplaceAllString(req.URL.Path, "${1}"+name)
	rewrittenRaw := ""
	if req.URL.RawPath != "" {
		rewrittenRaw = orgScopedPathPattern.ReplaceAllString(req.URL.RawPath, "${1}"+name)
	}

	if rewritten == req.URL.Path && rewrittenRaw == req.URL.RawPath {
		return next.RoundTrip(req)
	}

	cloned := req.Clone(req.Context())
	clonedURL := *req.URL
	cloned.URL = &clonedURL
	cloned.URL.Path = rewritten
	if req.URL.RawPath != "" {
		cloned.URL.RawPath = rewrittenRaw
	}
	cloned.RequestURI = ""
	return next.RoundTrip(cloned)
}

// NewHTTPClient builds an *http.Client with the api.name rewrite transport
// and, optionally, the debug curl-logging transport layered on top.
// Transport order (outer → inner):
//
//	APINameTransport → DebugTransport → http.DefaultTransport
//
// The debug transport logs after the path has been rewritten so the curl
// command shows the actual URL sent to the server.
func NewHTTPClient(apiName string, verbose bool) *http.Client {
	var inner http.RoundTripper = http.DefaultTransport
	if verbose {
		inner = &DebugTransport{Wrapped: inner}
	}
	return &http.Client{
		Transport: &APINameTransport{
			APIName: apiName,
			Wrapped: inner,
		},
	}
}

// DebugTransport wraps an http.RoundTripper and logs each outbound request
// as a curl command before executing it.
type DebugTransport struct {
	Wrapped http.RoundTripper
}

func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Print(requestAsCurl(req))
	return t.Wrapped.RoundTrip(req)
}

// requestAsCurl renders an *http.Request as a single-line curl command.
func requestAsCurl(req *http.Request) string {
	var b strings.Builder
	b.WriteString("\n[DEBUG] curl -sS")
	b.WriteString(" -X ")
	b.WriteString(req.Method)

	for key, vals := range req.Header {
		for _, v := range vals {
			if strings.EqualFold(key, "Authorization") {
				// Truncate token so the log is readable but not a full secret.
				if len(v) > 30 {
					v = v[:30] + "...<truncated>"
				}
			}
			fmt.Fprintf(&b, " -H %q", key+": "+v)
		}
	}

	fmt.Fprintf(&b, " %q\n", req.URL.String())
	return b.String()
}

// PrintSummary writes a brief item count line to w.
func PrintSummary(w io.Writer, resp *http.Response, count int) {
	meta := parsePaginationMeta(resp)
	if meta != nil && meta.Total > 0 {
		fmt.Fprintf(w, "%d items\n", meta.Total)
	} else {
		fmt.Fprintf(w, "%d items\n", count)
	}
}
