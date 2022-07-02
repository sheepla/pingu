<div align="right">

![golangci-lint](https://github.com/sheepla/pingu/actions/workflows/golangci-lint.yml/badge.svg)
![release](https://github.com/sheepla/pingu/actions/workflows/release.yml/badge.svg)

</div>


<div align="center">

# 🐧 pingu

`ping` command implementation in Go but with colorful output and pingu ascii art


![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)
[![Latest Release](https://img.shields.io/github/v/release/sheepla/pingu?style=flat-square)](https://github.com/sheepla/pingu/releases/latest)

</div>

## Features

- [x] Colorful and fun output.
- [x] Cross-platform support: Windows, macOS and Linux (also WSL)
- [x] It works with a single executable file, so it can be installed easily.
- [x] Surpports IPv4 and IPv6.

## Usage

Simply specify the target host name or IP address in the first argument e.g. `pingu github.com` or `pingu 13.114.40.48`.
You can change the number of transmissions by specifying the `-c` option.

```
Usage:
  pingu [OPTIONS] HOST

`ping` command but with pingu

Application Options:
  -c, --count=     Stop after <count> replies (default: 20)
  -P, --privilege  Enable privileged mode
  -V, --version    Show version

Help Options:
  -h, --help       Show this help message
```

## Installation

### Download executable binaries

You can download executable binaries from the latest release page.

> [![Latest Release](https://img.shields.io/github/v/release/sheepla/pingu?style=flat-square)](https://github.com/sheepla/pingu/releases/latest)

### Build from source

To build from source, clone this repository then run `make build` or `go install`. Develo*ping* on `go1.18.3 linux/amd64`.

Another way to install it if you have go in your machine just:

```sh
go install github.com/sheepla/pingu@latest
```

### WSL support

To run on WSL, give `CAP_NET_RAW` capability with the following command and start it with the `-P`, `--privilege` option.

```
sudo setcap cap_net_raw=+ep /path/to/pingu
pingu -P github.com
```

*c.f.*

> [Capability - Arch Linux Wiki](https://wiki.archlinux.org/title/Capabilities)

## LICENSE

[MIT](./LICENSE)

## Author

[Sheepla](https://github.com/sheepla)

