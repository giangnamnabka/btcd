txscript
========

[![Build Status](https://github.com/giangnamnabka/btcd/workflows/Build%20and%20Test/badge.svg)](https://github.com/giangnamnabka/btcd/actions)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://pkg.go.dev/github.com/giangnamnabka/btcd/txscript?status.png)](https://pkg.go.dev/github.com/giangnamnabka/btcd/txscript)

Package txscript implements the bitcoin transaction script language.  There is
a comprehensive test suite.

This package has intentionally been designed so it can be used as a standalone
package for any projects needing to use or validate bitcoin transaction scripts.

## Bitcoin Scripts

Bitcoin provides a stack-based, FORTH-like language for the scripts in
the bitcoin transactions.  This language is not turing complete
although it is still fairly powerful.  A description of the language
can be found at https://en.bitcoin.it/wiki/Script

## Installation and Updating

```bash
$ go get -u github.com/giangnamnabka/btcd/txscript
```

## Examples

* [Standard Pay-to-pubkey-hash Script](https://pkg.go.dev/github.com/giangnamnabka/btcd/txscript#example-PayToAddrScript)  
  Demonstrates creating a script which pays to a bitcoin address.  It also
  prints the created script hex and uses the DisasmString function to display
  the disassembled script.

* [Extracting Details from Standard Scripts](https://pkg.go.dev/github.com/giangnamnabka/btcd/txscript#example-ExtractPkScriptAddrs)  
  Demonstrates extracting information from a standard public key script.

* [Manually Signing a Transaction Output](https://pkg.go.dev/github.com/giangnamnabka/btcd/txscript#example-SignTxOutput)  
  Demonstrates manually creating and signing a redeem transaction.

## GPG Verification Key

All official release tags are signed by Conformal so users can ensure the code
has not been tampered with and is coming from the btcsuite developers.  To
verify the signature perform the following:

- Download the public key from the Conformal website at
  https://opensource.conformal.com/GIT-GPG-KEY-conformal.txt

- Import the public key into your GPG keyring:
  ```bash
  gpg --import GIT-GPG-KEY-conformal.txt
  ```

- Verify the release tag with the following command where `TAG_NAME` is a
  placeholder for the specific tag:
  ```bash
  git tag -v TAG_NAME
  ```

## License

Package txscript is licensed under the [copyfree](http://copyfree.org) ISC
License.
