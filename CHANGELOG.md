# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [v0.1.3] - 2019-09-20

### Fixed

- Fix travis deploy conidtional

## [v0.1.2] - 2019-09-17

### Changed

- Replace golint with golangci-lint
- Use goreleaser.

## [v0.1.1] - 2018-09-15

### Changed

- Replace Makefile with [magefile](https://github.com/magefile/mage)
- Switch from dep to go module

### Fixed

- Correctly build distribution executables

## v0.1.0 - 2018-06-11

### Added

- Generate 10, 4 word passphrase from ~8830 words, with ~52 bits of entropy.
- Basic support for multiple languages based on environment or command-line
  flag.
- Option to change the number of passphrases that are generated.
- Option to change the number of words in the passphrases.
- Pick from one of several capitalization strategies.
- Options to set minimum and maximum word length.
- Control the character used to separate words in the passphrase.
- Support config files for changing defaults for all command-line options.

[v0.1.3]: https://github.com/wfscheper/xkcdpwd/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/wfscheper/xkcdpwd/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/wfscheper/xkcdpwd/compare/v0.1.0...v0.1.1
[v0.1.1]: https://github.com/wfscheper/xkcdpwd/compare/4ec2e6...v0.1.0
