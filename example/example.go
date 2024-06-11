package example

import (
	"context"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"sync"
)

func getClientTestnet2() (rpc.Client, error) {
	ckbUrl := "https://testnet.ckb.dev/"
	indexerUrl := "https://testnet.ckb.dev/indexer"
	ckbUrl = "https://testnet.ckb.dev/"
	indexerUrl = "https://testnet.ckb.dev/"
	return rpc.DialWithIndexer(ckbUrl, indexerUrl)
}

func getNewDasCoreTestnet2() (*core.DasCore, error) {
	client, err := getClientTestnet2()
	if err != nil {
		return nil, err
	}

	env := core.InitEnvOpt(common.DasNetTypeTestnet2,
		common.DasContractNameConfigCellType,
		common.DasContractNameAccountCellType,
		common.DasContractNameDispatchCellType,
		//common.DasContractNameBalanceCellType,
		//common.DasContractNameAlwaysSuccess,
		//common.DasContractNameIncomeCellType,
		//common.DASContractNameSubAccountCellType,
		//common.DasContractNamePreAccountCellType,
		//common.DasContractNameReverseRecordRootCellType,
		//common.DasKeyListCellType,
		//common.DasContractNameDpCellType,
		common.DasContractNameDidCellType,
	)
	var wg sync.WaitGroup
	ops := []core.DasCoreOption{
		core.WithClient(client),
		core.WithDasContractArgs(env.ContractArgs),
		core.WithDasContractCodeHash(env.ContractCodeHash),
		core.WithDasNetType(common.DasNetTypeTestnet2),
		core.WithTHQCodeHash(env.THQCodeHash),
	}
	dc := core.NewDasCore(context.Background(), &wg, ops...)
	// contract
	dc.InitDasContract(env.MapContract)
	// config cell
	if err = dc.InitDasConfigCell(); err != nil {
		return nil, err
	}
	// so script
	if err = dc.InitDasSoScript(); err != nil {
		return nil, err
	}
	return dc, nil
}

func getClientMainNet() (rpc.Client, error) {
	ckbUrl := "http://127.0.0.1:8114"
	indexerUrl := "http://127.0.0.1:8116"
	ckbUrl = "https://mainnet.ckb.dev/"
	indexerUrl = "https://mainnet.ckb.dev/"
	return rpc.DialWithIndexer(ckbUrl, indexerUrl)
}

func getNewDasCoreMainNet() (*core.DasCore, error) {
	client, err := getClientMainNet()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	env := core.InitEnvOpt(common.DasNetTypeMainNet,
		common.DasContractNameConfigCellType,
		//common.DasContractNameAlwaysSuccess,
		//common.DasContractNameIncomeCellType,
		common.DasContractNameAccountCellType,
		//common.DasContractNameDispatchCellType,
		//common.DasContractNameAlwaysSuccess,
		//common.DASContractNameSubAccountCellType,
		//common.DasContractNameBalanceCellType,
		//common.DasContractNamePreAccountCellType,
		common.DasContractNameDpCellType,
	)
	ops := []core.DasCoreOption{
		core.WithClient(client),
		core.WithDasContractArgs(env.ContractArgs),
		core.WithDasContractCodeHash(env.ContractCodeHash),
		core.WithDasNetType(common.DasNetTypeMainNet),
		core.WithTHQCodeHash(env.THQCodeHash),
	}
	dc := core.NewDasCore(context.Background(), &wg, ops...)
	// contract
	dc.InitDasContract(env.MapContract)
	// config cell
	if err = dc.InitDasConfigCell(); err != nil {
		return nil, err
	}
	// so script
	if err = dc.InitDasSoScript(); err != nil {
		return nil, err
	}
	return dc, nil
}
