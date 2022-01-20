package ospkg

import (
	"os"

	ospkgs "github.com/system-transparency/stboot/ospkg"
)

func Create(out, label, url, kernel, initramfs, cmdline string) error {
	osp, err := ospkgs.CreateOSPackage(label, url, kernel, initramfs, cmdline)
	if err != nil {
		return err
	}

	archive, err := osp.ArchiveBytes()
	if err != nil {
		return err
	}
	if err := os.WriteFile(out+ospkgs.OSPackageExt, archive, 0666); err != nil {
		return err
	}

	descriptor, err := osp.DescriptorBytes()
	if err != nil {
		return err
	}
	if err := os.WriteFile(out+ospkgs.DescriptorExt, descriptor, 0666); err != nil {
		return err
	}

	return nil
}
