module system-transparency.org/stmgr

// We don't want to depend on golang version later than is available
// in debian's stable or backports repos.
go 1.19

require (
	github.com/gdamore/tcell/v2 v2.5.4
	github.com/rivo/tview v0.0.0-20230130130022-4a1b7a76c01c
	github.com/u-root/u-root v0.10.0
	system-transparency.org/stboot v0.0.0-20230130142012-033e4de02012
)

require (
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/golang/protobuf v1.4.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.10.6 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
