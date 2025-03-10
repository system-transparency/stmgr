package ospkg

import (
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"system-transparency.org/stboot/opts"
	ospkgs "system-transparency.org/stboot/ospkg"
	"system-transparency.org/stboot/stlog"
)

const (
	trustPolicyFile = "trust_policy.json"
	signingRootFile = "ospkg_signing_root.pem"
)

// VerifyTrustPolicy verifies an OS package using the provided path to
// a Trust policy directory
func VerifyTrustPolicy(trustPolicyDir, pkgPath string) error {
	stlog.Info("Using Trust policy directory %q", trustPolicyDir)
	now := time.Now()

	var threshold int
	trustPolicy, err := opts.ReadTrustPolicy(filepath.Join(trustPolicyDir, trustPolicyFile))
	if err != nil {
		return err
	}
	threshold = trustPolicy.SignatureThreshold

	rootCert, err := readRootCert(filepath.Join(trustPolicyDir, signingRootFile), now)
	if err != nil {
		return err
	}

	return verify(threshold, rootCert, now, pkgPath)
}

// VerifyRootCert verifies an OS package using the provided path to a
// root certificate.
func VerifyRootCert(rootCertPath, pkgPath string) error {
	stlog.Info("Using root certificate only: expecting all found signatures to be valid")
	now := time.Now()

	rootCert, err := readRootCert(rootCertPath, now)
	if err != nil {
		return err
	}

	return verify(0, rootCert, now, pkgPath)
}

func readRootCert(rootCertPath string, now time.Time) (*x509.Certificate, error) {
	certs, err := opts.ReadCertsFile(rootCertPath, now)
	if err != nil {
		return nil, err
	}
	if got := len(certs); got != 1 {
		return nil, fmt.Errorf("exactly one root certificate is expected, file contains %d", got)
	}
	return certs[0], nil
}

// verify verifies an OS package using the provided root certificate,
// requiring threshold number of valid signatures. If threshold is 0,
// all found signatures are expected to be valid.
func verify(threshold int, rootCert *x509.Certificate, now time.Time, pkgPath string) error {
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

	numSigs, valid, err := osp.Verify(rootCert, now)
	if err != nil {
		return err
	}

	if threshold == 0 {
		threshold = numSigs
	}

	if valid < threshold {
		return fmt.Errorf("not enough valid signatures: %d found, %d valid, %d required", numSigs, valid, threshold)
	}

	stlog.Info("Signatures: %d found, %d valid, %d required", numSigs, valid, threshold)

	return nil
}
