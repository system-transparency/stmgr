package ospkg

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ospkgs "system-transparency.org/stboot/ospkg"
	"system-transparency.org/stmgr/keygen"
)

var ErrInvalidSuffix = errors.New("invalid file extension")

const DefaultOutName = "system-transparency-os-package"

// Sign will sign an OS package using the provided path
// of the private ed25519 key and corresponding certificate.
func Sign(keyPath, certPath, pkgPath string) error {
	pkgPath, err := parsePkgPath(pkgPath)
	if err != nil {
		return err
	}

	archive, err := os.ReadFile(pkgPath + ospkgs.OSPackageExt)
	if err != nil {
		return err
	}

	descriptor, err := os.ReadFile(pkgPath + ospkgs.DescriptorExt)
	if err != nil {
		return err
	}

	osp, err := ospkgs.NewOSPackage(archive, descriptor)
	if err != nil {
		return err
	}

	key, err := keygen.LoadPEM(keyPath)
	if err != nil {
		return err
	}
	if key.Type != "PRIVATE KEY" {
		return fmt.Errorf("invalid key file, got type %q", key.Type)
	}
	priv, err := x509.ParsePKCS8PrivateKey(key.Bytes)
	if err != nil {
		return err
	}
	signer, ok := priv.(crypto.Signer)
	if !ok {
		return fmt.Errorf("invalid private key type: %T", priv)
	}

	cert, err := keygen.LoadPEM(certPath)
	if err != nil {
		return err
	}
	if cert.Type != "CERTIFICATE" {
		return fmt.Errorf("invalid cert file, got type %q", key.Type)
	}

	if err := osp.Sign(signer, cert.Bytes); err != nil {
		return err
	}

	signedDescriptor, err := osp.DescriptorBytes()
	if err != nil {
		return err
	}

	return os.WriteFile(pkgPath+ospkgs.DescriptorExt, signedDescriptor, defaultFilePerm)
}

func parsePkgPath(path string) (string, error) {
	if len(path) == 0 {
		return DefaultOutName, nil
	}

	if stat, err := os.Stat(path); err != nil { //nolint:nestif
		if dir := filepath.Dir(path); dir != "." {
			if _, err := os.Stat(dir); err != nil {
				return "", err
			}
		}
	} else {
		if stat.IsDir() {
			return filepath.Join(path, DefaultOutName), nil
		}
	}

	ext := filepath.Ext(path)
	switch ext {
	case "":
		return path, nil
	case ospkgs.OSPackageExt, ospkgs.DescriptorExt:
		return strings.TrimSuffix(path, ext), nil
	default:
		return "", fmt.Errorf("%w %q", ErrInvalidSuffix, ext)
	}
}
