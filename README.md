# xkcdpwd [![Build Status](https://travis-ci.org/wfscheper/xkcdpwd.svg?branch=master)](https://travis-ci.org/wfscheper/xkcdpwd) [![Coverage Status](https://coveralls.io/repos/github/wfscheper/xkcdpwd/badge.svg?branch=master)](https://coveralls.io/github/wfscheper/xkcdpwd?branch=master)

A passphrase generator in the style of XKCD comic **[Password Strength]**.

## Getting started

This project requires Go 1.11 to be installed, as it uses the new go module system for dependency management.
On OS X with Homebrew you can just run `brew install go`.
On Linux check with your distro's package manager.
Alternatively, or on Windows, download the Go source or binary from [golang.org](https://golang.org/dl/).

Running it then should be as simple as:

```console
git clone https://gitub.com/wfscheper/xkcdpwd.git
cd xkcdpwd
go get github.com/magefile/mage
mage
bin/xkcdpwd
```

## Testing

``mage test``

To generate coverage data:

``mage coverage``

## Similar projects

- [A pure JavaScript implementation with website](http://preshing.com/20110811/xkcd-password-generator/).
- [A Perl implementation](https://github.com/bbusschots/hsxkpasswd) with [a web-accessible UI](https://xkpasswd.net/s/).

[Password Strength]: https://xkcd.com/936/
