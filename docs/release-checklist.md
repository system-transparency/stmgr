# Releases checklist

This document is intended for maintainers that make releases.

## Making a release

  - [ ] The README.md and RELEASES.md files are up-to-date
  - [ ] Ensure tests pass. Make sure the stmgr dependencies in other
        repos (stboot, stprov, system-transparency) are updated to
        an stmgr release candidate version, and that the integration
        tests in those repos are passing.
  - [ ] Ensure that the stimages/test/getting-started.sh script uses
        relevant versions of stmgr and other tools. Check that the
        script works in podman, and that the introductory build guide
        is consistent with this script.
  - [ ] After finalizing the release documentation (in particular the
        NEWS file), create a new signed tag. Usually, this means
        incrementing the third number for the most recent tag that was
        used during our interoperability tests.
  - [ ] Send announcement email

## RELEASES-file

  - [ ] What in the repository is released and supported
  - [ ] The overall release process is described, e.g., where are releases
    announced, how often do we make releases, what type of releases, etc.
  - [ ] The expectation we as maintainers have on users is described
  - [ ] The expectations users can have on us as maintainers is
    described, e.g., what we intend to (not) break in the future or any
    relevant pointers on how we ensure that things are "working".

## NEWS-file

  - [ ] The previous NEWS entry is for the previous release
  - [ ] Explain what changed
  - [ ] Detailed instructions on how to upgrade on breaking changes
  - [ ] List interoperable repositories and tools, specify commits or tags
  - [ ] List implemented reference specifications, specify commits or tags

## Announcement email template

```
The ST team is happy to announce a new release of the stmgr programm,
tag v0.X.X, which succeeds the previous release at tag v0.Y.Y.  The
source code for this release is available from the git repository:

  git clone -b v0.X.X https://git.glasklar.is/system-transparency/core/stmgr.git

Or install using Go's tooling:

  go install system-transparency.org/stmgr@v0.X.X

The expectations and intended use of the stmgr program is documented
in the repository's RELEASES file.  This RELEASES file also contains
more information concerning the overall release process, see:

  https://git.glasklar.is/system-transparency/core/stmgr/-/blob/main/RELEASES.md

Learn about what's new in a release from the repository's NEWS file.  An
excerpt from the latest NEWS-file entry is listed below for convenience.

If you find any bugs, please report them on the System Transparency
discuss list or open an issue on GitLab in the stmgr repository:

  https://lists.system-transparency.org/mailman3/postorius/lists/st-discuss.lists.system-transparency.org/
  https://git.glasklar.is/system-transparency/core/stmgr/-/issues

Cheers,
The ST team

<COPY-PASTE EXCERPT OF LATEST NEWS FILE ENTRY HERE>
```
