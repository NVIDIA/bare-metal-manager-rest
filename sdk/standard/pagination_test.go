package standard

import (
	"context"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPaginationResponse(t *testing.T) {
	canonicalPaginationHeader := textproto.CanonicalMIMEHeaderKey(PaginationHeader)
	response := &http.Response{
		Header: http.Header{canonicalPaginationHeader: []string{"{\"pageNumber\":1,\"pageSize\":100,\"total\":46,\"orderBy\":null}"}},
	}

	pagination, err := GetPaginationResponse(context.Background(), response)
	require.NoError(t, err)

	assert.Equal(t, 1, pagination.PageNumber)
	assert.Equal(t, 100, pagination.PageSize)
	assert.Equal(t, 46, pagination.Total)
	assert.Nil(t, pagination.OrderBy)
}
