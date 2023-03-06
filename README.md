# Finance Market Data Downloader

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/regel/wsb)](https://goreportcard.com/report/github.com/regel/wsb)
[![Build](https://github.com/regel/wsb/actions/workflows/build.yaml/badge.svg)](https://github.com/regel/wsb/actions/workflows/build.yaml)
[![codecov](https://codecov.io/github/regel/wsb/coverage.svg)](https://codecov.io/gh/regel/wsb)

`wsb` is the the tool for downloading Yahoo! finance market data.

This go package aims to provide a reliable, threaded, and idiomatic way to download historical market data from Yahoo! Finance API and other finance data sources.

Although the Yahoo Finance API has officially been closed down, it does still work and it provides a free access to a vast number of stocks.

>Warning - The Yahoo Finance API could be removed or shut down at any point. You use this package at your own risk.

Other finance data sources supported in this package:

- [IEX Cloud](https://iexcloud.io/docs/api/): IEX Cloud is a platform that makes financial data and services accessible to everyone. There is a free tier for use during initial API exploration and application development. During registration you will receive security tokens required to access this API
- [CoinGecko](https://www.coingecko.com/): CoinGecko provides a comprehensive cryptocurrency API. See Crypto Data API Plans on their web site for more information. At the time of this writting, the free plan is limited at 50 calls/minute (varies)

## Backers :dart: :heart_eyes:

Thank you to all our backers! 🙏 [[Become a backer](https://opencollective.com/gh-regel#backer)]

<a href="https://opencollective.com/gh-regel#backers" target="_blank"><img src="https://opencollective.com/gh-regel/backers.svg?width=890"></a>

## Sponsors :whale:

Support this project by becoming a sponsor. Your logo will show up here with a
link to your website. [[Become a
sponsor](https://opencollective.com/gh-regel#sponsor)]

## Installation

### Binary Distribution

Download the release distribution for your OS from the [Releases](https://github.com/regel/wsb/releases) page.

Unpack the `wsb` binary, add it to your PATH, and you are good to go!

### Docker Image

A Docker image is available at `https://github.com/regel/wsb/pkgs/container/wsb` with list of
available tags [here](https://github.com/regel/wsb/pkgs/container/wsb/versions).

### Homebrew

```console
$ brew tap regel/wsb
$ brew install wsb
```

## Usage

See documentation for individual commands:

* [wsb](doc/wsb.md)
* [wsb version](doc/wsb_version.md)
* [wsb chart](doc/wsb_chart.md)
* [wsb hold](doc/wsb_hold.md)

## Configuration

`wsb` is a command-line application.

All command-line flags can also be set via environment variables or config file.
Environment variables must be prefixed with `WSB_`.
Underscores must be used instead of hyphens.

CLI flags, environment variables, and a config file can be mixed.
The following order of precedence applies:

1. CLI flags
1. Environment variables
1. Config file

### Examples

Pulling historic price data for Bitcoin and [Cardano](https://cardano.org/) cryptocurrencies:

```
wsb chart --provider coingecko --tickers bitcoin,cardano --from 2021-02-01 --to 2021-04-01
```

Output:

```
+---------------------+----------+----------+----------+----------+--------+
|        DATE         |   OPEN   |   HIGH   |   LOW    |  CLOSE   | VOLUME |
+---------------------+----------+----------+----------+----------+--------+
| 2021-02-03T00:00:00 | 33064.79 | 35485.99 | 33064.79 | 35485.99 |      0 |
| 2021-02-07T00:00:00 | 37494.72 | 39279.41 | 36816.51 | 39279.41 |      0 |
| 2021-02-11T00:00:00 | 38833.34 | 46569.56 | 38833.34 | 44848.69 |      0 |
| 2021-02-15T00:00:00 | 47815.96 | 48607.87 | 46941.29 | 48607.87 |      0 |
| 2021-02-19T00:00:00 | 47898.49 | 52143.68 | 47898.49 | 51733.08 |      0 |
| 2021-02-23T00:00:00 | 56038.73 | 57669.30 | 54410.86 | 54410.86 |      0 |
| 2021-02-27T00:00:00 | 48691.89 | 49849.38 | 46551.49 | 46551.49 |      0 |
| 2021-02-28T00:00:00 | 46653.53 | 49787.34 | 44970.16 | 48532.24 |      0 |
| 2021-03-07T00:00:00 | 50577.46 | 50577.46 | 48727.45 | 49019.37 |      0 |
| 2021-03-11T00:00:00 | 51313.09 | 56020.49 | 51313.09 | 56020.49 |      0 |
| 2021-03-15T00:00:00 | 57788.87 | 61315.20 | 57353.86 | 59428.97 |      0 |
| 2021-03-19T00:00:00 | 55805.33 | 59014.93 | 55805.33 | 57922.41 |      0 |
| 2021-03-23T00:00:00 | 58243.27 | 58376.16 | 54370.14 | 54370.14 |      0 |
| 2021-03-27T00:00:00 | 54584.87 | 55033.10 | 51416.91 | 55033.10 |      0 |
| 2021-03-31T00:00:00 | 55832.42 | 58668.63 | 55728.10 | 58668.63 |      0 |
+---------------------+----------+----------+----------+----------+--------+
History of 'bitcoin'.
+---------------------+------+------+------+-------+--------+
|        DATE         | OPEN | HIGH | LOW  | CLOSE | VOLUME |
+---------------------+------+------+------+-------+--------+
| 2021-02-03T00:00:00 | 0.35 | 0.43 | 0.35 |  0.43 |      0 |
| 2021-02-07T00:00:00 | 0.44 | 0.63 | 0.44 |  0.63 |      0 |
| 2021-02-11T00:00:00 | 0.66 | 0.92 | 0.66 |  0.92 |      0 |
| 2021-02-15T00:00:00 | 0.93 | 0.93 | 0.85 |  0.85 |      0 |
| 2021-02-19T00:00:00 | 0.86 | 0.91 | 0.86 |  0.91 |      0 |
| 2021-02-23T00:00:00 | 0.93 | 1.13 | 0.93 |  1.10 |      0 |
| 2021-02-27T00:00:00 | 0.96 | 1.25 | 0.96 |  1.25 |      0 |
| 2021-02-28T00:00:00 | 1.34 | 1.34 | 1.23 |  1.23 |      0 |
| 2021-03-07T00:00:00 | 1.22 | 1.22 | 1.12 |  1.13 |      0 |
| 2021-03-11T00:00:00 | 1.14 | 1.19 | 1.12 |  1.14 |      0 |
| 2021-03-15T00:00:00 | 1.13 | 1.13 | 1.04 |  1.06 |      0 |
| 2021-03-19T00:00:00 | 1.03 | 1.38 | 1.03 |  1.24 |      0 |
| 2021-03-23T00:00:00 | 1.30 | 1.30 | 1.11 |  1.11 |      0 |
| 2021-03-27T00:00:00 | 1.13 | 1.21 | 1.06 |  1.21 |      0 |
| 2021-03-31T00:00:00 | 1.18 | 1.21 | 1.18 |  1.21 |      0 |
+---------------------+------+------+------+-------+--------+
History of 'cardano'.
```

The following example show various way of configuring the same thing:

#### CLI

    wsb chart --tickers AAPL,GME --from "2021-01-01"

#### Environment Variables

    export WSB_TICKERS=AAPL,GME

    wsb chart --from "2021-01-01"

#### Config File

`config.yaml`:

```yaml
tickers:
  - AAPL
  - GME
```

#### Config Usage

    wsb chart --config config.yaml --from "2021-01-01"


`wsb` supports any format [Viper](https://github.com/spf13/viper) can read, i. e. JSON, TOML, YAML, HCL, and Java properties files.

Notice that if no config file is specified, then `wsb.yaml` (or any of the supported formats) is loaded from the current directory, `$HOME/.wsb`, or `/etc/wsb`, in that order, if found.

## Building from Source

`wsb` is built using Go 1.13 or higher.

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

