package remote_sign

import (
	"context"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
)

type reqParam struct {
	Errno  int         `json:"errno"`
	Errmsg interface{} `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type RemoteSignClient struct {
	ctx    context.Context
	client rpc.Client
}

func NewRemoteSignClient(ctx context.Context, apiUrl string) (*RemoteSignClient, error) {
	client, err := rpc.Dial(apiUrl)
	if err != nil {
		return nil, err
	}

	return &RemoteSignClient{
		ctx:    ctx,
		client: client,
	}, nil
}

const (
	SignMethodEvm  string = "wallet_eTHSignMsg"
	SignMethodTron string = "wallet_tronSignMsg"
	SignMethodCkb  string = "wallet_cKBSignMsg"
)

func (r *RemoteSignClient) Client() rpc.Client {
	return r.client
}
