# stmgr manual

stmgr is a tool to create, sign and manage System Transparency boot
images, OS packages and related configuration files. It offers several
subcommands: `ospkg`, `keygen`, `uki`, `trustpolicy` and `hostconfig`.
This manual covers the most important features. All subcommands accept
the `-h` option which provides more information on defaults and optional
arguments.

## Specifying the keys to use

Public keys, and corresponding private keys, are used both for signing
OS packages (`stmgr ospkg sign` command) and for signing certificates
(`stmgr keygen certificate` command). Only Ed25519 keys are supported.
When specifying a public key, you provide a file with a public Ed25519
key in either PKIX PEM format, or OpenSSH single-line public key format.
Certificates use x.509 PEM format.

For private keys, you can specify an unencrypted private key in either
PKIX or OpenSSH format. The recommended way, however, is to pass a
*public* key in OpenSSH format, in which case stmgr will use ssh-agent,
based on the `$SSH_AUTH_SOCK` environment variable, to access the
corresponding signing key.

## The stmgr ospkg command

System transparency OS packages are defined by the [OS package][]
specification. The `stmgr ospkg` command has two subcommands operating
on OS packages: `create` and `sign`. To create an OS package use

```
stmgr ospkg create [OPTIONS] -cmdline STRING -initramfs FILENAME -kernel FILENAME -out FILENAME [-url OSPKG-URL]
```

This command creates two files, a .zip archive and a .json descriptor
file. The `-url` option is required for network boot. It is included in
the descriptor file and specifies from where the OS package should be
downloaded at boot time.

To command for signing an OS package is

```
stmgr ospkg sign -cert X509-CERT -key PRIVATE-KEY -ospkg FILE
```

Both the archive .zip and the descriptor .json files are needed, the
`-ospkg` flag takes the name of either file. The certificate and a
corresponding signature are added to the descriptor file. The `-key`
option specifies the corresponding signing key, possibly with access via
ssh-agent, as described above.

[OS package]: https://docs.system-transparency.org/docs/reference/data-structures/os_package/

## The stmgr keygen command

There's only one subcommand, which is used to create certificates, and
optionally generate a corresponding keypair. There are defaults for the
file name arguments, see `stmgr keygen certificate` for details. To
create a self-signed root certificate:

```
stmgr keygen certificate -isCA [-rootKey FILENAME] [-certOut FILENAME] [-keyOut FILENAME]
```

The `-rootKey` option specifies a signing key to use, possibly with
access via ssh-agent, as described above. If not specified, a new
keypair is generated, and the private key is written to the file
specified with `-keyOut`.

To create a leaf signing certificate:

```
stmgr keygen certificate [-rootCert FILENAME] [-rootKey FILENAME] [-certOut FILENAME] [-keyOut FILENAME] [-leafKey FILENAME]
```

The `-rootCert` and `-rootKey` specify the CA root and corresponding
signing key. `-leafKey` specifies the public key to certify, if not
provided, a new keypair is generated, and the private key is written to
the file specified with `-keyOut`.

## The stmgr uki command

This command is used to create a Unified Kernel Image (UKI) that is
bootable directly by UEFI firmware. Essentially, a UKI is kernel,
initramfs and command line packages into a UEFI PE executable. This
command is used for packaging the stboot executable, trust policy, and
other related files. Inputs are similar to those for `ospkg create`, but
for a different purpose and with different output format.

```
stmgr uki create -cmdline STRING [-format iso|uki] -initramfs FILENAME -kernel FILENAME -out FILENAME
```

The default output format is `iso`, and means that the UKI is wrapped in
a bootable CDROM image. To get just the uki, pass `-format uki`.

The UKI, a PE executable, can optionally be signed for Secure Boot.  Use
the flags `-signkey` and `-signcert` to set the file names to a private
key and its corresponding certificate, both in PEM format.

## The stmgr trustpolicy and host config commands

These commands can be used to validate syntax and contents of [host
config][] and [trust policy][] configuration files, respectively. They
take the contents of the configuration (not a filename!) on the command
line.

```
stmgr hostconfig check JSON-DATA
stmgr trustpolicy check JSON-DATA
```

[trust policy]: https://docs.system-transparency.org/docs/reference/data-structures/trust_policy/
[host config]: https://docs.system-transparency.org/docs/reference/data-structures/host_configuration/
