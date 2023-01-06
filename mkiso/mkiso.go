package mkiso

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	diskfs "github.com/diskfs/go-diskfs"
	diskpkg "github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/diskfs/go-diskfs/partition/gpt"
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

func mkvfat(out, binary, config string) error {
	var espSize int64
	for _, file := range []string{binary, config} {
		if file != "" {
			fi, err := os.Stat(file)
			if err != nil {
				return err
			}
			espSize += fi.Size()
		}
	}

	var (
		align1MiBMask    uint64 = (1<<44 - 1) << 20
		partSize         int64  = int64(uint64(espSize) & align1MiBMask)
		diskSize         int64  = partSize + 5*1024*1024
		blkSize          int64  = 512
		partitionStart   int64  = 2048
		partitionSectors int64  = partSize / blkSize
		partitionEnd     int64  = partitionSectors - partitionStart + 1
	)

	disk, err := diskfs.Create(out, diskSize, diskfs.Raw)
	if err != nil {
		return fmt.Errorf("failed to create disk file: %w", err)
	}

	table := &gpt.Table{
		Partitions: []*gpt.Partition{
			{
				Start: uint64(partitionStart),
				End:   uint64(partitionEnd),
				Type:  gpt.EFISystemPartition,
				Name:  "EFI System"},
		},
	}

	err = disk.Partition(table)
	if err != nil {
		return fmt.Errorf("failed to create partitiont table: %w", err)
	}

	spec := diskpkg.FilesystemSpec{Partition: 0, FSType: filesystem.TypeFat32}
	fs, err := disk.CreateFilesystem(spec)
	if err != nil {
		return fmt.Errorf("failed to create filesystem")
	}

	if err := writeDiskFs(fs, binary, "/EFI/BOOT/BOOTX64.EFI"); err != nil {
		return fmt.Errorf("failed to write kernel: %w", err)
	}

	if config != "" {
		if err := writeDiskFs(fs, config, "/host_config.json"); err != nil {
			return fmt.Errorf("failed to write host config: %w", err)
		}
	}
	return nil
}

func mkiso(out, vfat string) error {
	fi, err := os.Stat(vfat)
	if err != nil {
		return err
	}
	size := fi.Size()
	size = size + 5*1024*1024 // disk padding
	iso, err := diskfs.Create(out, size, diskfs.Raw)
	if err != nil {
		return err
	}
	iso.LogicalBlocksize = 2048
	fs, err := iso.CreateFilesystem(diskpkg.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeISO9660,
		VolumeLabel: "stboot",
	})
	if err != nil {
		return err
	}
	// This avoids an issue where path.Base in go-diskfs gives us a sigsegv
	vfatName := filepath.Join("vfat", filepath.Base(vfat))
	if err := writeDiskFs(fs, vfat, vfatName); err != nil {
		return fmt.Errorf("failed to write file %s to ISO: %w", vfat, err)
	}
	diskImage, ok := fs.(*iso9660.FileSystem)
	if !ok {
		return fmt.Errorf("not an iso9660 filesystem")
	}
	options := iso9660.FinalizeOptions{
		VolumeIdentifier: "stboot",
		ElTorito: &iso9660.ElTorito{
			BootCatalog: "/BOOT.CAT",
			Entries: []*iso9660.ElToritoEntry{
				{
					Platform:  iso9660.EFI,
					Emulation: iso9660.NoEmulation,
					BootFile:  vfatName,
				},
			},
		},
	}
	if err = diskImage.Finalize(options); err != nil {
		return err
	}
	return nil
}

func Create(args []string) error {
	mkisoCmd := flag.NewFlagSet("mkiso", flag.ExitOnError)
	mkosiOut := mkisoCmd.String("out", "stmgr.iso", "ISO output path (default: stmgr.iso)")
	mkosiKernel := mkisoCmd.String("kernel", "", "kernel or EFI binary to boot")
	mkosiConfig := mkisoCmd.String("config", "host_config.json", "stboot host_configuration (optional)")
	mkosiForce := mkisoCmd.Bool("force", false, "remove existing files (default: false)")

	if err := mkisoCmd.Parse(args); err != nil {
		return err
	}

	if *mkosiForce {
		os.Remove(*mkosiOut)
	}

	if *mkosiKernel == "" {
		return fmt.Errorf("no kernel specified")
	}

	// We care about the name, not the file. Create the file, delete it and use it's name
	f, err := os.CreateTemp("/var/tmp", "stmgr.*.vfat")
	if err != nil {
		return fmt.Errorf("failed to create temporary vfat file: %w", err)
	}
	f.Close()
	os.RemoveAll(f.Name())

	if err := mkvfat(f.Name(), *mkosiKernel, *mkosiConfig); err != nil {
		return fmt.Errorf("failed to make vfat partition: %w", err)
	}

	if err := mkiso(*mkosiOut, f.Name()); err != nil {
		return fmt.Errorf("failed to make iso: %w", err)
	}
	return nil
}
