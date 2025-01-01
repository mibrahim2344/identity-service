package password

import (
	"crypto/rand"
	"fmt"
	"io"
)

// EntropyProvider defines the interface for providing entropy
type EntropyProvider interface {
	io.Reader
}

// CryptoEntropyProvider implements EntropyProvider using crypto/rand
type CryptoEntropyProvider struct{}

// Read implements io.Reader using crypto/rand
func (p *CryptoEntropyProvider) Read(b []byte) (n int, err error) {
	n, err = rand.Read(b)
	if err != nil {
		return 0, fmt.Errorf("failed to read random bytes: %w", err)
	}
	return n, nil
}

// MockEntropyProvider implements EntropyProvider for testing
type MockEntropyProvider struct {
	Data []byte
	Pos  int
}

// Read implements io.Reader for testing
func (p *MockEntropyProvider) Read(b []byte) (n int, err error) {
	if p.Pos >= len(p.Data) {
		return 0, io.EOF
	}
	n = copy(b, p.Data[p.Pos:])
	p.Pos += n
	return n, nil
}
