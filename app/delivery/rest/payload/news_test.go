package payload

import (
	"testing"

	"github.com/financial_advisor/app/errors"
	"github.com/stretchr/testify/require"
)

func TestNewsCreateRequest_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload NewsCreateRequest
		wantErr error
	}{
		{
			name: "valid payload",
			payload: NewsCreateRequest{
				URL:      "https://example.com/news/sample-news-title",
				Category: "finance",
			},
		},
		{
			name:    "violate required tag",
			payload: NewsCreateRequest{},
			wantErr: errors.SystemErrors{
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for URL with tag required: ",
				),
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for Category with tag required: ",
				),
			},
		},
		{
			name: "violate url tag",
			payload: NewsCreateRequest{
				URL:      "invalid-url",
				Category: "finance",
			},
			wantErr: errors.SystemErrors{
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for URL with tag url: invalid-url",
				),
			},
		},
		{
			name: "violate oneof tag",
			payload: NewsCreateRequest{
				URL:      "https://example.com/news/sample-news-title",
				Category: "invalid category",
			},
			wantErr: errors.SystemErrors{
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for Category with tag oneof: invalid category",
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.payload.Validate()

			require.Equal(t, tt.wantErr, err)
		})
	}
}
