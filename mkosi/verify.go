package mkosi

import (
	"fmt"
	"os"

	"github.com/foxboron/go-uefi/efi/pecoff"
	"github.com/foxboron/go-uefi/efi/pkcs7"
	"github.com/foxboron/go-uefi/efi/util"
)

func Verify(certPath, uki string) error {
	ukiFile, err := os.ReadFile(uki)
	if err != nil {
		return err
	}
	x509Cert, err := util.ReadCertFromFile(certPath)
	if err != nil {
		return err
	}
	sigs, err := pecoff.GetSignatures(ukiFile)
	if err != nil {
		return fmt.Errorf("%s: %w", ukiFile, err)
	}
	if len(sigs) == 0 {
		return fmt.Errorf("no signatures found")
	}
	var noMatch uint
	for _, signature := range sigs {
		ok, err := pkcs7.VerifySignature(x509Cert, signature.Certificate)
		if err != nil {
			return err
		} else if !ok {
			noMatch++
		}
	}
	if noMatch > 0 {
		return fmt.Errorf("%d signatures matched out of %d", noMatch, len(sigs))
	}
	return nil
}
