# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

### Changed

- Replace Makefile with [magefile](https://github.com/magefile/mage)
- Switch from dep to go module

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
