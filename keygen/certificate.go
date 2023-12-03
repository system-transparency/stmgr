package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"system-transparency.org/stboot/stlog"
)

var (
	ErrNoRootCert = errors.New("missing rootCert")
	ErrNoRootKey  = errors.New("missing rootKey")
)

const (
	DefaultCertName     = "cert.pem"
	DefaultRootCertName = "rootcert.pem"
	DefaultKeyName      = "key.pem"
	DefaultRootKeyName  = "rootkey.pem"
	serialNumberRange   = 128
)

// CertificateArgs is a list of arguments
// that's passed to Certificate().
type CertificateArgs struct {
	IsCa         bool
	RootCertPath string
	RootKeyPath  string
	NotBefore    time.Time
	NotAfter     time.Time
	CertOut      string
	KeyOut       string
}

// Certificate is used to create a new certificate and private
// key to sign OS packages with.
func Certificate(args *CertificateArgs) error {
	// Check if the provided args are sane
	if err := checkArgs(args); err != nil {
		return err
	}

	// Evaluate the path for the private key.
	keyOut, err := parseKeyPath(args.IsCa, args.KeyOut)
	if err != nil {
		return err
	}

	// Evaluate the path for the certificate.
	args.CertOut, err = parseCertPath(args.IsCa, args.CertOut)
	if err != nil {
		return err
	}

	var (
		newCert *x509.Certificate
		newKey  *rsa.PrivateKey
	)

	if len(args.RootCertPath) == 0 {
		// Create a self-signed certificate.
		newCert, newKey, err = newCertWithKey(nil, nil, args.NotBefore, args.NotAfter)
		if err != nil {
			return err
		}
	} else {
		rootCert, rootKey, err := parseCaFiles(args.RootCertPath, args.RootKeyPath)
		if err != nil {
			return err
		}

		// Create a certificate signed by a root certificate.
		newCert, newKey, err = newCertWithKey(rootCert, rootKey, args.NotBefore, args.NotAfter)
		if err != nil {
			return err
		}
	}

	return writeToDisk(newCert, newKey, args.CertOut, keyOut)
}

func checkArgs(args *CertificateArgs) error {
	switch {
	case args.IsCa && (len(args.RootCertPath) != 0 || len(args.RootKeyPath) != 0):
		stlog.Warn("isCa specified, will ignore rootKey and rootCert")

		return nil

	case len(args.RootCertPath) == 0 && len(args.RootKeyPath) != 0:
		return ErrNoRootCert

	case len(args.RootKeyPath) == 0 && len(args.RootCertPath) != 0:
		return ErrNoRootKey

	default:
		return nil
	}
}

// This function makes sure the on-disk format of the key and certificate are correct.
func writeToDisk(cert *x509.Certificate, key *rsa.PrivateKey, certOut, keyOut string) error {
	marshaledKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	certBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	if err := WritePEM(&certBlock, certOut); err != nil {
		return err
	}

	keyBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: marshaledKey,
	}
	if err := WritePEM(&keyBlock, keyOut); err != nil {
		return err
	}

	return nil
}

func parseCaFiles(rootCertPath, rootKeyPath string) (*x509.Certificate, *interface{}, error) {
	rootCertPath, err := filepath.Abs(rootCertPath)
	if err != nil {
		return nil, nil, err
	}

	rootCertBlock, err := LoadPEM(rootCertPath)
	if err != nil {
		return nil, nil, err
	}

	rootCert, err := x509.ParseCertificate(rootCertBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	rootKeyPath, err = filepath.Abs(rootKeyPath)
	if err != nil {
		return nil, nil, err
	}

	rootKeyBlock, err := LoadPEM(rootKeyPath)
	if err != nil {
		return nil, nil, err
	}

	rootKey, err := x509.ParsePKCS8PrivateKey(rootKeyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return rootCert, &rootKey, nil
}

func parseKeyPath(isCA bool, path string) (string, error) {
	if len(path) == 0 {
		if isCA {
			return DefaultRootKeyName, nil
		}

		return DefaultKeyName, nil
	}

	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		return "", err
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func parseCertPath(isCA bool, path string) (string, error) {
	if len(path) == 0 {
		if isCA {
			return DefaultRootCertName, nil
		}

		return DefaultCertName, nil
	}

	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		return "", err
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

// This is the core function to create a certificate and the corresponding rsa key.
// It uses Golangs std crypto package and uses a cryptographically secure (enough) source
// of entropy, namely /dev/urandom on Linux.
func newCertWithKey(rootCert *x509.Certificate, rootKey *interface{}, notBefore, notAfter time.Time) (*x509.Certificate, *rsa.PrivateKey, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), serialNumberRange)

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		NotBefore:    notBefore,
		NotAfter:     notAfter,
	}

	newPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	var certBytes []byte

	if rootCert == nil || rootKey == nil {
		template.KeyUsage |= x509.KeyUsageCertSign
		template.BasicConstraintsValid = true
		template.IsCA = true

		certBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, newPriv.Public(), newPriv)
		if err != nil {
			return nil, nil, err
		}
	} else {
		certBytes, err = x509.CreateCertificate(rand.Reader, &template, rootCert, newPriv.Public(), *rootKey)
		if err != nil {
			return nil, nil, err
		}
	}

	newCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, err
	}

	return newCert, newPriv, nil
}
