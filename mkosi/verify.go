package mkosi

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/foxboron/go-uefi/efi/pecoff"
	"github.com/foxboron/go-uefi/efi/pkcs7"
	"github.com/foxboron/go-uefi/efi/util"
	"system-transparency.org/stboot/trust"
)

func Verify(certPath, uki, truspolicy string) error {
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
	trustpolicyFile, err := os.Open(truspolicy)
	if err != nil {
		return err
	}
	var policy trust.Policy
	tp := json.NewDecoder(trustpolicyFile)
	if err := tp.Decode(&policy); err != nil {
		return err
	}
	if policy.SignatureThreshold == 0 {
		return fmt.Errorf("signature threshold is 0")
	}
	if len(sigs) < policy.SignatureThreshold {
		return fmt.Errorf("signature threshold not met")
	}
	var match int
	for _, signature := range sigs {
		ok, err := pkcs7.VerifySignature(x509Cert, signature.Certificate)
		if err != nil {
			return err
		} else if ok {
			match++
		}
	}
	if policy.SignatureThreshold > match {
		return fmt.Errorf("signature threshold not met, out of %d signatures only %d matched with a matching treshold of %d", len(sigs), match, policy.SignatureThreshold)
	}
	return nil
}
