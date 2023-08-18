package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/schollz/progressbar/v3"

	"github.com/spf13/cobra"
)

var lowContrastBackground = lipgloss.AdaptiveColor{Light: "255", Dark: "234"}
var lowContrastForeground = lipgloss.AdaptiveColor{Light: "234", Dark: "255"}

var highContrastForeground = lipgloss.AdaptiveColor{Light: "238", Dark: "251"}
var highContrastBackground = lipgloss.AdaptiveColor{Light: "251", Dark: "241"}

var (
	RootCmd = &cobra.Command{
		Use: "stmgr",
	}

	kernelCmd = &cobra.Command{
		Use: `kernel [URL | PATH | VERSION | GIT-REF]
               [-c,--config PATH]
               [--menu]
               [-n,--name NAME]`,
		DisableFlagsInUseLine: true,
		Short:                 "Build and list Linux kernels for stboot and OS packages.",
		Run:                   doKernel,
		Args:                  cobra.MaximumNArgs(1),
	}

	kernelConfig string
	kernelMenu   bool
	kernelName   string

	stbootCmd = &cobra.Command{
		Use: `stboot [URL | PATH | VERSION | GIT-REF]
		       [-k,--kernel NAME | URL | PATH | VERSION | GIT-REF]
               [-u,--uroot URL | PATH | GIT-REF]
               [-n,--name NAME]
               [-a,--add URL | PATH=[PATH]]
               [-c,--host-config PATH]
               [-p,--trust-policy PATH]
               [-d,--secure-data PATH]
               [-s,--shim NAME | URL | PATH | VERSION | GIT-REF]
               [--no-shim]
               [--iso | --file]`,
		DisableFlagsInUseLine: true,
		Short:                 "Build and list stboot.",
		Run:                   doStboot,
		Args:                  cobra.MaximumNArgs(1),
	}

	stbootKernel      string
	stbootUroot       string
	stbootName        string
	stbootAdditions   []string
	stbootHostCfg     string
	stbootTrustPolicy string
	stbootSecureData  string
	stbootShim        string
	stbootNoShim      bool
	stbootOuptutIso   string
	stbootOutputFile  string

	osPkgCmd = &cobra.Command{
		Use:   "ospkg",
		Short: "Build, sign and list OS packages.",
		Run:   doOspkg,
	}

	osPkgCreateCmd = &cobra.Command{
		Use: `create [DIRECTORY | DEBOS-FILE | DOCKER-IMAGE]
                     [-k,--kernel NAME | URL | PATH | VERSION | GIT-REF]
                     [-n,--name NAME]
                     [-s,--sign KEY-SPEC]
                     [-o,--output PATH]
                     [-u,--upload REPOSITORY]`,
		DisableFlagsInUseLine: true,
		Short:                 "Build new OS packages.",
		Run:                   doOspkgCreate,
		Args:                  cobra.MaximumNArgs(1),
	}

	osPkgCreateKernel string
	osPkgCreateName   string
	osPkgCreateSign   string
	osPkgCreateOutput string
	osPkgCreateUpload string

	osPkgSignCmd = &cobra.Command{
		Use: `sign [--zip PATH]
                   [--json PATH]
                   [-s,--sign KEY-SPEC]`,
		DisableFlagsInUseLine: true,
		Short:                 "Sign OS packages.",
		Run:                   doOspkgCreate,
		Args:                  cobra.MaximumNArgs(0),
	}

	osPkgSignZip  string
	osPkgSignJson string
	osPkgSignSign string

	repositoryCmd = &cobra.Command{
		Use: `repository [-s,--stboot URL | PATH | GIT-REF]
                   [--log | --no-log]
        repository [-p,--ospkg URL | PATH | GIT-REF]
                   [--log | --no-log]`,
		DisableFlagsInUseLine: true,
		Short:                 "List repositories and upload new OS packages and stboot images.",
		Run:                   doRepository,
	}

	repositoryStboot string
	repositoryOspkg  string
	repositoryLog    bool
	repositoryNoLog  bool

	keyCmd = &cobra.Command{
		Use:                   "key",
		DisableFlagsInUseLine: true,
		Short:                 "Create and list keys.",
		Run:                   doKey,
	}

	keyCreateCmd = &cobra.Command{
		Use: `create [-n,--name KEY-SPEC]
                   [--certify KEY-SPEC]
                   [--self-certify]
                   [--ca | --no-ca]`,
		DisableFlagsInUseLine: true,
		Short:                 "Create new keys.",
		Run:                   doKeyCreate,
	}

	keyCreateName        string
	keyCreateCertify     string
	keyCreateSelfCertify bool
	keyCreateCA          bool
	keyCreateNoCA        bool

	keyExportCmd = &cobra.Command{
		Use: `export [KEY-SPEC]
                   [-o,--output PATH]
                   [--private | --no-private]`,
		DisableFlagsInUseLine: true,
		Short:                 "Export keys.",
		Run:                   doKeyExport,
	}

	keyExportOutput    string
	keyExportPrivate   bool
	keyExportNoPrivate bool

	machineCmd = &cobra.Command{
		Use: `machine [HOSTNAME | IP]
                [--vm]
                [-u,--user USERNAME]
                [-p,--password PASSWD]
                [-s,--stboot PATH | VERSION | GIT-REF]
                [-o,--ospkg PATH | NAME]
                [--zip PATH] [--json PATH]
                [--manual-provision]
                [-c,--host-config PATH]
                [-p,--trust-policy PATH]
                [-d,--secure-data PATH]
                [--reseal]`,
		DisableFlagsInUseLine: true,
		Short:                 "Provision machines.",
		Run:                   doMachine,
	}

	machineVM              bool
	machineUser            string
	machinePassword        string
	machineStboot          string
	machineOspkg           string
	machineZip             string
	machineJson            string
	machineManualProvision bool
	machineHostCfg         string
	machineTrustPolicy     string
	machineSecureData      string
	machineReseal          bool
)

