package ospkg

import (
	"io/fs"
	"os"

	ospkgs "github.com/system-transparency/stboot/ospkg"
)

const defaultFilePerm fs.FileMode = 0o600

type CreateArgs struct {
	OutPath   string
	Label     string
	URL       string
	Kernel    string
	Initramfs string
	Cmdline   string
}

// Create will create a new OS package using the
// stboot/ospkg package. If no errors occur, it will
// be written to disk using the provided path or
// default values.
func Create(args *CreateArgs) error {
	args, err := checkArgs(args)
	if err != nil {
		return err
	}

	osp, err := ospkgs.CreateOSPackage(
		args.Label,
		args.URL,
		args.Kernel,
		args.Initramfs,
		args.Cmdline,
	)
	if err != nil {
		return err
	}

	archive, err := osp.ArchiveBytes()
	if err != nil {
		return err
	}

	if err := os.WriteFile(args.OutPath+ospkgs.OSPackageExt, archive, defaultFilePerm); err != nil {
		return err
	}

	descriptor, err := osp.DescriptorBytes()
	if err != nil {
		return err
	}

	if err := os.WriteFile(args.OutPath+ospkgs.DescriptorExt, descriptor, defaultFilePerm); err != nil {
		return err
	}

	return nil
}

func checkArgs(args *CreateArgs) (*CreateArgs, error) {
	var err error

	args.OutPath, err = parsePkgPath(args.OutPath)
	if err != nil {
		return nil, err
	}

	return args, nil
}
