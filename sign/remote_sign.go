package sign

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/address"
)

func RemoteSignNew(c *Client, addr string) HandleSignCkbMessage {
	return func(message string) ([]byte, error) {
		log.Info("RemoteSign:", message)
		return c.SignCkbMessage(addr, message)
	}
}

func RemoteSign(c *Client, net common.DasNetType, args string) HandleSignCkbMessage {
	return func(message string) ([]byte, error) {
		log.Info("RemoteSign:", message)
		addr, err := GenerateAddressByArgs(net, args)
		if err != nil {
			return nil, fmt.Errorf("address.Generate err: %s", err.Error())
		}
		return c.SignCkbMessage(addr, message)
	}
}

func GenerateAddressByArgs(net common.DasNetType, args string) (string, error) {
	serverLock := common.GetNormalLockScript(args)
	netMode := address.Testnet
	if net == common.DasNetTypeMainNet {
		netMode = address.Mainnet
	}
	return common.ConvertScriptToAddress(netMode, serverLock)
}
