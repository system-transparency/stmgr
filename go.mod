module system-transparency.org/stmgr

// We don't want to depend on golang version later than is available
// in debian's stable or backports repos.
go 1.19

require (
	github.com/diskfs/go-diskfs v1.3.0
	github.com/foxboron/go-uefi v0.0.0-20230808201820-18b9ba9cd4c3
	sigsum.org/sigsum-go v0.7.2
	system-transparency.org/stboot v0.3.4
)

require (
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/frankban/quicktest v1.14.5 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.10.6 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/xattr v0.4.1 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/u-root/u-root v0.10.0 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20211118161826-650dca95af54 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/djherbis/times.v1 v1.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
