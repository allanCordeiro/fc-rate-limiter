package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpirationTimeByToken(t *testing.T) {

	tests := []struct {
		name           string
		key            string
		expectedExpire int
	}{
		{
			name:           "given a valid key when search its expire time should return properly",
			key:            "ccd42b6a-6e64-410d-9122-372d922858f2",
			expectedExpire: 10,
		},
		{
			name:           "given an invalid key when search its expire time should return 0",
			key:            "an-invalid-key",
			expectedExpire: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expire, err := GetTokenExpirationParam("test/tokens.json", test.key)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedExpire, expire)
		})
	}
}
