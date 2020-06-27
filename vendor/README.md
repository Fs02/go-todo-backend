# vendor

It's recommended to commit package dependency for a non library/package project, Because golang uses a decentralized dependency management, so there's no telling when a repo might be deleted or renamed, so it's better to have it vendored in your package.

when using go mod, you can use `go mod vendor` command to update the vendor folder.