func init() {
	RootCmd.AddCommand(kernelCmd)
	kernelCmd.Flags().StringVarP(&kernelConfig, "config", "c", "", "Kernel config file to use")
	kernelCmd.Flags().BoolVarP(&kernelMenu, "menu", "", false, "Run kconfig menuconfig before building")
	kernelCmd.Flags().StringVarP(&kernelName, "name", "n", "", "Shorthand name for kernel")

	RootCmd.AddCommand(stbootCmd)
	stbootCmd.Flags().StringVarP(&stbootKernel, "kernel", "k", "default", "Kernel to use")
	stbootCmd.Flags().StringVarP(&stbootUroot, "uroot", "u", "1.10.1", "u-root to use")
	stbootCmd.Flags().StringVarP(&stbootName, "name", "n", "", "Shorthand name for stboot")
	stbootCmd.Flags().StringSliceVarP(&stbootAdditions, "add", "a", nil, "Add a file to the stboot image")
	stbootCmd.Flags().StringVarP(&stbootHostCfg, "host-config", "c", "", "Host configuration file to use")
	stbootCmd.Flags().StringVarP(&stbootTrustPolicy, "trust-policy", "p", "", "Trust policy file to use")
	stbootCmd.Flags().StringVarP(&stbootSecureData, "secure-data", "d", "", "Secure data file to use")
	stbootCmd.Flags().StringVarP(&stbootShim, "shim", "s", "", "Shim to use")
	stbootCmd.Flags().BoolVarP(&stbootNoShim, "no-shim", "", false, "Don't include a Shim")
	stbootCmd.Flags().StringVarP(&stbootOuptutIso, "iso", "", "", "Build an ISO image")
	stbootCmd.Flags().StringVarP(&stbootOutputFile, "file", "", "", "Build a file image")

	RootCmd.AddCommand(osPkgCmd)
	osPkgCmd.AddCommand(osPkgCreateCmd)
	osPkgCreateCmd.Flags().StringVarP(&osPkgCreateKernel, "kernel", "k", "", "Kernel to use")
	osPkgCreateCmd.Flags().StringVarP(&osPkgCreateName, "name", "n", "", "Shorthand name for OS package")
	osPkgCreateCmd.Flags().StringVarP(&osPkgCreateSign, "sign", "s", "", "Sign the OS package")
	osPkgCreateCmd.Flags().StringVarP(&osPkgCreateOutput, "output", "o", "", "Output directory for OS package")
	osPkgCreateCmd.Flags().StringVarP(&osPkgCreateUpload, "upload", "u", "", "Upload the OS package")

	osPkgCmd.AddCommand(osPkgSignCmd)
	osPkgSignCmd.Flags().StringVarP(&osPkgSignZip, "zip", "", "", "Zip file to sign")
	osPkgSignCmd.Flags().StringVarP(&osPkgSignJson, "json", "", "", "JSON file to sign")
	osPkgSignCmd.Flags().StringVarP(&osPkgSignSign, "sign", "s", "", "Sign the OS package")

	RootCmd.AddCommand(repositoryCmd)
	repositoryCmd.Flags().StringVarP(&repositoryStboot, "stboot", "s", "", "Stboot to use")
	repositoryCmd.Flags().StringVarP(&repositoryOspkg, "ospkg", "o", "", "OS package to use")
	repositoryCmd.Flags().BoolVarP(&repositoryLog, "log", "l", false, "Transparency log")
	repositoryCmd.Flags().BoolVarP(&repositoryNoLog, "no-log", "", true, "Don't transparency log")

	RootCmd.AddCommand(keyCmd)
	keyCmd.AddCommand(keyCreateCmd)
	keyCreateCmd.Flags().StringVarP(&keyCreateName, "name", "n", "", "Shorthand name for key")
	keyCreateCmd.Flags().StringVarP(&keyCreateCertify, "certify", "c", "", "Certify the key")
	keyCreateCmd.Flags().BoolVarP(&keyCreateSelfCertify, "self-certify", "s", false, "Self certify the key")
	keyCreateCmd.Flags().BoolVarP(&keyCreateCA, "ca", "", false, "CA to use")
	keyCreateCmd.Flags().BoolVarP(&keyCreateNoCA, "no-ca", "", true, "No CA to use")

	keyCmd.AddCommand(keyExportCmd)
	keyExportCmd.Flags().StringVarP(&keyExportOutput, "output", "o", "", "Output file for key")
	keyExportCmd.Flags().BoolVarP(&keyExportPrivate, "private", "p", false, "Export private key")
	keyExportCmd.Flags().BoolVarP(&keyExportNoPrivate, "no-private", "", true, "Don't export private key")

	RootCmd.AddCommand(machineCmd)
	machineCmd.Flags().BoolVarP(&machineVM, "vm", "", false, "Provision a VM")
	machineCmd.Flags().StringVarP(&machineUser, "user", "u", "", "Username to use")
	machineCmd.Flags().StringVarP(&machinePassword, "password", "p", "", "Password to use")
	machineCmd.Flags().StringVarP(&machineStboot, "stboot", "s", "", "Stboot to use")
	machineCmd.Flags().StringVarP(&machineOspkg, "ospkg", "o", "", "OS package to use")
	machineCmd.Flags().StringVarP(&machineZip, "zip", "z", "", "Zip file to use")
	machineCmd.Flags().StringVarP(&machineJson, "json", "j", "", "JSON file to use")
	machineCmd.Flags().BoolVarP(&machineManualProvision, "manual", "m", false, "Manual provision")
	machineCmd.Flags().StringVarP(&machineHostCfg, "host-config", "c", "", "Host configuration file to use")
	machineCmd.Flags().StringVarP(&machineTrustPolicy, "trust-policy", "t", "", "Trust policy file to use")
	machineCmd.Flags().StringVarP(&machineSecureData, "secure-data", "d", "", "Secure data file to use")
	machineCmd.Flags().BoolVarP(&machineReseal, "reseal", "r", false, "Reseal file to use")
}

