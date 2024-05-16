module system-transparency.org/stmgr

// We don't want to depend on golang version later than is available
// in debian's stable or backports repos.
go 1.19

require (
	github.com/diskfs/go-diskfs v1.3.0
	github.com/foxboron/go-uefi v0.0.0-20230808201820-18b9ba9cd4c3
	sigsum.org/sigsum-go v0.7.2
	system-transparency.org/stboot v0.4.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/frankban/quicktest v1.14.5 // indirect
	github.com/google/go-tpm v0.9.1-0.20230914180155-ee6cbcd136f8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/xattr v0.4.1 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/u-root/u-root v0.12.0 // indirect
	github.com/u-root/uio v0.0.0-20230305220412-3e8cd9d6bf63 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/vishvananda/netlink v1.2.1-beta.2 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/djherbis/times.v1 v1.2.0 // indirect
)
