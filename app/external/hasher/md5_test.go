package hasher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMD5(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &md5Hashing{}, NewMD5())
}

func Test_md5Hashing_Hash(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		"d0adc7d0d152c10454693682f663c36a",
		fmt.Sprintf("%x", (&md5Hashing{}).Hash("heelo")),
	)
}
