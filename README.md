# STMGR

[![CI](https://github.com/system-transparency/stmgr/actions/workflows/go.yml/badge.svg)](https://github.com/system-transparency/stmgr/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/system-transparency/stmgr/branch/main/graph/badge.svg)](https://codecov.io/gh/system-transparency/stmgr)

Next generation of stmanager in a standalone repo instead of being closely tied into [stboot](https://github.com/system-transparency/stboot).
This tool is used to create or sign [OS Packages](https://github.com/system-transparency/system-transparency#OS-Package), provision nodes for usage with stboot and has several other features to ease the usage of [system-transparency](https://github.com/system-transparency/system-transparency)

---

## Requirements

Go version 1.17 or higher.

---

## Installation instructions

Either run `go install github.com/system-transparency/stmgr@latest` or clone the repo and run `go build`.

---

## Usage

This tool requires to be invoked with a command and corresponding subcommand.
That means an example invokation looks like this: `stmgr ospkg create`.
In that example, `ospkg` is the command and `create` is the subcommand.
The best way to find out about all the commands and their subcommands is to run `stmgr -help` and follow the printed out instructions for further info.
