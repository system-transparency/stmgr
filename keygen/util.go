package keygen

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	sigsumCrypto "sigsum.org/sigsum-go/pkg/crypto"
	"sigsum.org/sigsum-go/pkg/key"
)

var (
	ErrNoPEMBlock = errors.New("no PEM block found")
	ErrTrailing   = errors.New("unexpected trailing data after PEM block")
)

const defaultFilePerm fs.FileMode = 0o600

// Sigsum uses its own (simpler) Signer interface. Wrap it in a type
// implementing ths stdlib crypto.Signer interface.
type sigsumSigner struct {
	sss sigsumCrypto.Signer
}

func (s sigsumSigner) Sign(_ io.Reader, msg []byte, _ crypto.SignerOpts) ([]byte, error) {
	sig, err := s.sss.Sign(msg)
	if err != nil {
		return nil, err
	}
	return sig[:], nil
}

func (s sigsumSigner) Public() crypto.PublicKey {
	pub := s.sss.Public()
	return ed25519.PublicKey(pub[:])
}

// Loads a private key file, either x509 style, or an OpenSSH public
// key file where private key is accessed using ssh-agent.
func LoadPrivateKey(fileName string) (crypto.Signer, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if bytes.HasPrefix(data, []byte("ssh-ed25519 ")) {
		// Attempt decoding as OpenSSH pubkey. Sigsum's key
		// package supports reading OpenSSH private keys, as
		// well as OpenSSH public keys where the corresponding
		// private key is available via ssh-agent.
		signer, err := key.ParsePrivateKey(string(data))
		if err != nil {
			return nil, err
		}
		return sigsumSigner{signer}, nil
	}
	block, rest := pem.Decode(data)
	if block == nil {
		return nil, ErrNoPEMBlock
	}
	if len(rest) != 0 {
		return nil, ErrTrailing
	}

	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key file, PEM type: %q", block.Type)
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	signer, ok := priv.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("invalid private key type: %T", priv)
	}
	return signer, nil
}

// Loads a PEM coded x509 certificate, without decoding the DER blob.
func LoadCertBytes(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, rest := pem.Decode(bytes)
	if block == nil {
		return nil, ErrNoPEMBlock
	}

	if len(rest) != 0 {
		return nil, ErrTrailing
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid cert file, got type %q", block.Type)
	}

	return block.Bytes, nil
}

// WritePEM writes the pem.Block data to a PEM formatted
// file to the specified path.
func WritePEM(block *pem.Block, path string) error {
	if block == nil {
		return ErrNoPEMBlock
	}

	return os.WriteFile(path, pem.EncodeToMemory(block), defaultFilePerm)
}
