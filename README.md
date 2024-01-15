# STMGR

Next generation of stmanager in a standalone repo instead of being closely tied into [stboot](https://git.glasklar.is/system-transparency/core/stboot).
This tool is used to create or sign [OS Packages](https://git.glasklar.is/system-transparency/core/system-transparency#os-package), provision nodes for usage with stboot and has several other features to ease the usage of [system-transparency](https://git.glasklar.is/system-transparency/core/system-transparency)

---

## Requirements

Go version 1.17 or higher.

---

## Installation instructions

Either run `go install system-transparency.org/stmgr@latest` or clone the repo and run `go build`.

---

## Usage

This tool requires to be invoked with a command and corresponding subcommand.
That means an example invokation looks like this: `stmgr ospkg create`.
In that example, `ospkg` is the command and `create` is the subcommand.
The best way to find out about all the commands and their subcommands is to run `stmgr -help` and follow the printed out instructions for further info.
