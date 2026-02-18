package standard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	PaginationHeader = "x-pagination"
)

// PaginationResponse is the response contained in the x-pagination header of http response
type PaginationResponse struct {
	PageNumber int
	PageSize   int
	Total      int
	OrderBy    *string
}

// GetPaginationResponse extracts the pagination response from the JSON contained in the x-pagination header
func GetPaginationResponse(ctx context.Context, response *http.Response) (*PaginationResponse, error) {
	pagination := &PaginationResponse{}

	paginationHeader := response.Header.Get(PaginationHeader)
	if paginationHeader == "" {
		return nil, fmt.Errorf("pagination header not found in response")
	}

	err := json.Unmarshal([]byte(paginationHeader), pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal pagination header: %v", err)
	}

	return pagination, nil
}
