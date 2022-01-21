package keygen

import (
	"encoding/pem"
	"errors"
	"io/fs"
	"os"
)

const defaultFilePerm fs.FileMode = 0o600

func LoadPEM(path string) (*pem.Block, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, rest := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}

	if len(rest) != 0 {
		return nil, errors.New("unexpected trailing data after PEM block")
	}

	return block, nil
}

func WritePEM(block *pem.Block, path string) error {
	return os.WriteFile(path, pem.EncodeToMemory(block), defaultFilePerm)
}
