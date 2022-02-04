package keygen

import (
	"encoding/pem"
	"errors"
	"io/fs"
	"os"
)

var (
	ErrNoPEMBlock = errors.New("no PEM block found")
	ErrTrailing   = errors.New("unexpected trailing data after PEM block")
)

const defaultFilePerm fs.FileMode = 0o600

// LoadPEM reads a PEM formatted file from disk
// and returns a pointer to the pem.Block data.
func LoadPEM(path string) (*pem.Block, error) {
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

	return block, nil
}

// WritePEM writes the pem.Block data to a PEM formatted
// file to the specified path.
func WritePEM(block *pem.Block, path string) error {
	if block == nil {
		return ErrNoPEMBlock
	}

	return os.WriteFile(path, pem.EncodeToMemory(block), defaultFilePerm)
}
