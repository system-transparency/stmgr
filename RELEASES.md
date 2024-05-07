# Releases of stmgr

## What is being released?

The following program is released and supported:

  - `./stmgr`

New releases are announced on the System Transparency [announce list][].
What changed in each release is documented in a [NEWS file](./NEWS). The
NEWS file also specifies which other System Transparency components are
known to be interoperable, as well as which reference specifications are
being implemented.

Note that a release is simply a signed git tag specified on our mailing
list, accessed from the [stmgr repository][]. To verify tag signatures,
get the `allowed-ST-release-signers` file published at [signing keys][],
and verify the tag `vX.Y.Z` using the command
```
git -c gpg.format=ssh -c gpg.ssh.allowedSignersFile=allowed-ST-release-signers \
  tag --verify vX.Y.Z
```
If desired, the config settings above can be stored more permanently using
`git config`.

The stmgr Go module is **not** considered stable before a v1.0.0 release.  By
the terms of the LICENSE file you are free to use this code "as is" in almost
any way you like, but for now, we support its use _only_ via the above program.
We don't aim to provide any backwards-compatibility for internal interfaces.

[announce list]: https://lists.system-transparency.org/mailman3/postorius/lists/st-announce.lists.system-transparency.org/
[stmgr repository]: https://git.glasklar.is/system-transparency/core/stmgr/
[signing keys]: https://www.system-transparency.org/keys

## What release cycle is used?

We make feature releases when something new is ready.  As a rule of thumb,
feature releases will not happen more often than once per month.

In case critical bugs are discovered, we intend to provide bug-fix-only updates
for the latest release in a timely manner.  Backporting bug-fixes to older
releases than the latest one will be considered on a case-by-case basis.  Such
consideration is most likely if the latest feature release is very recent and
upgrading to it is particularly disruptive due to the changes that it brings.

## Upgrading

We strive to make stmgr upgrades easy and well-documented. Any complications,
e.g., chages to command line flags, will be clearly outlined in the [NEWS
file](./NEWS). Pay close attention to the "Incompatible changes" section before
upgrading to a new version.

## Expected changes in upcoming releases

  - The command line interface is expected to be overhauled.
  - Transition to new signature format that's compatible with Sigsum,
    likely with changes to the OS package format.
  - Any changes to the System Transparency reference specifications will be
implemented.  This could for example affect the format of configuration
files such as host configuration or trust policy.

