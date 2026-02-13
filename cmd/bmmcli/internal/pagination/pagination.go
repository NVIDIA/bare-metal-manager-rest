// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pagination

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// MaxPageSize matches the OpenAPI maximum pageSize constraint for list endpoints.
	MaxPageSize int32 = 100

	paginationHeaderName = "X-Pagination"
)

type header struct {
	PageNumber int32 `json:"pageNumber"`
	PageSize   int32 `json:"pageSize"`
	Total      int32 `json:"total"`
}

// Summary describes pagination details for user-facing list output.
type Summary struct {
	FetchedCount int
	PageSize     int32
	TotalCount   int32
	PageCount    int32
	HeaderFound  bool
	PageLimitHit bool
}

// FetchAllPages requests all pages for list APIs and returns combined results.
// It requests MaxPageSize on each call and keeps paging until all items are fetched.
func FetchAllPages[T any](fetch func(pageNumber, pageSize int32) ([]T, *http.Response, error)) ([]T, *http.Response, error) {
	pageNumber := int32(1)
	pageSize := MaxPageSize
	all := make([]T, 0)
	var lastResp *http.Response

	for {
		pageItems, resp, err := fetch(pageNumber, pageSize)
		lastResp = resp
		if err != nil {
			return nil, lastResp, err
		}

		all = append(all, pageItems...)

		if done, nextPage, ok := nextPageFromHeader(resp, int32(len(all)), pageNumber, pageSize); ok {
			if done {
				return all, lastResp, nil
			}
			pageNumber = nextPage
			continue
		}

		// Fallback when no pagination header is present.
		if len(pageItems) < int(pageSize) {
			return all, lastResp, nil
		}
		pageNumber++
	}
}

func nextPageFromHeader(resp *http.Response, totalFetched, currentPage, currentPageSize int32) (done bool, nextPage int32, ok bool) {
	if resp == nil {
		return false, 0, false
	}

	raw := resp.Header.Get(paginationHeaderName)
	if raw == "" {
		return false, 0, false
	}

	var h header
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		return false, 0, false
	}

	if h.Total <= 0 || totalFetched >= h.Total {
		return true, 0, true
	}

	page := h.PageNumber
	if page <= 0 {
		page = currentPage
	}
	size := h.PageSize
	if size <= 0 {
		size = currentPageSize
	}

	if page*size >= h.Total {
		return true, 0, true
	}

	next := page + 1
	if next <= currentPage {
		next = currentPage + 1
	}
	return false, next, true
}

// BuildSummary builds pagination summary metadata from the final response header.
func BuildSummary(resp *http.Response, fetchedCount int) Summary {
	s := Summary{
		FetchedCount: fetchedCount,
		PageSize:     MaxPageSize,
		TotalCount:   int32(fetchedCount),
	}
	if s.TotalCount < 0 {
		s.TotalCount = 0
	}

	if h, ok := parseHeader(resp); ok {
		s.HeaderFound = true
		if h.PageSize > 0 {
			s.PageSize = h.PageSize
		}
		if h.Total > 0 {
			s.TotalCount = h.Total
		}
	}

	s.PageCount = pageCountForTotal(s.TotalCount, s.PageSize)
	s.PageLimitHit = s.PageCount > 1
	return s
}

// PrintSummary writes pagination details for list output.
func PrintSummary(w io.Writer, resp *http.Response, fetchedCount int) {
	if w == nil {
		return
	}
	s := BuildSummary(resp, fetchedCount)
	if s.PageLimitHit {
		fmt.Fprintf(w, "INFO: pageSize limit hit; fetched all %d items across %d pages (pageSize=%d, total=%d)\n", s.FetchedCount, s.PageCount, s.PageSize, s.TotalCount)
		return
	}
	fmt.Fprintf(w, "INFO: fetched %d items in 1 page (pageSize=%d, total=%d)\n", s.FetchedCount, s.PageSize, s.TotalCount)
}

func parseHeader(resp *http.Response) (header, bool) {
	if resp == nil {
		return header{}, false
	}

	raw := resp.Header.Get(paginationHeaderName)
	if raw == "" {
		return header{}, false
	}

	var h header
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		return header{}, false
	}
	return h, true
}

func pageCountForTotal(total, pageSize int32) int32 {
	if pageSize <= 0 {
		pageSize = MaxPageSize
	}
	if total <= 0 {
		return 1
	}
	pages := total / pageSize
	if total%pageSize != 0 {
		pages++
	}
	if pages <= 0 {
		return 1
	}
	return pages
}
