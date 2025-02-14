package ospkg

import (
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

	signer, err := keygen.LoadPrivateKey(keyPath)
	if err != nil {
		return err
	}
	cert, err := keygen.LoadCertBytes(certPath)
	if err != nil {
		return err
	}
	if err := osp.Sign(signer, cert); err != nil {
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

	if stat, err := os.Stat(path); err != nil {
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
