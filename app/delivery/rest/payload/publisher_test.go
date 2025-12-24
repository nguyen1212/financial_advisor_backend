package payload

import (
	"testing"

	"github.com/financial_advisor/app/errors"
	"github.com/stretchr/testify/require"
)

func TestPublisherCreateRequest_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload PublisherCreateRequest
		wantErr error
	}{
		{
			name: "valid payload",
			payload: PublisherCreateRequest{
				Name:        "Financial Times",
				Domain:      "https://www.ft.com",
				Description: "Leading financial news publication",
			},
		},
		{
			name: "valid payload without description",
			payload: PublisherCreateRequest{
				Name:   "Reuters",
				Domain: "https://www.reuters.com",
			},
		},
		{
			name:    "violate required tag",
			payload: PublisherCreateRequest{},
			wantErr: errors.SystemErrors{
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for Name with tag required: ",
				),
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for Domain with tag required: ",
				),
			},
		},
		{
			name: "violate url tag",
			payload: PublisherCreateRequest{
				Name:   "Invalid Publisher",
				Domain: "invalid-url",
			},
			wantErr: errors.SystemErrors{
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for Domain with tag url: invalid-url",
				),
			},
		},
		{
			name: "missing name only",
			payload: PublisherCreateRequest{
				Domain:      "https://example.com",
				Description: "Test description",
			},
			wantErr: errors.SystemErrors{
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for Name with tag required: ",
				),
			},
		},
		{
			name: "missing domain only",
			payload: PublisherCreateRequest{
				Name:        "Test Publisher",
				Description: "Test description",
			},
			wantErr: errors.SystemErrors{
				errors.NewErrorBadRequest(
					errors.ErrorCodeBadRequest,
					"Invalid value for Domain with tag required: ",
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

func TestPublisherCreateRequest_ToDTO(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload PublisherCreateRequest
		want    struct {
			Name        string
			Domain      string
			Description string
		}
	}{
		{
			name: "full payload conversion",
			payload: PublisherCreateRequest{
				Name:        "Bloomberg",
				Domain:      "https://www.bloomberg.com",
				Description: "Business and financial news",
			},
			want: struct {
				Name        string
				Domain      string
				Description string
			}{
				Name:        "Bloomberg",
				Domain:      "https://www.bloomberg.com",
				Description: "Business and financial news",
			},
		},
		{
			name: "payload without description",
			payload: PublisherCreateRequest{
				Name:   "Wall Street Journal",
				Domain: "https://www.wsj.com",
			},
			want: struct {
				Name        string
				Domain      string
				Description string
			}{
				Name:        "Wall Street Journal",
				Domain:      "https://www.wsj.com",
				Description: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dto := tt.payload.ToDTO()

			require.Equal(t, tt.want.Name, dto.Name)
			require.Equal(t, tt.want.Domain, dto.Domain)
			require.Equal(t, tt.want.Description, dto.Description)
		})
	}
}

func TestPublishersFindRequest_ToDTO(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload PublishersFindRequest
		want    struct {
			Page int
			Size int
		}
	}{
		{
			name: "with custom paging",
			payload: PublishersFindRequest{
				Paging: Paging{
					Page: 2,
					Size: 50,
				},
			},
			want: struct {
				Page int
				Size int
			}{
				Page: 2,
				Size: 50,
			},
		},
		{
			name: "with zero values - should use defaults",
			payload: PublishersFindRequest{
				Paging: Paging{
					Page: 0,
					Size: 0,
				},
			},
			want: struct {
				Page int
				Size int
			}{
				Page: 1,
				Size: 30,
			},
		},
		{
			name: "with only page set",
			payload: PublishersFindRequest{
				Paging: Paging{
					Page: 3,
					Size: 0,
				},
			},
			want: struct {
				Page int
				Size int
			}{
				Page: 3,
				Size: 30,
			},
		},
		{
			name: "with only size set",
			payload: PublishersFindRequest{
				Paging: Paging{
					Page: 0,
					Size: 100,
				},
			},
			want: struct {
				Page int
				Size int
			}{
				Page: 1,
				Size: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dto := tt.payload.ToDTO()

			require.Equal(t, tt.want.Page, dto.Paging.Page)
			require.Equal(t, tt.want.Size, dto.Paging.Size)
		})
	}
}