package eval

import (
	"flag"

	"github.com/system-transparency/stmgr/keygen"
	"github.com/system-transparency/stmgr/log"
)

func KeygenCertificate(args []string) error {
	certificateCmd := flag.NewFlagSet("keygen", flag.ExitOnError)
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

	if err := certificateCmd.Parse(args); err != nil {
		return err
	}

	setLoglevel(*certificateLogLevel)

	certificateCmd.Visit(func(f *flag.Flag) {
		log.Debugf("Registered flag %q", f)
	})

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
