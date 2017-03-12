asciist
=======

[![build status](https://img.shields.io/travis/enaeseth/asciist.svg)](https://travis-ci.org/enaeseth/asciist)
[![Go Report Card](https://goreportcard.com/badge/github.com/enaeseth/asciist)](https://goreportcard.com/report/github.com/enaeseth/asciist)
[![GoDoc](https://godoc.org/github.com/enaeseth/asciist?status.svg)](https://godoc.org/github.com/enaeseth/asciist)

*a toy service to turn images into ASCII art*

![asciist action shot](./asciist.gif)

Usage
-----

### Install

asciist is a Go package, and installation is simple with the Go toolchain:

```sh
go install github.com/enaeseth/asciist
```

### Run

For detailed usage information, run `asciist help`.

#### Service

To start the asciist service, run `asciist serve`:

```
usage: asciist serve [<flags>]

Run the service

Flags:
  --debug       Use debugging mode
  --host=HOST   Interface to listen on
  --port=27244  Port to listen on
```

#### Client

Once the service is running, you can use the built-in client to convert images with `asciist convert`. `convert` accepts a file to convert as an argument, or reads the image from stdin if no argument is provided.

By default, `convert` will output ASCII art that is the same width as your terminal. You can use a specific width with the `-w/--width` option.

```sh
$ curl -sL https://git.io/vyKUO | asciist convert -w 20
+++****####%%%@@@@@@
==++++***####%%%%@@@
====++++***####%%%%@
--====++++***####%%%
:---====++++***####%
:::---====++++***###
.::::---====++++***#
...::::---====++++**
 ....::::---====++++
    ...::::---====++
```

#### Testing

To convert an image to ASCII art locally (without using the service), run `asciist test` instead of `convert`.

### Unit Tests

asciist includes tests for its image conversion and service packages, and they can be run with `go test ./...`:

```sh
$ (cd "$GOPATH/src/github.com/enaeseth/asciist" && go test ./...)
?   github.com/enaeseth/asciist            [no test files]
?   github.com/enaeseth/asciist/client     [no test files]
ok  github.com/enaeseth/asciist/convert    0.043s
?   github.com/enaeseth/asciist/fixture    [no test files]
ok  github.com/enaeseth/asciist/service    0.062s
```

Protocol
--------

asciist accepts POST requests to the root URL (`/`). It expects JSON requests and produces JSON responses.

### Request

```json
{
  "width": 80,
  "image": "[base64-encoded image file]"
}
```

### Response

#### 200 OK

```json
{
  "art": "[ASCII art, no trailing newline]"
}
```

#### 400 Bad Request

```json
{
  "error": "[brief error message]"
}
```

Notes
-----

### Dependencies

asciist depends on these third-party Go packages:

#### [gin][gin]

Gin is a fast and popular HTTP framework for Go with excellent test coverage and a stable API.

#### [nfnt/resize][resize]

A simple, stable, popular pure-Go library for performing image resizing.

#### [mattn/go-isatty][isatty]

A small cross-platform (Windows, macOS, Linux, BSD, Solaris) library for determining whether `stdout` is a terminal.

[gin]: https://github.com/gin-gonic/gin#readme
[resize]: https://github.com/nfnt/resize
[isatty]: https://github.com/mattn/go-isatty