func doOspkg(cmd *cobra.Command, args []string) {
	headers := []string{"", "Kernel", "Size", "Signers", "Repository"}
	rows := [][]string{
		{"my-debian", "linux", "1.2 Gb", "kai,zaolin", "my-repo"},
		{"mystvmm", "my-kernel", "304 Mb", "zaolin", "(local)"},
		{"test", "stboot", "123 Mb", "(none)", "(local)"},
	}

	printTable(headers, rows)
}

func doRepository(cmd *cobra.Command, args []string) {
}

func doKey(cmd *cobra.Command, args []string) {
}

func doKeyCreate(cmd *cobra.Command, args []string) {
}

func doKeyExport(cmd *cobra.Command, args []string) {
}

func doMachine(cmd *cobra.Command, args []string) {
	if len(args) == 0 && machineVM == false {
		doMachineList()
	} else if len(args) > 0 && machineVM == true {
		fmt.Println("Error: cannot specify both a machine name and --vm")
	} else {
		if machineStboot == "" {
			machineStboot = "default"
		}
		if machineOspkg == "" {
			machineOspkg = "demo"
		}

		if len(args) == 0 && machineVM == true {
			doMachineVirtual(machineStboot, machineOspkg)
		} else {
			doMachineProvision(args[0], machineUser, machinePassword, machineStboot, machineOspkg, machineManualProvision)
		}
	}
}

func doMachineList() {
	headers := []string{"", "Hostname", "stboot", "ospkg", ""}
	rows := [][]string{
		{" v", "x11ssh.local", "1.2.0", "my-debian", ""},
		{" o", "10.0.1.2", "1.1.0 (outdated)", "my-debian", ""},
		{" o", "example.com", "my-stboot", "mystvmm", ""},
		{" x", "  www.example.com", "1.2.0", "my-debian", "ospkg mismatch"},
		{" v", "  mail.example.com", "1.2.1", "my-debian", ""},
		{"", "10.0.1.3", "(none)", "(none)", ""},
	}

	printTable(headers, rows)
}

func doMachineVirtual(stboot string, ospkg string) {
	fmt.Printf("Starting fresh virtual machine with stboot %s and ospkg %s\n", stboot, ospkg)

	// starting swtpm
	// starting http srv
	// starting vm

	// show boot
	// drop into shell
}

func doMachineProvision(name string, user string, password string, stboot string, ospkg string, manual bool) {
	fmt.Printf("Provisioning machine %s with stboot %s and ospkg %s\n", name, stboot, ospkg)

	// upload ospkg

	// if manual
	//     download ubuntu iso

	// connect to bmc

	// if manual
	//     mount iso
	//     start

	// upload stboot
	// mount stboot iso

	// wait for boot
	// attest
}

func doOspkgCreate(cmd *cobra.Command, args []string) {
}

func doOspkgSign(cmd *cobra.Command, args []string) {
}

func doStboot(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		doStbootList()
	} else {
		if stbootName == "" {
			stbootName = args[0]
		}
		doStbootBuild(args[0], stbootKernel, stbootUroot, stbootName)
	}
}

