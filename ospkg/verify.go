package ospkg

import (
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"system-transparency.org/stboot/opts"
	ospkgs "system-transparency.org/stboot/ospkg"
	"system-transparency.org/stboot/stlog"
)

// Verify verifies an OS package using the provided path to a root
// certificate.
func Verify(rootCertPath, trustPolicyPath, pkgPath string) error {
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

	now := time.Now().Truncate(time.Second)
	var rootCert *x509.Certificate

	if certs, err := opts.ReadCertsFile(rootCertPath, now); err != nil {
		return err
	} else {
		if got := len(certs); got != 1 {
			return fmt.Errorf("exactly one root certificate is expected, file contains %d", got)
		}
		rootCert = certs[0]
	}

	numSigs, valid, err := osp.Verify(rootCert, now)
	if err != nil {
		return err
	}

	var threshold int
	if trustPolicyPath == "" {
		stlog.Info("No Trust policy file set: expecting all found signatures to be valid")
		threshold = numSigs
	} else {
		stlog.Info("Reading Trust policy file %q", trustPolicyPath)
		trustPolicy, err := opts.ReadTrustPolicy(trustPolicyPath)
		if err != nil {
			return err
		}
		threshold = trustPolicy.SignatureThreshold
	}

	if valid < threshold {
		return fmt.Errorf("not enough valid signatures: %d found, %d valid, %d required", numSigs, valid, threshold)
	}

	fmt.Fprintf(os.Stderr, "Signatures: %d found, %d valid, %d required\n", numSigs, valid, threshold)

	return nil
}
