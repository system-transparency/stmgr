package uki

import (
	"bytes"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/foxboron/go-uefi/authenticode"
	"system-transparency.org/stmgr/keygen"
)

func writeDiskFs(fs filesystem.FileSystem, file, diskPath string) error {
	if path := filepath.Dir(diskPath); path != "/" {
		if err := fs.Mkdir(path); err != nil {
			return err
		}
	}

	rw, err := fs.OpenFile(diskPath, os.O_CREATE|os.O_RDWR)
	if err != nil {
		return fmt.Errorf("failed to make %s on the disk image", diskPath)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read %s", filepath.Base(file))
	}

	_, err = rw.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write %s", filepath.Base(file))
	}

	return nil
}

func createTempFilename() (string, error) {
	// Go only allows us to create templated filenames if we make one then delete
	// it.
	f, err := os.CreateTemp("", "stmgr.*.vfat")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary vfat file: %w", err)
	}

	f.Close()
	os.RemoveAll(f.Name())

	return f.Name(), nil
}

const (
	formatIso = "iso"
	formatUki = "uki"
)

func Create(args []string) error {
	ukiCmd := flag.NewFlagSet("uki", flag.ExitOnError)
	out := ukiCmd.String("out", "stmgr", "output path with format as suffix (default: stmgr)")
	initramfs := ukiCmd.String("initramfs", "/tmp/initramfs.linux_amd64.cpio", "initramfs of a system, commonly from u-root")
	cmdline := ukiCmd.String("cmdline", "", "additional cmdline options for the kernel")
	osrelease := ukiCmd.String("osrelease", "", "os-release file for the uki")
	kernel := ukiCmd.String("kernel", "", "kernel or EFI binary to boot")
	force := ukiCmd.Bool("force", false, "remove existing files (default: false)")
	format := ukiCmd.String("format", "iso", "comma separated list of output formats (iso, uki)")
	stub := ukiCmd.String("stub", "", "UKI stub location (defaults to an embedded stub)")
	sbat := ukiCmd.String("sbat", "", "SBAT metadata")
	appendSbat := ukiCmd.Bool("append-sbat", false, "Append SBAT metadata to the existing section (default: false)")
	signCert := ukiCmd.String("signcert", "", "Certificate corresponding to the private key (a file in PEM format)")
	signKey := ukiCmd.String("signkey", "", "Private key for signing the uki for Secure Boot (a file in PEM format)")

	if err := ukiCmd.Parse(args); err != nil {
		return err
	}

	if ukiCmd.NArg() > 0 {
		return errors.New("unexpected positional argument")
	}

	formats := strings.Split(*format, ",")
	outputIso := false
	outputUki := false
	for _, f := range formats {
		switch f {
		case formatIso:
			outputIso = true
		case formatUki:
			outputUki = true
		case "":
		default:
			return fmt.Errorf("format list can only contain iso or uki")
		}
	}

	if !outputIso && !outputUki {
		return fmt.Errorf("no output format specified")
	}

	ukiFilename := *out
	isoFilename := *out
	if !strings.HasSuffix(isoFilename, ".iso") {
		isoFilename = fmt.Sprintf("%s.iso", isoFilename)
	}
	if !strings.HasSuffix(ukiFilename, ".uki") {
		ukiFilename = fmt.Sprintf("%s.uki", ukiFilename)
	}

	if *force && outputIso {
		os.Remove(isoFilename)
	}
	if *force && outputUki {
		os.Remove(ukiFilename)
	}
	if !outputUki {
		// File we write for the UKI
		stmgrUkiTmpfile, err := os.CreateTemp("", "stmgr-uki.*.efi")
		if err != nil {
			return fmt.Errorf("failed to make temporary file for the UKI")
		}
		defer os.Remove(stmgrUkiTmpfile.Name())
		ukiFilename = stmgrUkiTmpfile.Name()
	}

	// Require both or none of these flags (XOR)
	if (*signCert != "") != (*signKey != "") {
		return fmt.Errorf("both -signcert and -signkey are required for signing UKI")
	}

	uki := &UKI{}

	defer uki.Cleanup()

	if err := uki.SetCmdline(*cmdline); err != nil {
		return fmt.Errorf("failed setting cmdline: %w", err)
	}

	if err := uki.SetKernel(*kernel); err != nil {
		return fmt.Errorf("failed setting kernel: %w", err)
	}

	if err := uki.SetOSRelease(*osrelease); err != nil {
		return fmt.Errorf("failed setting os-release: %w", err)
	}

	if err := uki.SetInitramfs(*initramfs); err != nil {
		return fmt.Errorf("failed setting initramfs")
	}

	// SBAT section is optional
	uki.SetSBAT(*sbat, *appendSbat)

	if err := generateUKI(uki, *stub, ukiFilename); err != nil {
		return fmt.Errorf("failed to write UKI: %w", err)
	}

	if (*signKey != "") && (*signCert != "") {
		if err := signPE(*signKey, *signCert, ukiFilename); err != nil {
			return fmt.Errorf("failed to sign UKI/PE: %w", err)
		}
	}

	if outputIso {
		return toISO(ukiFilename, isoFilename)
	}

	return nil
}

