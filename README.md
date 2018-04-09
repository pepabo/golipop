<p align="center"><img src="https://raw.githubusercontent.com/pepabo/lolp.rb/images/lolipop-logo-by-gmo-pepabo.png" width="300" alt="lolp" /></p><p align="center"><strong>lolp</strong>: A Lolipop! Managed Cloud API client library for Go</p> <br /> <br />

[![Travis](https://img.shields.io/travis/pepabo/golipop.svg?style=flat-square)][travis]
[![GitHub release](http://img.shields.io/github/release/pepabo/golipop.svg?style=flat-square)][release]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[travis]: https://travis-ci.org/pepabo/golipop
[release]: https://github.com/pepabo/golipop/releases
[license]: https://github.com/pepabo/golipop/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/pepabo/golipop

Installation
------------

To install, use `go get`:

```sh
$ go get -d github.com/pepabo/golipop
```

Usage
-----

As CLI:

```sh
$ eval $(lolp login -u <your@example.com> -p <your_password>)
$ lolp project create -k rails -s foobar -d password:<your_password>
foobar.lolipop.io
$ lolp project foobar
ID                 "cdd32ae5-c118-4fc9-b9d6-ea5ad18f3737"
Kind               "rails"
Domain             "foobar.lolipop.io"
...
```

As library:

```go
client := lolp.DefaultClinet()
token, err := client.Login("your@example.com", "your_password")
if err != nil {
  panic(err)
}
p := &ProjectNew{
  Kind: "rails",
  Database: map[string]interface{} {"password": "********"},
}
project, err := client.CreateProject(p)
if err != nil {
  panic(err)
}
```

Contribution
------------

1. Fork ([https://github.com/pepabo/golipop/fork](https://github.com/pepabo/golipop/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

Author
------

[linyows](https://github.com/linyows)
