package uki

import (
	"fmt"
	"os"
	"path/filepath"

	diskfs "github.com/diskfs/go-diskfs"
	diskpkg "github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/diskfs/go-diskfs/partition/gpt"
)

func mkvfat(out, binary string) error {
	var espSize int64

	for _, file := range []string{binary} {
		if file != "" {
			fi, err := os.Stat(file)
			if err != nil {
				return err
			}

			espSize += fi.Size()
		}
	}

	// Smallest possible FAT32 partition
	var minVfatSize int64 = 33548800
	if espSize < minVfatSize {
		espSize = minVfatSize
	}

	var (
		align1MiBMask    uint64 = (1<<44 - 1) << 20
		blkSize          int64  = 512
		partitionStart   int64  = 2048
		partSize                = int64(uint64(espSize) & align1MiBMask)
		diskSize                = partSize + 5*1024*1024
		partitionSectors        = partSize / blkSize
		partitionEnd            = partitionSectors - partitionStart + 1
	)

	disk, err := diskfs.Create(out, diskSize, diskfs.Raw, diskfs.SectorSize512)
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

	return nil
}

func mkiso(out, vfat string) error {
	fi, err := os.Stat(vfat)
	if err != nil {
		return err
	}

	size := fi.Size()
	size += 5 * 1024 * 1024 // disk padding
	iso, err := diskfs.Create(out, size, diskfs.Raw, diskfs.SectorSize512)

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
