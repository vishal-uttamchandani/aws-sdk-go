package credentials

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type stubProvider struct {
	creds   Value
	expired bool
	err     error
}

func (s *stubProvider) Retrieve() (Value, error) {
	s.expired = false
	return s.creds, s.err
}
func (s *stubProvider) IsExpired() bool {
	return s.expired
}

func TestCredentialsGet(t *testing.T) {
	c := NewCredentials(&stubProvider{
		creds: Value{
			AccessKeyID:     "AKID",
			SecretAccessKey: "SECRET",
			SessionToken:    "",
		},
		expired: true,
	})

	creds, err := c.Get()
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, "AKID", creds.AccessKeyID, "Expect access key ID to match")
	assert.Equal(t, "SECRET", creds.SecretAccessKey, "Expect secret access key to match")
	assert.Empty(t, creds.SessionToken, "Expect session token to be empty")
}

func TestCredentialsGetWithError(t *testing.T) {
	c := NewCredentials(&stubProvider{err: errors.New("provider error"), expired: true})

	_, err := c.Get()
	assert.Equal(t, "provider error", err.Error(), "Expected provider error")
}

func TestCredentialsExpire(t *testing.T) {
	stub := &stubProvider{}
	c := NewCredentials(stub)

	assert.False(t, c.IsExpired(), "Expected not to be expired")
	c.Expire()
	assert.True(t, c.IsExpired(), "Expected to be expired")

	c.forceRefresh = false
	assert.False(t, c.IsExpired(), "Expected not to be expired")

	stub.expired = true
	assert.True(t, c.IsExpired(), "Expected to be expired")
}
