package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebDomain_Constants(t *testing.T) {
	tests := []struct {
		name     string
		domain   WebDomain
		expected string
	}{
		{
			name:     "VnExpress domain",
			domain:   WebDomainVnExpress,
			expected: "vnexpress.net",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.domain))
		})
	}
}

