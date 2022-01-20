package keygen

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultCertName     = "cert.pem"
	DefaultRootCertName = "rootcert.pem"
	DefaultKeyName      = "key.pem"
	DefaultRootKeyName  = "rootkey.pem"
)

func Certificate(isCa bool, rootCertPath, rootKeyPath, validFrom, validUntil, certOut, keyOut string) error {
	keyOut, err := parseKeyPath(isCa, keyOut)
	if err != nil {
		return err
	}

	certOut, err = parseCertPath(isCa, certOut)
	if err != nil {
		return err
	}

	notBefore, err := parseValidFrom(validFrom)
	if err != nil {
		return err
	}

	notAfter, err := parseValidUntil(validUntil)
	if err != nil {
		return err
	}

	var newCert *x509.Certificate
	var newKey ed25519.PrivateKey
	if rootCertPath == "" || rootKeyPath == "" {
		newCert, newKey, err = newCertWithED25519Key(nil, nil, notBefore, notAfter)
		if err != nil {
			return err
		}
	} else {
		rootCertBlock, err := LoadPEM(rootKeyPath)
		if err != nil {
			return err
		}
		rootCert, err := x509.ParseCertificate(rootCertBlock.Bytes)
		if err != nil {
			return err
		}

		rootKeyBlock, err := LoadPEM(rootKeyPath)
		if err != nil {
			return err
		}
		rootKey, err := x509.ParsePKCS8PrivateKey(rootKeyBlock.Bytes)
		if err != nil {
			return err
		}

		newCert, newKey, err = newCertWithED25519Key(rootCert, &rootKey, notBefore, notAfter)
		if err != nil {
			return err
		}
	}

	key, err := x509.MarshalPKCS8PrivateKey(newKey)
	if err != nil {
		return err
	}

	certBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCert.Raw,
	}
	if err := WritePEM(&certBlock, certOut); err != nil {
		return err
	}

	keyBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: key,
	}
	if err := WritePEM(&keyBlock, keyOut); err != nil {
		return err
	}

	return nil
}

func parseKeyPath(isCA bool, path string) (string, error) {
	if path == "" {
		if isCA {
			return DefaultRootKeyName, nil
		}
		return DefaultKeyName, nil
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
		return "", err
	}
	return path, nil
}

func parseCertPath(isCA bool, path string) (string, error) {
	if path == "" {
		if isCA {
			return DefaultRootCertName, nil
		}
		return DefaultCertName, nil
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
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
		return time.Now().Add(72 * time.Hour), nil
	}
	return time.Parse(time.RFC822, date)
}

func newCertWithED25519Key(rootCert *x509.Certificate, rootKey *interface{}, notBefore, notAfter time.Time) (*x509.Certificate, ed25519.PrivateKey, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
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
		certBytes, err = x509.CreateCertificate(rand.Reader, &template, rootCert, newPub, newPriv)
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
