# About

A small chore log.

# Release Process

1. Generate with prod tags.

    $ go generate -tags prod .

2. Commit and tag

    $ git tag -s -m "Releasing v0.0.0" v0.0.0

3. Goreleaser

    $ 