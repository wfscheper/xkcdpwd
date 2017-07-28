# xkcdpwd

A pass phrase generator in the style of XKCD comic **[Password Strength]**

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just
run `brew install go`. On Linux check with your distro's package manager.
Alternatively, or on Windows, download the go source or binary from
https://golang.org/dl/.

Running it then should be as simple as:

```console
$ make
$ bin/xkcdpwd
```

## Testing

``make test``

To generate coverage data

``make test-coverage``

## Similar projects

- http://preshing.com/20110811/xkcd-password-generator/
  A pure javascript implementation with website.
- https://github.com/bbusschots/hsxkpasswd
  A perl implementation, with a web-accessible UI at https://xkpasswd.net/s/

[Password Strength]: https://xkcd.com/936/