func doStbootList() {
	headers := []string{"", "stboot", "u-root", "shim", "Options"}
	rows := [][]string{
		{"demo", "default", "1.10.1", "1.0", "H"},
		{"default", "5.10.1", "1.10.1", "custom", "+3 files, T"},
		{"my-stboot", "6.0.1", "fd865a4", "1.0", "S, T"},
	}

	printTable(headers, rows)
}

func doStbootBuild(stboot string, kernel string, uroot string, name string) {
	newKernel := true
	for _, k := range kernels {
		newKernel = newKernel && k[0] != kernel
	}
	fmt.Printf("Building %s\n", name)
	fmt.Printf("  Linux: %s", kernel)
	if newKernel {
		fmt.Printf(" (new)\n")
	} else {
		fmt.Printf("\n")
	}
	fmt.Printf("  stboot: %s\n", stboot)
	fmt.Printf("  u-root: %s\n\n", uroot)

	if newKernel {
		doKernelBuild(kernel, "default", false, kernel)
	}

	cloneDesc := "Cloning u-root %s repository..."
	checkoutDesc := fmt.Sprintf("Checking out u-root %s...", uroot)

	printProgress([]string{cloneDesc, checkoutDesc}, 2*time.Second)
	printCompilation(goOutput, 3*time.Second)

	fmt.Printf("Building stboot image...\n")
	time.Sleep(2 * time.Second)
	fmt.Printf("Signing stboot image...\n")
	fmt.Printf("Done.\n")
}

func doKernel(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		doKernelList()
	} else {
		if kernelName == "" {
			kernelName = args[0]
		}
		if kernelConfig == "" {
			kernelConfig = "default"
		}
		doKernelBuild(args[0], kernelConfig, kernelMenu, kernelName)
	}
}

func doKernelBuild(kernelSpec string, config string, menu bool, name string) {
	fmt.Printf("Building %s from Linux %s with %s configuration\n", name, kernelSpec, config)

	downDesc := fmt.Sprintf("Downloading Linux %s tarball...", kernelSpec)
	extractDesc := fmt.Sprintf("Extracting tarball...")

	printProgress([]string{downDesc, extractDesc}, 5*time.Second)
	printCompilation(kernelOutput, 5*time.Second)
	fmt.Printf("make: Leaving directory '/home/%s/.local/stmgr/kernel/%s'\n", os.Getenv("USER"), name)
	fmt.Println("Done")
}

func printCompilation(lines []string, duration time.Duration) {
	for _, line := range lines {
		fmt.Println(line)
		time.Sleep(duration / time.Duration(len(lines)))
	}
}

func printProgress(desc []string, duration time.Duration) {
	width := 0
	for _, line := range desc {
		if lipgloss.Width(line) > width {
			width = lipgloss.Width(line)
		}
	}

	style := lipgloss.NewStyle().Width(width).Align(lipgloss.Left)
	bar := progressbar.Default(100, style.Render(""))

	for _, line := range desc {
		bar.Describe(style.Render(line))
		for i := 0; i < 100/len(desc); i++ {
			bar.Add(1)
			time.Sleep(duration / time.Duration(len(desc)*100))
		}
	}
	bar.Finish()
}

func printTable(headers []string, rows [][]string) {
	measure := make([]int, len(headers))
	for i, h := range headers {
		measure[i] = lipgloss.Width(h)
		for _, r := range rows {
			if len(r[i]) > measure[i] {
				measure[i] = lipgloss.Width(r[i])
			}
		}
	}
	columns := make([]lipgloss.Style, len(headers))
	for i := range headers {
		columns[i] = lipgloss.NewStyle().Width(measure[i]).MarginRight(3)
	}

	head := make([]string, len(headers))
	for i, h := range headers {
		head[i] = columns[i].Copy().Bold(true).Foreground(highContrastForeground).Render(h)
	}
	fmt.Printf("%s\n", lipgloss.JoinHorizontal(lipgloss.Top, head...))
	for _, r := range rows {
		for j, c := range r {
			fmt.Printf("%s", columns[j].Copy().Render(c))
		}
		fmt.Printf("\n")
	}
}

var kernels = [][]string{
	{"default", "5.10.1", "(default)"},
	{"my-kernel", "6.0.1", "custom"},
	{"stboot", "5.10.1", "custom"},
}

func doKernelList() {
	headers := []string{"", "Version", "Config"}
	printTable(headers, kernels)
}

