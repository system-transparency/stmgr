package keygen

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"git.glasklar.is/system-transparency/core/stmgr/log"
)

var (
	ErrNoRootCert = errors.New("missing rootCert")
	ErrNoRootKey  = errors.New("missing rootKey")
)

const (
	DefaultCertName      = "cert.pem"
	DefaultRootCertName  = "rootcert.pem"
	DefaultKeyName       = "key.pem"
	DefaultRootKeyName   = "rootkey.pem"
	defaultValidDuration = 72 * time.Hour
	serialNumberRange    = 128
)

// CertificateArgs is a list of arguments
// that's passed to Certificate().
type CertificateArgs struct {
	IsCa         bool
	RootCertPath string
	RootKeyPath  string
	NotBefore    string
	NotAfter     string
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

	// Evaluate the certificate validity date.
	notBefore, err := parseValidFrom(args.NotBefore)
	if err != nil {
		return err
	}

	// Evaluate the certificate expiration date.
	notAfter, err := parseValidUntil(args.NotAfter)
	if err != nil {
		return err
	}

	var (
		newCert *x509.Certificate
		newKey  ed25519.PrivateKey
	)

	if len(args.RootCertPath) == 0 {
		// Create a self-signed certificate.
		newCert, newKey, err = newCertWithKey(nil, nil, notBefore, notAfter)
		if err != nil {
			return err
		}
	} else {
		rootCert, rootKey, err := parseCaFiles(args.RootCertPath, args.RootKeyPath)
		if err != nil {
			return err
		}

		// Create a certificate signed by a root certificate.
		newCert, newKey, err = newCertWithKey(rootCert, rootKey, notBefore, notAfter)
		if err != nil {
			return err
		}
	}

	return writeToDisk(newCert, newKey, args.CertOut, keyOut)
}

func checkArgs(args *CertificateArgs) error {
	switch {
	case args.IsCa && (len(args.RootCertPath) != 0 || len(args.RootKeyPath) != 0):
		log.Warn("isCa specified, will ignore rootKey and rootCert")

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
func writeToDisk(cert *x509.Certificate, key ed25519.PrivateKey, certOut, keyOut string) error {
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

func parseValidFrom(date string) (time.Time, error) {
	if len(date) == 0 {
		return time.Now(), nil
	}

	return time.Parse(time.RFC822, date)
}

func parseValidUntil(date string) (time.Time, error) {
	if len(date) == 0 {
		return time.Now().Add(defaultValidDuration), nil
	}

	return time.Parse(time.RFC822, date)
}

// This is the core function to create a certificate and the corresponding ed25519 key.
// It uses Golangs std crypto package and uses a cryptographically secure (enough) source
// of entropy, namely /dev/urandom on Linux.
func newCertWithKey(rootCert *x509.Certificate, rootKey *interface{}, notBefore, notAfter time.Time) (*x509.Certificate, ed25519.PrivateKey, error) {
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

	newPub, newPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	var certBytes []byte

	if rootCert == nil || rootKey == nil {
		template.KeyUsage |= x509.KeyUsageCertSign
		template.BasicConstraintsValid = true
		template.IsCA = true

		certBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, newPub, newPriv)
		if err != nil {
			return nil, nil, err
		}
	} else {
		certBytes, err = x509.CreateCertificate(rand.Reader, &template, rootCert, newPub, *rootKey)
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
