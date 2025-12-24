package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MessageType_String(t *testing.T) {
	tests := []struct {
		name string
		t    MessageType
		want string
	}{
		{
			name: "WebScrapper",
			t:    MessageTypeWebScrapper,
			want: "web_scrapper",
		},
		{
			name: "Unknown",
			t:    MessageType(999),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.t.String())
		})
	}
}