func ToISO(args []string) error {
	cmd := flag.NewFlagSet("iso", flag.ExitOnError)
	inFilename := cmd.String("in", "", "filename of an input UKI to format as an ISO")
	outFilename := cmd.String("out", "", "where to store output ISO (default: <INPUT-NAME>.iso)")
	if err := cmd.Parse(args); err != nil {
		return err
	}
	if *inFilename == "" {
		return fmt.Errorf("missing required option: -in")
	}
	if *outFilename == "" {
		*outFilename = strings.TrimSuffix(*inFilename, ".uki")
		*outFilename += ".iso"
	}
	return toISO(*inFilename, *outFilename)
}

func toISO(ukiFilename, isoFilename string) error {
	tmpfilename, err := createTempFilename()
	if err != nil {
		return fmt.Errorf("failed to make temporary filename: %w", err)
	}
	defer os.Remove(tmpfilename)

	if err := mkvfat(tmpfilename, ukiFilename); err != nil {
		return fmt.Errorf("failed to make vfat partition: %w", err)
	}
	if err := mkiso(isoFilename, tmpfilename); err != nil {
		return fmt.Errorf("failed to make iso: %w", err)
	}
	return nil
}

func signPE(keyFileName, certFileName, peFileName string) error {
	peData, err := os.ReadFile(peFileName)
	if err != nil {
		return fmt.Errorf("ReadFile failed: %w", err)
	}

	pe, err := authenticode.Parse(bytes.NewReader(peData))
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	if sigs, err := pe.Signatures(); err != nil {
		return fmt.Errorf("signatures failed: %w", err)
	} else if len(sigs) > 0 {
		return fmt.Errorf("PE is already signed")
	}

	signer, err := keygen.LoadPrivateKey(keyFileName)
	if err != nil {
		return fmt.Errorf("LoadPrivateKey failed: %w", err)
	}
	certDER, err := keygen.LoadCertBytes(certFileName)
	if err != nil {
		return fmt.Errorf("LoadCertBytes failed: %w", err)
	}
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return fmt.Errorf("invalid x509 certificate: %w", err)
	}

	// pe.Sign both signs and returns the signature
	if _, err := pe.Sign(signer, cert); err != nil {
		return fmt.Errorf("sign failed: %w", err)
	}

	info, err := os.Stat(peFileName)
	if err != nil {
		return fmt.Errorf("stat failed: %w", err)
	}

	// pe.Bytes returns the signed PE
	if err = os.WriteFile(peFileName, pe.Bytes(), info.Mode()); err != nil {
		return fmt.Errorf("WriteFile failed: %w", err)
	}

	return nil
}
