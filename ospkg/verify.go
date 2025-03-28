package ospkg

import (
	"crypto/x509"
	"os"
	"path/filepath"
	"time"

	sigsumCrypto "sigsum.org/sigsum-go/pkg/crypto"
	"system-transparency.org/stboot/opts"
	"system-transparency.org/stboot/ospkg"
	"system-transparency.org/stboot/stlog"
	"system-transparency.org/stboot/trust"
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

	trustPolicy, err := opts.ReadTrustPolicy(filepath.Join(trustPolicyDir, trustPolicyFile))
	if err != nil {
		return err
	}

	rootCerts, err := opts.ReadCertsFile(filepath.Join(trustPolicyDir, signingRootFile), now)
	if err != nil {
		return err
	}

	return verify(rootCerts, trustPolicy, now, pkgPath)
}

// VerifyRootCerts verifies an OS package using the provided path to a
// file containing root certificate(s).
func VerifyRootCerts(rootCertsPath, pkgPath string) error {
	stlog.Info("Using root certificate(s) only: expecting all found signatures to be valid")
	now := time.Now()

	rootCerts, err := opts.ReadCertsFile(rootCertsPath, now)
	if err != nil {
		return err
	}

	return verify(rootCerts, &trust.Policy{SignatureThreshold: ospkg.SignatureThresholdAll}, now, pkgPath)
}

// verify verifies an OS package using the provided root
// certificate(s) and Trust Policy
func verify(rootCerts *x509.CertPool, trustPolicy *trust.Policy, now time.Time, pkgPath string) error {
	pkgPath, err := parsePkgPath(pkgPath)
	if err != nil {
		return err
	}

	r, err := os.Open(pkgPath + ospkg.OSPackageExt)
	if err != nil {
		return err
	}
	hash, err := sigsumCrypto.HashFile(r)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(pkgPath + ospkg.DescriptorExt)
	if err != nil {
		return err
	}
	descriptor, err := ospkg.DescriptorFromBytes(b)
	if err != nil {
		return err
	}

	return descriptor.Verify(rootCerts, nil, trustPolicy, &hash, now)
}
