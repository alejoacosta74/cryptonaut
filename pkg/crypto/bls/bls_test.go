package bls

import (
	"testing"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignAndVerifyMessage(t *testing.T) {
	bls.Init(bls.BLS12_381)

	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{name: "Test 1", message: "Hello, world!", want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sec bls.SecretKey
			sec.SetByCSPRNG()

			signature, err := SignMessage(&sec, test.message)
			require.NoError(t, err)

			pub := sec.GetPublicKey()

			isValid, err := VerifyMessage(pub, test.message, signature)
			require.NoError(t, err)
			assert.Equal(t, test.want, isValid)
		})
	}
}
