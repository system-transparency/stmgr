package ospkg

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"fmt"
	"os"

	"sigsum.org/sigsum-go/pkg/crypto"
	"sigsum.org/sigsum-go/pkg/proof"
	"sigsum.org/sigsum-go/pkg/types"

	ospkgs "system-transparency.org/stboot/ospkg"
	"system-transparency.org/stmgr/keygen"
)

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
