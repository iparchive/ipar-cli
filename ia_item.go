package ipar

import (
	"github.com/multiformats/go-multiaddr"
)

var (
	LocalIpfsEndpoint multiaddr.Multiaddr
)

func init() {
	var err error
	LocalIpfsEndpoint, err = multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		panic(err)
	}
}
