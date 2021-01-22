(Adapted from [Ginkgo's RELEASING.md](https://github.com/onsi/ginkgo/blob/v1.5.0/RELEASING.md) Copyright (c) 2013-2014 Onsi Fakhouri)

A `satcon-client-go` release is a tagged git sha and a GitHub release. Please follow best practices for naming/numbering versions as given in https://blog.golang.org/publishing-go-modules.

To cut a release:

1. Run `go mod tidy` and run all tests (e.g. `ginkgo ./...` from top directory). Commit any changes.
1. Ensure CHANGELOG.md is up to date.
    - Check dependency changes since last release: `git diff vX.X.X HEAD -- go.mod`
    - Use `git log --pretty=format:'- %s [%h]' HEAD...vX.X.X` to list all the commits since the last release
    - Categorize the changes into
        - Breaking Changes (requires a major version)
        - New Features (minor version)
        - Fixes (fix version)
        - Maintenance (which in general should not be mentioned in `CHANGELOG.md` as they have no user impact)
1. Create a commit with the version number as the commit message (e.g. `v1.3.0`)
1. Tag the commit with the version number as the tag name (e.g. `v1.3.0`)
1. Push the commit and tag to GitHub
1. Create a new [GitHub release](https://help.github.com/articles/creating-releases/) with the version number as the tag  (e.g. `v1.3.0`).  List the key changes in the release notes.