// Copyright 2022 the System Transparency Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sign

import (
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/system-transparency/stboot/ospkg"
)

const (
	DefaultOutName = "system-transparency-os-package"
)

func Run(keyPath, certPath, pkgPath string) error {
	pkgPath, err := parsePkgPath(pkgPath)
	if err != nil {
		return err
	}
	archive, err := os.ReadFile(pkgPath + ospkg.OSPackageExt)
	if err != nil {
		return err
	}
	descriptor, err := os.ReadFile(pkgPath + ospkg.DescriptorExt)
	if err != nil {
		return err
	}
	osp, err := ospkg.NewOSPackage(archive, descriptor)
	if err != nil {
		return err
	}

	key, err := loadPEM(keyPath)
	if err != nil {
		return err
	}
	cert, err := loadPEM(certPath)
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

	return os.WriteFile(pkgPath+ospkg.DescriptorExt, signedDescriptor, 0666)
}

func parsePkgPath(path string) (string, error) {
	if path == "" {
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
	case ospkg.OSPackageExt, ospkg.DescriptorExt:
		return strings.TrimSuffix(path, ext), nil
	default:
		return "", fmt.Errorf("invalid file extension %q", ext)
	}
}

func loadPEM(path string) (*pem.Block, error) {
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
