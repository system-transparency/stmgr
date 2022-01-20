package keygen

import (
	"encoding/pem"
	"errors"
	"os"
)

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
	pemBytes := pem.EncodeToMemory(block)

	return os.WriteFile(path, pemBytes, 0o600)
}
