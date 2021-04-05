# Finance Market Data Downloader

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/regel/tinkerbell)](https://goreportcard.com/report/github.com/regel/tinkerbell)
[![Build](https://github.com/regel/tinkerbell/actions/workflows/build.yaml/badge.svg)](https://github.com/regel/tinkerbell/actions/workflows/build.yaml)
[![codecov](https://codecov.io/github/regel/tinkerbell/coverage.svg)](https://codecov.io/gh/regel/tinkerbell)

`tb` (tinkerbell) is the the tool for downloading Yahoo! finance market data.

This go package aims to provide a reliable, threaded, and idiomatic way to download historical market data from Yahoo! Finance API and other finance data sources.

Although the Yahoo Finance API has officially been closed down, it does still work and it provides a free access to a vast number of stocks.

>Warning - The Yahoo Finance API could be removed or shut down at any point. You use this package at your own risk.

Other finance data sources supported in this package:

- [IEX Cloud](https://iexcloud.io/docs/api/): IEX Cloud is a platform that makes financial data and services accessible to everyone. There is a free tier for use during initial API exploration and application development. During registration you will receive security tokens required to access this API

## Installation

### Binary Distribution

Download the release distribution for your OS from the Releases page:

https://github.com/regel/tinkerbell/releases

Unpack the `tb` binary, add it to your PATH, and you are good to go!

### Docker Image

A Docker image is available at `https://hub.docker.com/r/regel/tinkerbell` with list of
available tags [here](https://hub.docker.com/r/regel/tinkerbell/tags).

### Homebrew

```console
$ brew tap regel/tinkerbell
$ brew install tinkerbell
```

## Usage

See documentation for individual commands:

* [tb](doc/tb.md)
* [tb version](doc/tb_version.md)
* [tb chart](doc/tb_chart.md)
* [tb hold](doc/tb_hold.md)

## Configuration

`tb` is a command-line application.

All command-line flags can also be set via environment variables or config file.
Environment variables must be prefixed with `TB_`.
Underscores must be used instead of hyphens.

CLI flags, environment variables, and a config file can be mixed.
The following order of precedence applies:

1. CLI flags
1. Environment variables
1. Config file

### Examples

The following example show various way of configuring the same thing:

#### CLI

    tb chart --tickers AAPL,GME --from "2021-01-01"

#### Environment Variables

    export TB_TICKERS=AAPL,GME

    tb chart --from "2021-01-01"

#### Config File

`config.yaml`:

```yaml
tickers:
  - AAPL
  - GME
```

#### Config Usage

    tb chart --config config.yaml --from "2021-01-01"


`tb` supports any format [Viper](https://github.com/spf13/viper) can read, i. e. JSON, TOML, YAML, HCL, and Java properties files.

Notice that if no config file is specified, then `tb.yaml` (or any of the supported formats) is loaded from the current directory, `$HOME/.tb`, or `/etc/tb`, in that order, if found.

## Building from Source

`tb` is built using Go 1.13 or higher.

It uses [Goreleaser](https://goreleaser.com/) under the covers.

To build:

```
goreleaser build --rm-dist --snapshot
```

### Known issues

On MacOS the `boring` build fails. Comment the lines in `.goreleaser.yml` to disable this build.

```
runtime/cgo(__TEXT/__text): relocation target x_cgo_inittls not defined
```