func main() {
	ctx := context.Background()
	err := RootCmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var kernelOutput = []string{
	`  CC      net/ipv4/netfilter.o`,
	`  CC      drivers/fake-pmem/fake-pmem.o`,
	`  CC      drivers/acpi/acpica/psopcode.o`,
	`  CC      drivers/firmware/efi/libstub/skip_spaces.o`,
	`  CC      drivers/firmware/efi/libstub/lib-cmdline.o`,
	`  CC      lib/crc16.o`,
	`  CC      drivers/acpi/acpica/psopinfo.o`,
	`  CC      drivers/firmware/efi/memmap.o`,
	`  CC      drivers/acpi/evged.o`,
	`  CC      net/ipv4/inet_diag.o`,
	`  CC      net/ipv4/tcp_diag.o`,
	`  CC      drivers/md/dm-builtin.o`,
	`  CC      drivers/firmware/efi/libstub/lib-ctype.o`,
	`  AR      drivers/nvmem/built-in.a`,
	`  CC      lib/crc-itu-t.o`,
	`  CC      drivers/md/dm-crypt.o`,
	`  CC      net/ipv6/xfrm6_policy.o`,
	`  CC      drivers/firmware/efi/libstub/alignedmem.o`,
	`  CC      drivers/acpi/sysfs.o`,
	`  CC      drivers/net/ethernet/intel/i40e/i40e_client.o`,
	`  CC      drivers/net/ethernet/intel/i40e/i40e_virtchnl_pf.o`,
	`  CC      net/ipv6/xfrm6_state.o`,
	`  CC      fs/userfaultfd.o`,
	`  CC      net/ipv4/udp_diag.o`,
	`  CC      net/ipv6/xfrm6_input.o`,
	`  CC      drivers/firmware/efi/libstub/relocate.o`,
	`  AR      drivers/fake-pmem/built-in.a`,
	`  CC      drivers/firmware/efi/esrt.o`,
	`  CC      drivers/acpi/acpica/psparse.o`,
	`  CC      net/ipv4/raw_diag.o`,
	`  CC      drivers/firmware/efi/runtime-map.o`,
	`  CC      drivers/acpi/acpica/psscope.o`,
	`  CC      drivers/acpi/acpica/pstree.o`,
	`  CC      drivers/firmware/efi/libstub/vsprintf.o`,
	`  CC      drivers/acpi/acpica/psutils.o`,
	`  CC      drivers/acpi/property.o`,
	`  CC      net/ipv6/xfrm6_output.o`,
	`  HOSTCC  lib/gen_crc32table`,
	`  CC      lib/xxhash.o`,
	`  AR      drivers/hid/built-in.a`,
	`  CC      net/ipv6/xfrm6_protocol.o`,
	`  CC      net/ipv6/netfilter.o`,
	`  CC      fs/binfmt_elf.o`,
	`  CC      net/ipv6/proc.o`,
	`  CC      fs/mbcache.o`,
	`  CC      drivers/acpi/acpi_cmos_rtc.o`,
	`  CC      fs/posix_acl.o`,
	`  CC      drivers/firmware/efi/libstub/x86-stub.o`,
	`  CC      drivers/firmware/efi/runtime-wrappers.o`,
	`  CC      drivers/acpi/acpica/pswalk.o`,
	`  CC      drivers/acpi/acpica/psxface.o`,
	`  CC      drivers/acpi/acpica/rsaddr.o`,
	`  CC      drivers/firmware/efi/sysfb_efi.o`,
	`  CC      lib/genalloc.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_flex_pipe.o`,
	`  STUBCPY drivers/firmware/efi/libstub/efi-stub-helper.stub.o`,
	`  STUBCPY drivers/firmware/efi/libstub/file.stub.o`,
	`  CC      net/ipv6/sit.o`,
	`  STUBCPY drivers/firmware/efi/libstub/gop.stub.o`,
	`  CC      drivers/firmware/efi/earlycon.o`,
	`  CC      fs/drop_caches.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_flow.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_idc.o`,
	`  CC      lib/percpu_counter.o`,
	`  CC      lib/syscall.o`,
	`  STUBCPY drivers/firmware/efi/libstub/lib-cmdline.stub.o`,
	`  CC      net/ipv4/tcp_cubic.o`,
	`  STUBCPY drivers/firmware/efi/libstub/lib-ctype.stub.o`,
	`  CC      net/ipv6/addrconf_core.o`,
	`  CC      net/ipv6/exthdrs_core.o`,
	`  CC      drivers/net/ethernet/intel/i40e/i40e_xsk.o`,
	`  CC      drivers/acpi/acpica/rscalc.o`,
	`  CC      drivers/acpi/acpica/rscreate.o`,
	`  CC      drivers/acpi/x86/apple.o`,
	`  CC      net/ipv4/tcp_bpf.o`,
	`  CC      drivers/acpi/acpica/rsdumpinfo.o`,
	`  CC      net/ipv4/udp_bpf.o`,
	`  CC      net/ipv6/ip6_checksum.o`,
	`  CC      net/ipv6/ip6_icmp.o`,
	`  CC      drivers/acpi/acpica/rsinfo.o`,
	`  CC      net/ipv6/output_core.o`,
	`  CC      drivers/acpi/x86/utils.o`,
	`  CC      net/ipv6/protocol.o`,
	`  CC      net/ipv4/xfrm4_policy.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_devlink.o`,
	`  CC      net/ipv4/xfrm4_state.o`,
	`  CC      net/ipv6/ip6_offload.o`,
	`  CC      lib/errname.o`,
	`  CC      drivers/acpi/x86/s2idle.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_fw_update.o`,
	`  CC      drivers/acpi/acpi_lpat.o`,
	`  CC      lib/nlattr.o`,
	`  CC      lib/cpu_rmap.o`,
	`  CC      net/ipv6/tcpv6_offload.o`,
	`  CC      fs/fhandle.o`,
	`  CC      drivers/acpi/acpica/rsio.o`,
	`  CC      lib/dynamic_queue_limits.o`,
	`  CC      lib/glob.o`,
	`  CC      lib/strncpy_from_user.o`,
	`  CC      lib/strnlen_user.o`,
	`  CC      net/ipv6/exthdrs_offload.o`,
	`  CC      drivers/acpi/acpica/rsirq.o`,
	`  CC      drivers/acpi/acpica/rslist.o`,
	`  CC      drivers/acpi/acpica/rsmemory.o`,
	`  CC      lib/net_utils.o`,
	`  CC      net/ipv6/inet6_hashtables.o`,
	`  CC      drivers/acpi/acpi_lpit.o`,
	`  STUBCPY drivers/firmware/efi/libstub/mem.stub.o`,
	`  STUBCPY drivers/firmware/efi/libstub/pci.stub.o`,
	`  CC      lib/sg_pool.o`,
	`  STUBCPY drivers/firmware/efi/libstub/random.stub.o`,
	`  CC      net/ipv6/mcast_snoop.o`,
	`  CC      lib/ucs2_string.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_lag.o`,
	`  CC      drivers/acpi/prmt.o`,
	`  STUBCPY drivers/firmware/efi/libstub/randomalloc.stub.o`,
	`  CC      drivers/acpi/button.o`,
	`  STUBCPY drivers/firmware/efi/libstub/relocate.stub.o`,
	`  AR      drivers/firmware/efi/built-in.a`,
	`  STUBCPY drivers/firmware/efi/libstub/secureboot.stub.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_ethtool.o`,
	`  CC      drivers/acpi/fan.o`,
	`  STUBCPY drivers/firmware/efi/libstub/skip_spaces.stub.o`,
	`  STUBCPY drivers/firmware/efi/libstub/tpm.stub.o`,
	`  STUBCPY drivers/firmware/efi/libstub/vsprintf.stub.o`,
	`  STUBCPY drivers/firmware/efi/libstub/x86-stub.stub.o`,
	`  CC      drivers/acpi/acpica/rsmisc.o`,
	`  STUBCPY drivers/firmware/efi/libstub/alignedmem.stub.o`,
	`  CC      lib/sbitmap.o`,
	`  AR      drivers/firmware/efi/libstub/lib.a`,
	`  CC      net/ipv4/xfrm4_input.o`,
	`  AR      drivers/firmware/built-in.a`,
	`  CC      net/ipv4/xfrm4_output.o`,
	`  CC      drivers/acpi/acpica/rsserial.o`,
	`  CC      drivers/acpi/acpica/rsutils.o`,
	`  CC      drivers/acpi/processor_driver.o`,
	`  CC      drivers/acpi/acpica/rsxface.o`,
	`  AR      lib/lib.a`,
	`  GEN     lib/crc32table.h`,
	`  CC      lib/crc32.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_virtchnl_allowlist.o`,
	`  CC      drivers/acpi/acpica/tbdata.o`,
	`  CC      drivers/acpi/processor_idle.o`,
	`  CC      drivers/acpi/processor_throttling.o`,
	`  CC      drivers/acpi/acpica/tbfadt.o`,
	`  CC      drivers/acpi/processor_thermal.o`,
	`  CC      drivers/acpi/acpica/tbfind.o`,
	`  CC      drivers/acpi/acpica/tbinstal.o`,
	`  AR      fs/built-in.a`,
	`  CC      drivers/net/ethernet/intel/ice/ice_virtchnl_pf.o`,
	`  CC      drivers/acpi/acpica/tbprint.o`,
	`  CC      net/ipv4/xfrm4_protocol.o`,
	`  CC      drivers/acpi/acpica/tbutils.o`,
	`  CC      drivers/acpi/acpica/tbxface.o`,
	`  CC      drivers/acpi/acpica/tbxfload.o`,
	`  CC      drivers/acpi/acpica/tbxfroot.o`,
	`  CC      drivers/acpi/acpica/utaddress.o`,
	`  CC      drivers/acpi/acpica/utalloc.o`,
	`  CC      drivers/acpi/processor_perflib.o`,
	`  CC      drivers/acpi/acpica/utascii.o`,
	`  CC      drivers/acpi/container.o`,
	`  CC      drivers/acpi/thermal.o`,
	`  CC      drivers/acpi/acpica/utbuffer.o`,
	`  CC      drivers/acpi/acpi_memhotplug.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_sriov.o`,
	`  AR      drivers/md/built-in.a`,
	`  CC      drivers/acpi/acpica/utcopy.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_virtchnl_fdir.o`,
	`  CC      drivers/acpi/acpica/utexcep.o`,
	`  CC      drivers/acpi/ioapic.o`,
	`  CC      drivers/acpi/acpica/utdebug.o`,
	`  CC      drivers/acpi/acpica/utdecode.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_ptp.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_ptp_hw.o`,
	`  CC      drivers/acpi/acpica/utdelete.o`,
	`  CC      drivers/acpi/cppc_acpi.o`,
	`  CC      drivers/acpi/acpica/uterror.o`,
	`  CC      drivers/acpi/acpica/uteval.o`,
	`  CC      drivers/net/ethernet/intel/ice/ice_arfs.o`,
	`  CC      drivers/acpi/acpica/utglobal.o`,
	`  AR      lib/built-in.a`,
	`  CC      drivers/acpi/acpica/uthex.o`,
	`  CC      drivers/acpi/acpica/utids.o`,
	`  CC      drivers/acpi/spcr.o`,
	`  CC      drivers/acpi/acpica/utinit.o`,
	`  CC      drivers/acpi/acpica/utlock.o`,
	`  CC      drivers/acpi/acpica/utmath.o`,
	`  CC      drivers/acpi/acpica/utmisc.o`,
	`  CC      drivers/acpi/acpica/utmutex.o`,
	`  CC      drivers/acpi/acpica/utnonansi.o`,
	`  CC      drivers/acpi/acpica/utosi.o`,
	`  CC      drivers/acpi/acpica/utobject.o`,
	`  CC      drivers/acpi/acpica/utownerid.o`,
	`  CC      drivers/acpi/acpica/utpredef.o`,
	`  CC      drivers/acpi/acpica/utresdecode.o`,
	`  CC      drivers/acpi/acpica/utresrc.o`,
	`  CC      drivers/acpi/acpica/utstate.o`,
	`  CC      drivers/acpi/acpica/utstring.o`,
	`  CC      drivers/acpi/acpica/utstrsuppt.o`,
	`  CC      drivers/acpi/acpica/utstrtoul64.o`,
	`  CC      drivers/acpi/acpica/utxface.o`,
	`  CC      drivers/acpi/acpica/utxfinit.o`,
	`  CC      drivers/acpi/acpica/utxferror.o`,
	`  CC      drivers/acpi/acpica/utxfmutex.o`,
	`  AR      net/ipv6/built-in.a`,
	`  AR      drivers/net/ethernet/intel/i40e/built-in.a`,
	`  AR      drivers/acpi/acpica/built-in.a`,
	`  AR      drivers/acpi/built-in.a`,
	`  AR      net/ipv4/built-in.a`,
	`  AR      net/built-in.a`,
	`  AR      drivers/net/ethernet/intel/ice/built-in.a`,
	`  AR      drivers/net/ethernet/intel/built-in.a`,
	`  AR      drivers/net/ethernet/built-in.a`,
	`  AR      drivers/net/built-in.a`,
	`  AR      drivers/built-in.a`,
	`  GEN     .version`,
	`  CHK     include/generated/compile.h`,
	`  LD      vmlinux.o`,
	`  MODPOST vmlinux.symvers`,
	`  MODINFO modules.builtin.modinfo`,
	`  GEN     modules.builtin`,
	`  LD      vmlinux`,
	`  SORTTAB vmlinux`,
	`  SYSMAP  System.map`,
	`  CC      arch/x86/boot/a20.o`,
	`  AS      arch/x86/boot/bioscall.o`,
	`  CC      arch/x86/boot/cmdline.o`,
	`  AS      arch/x86/boot/copy.o`,
	`  HOSTCC  arch/x86/boot/mkcpustr`,
	`  CC      arch/x86/boot/cpuflags.o`,
	`  CC      arch/x86/boot/cpucheck.o`,
	`  CC      arch/x86/boot/early_serial_console.o`,
	`  CC      arch/x86/boot/edd.o`,
	`  CC      arch/x86/boot/main.o`,
	`  CC      arch/x86/boot/memory.o`,
	`  CC      arch/x86/boot/pm.o`,
	`  AS      arch/x86/boot/pmjump.o`,
	`  CC      arch/x86/boot/printf.o`,
	`  CC      arch/x86/boot/regs.o`,
	`  CC      arch/x86/boot/string.o`,
	`  CC      arch/x86/boot/tty.o`,
	`  CC      arch/x86/boot/video.o`,
	`  CC      arch/x86/boot/video-mode.o`,
	`  CC      arch/x86/boot/version.o`,
	`  CC      arch/x86/boot/video-vga.o`,
	`  CC      arch/x86/boot/video-vesa.o`,
	`  CC      arch/x86/boot/video-bios.o`,
	`  HOSTCC  arch/x86/boot/tools/build`,
	`  CPUSTR  arch/x86/boot/cpustr.h`,
	`  CC      arch/x86/boot/cpu.o`,
	`  LDS     arch/x86/boot/compressed/vmlinux.lds`,
	`  AS      arch/x86/boot/compressed/kernel_info.o`,
	`  AS      arch/x86/boot/compressed/head_64.o`,
	`  VOFFSET arch/x86/boot/compressed/../voffset.h`,
	`  CC      arch/x86/boot/compressed/string.o`,
	`  CC      arch/x86/boot/compressed/cmdline.o`,
	`  CC      arch/x86/boot/compressed/error.o`,
	`  OBJCOPY arch/x86/boot/compressed/vmlinux.bin`,
	`  HOSTCC  arch/x86/boot/compressed/mkpiggy`,
	`  RELOCS  arch/x86/boot/compressed/vmlinux.relocs`,
	`  CC      arch/x86/boot/compressed/cpuflags.o`,
	`  CC      arch/x86/boot/compressed/kaslr.o`,
	`  CC      arch/x86/boot/compressed/ident_map_64.o`,
	`  CC      arch/x86/boot/compressed/idt_64.o`,
	`  AS      arch/x86/boot/compressed/idt_handlers_64.o`,
	`  AS      arch/x86/boot/compressed/mem_encrypt.o`,
	`  CC      arch/x86/boot/compressed/pgtable_64.o`,
	`  CC      arch/x86/boot/compressed/acpi.o`,
	`  XZKERN  arch/x86/boot/compressed/vmlinux.bin.xz`,
	`  CC      arch/x86/boot/compressed/misc.o`,
	`  MKPIGGY arch/x86/boot/compressed/piggy.S`,
	`  AS      arch/x86/boot/compressed/piggy.o`,
	`  LD      arch/x86/boot/compressed/vmlinux`,
	`  ZOFFSET arch/x86/boot/zoffset.h`,
	`  OBJCOPY arch/x86/boot/vmlinux.bin`,
	`  AS      arch/x86/boot/header.o`,
	`  LD      arch/x86/boot/setup.elf`,
	`  OBJCOPY arch/x86/boot/setup.bin`,
	`  BUILD   arch/x86/boot/bzImage`,
	`Kernel: arch/x86/boot/bzImage is ready  (#1)`,
}
var goOutput = []string{
	`go: downloading system-transparency.org/stmgr v0.2.1`,
	`go: downloading github.com/diskfs/go-diskfs v1.3.0`,
	`go: downloading system-transparency.org/stboot v0.2.0`,
	`go: downloading github.com/gdamore/tcell/v2 v2.5.4`,
	`go: downloading github.com/rivo/tview v0.0.0-20230130130022-4a1b7a76c01c`,
	`go: downloading github.com/u-root/u-root v0.10.0`,
	`go: downloading git.glasklar.is/system-transparency/core/stauth v0.0.0-20230621112137-6e1b46d9f57b`,
	`go: downloading github.com/lucasb-eyer/go-colorful v1.2.0`,
	`go: downloading golang.org/x/term v0.0.0-20210927222741-03fcf44c2211`,
	`go: downloading github.com/mattn/go-runewidth v0.0.14`,
	`go: downloading golang.org/x/text v0.5.0`,
	`go: downloading golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f`,
	`go: downloading github.com/gdamore/encoding v1.0.0`,
	`go: downloading github.com/vishvananda/netlink v1.1.1-0.20211118161826-650dca95af54`,
	`go: downloading github.com/rivo/uniseg v0.4.2`,
	`go: downloading github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f`,
	`go: downloading github.com/rs/zerolog v1.29.1`,
	`go: downloading github.com/google/go-tpm v0.3.3`,
	`go: downloading github.com/spf13/cobra v1.1.3`,
	`go: downloading google.golang.org/protobuf v1.30.0`,
	`go: downloading github.com/golang/protobuf v1.5.0`,
	`go: downloading github.com/go-chi/chi/v5 v5.0.8`,
	`go: downloading github.com/diskfs/go-diskfs v1.3.0`,
	`go: downloading github.com/cavaliergopher/cpio v1.0.1`,
	`go: downloading github.com/saferwall/pe v1.4.2`,
	`go: downloading github.com/google/uuid v1.1.2`,
	`go: downloading github.com/mattn/go-colorable v0.1.12`,
	`go: downloading github.com/spf13/pflag v1.0.5`,
	`go: downloading golang.org/x/sys v0.8.0`,
	`go: downloading github.com/mattn/go-isatty v0.0.19`,
	`go: downloading gopkg.in/djherbis/times.v1 v1.2.0`,
	`go: downloading github.com/sirupsen/logrus v1.7.0`,
	`go: downloading github.com/google/uuid v1.3.0`,
	`go: downloading github.com/ulikunitz/xz v0.5.10`,
	`go: downloading github.com/pierrec/lz4 v2.3.0+incompatible`,
	`go: downloading github.com/pkg/xattr v0.4.1`,
	`go: downloading github.com/sirupsen/logrus v1.7.0`,
	`go: downloading gopkg.in/djherbis/times.v1 v1.2.0`,
	`go: downloading github.com/ulikunitz/xz v0.5.10`,
	`go: downloading github.com/pierrec/lz4 v2.3.0+incompatible`,
	`go: downloading github.com/pkg/xattr v0.4.1`,
	`go: downloading github.com/edsrzf/mmap-go v1.1.0`,
	`go: downloading go.mozilla.org/pkcs7 v0.0.0-20210826202110-33d05740a352`,
	`go: downloading golang.org/x/text v0.9.0`,
	`go: downloading github.com/pierrec/lz4/v4 v4.1.14`,
	`go: downloading github.com/klauspost/compress v1.10.6`,
	`go: downloading github.com/dustin/go-humanize v1.0.0`,
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
