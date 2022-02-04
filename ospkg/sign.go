package ospkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ospkgs "github.com/system-transparency/stboot/ospkg"
	"github.com/system-transparency/stmgr/keygen"
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

	cert, err := keygen.LoadPEM(certPath)
	if err != nil {
		return err
	}

	if err := osp.Sign(key, cert); err != nil {
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
