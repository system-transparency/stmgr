package eval

import (
	"flag"

	"system-transparency.org/stboot/stlog"
	"system-transparency.org/stmgr/keygen"
)

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
	certificateIsCA := certificateCmd.Bool("isCA", false, "Generate self signed root certificate.")
	certificateValidFrom := certificateCmd.String("validFrom", "", "Date formatted as RFC822."+
		" Defaults to time of creation.")
	certificateValidUntil := certificateCmd.String("validUntil", "", "Date formatted as RFC822."+
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

	// Call function with parsed flags
	return keygen.Certificate(
		&keygen.CertificateArgs{
			IsCa:         *certificateIsCA,
			RootCertPath: *certificateRootCert,
			RootKeyPath:  *certificateRootKey,
			NotBefore:    *certificateValidFrom,
			NotAfter:     *certificateValidUntil,
			CertOut:      *certificateCertOut,
			KeyOut:       *certificateKeyOut,
		},
	)
}
