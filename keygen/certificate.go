package keygen

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
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
	IsCa           bool
	IssuerCertFile string // Empty, for creating a self-signed cert.
	IssuerKeyFile  string // Private root CA signing key.
	LeafKeyFile    string // Public key
	NotBefore      time.Time
	NotAfter       time.Time
	CertOut        string
	KeyOut         string
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
		newCert []byte
		newKey  crypto.Signer
	)

	if len(args.IssuerCertFile) == 0 {
		var err error
		var signer crypto.Signer
		// Create a self-signed certificate.
		if len(args.IssuerKeyFile) > 0 {
			signer, err = LoadPrivateKey(args.IssuerKeyFile)
		} else {
			_, signer, err = ed25519.GenerateKey(rand.Reader)
			newKey = signer
		}
		if err != nil {
			return err
		}
		newCert, err = newCaCert(signer, args.NotBefore, args.NotAfter)
		if err != nil {
			return err
		}
	} else {
		rootCert, rootKey, err := parseCaFiles(args.IssuerCertFile, args.IssuerKeyFile)
		if err != nil {
			return err
		}
		var leafPublicKey crypto.PublicKey

		if len(args.LeafKeyFile) > 0 {
			leafPublicKey, err = LoadPublicKey(args.LeafKeyFile)
		} else {
			leafPublicKey, newKey, err = ed25519.GenerateKey(rand.Reader)
		}
		if err != nil {
			return err
		}
		// Create a certificate signed by a root certificate.
		newCert, err = newSigningCert(rootCert, rootKey, leafPublicKey, args.NotBefore, args.NotAfter)
		if err != nil {
			return err
		}
	}

	if newKey != nil {
		if err := writeKey(newKey, keyOut); err != nil {
			return err
		}
	}
	return writeCert(newCert, args.CertOut)
}

func checkArgs(args *CertificateArgs) error {
	if args.IsCa {
		if len(args.IssuerCertFile) != 0 {
			stlog.Warn("isCA specified, will ignore rootCert")
			args.IssuerCertFile = ""
		}
		return nil
	}
	// For generating non-CA certs, either both key and cert must
	// be provided, or none (in which case default filenames are
	// used).
	if len(args.IssuerCertFile) == 0 && len(args.IssuerKeyFile) != 0 {
		return ErrNoRootCert
	}

	if len(args.IssuerKeyFile) == 0 && len(args.IssuerCertFile) != 0 {
		return ErrNoRootKey
	}
	return nil
}

func writeCert(cert []byte, certOut string) error {
	return WritePEM(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}, certOut)
}

// This function makes sure the on-disk format of the key and certificate are correct.
func writeKey(key crypto.Signer, keyOut string) error {
	marshaledKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	return WritePEM(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: marshaledKey,
	}, keyOut)
}

func parseCaFiles(rootCertPath, rootKeyPath string) (*x509.Certificate, crypto.Signer, error) {
	rootCertPath, err := filepath.Abs(rootCertPath)
	if err != nil {
		return nil, nil, err
	}

	rootCertDER, err := LoadCertBytes(rootCertPath)
	if err != nil {
		return nil, nil, err
	}

	rootCert, err := x509.ParseCertificate(rootCertDER)
	if err != nil {
		return nil, nil, err
	}

	rootKeyPath, err = filepath.Abs(rootKeyPath)
	if err != nil {
		return nil, nil, err
	}

	rootKey, err := LoadPrivateKey(rootKeyPath)
	if err != nil {
		return nil, nil, err
	}

	return rootCert, rootKey, nil
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

func randomSerial() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), serialNumberRange)

	return rand.Int(rand.Reader, serialNumberLimit)
}

func hashPublicKey(pub crypto.PublicKey) (string, error) {
	ed25519pub, ok := pub.(ed25519.PublicKey)
	if !ok {
		return "", fmt.Errorf("not ed25519")
	}
	hash := sha256.Sum256(ed25519pub)
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// Create a new self-signed certificate.
func newCaCert(signer crypto.Signer, notBefore, notAfter time.Time) ([]byte, error) {
	keyHash, err := hashPublicKey(signer.Public())
	if err != nil {
		return nil, err
	}
	serialNumber, err := randomSerial()
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		Issuer:                pkix.Name{CommonName: keyHash}, // anything that is unique
		Subject:               pkix.Name{CommonName: keyHash}, // anything that is unique
		SerialNumber:          serialNumber,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
	}
	return x509.CreateCertificate(rand.Reader, &template, &template, signer.Public(), signer)
}

// Creates a new signing certificate, signed by the CA key.
func newSigningCert(caCert *x509.Certificate, caSigner crypto.Signer, leafPublicKey crypto.PublicKey, notBefore, notAfter time.Time) ([]byte, error) {
	keyHash, err := hashPublicKey(leafPublicKey)
	if err != nil {
		return nil, err
	}
	serialNumber, err := randomSerial()
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		Issuer:       caCert.Subject,
		Subject:      pkix.Name{CommonName: keyHash}, // anything that is unique
		SerialNumber: serialNumber,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		NotBefore:    notBefore,
		NotAfter:     notAfter,
	}
	return x509.CreateCertificate(rand.Reader, &template, caCert, leafPublicKey, caSigner)
}
