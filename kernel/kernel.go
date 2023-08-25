package kernel

import (
	"context"
	_ "embed"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile"
)

//go:embed linux.yml
var linuxTaskfileYaml []byte

//go:embed linuxboot.defconfig
var linuxbootDefconfig []byte

func runTasks(ctx context.Context, vars map[string]string, tasks ...string) error {
	tmpdir, err := os.MkdirTemp("", "stmgr-task-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	err = ioutil.WriteFile(path.Join(tmpdir, "Taskfile.yml"), linuxTaskfileYaml, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(tmpdir, "linuxboot.defconfig"), linuxbootDefconfig, 0644)
	if err != nil {
		return err
	}

	exec := task.Executor{
		Taskfile:    nil,
		Dir:         tmpdir,
		TempDir:     tmpdir,
		Entrypoint:  "Taskfile.yml",
		Force:       false,
		ForceAll:    false,
		Watch:       false,
		Verbose:     false,
		Silent:      false,
		AssumeYes:   false,
		Dry:         false,
		Summary:     false,
		Parallel:    false,
		Color:       false,
		Concurrency: 0,
		Interval:    0,
		AssumesTerm: false,
		Stdin:       nil,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Logger:      nil,
		Compiler:    nil,
		Output:      nil,
		OutputStyle: taskfile.Output{},
		TaskSorter:  nil,
	}
	err = exec.Setup()
	if err != nil {
		return err
	}
	for k, v := range vars {
		exec.Taskfile.Vars.Set(k, taskfile.Var{Static: v})
	}

	for _, taskName := range tasks {
		call := taskfile.Call{
			Task:   taskName,
			Direct: true,
		}
		err = exec.RunTask(ctx, call)
		if err != nil {
			return err
		}
	}

	return nil
}

func Build(ctx context.Context, version string, name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	vars := map[string]string{}
	vars["KERNEL_VERSION"] = "5.10.1"
	vars["KERNEL_DIR"] = path.Join(home, ".local/stmgr/build/linux-5.10.1")
	vars["LINUXBOOT_KERNEL_CONFIG"] = "linuxboot.defconfig"

	return runTasks(ctx, vars, "fetch", "unpack", "config", "build")
}
