[![Go](https://github.com/sha1n/clib/actions/workflows/go.yml/badge.svg)](https://github.com/sha1n/clib/actions/workflows/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sha1n/clib)
[![Go Report Card](https://goreportcard.com/badge/sha1n/clib)](https://goreportcard.com/report/sha1n/clib) 
[![Release](https://img.shields.io/github/release/sha1n/clib.svg?style=flat-square)](https://github.com/sha1n/clib/releases)
![GitHub all releases](https://img.shields.io/github/downloads/sha1n/clib/total)
[![Release Drafter](https://github.com/sha1n/clib/actions/workflows/release-drafter.yml/badge.svg)](https://github.com/sha1n/clib/actions/workflows/release-drafter.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# clib

## Before anything else
```bash
git clone git@github.com:<owner>/<repo>.git
cd <repo>
./init.sh <owner> <repo>
```

## Features

- Init script for setup
- Makefile
- Workflows
  - Go build + coverage - [go.yml](/.github/workflows/go.yml)
  - Go report card - [go-report-card.yml](/.github/workflows/go-report-card.yml)
  - Release Drafter - [release-drafter.yml](/.github/workflows/release-drafter.yml)
  - Dependabot App - [dependabot.yml](/.github/dependabot.yml)
- Jekyll site setup with the [Cayman](https://github.com/pages-themes/cayman) theme (and some color overrides)
- .travis.yml for Go
