# stmgr manual

stmgr is a tool to create, sign and manage System Transparency boot
images, OS packages and related configuration files. It offers several
subcommands: `ospkg`, `keygen`, `uki`, `trustpolicy` and `hostconfig`.
This manual covers the most important features. All subcommands accept
the `-h` option which provides more information on defaults and optional
arguments.

## Specifying the keys to use

Keys are used for signing OS packages and certificates.  Use the `stmgr
ospkg sign` and `stmgr keygen certificate` commands.  Only Ed25519 keys
are supported.

When specifying a public key, provide a file with a public Ed25519 key
in either PKIX PEM format or OpenSSH single-line public key format.
Certificates use X.509 PEM format.

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

This command creates two files, a `.zip` archive and a `.json`
descriptor file. The `-url` option is required for network boot. It is
included in the descriptor file and specifies from where the OS package
should be downloaded at boot time.

The command for signing an OS package is

```
stmgr ospkg sign -cert FILENAME -key FILENAME -ospkg FILENAME
```

Both the archive `.zip` and the descriptor `.json` files are needed; the
`-ospkg` flag takes the name of either file. The certificate and a
corresponding signature are added to the descriptor file. The `-key`
option specifies the corresponding signing key, possibly with access via
ssh-agent, as described above.

To instead use a [Sigsum][] signature, first submit the OS package zip
file to a Sigsum log using the `sigsum-submit` tool. On success, this
produces a proof of logging, e.g., `ospkg.zip.proof`. This proof can
be attached to the OS package using

```
stmgr ospkg sigsum -cert FILENAME -proof FILENAME -ospkg FILENAME
```

Both the archive `.zip` and the descriptor `.json` files are needed; the
`-ospkg` flag takes the name of either file. The certificate and
Sigsum proof are added to the descriptor file. The certificate's
subject public key must correspond to the leaf signature and keyhash
in the Sigsum proof.

To verify the signatures of an OS package use:

```
stmgr ospkg verify [-trustPolicy DIRECTORY | -rootCerts FILENAME] -ospkg FILENAME
```

Both the archive `.zip` and the descriptor `.json` files are needed;
the `-ospkg` flag takes the name of either file. Signing root
certificate(s) are used to verify the signatures found in the
descriptor file.

The `-trustPolicy` flag is used to pass a [Trust policy][] directory,
where `trust_policy.json`, `ospkg_signing_root.pem`, and (optional)
`ospkg_trust_policy` will be read. If the number of valid signatures
is smaller than the `ospkg_signature_threshold` set in the JSON file,
then verification fails and the program exits with a non-zero status
code.

Alternatively, the `-rootCerts` flag can be used to only pass a file
containing signing root certificate(s). In this case *all* signatures
found in the descriptor file must be valid, otherwise verification
fails and the program exits with a non-zero status code. This means
that if there are no signatures present, verification will always
succeed.

[OS package]: https://git.glasklar.is/system-transparency/project/docs/-/blob/v0.5.2/content/docs/reference/os_package.md
[Trust policy]: https://git.glasklar.is/system-transparency/project/docs/-/blob/v0.5.2/content/docs/reference/trust_policy.md
[Sigsum]: https://www.sigsum.org

## The stmgr keygen command

There's only one subcommand, which is used to create certificates, and
optionally to generate a corresponding key-pair. There are defaults for
the file name arguments, see `stmgr keygen certificate` for details. To
create a self-signed root certificate:

```
stmgr keygen certificate -isCA [-rootKey FILENAME] [-certOut FILENAME] [-keyOut FILENAME]
```

The `-rootKey` option specifies a signing key to use, possibly with
access via ssh-agent, as described above. If not specified, a new
key-pair is generated, and the private key is written to the file
specified with `-keyOut`.

To create a leaf signing certificate:

```
stmgr keygen certificate [-rootCert FILENAME] [-rootKey FILENAME] [-certOut FILENAME] [-keyOut FILENAME] [-leafKey FILENAME]
```

The `-rootCert` and `-rootKey` specify the CA root and corresponding
signing key. `-leafKey` specifies the public key to certify, if not
provided, a new key-pair is generated, and the private key is written to
the file specified with `-keyOut`.

## The stmgr uki command

This command is used to create a Unified Kernel Image (UKI) that is
bootable directly by UEFI firmware. Essentially, a UKI is a kernel,
an initramfs and a command line packaged into a UEFI PE executable. This
command is used for packaging the stboot executable, trust policy, and
other related files. Inputs are similar to those of `ospkg create`, but
for a different purpose and with a different output format.

```
stmgr uki create -cmdline STRING [-format iso|uki] -initramfs FILENAME -kernel FILENAME -out FILENAME
```

The default output format is `iso`, and means that the UKI is wrapped in
a bootable CDROM image. To get just the UKI, pass `-format uki`.

The UKI (a PE executable) can optionally be signed for Secure Boot.  Use
the flags `-signkey` and `-signcert` to set the file names to a private
key and its corresponding certificate, both in PEM format.  Because
Secure Boot does not support Ed25519, RSA keys are required here.

If you Secure Boot sign the created UKI without stmgr, or if you prefer to
create a UKI and format it as an ISO in separate steps, then use the `to-iso`
subcommand:

```
stmgr uki to-iso -in FILENAME [-out FILENAME]
```

The output filename defaults to the input filename with a `.iso` suffix.

## The stmgr trustpolicy and host config commands

These commands can be used to validate syntax and contents of [host
config][] and [trust policy][] configuration files, respectively. They
take the contents of the configuration (not a filename!) on the command
line.

```
stmgr hostconfig check JSON-DATA
stmgr trustpolicy check JSON-DATA
```

[trust policy]: https://git.glasklar.is/system-transparency/project/docs/-/blob/v0.5.2/content/docs/reference/trust_policy.md
[host config]: https://git.glasklar.is/system-transparency/project/docs/-/blob/v0.5.2/content/docs/reference/host_configuration.md
