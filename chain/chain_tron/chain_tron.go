package chain_tron

import (
	"context"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/grpc"
	"strings"
	"unicode"
)

type ChainTron struct {
	Ctx    context.Context
	Client api.WalletClient
}

func NewChainTron(ctx context.Context, node string) (*ChainTron, error) {
	conn, err := grpc.DialContext(ctx, node, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &ChainTron{
		Ctx:    ctx,
		Client: api.NewWalletClient(conn),
	}, nil
}

func GetMemo(s []byte) string {
	str := make([]rune, 0, len(s))
	for _, v := range string(s) {
		if unicode.IsControl(v) {
			continue
		}
		str = append(str, v)
	}
	return strings.Replace(string(str), " ", "", 1)
}
