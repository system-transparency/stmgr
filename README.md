# ST manager

ST manager can be used to create and sign [OS Packages][] and more to ease the usage of [System Transparency][].

[OS Packages](https://git.glasklar.is/system-transparency/core/system-transparency#os-package)
[System Transparency](https://git.glasklar.is/system-transparency/core/system-transparency)

---

## Requirements

Go version 1.17 or higher.

---

## Installation instructions

Either run `go install system-transparency.org/stmgr@latest` or clone the repo and run `go build`.

---

## Usage

stmgr is invoked with a command and a corresponding subcommand, for example `stmgr ospkg create`.
In this example, `ospkg` is the command and `create` is the subcommand.
The best way to find out about all the commands and their subcommands is to run `stmgr -help` and follow the instructions for further info.
