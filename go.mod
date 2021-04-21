module github.com/giangnamnabka/btcd

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd
	github.com/btcsuite/goleveldb v1.0.0
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/davecgh/go-spew v1.1.1
	github.com/decred/dcrd/lru v1.0.0
	github.com/giangnamnabka/btcutil v1.0.3-0.20210421053153-55ad5ef66f21
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
)

// replace github.com/giangnamnabka/btcutil => github.com/giangnamnabka/btcutil v1.0.3-0.20210420155121-8b03ebc257c8

replace github.com/giangnamnabka/btcutil v1.0.3-0.20210421053153-55ad5ef66f21 => /home/tinhnguyenhn/anyswap/20210421_2/btcutil

go 1.14
