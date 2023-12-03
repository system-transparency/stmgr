package mkosi

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/foxboron/go-uefi/efi/pecoff"
	"system-transparency.org/stmgr/keygen"
)

func Sign(keyPath, certPath, uki string) error {
	ukiFile, err := os.ReadFile(uki)
	if err != nil {
		return err
	}
	ctx := pecoff.PECOFFChecksum(ukiFile)
	certPem, err := keygen.LoadPEM(certPath)
	if err != nil {
		return err
	}
	cert, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse cert: %w", err)
	}
	keyPem, err := keygen.LoadPEM(keyPath)
	if err != nil {
		return err
	}
	key, err := x509.ParsePKCS8PrivateKey(keyPem.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse key: %w", err)
	}
	signer, ok := key.(crypto.Signer)
	if !ok {
		return fmt.Errorf("key is not a signer")
	}
	signature, err := pecoff.CreateSignature(ctx, cert, signer)
	if err != nil {
		return err
	}
	newUKI, err := pecoff.AppendToBinary(ctx, signature)
	if err != nil {
		return err
	}
	return os.WriteFile(uki, newUKI, 0644)
}
