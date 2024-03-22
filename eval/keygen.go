package eval

import (
	"flag"
	"fmt"
	"time"

	"system-transparency.org/stboot/stlog"
	"system-transparency.org/stmgr/keygen"
)

const (
	defaultValidDuration = 72 * time.Hour
)

func parseDate(date string, defaultDate time.Time) (time.Time, error) {
	if len(date) == 0 {
		return defaultDate, nil
	}

	return time.Parse(time.RFC3339, date)
}

// KeygenCertificate takes arguments like os.Args as a string array
// and maps them to their corresponding flags using the std flag
// package. It then calls keygen.Certificate after they are parsed.
func KeygenCertificate(args []string) error {
	// Create a custom flag set and register flags
	certificateCmd := flag.NewFlagSet("certificate", flag.ExitOnError)
	certificateRootCert := certificateCmd.String("rootCert", "", "Root cert in PEM format to sign the new certificate."+
		" Ignored if -isCA is set.")
	certificateRootKey := certificateCmd.String("rootKey", "", "Root key in PEM format to sign the new certificate."+
		" Ignored if -isCA is set.")
	certificateSubjectKey := certificateCmd.String("subjectKey", "", "public key to certify, in PEM or OpenSSH format.")
	certificateIsCA := certificateCmd.Bool("isCA", false, "Generate self signed root certificate.")
	certificateValidFrom := certificateCmd.String("validFrom", "", "Date formatted as RFC3339."+
		" Defaults to time of creation.")
	certificateValidUntil := certificateCmd.String("validUntil", "", "Date formatted as RFC3339."+
		" Defaults to time of creation + 72h.")
	certificateCertOut := certificateCmd.String("certOut", "", "Output certificate file."+
		" Defaults to cert.pem or rootcert.pem is -isCA is set.")
	certificateKeyOut := certificateCmd.String("keyOut", "", "Output key file."+
		" Defaults to key.pem or rootkey.pem if -isCA is set.")
	certificateLogLevel := certificateCmd.String("loglevel", "", "Set loglevel to any of debug, info, warn, error (default) and panic.")

	// Parse which flags are provided to the function
	if err := certificateCmd.Parse(args); err != nil {
		return err
	}

	// Adjust loglevel
	setLoglevel(*certificateLogLevel)

	// Print the successfully parsed flags in debug level
	certificateCmd.Visit(func(f *flag.Flag) {
		stlog.Debug("Registered flag %q", f)
	})

	now := time.Now()

	notBefore, err := parseDate(*certificateValidFrom, now)
	if err != nil {
		return fmt.Errorf("invalid validFrom date %q: %w", *certificateValidFrom, err)
	}

	notAfter, err := parseDate(*certificateValidUntil, now.Add(defaultValidDuration))
	if err != nil {
		return fmt.Errorf("invalid validUntil date %q: %w", *certificateValidUntil, err)
	}

	// Call function with parsed flags
	return keygen.Certificate(
		&keygen.CertificateArgs{
			IsCa:           *certificateIsCA,
			IssuerCertFile: *certificateRootCert,
			IssuerKeyFile:  *certificateRootKey,
			SubjectKeyFile: *certificateSubjectKey,
			NotBefore:      notBefore,
			NotAfter:       notAfter,
			CertOut:        *certificateCertOut,
			KeyOut:         *certificateKeyOut,
		},
	)
}
