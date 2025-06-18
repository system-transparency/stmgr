module system-transparency.org/stmgr

// We don't want to depend on golang version later than is available
// in debian's stable or backports repos.
go 1.23.0

require (
	github.com/diskfs/go-diskfs v1.3.0
	github.com/foxboron/go-uefi v0.0.0-20250207204325-69fb7dba244f
	sigsum.org/sigsum-go v0.11.2
	system-transparency.org/stboot v0.6.1
)

require (
	filippo.io/age v1.2.1 // indirect
	github.com/dchest/safefile v0.0.0-20151022103144-855e8d98f185 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/pborman/getopt/v2 v2.1.0 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/xattr v0.4.9 // indirect
	github.com/sirupsen/logrus v1.9.4-0.20230606125235-dd1b4c2e81af // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/u-root/u-root v0.14.0 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/vishvananda/netlink v1.3.0 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/djherbis/times.v1 v1.2.0 // indirect
)
