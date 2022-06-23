<div align="right">

![golangci-lint](https://github.com/sheepla/pingu/actions/workflows/golangci-lint.yml/badge.svg)
![release](https://github.com/sheepla/pingu/actions/workflows/release.yml/badge.svg)

</div>


<div align="center">

# 🐧 pingu

`ping` command implementation but with pingu ascii art


![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)
[![Latest Release](https://img.shields.io/github/v/release/sheepla/pingu?style=flat-square)](https://github.com/sheepla/pingu/releases/latest)

</div>

<div align="center">

<img src="https://user-images.githubusercontent.com/62412884/173192079-a6229925-2575-476c-9bb1-b5f0c6e2dbfc.png" width="70%" alt="screenshot" />

</div>

## Usage

Simply specify the target host name or IP address in the first argument e.g. `pingu github.com` or `pingu 13.114.40.48`.
You can change the number of transmissions by specifying the `-c` option.

```
Usage:
  pingu [OPTIONS] HOST

`ping` command but with pingu

Application Options:
  -V, --version  Show version
  -c, --count=   Stop after <count> replies (default: 20)

Help Options:
  -h, --help     Show this help message
```

## Installation

You can download executable binaries from the latest release page.

> [![Latest Release](https://img.shields.io/github/v/release/sheepla/pingu?style=flat-square)](https://github.com/sheepla/pingu/releases/latest)

To build from source, clone this repository then run `make build` or `go install`. Develo*ping* on `go1.18.3 linux/amd64`.

Another way to install it if you have go in your machine just
```
go install github.com/sheepla/pingu@latest
```

## LICENSE

[MIT](./LICENSE)

## Author

[Sheepla](https://github.com/sheepla)

