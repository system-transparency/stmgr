package ospkg

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sigsum.org/sigsum-go/pkg/crypto"
	"sigsum.org/sigsum-go/pkg/proof"
	"sigsum.org/sigsum-go/pkg/types"

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

// AddSigsumProof attaches a Sigsum proof and a corresponding
// certificate for the sigsum submission key.
func AddSigsumProof(proofPath, certPath, pkgPath string) error {
	pkgPath, err := parsePkgPath(pkgPath)
	if err != nil {
		return err
	}

	archive, err := os.Open(pkgPath + ospkgs.OSPackageExt)
	if err != nil {
		return err
	}
	defer archive.Close()
	archiveHash, err := crypto.HashFile(archive)
	if err != nil {
		return err
	}

	descriptorBytes, err := os.ReadFile(pkgPath + ospkgs.DescriptorExt)
	if err != nil {
		return err
	}
	descriptor, err := ospkgs.DescriptorFromBytes(descriptorBytes)
	if err != nil {
		return err
	}

	sigsumProofBytes, err := os.ReadFile(proofPath)
	if err != nil {
		return err
	}
	var sigsumProof proof.SigsumProof
	if err := sigsumProof.FromASCII(bytes.NewReader(sigsumProofBytes)); err != nil {
		return err
	}

	cert, err := keygen.LoadCertBytes(certPath)
	if err != nil {
		return err
	}
	publicKey, err := ed25519PublicKeyFromCert(cert)
	if err != nil {
		return err
	}

	// Verify leaf signature, without examining the rest of the proof.
	if got, want := crypto.HashBytes(publicKey[:]), sigsumProof.Leaf.KeyHash; got != want {
		return fmt.Errorf("public key mismatch, certificate key hash: %x, proof key hash: %x",
			got, want)
	}
	if !types.VerifyLeafMessage(&publicKey, archiveHash[:], &sigsumProof.Leaf.Signature) {
		return fmt.Errorf("invalid leaf signature in sigsum proof")
	}

	if err := descriptor.AddSignature(cert, sigsumProofBytes); err != nil {
		return err
	}

	signedDescriptor, err := descriptor.Bytes()
	if err != nil {
		return err
	}

	return os.WriteFile(pkgPath+ospkgs.DescriptorExt, signedDescriptor, defaultFilePerm)
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
	case ospkgs.OSPackageExt, ospkgs.DescriptorExt:
		return strings.TrimSuffix(path, ext), nil
	default:
		return "", fmt.Errorf("%w %q", ErrInvalidSuffix, ext)
	}
}

func ed25519PublicKeyFromCert(certDER []byte) (crypto.PublicKey, error) {
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return crypto.PublicKey{}, err
	}
	publicKey, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok {
		return crypto.PublicKey{}, fmt.Errorf("invalid public key type: %T", cert.PublicKey)
	}
	var ed25519PublicKey crypto.PublicKey
	copy(ed25519PublicKey[:], publicKey)

	return ed25519PublicKey, nil
}
