package mkosi

type BuildArgs struct {
	RootPassword  string
	SSH           bool
	Hostname      string
	Cmdline       string
	Locale        string
	Keymap        string
	TimeZone      string
	ExtraPackages string
}

func Build(params *BuildArgs) error {
	var args []string
	if params.RootPassword != "" {
		args = append(args, "--root-password", params.RootPassword)
	}
	if params.Hostname != "" {
		args = append(args, "--hostname", params.Hostname)
	}
	if params.Cmdline != "" {
		args = append(args, "--kernel-command-line", params.Cmdline)
	}
	if params.Locale != "" {
		args = append(args, "--locale", params.Locale)
	}
	if params.Keymap != "" {
		args = append(args, "--keymap", params.Keymap)
	}
	if params.TimeZone != "" {
		args = append(args, "--timezone", params.TimeZone)
	}
	if params.ExtraPackages != "" {
		args = append(args, "-p", params.ExtraPackages)
	}
	if params.SSH {
		args = append(args, "--ssh")
	}
	return RunMkosi(args)
}
