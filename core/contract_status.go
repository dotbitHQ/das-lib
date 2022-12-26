package core

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
)

var ContractStatusMapMainNet = map[common.DasContractName]common.ContractStatus{
	common.DasContractNameApplyRegisterCellType: {Version: "1.7.0"},
	common.DasContractNamePreAccountCellType:    {Version: "1.7.0"},
	common.DasContractNameProposalCellType:      {Version: "1.7.0"},
	common.DasContractNameConfigCellType:        {Version: "1.7.0"},
	common.DasContractNameAccountCellType:       {Version: "1.7.0"},
	common.DasContractNameAccountSaleCellType:   {Version: "1.7.0"},
	common.DASContractNameSubAccountCellType:    {Version: "1.7.0"},
	common.DASContractNameOfferCellType:         {Version: "1.7.0"},
	common.DasContractNameBalanceCellType:       {Version: "1.7.0"},
	common.DasContractNameIncomeCellType:        {Version: "1.7.0"},
	common.DasContractNameReverseRecordCellType: {Version: "1.7.0"},
	common.DASContractNameEip712LibCellType:     {Version: "1.7.0"},
}

var ContractStatusMapTestNet = map[common.DasContractName]common.ContractStatus{
	common.DasContractNameApplyRegisterCellType: {Version: "1.7.0"},
	common.DasContractNamePreAccountCellType:    {Version: "1.7.0"},
	common.DasContractNameProposalCellType:      {Version: "1.7.0"},
	common.DasContractNameConfigCellType:        {Version: "1.7.0"},
	common.DasContractNameAccountCellType:       {Version: "1.7.0"},
	common.DasContractNameAccountSaleCellType:   {Version: "1.7.0"},
	common.DASContractNameSubAccountCellType:    {Version: "1.7.0"},
	common.DASContractNameOfferCellType:         {Version: "1.7.0"},
	common.DasContractNameBalanceCellType:       {Version: "1.7.0"},
	common.DasContractNameIncomeCellType:        {Version: "1.7.0"},
	common.DasContractNameReverseRecordCellType: {Version: "1.7.0"},
	common.DASContractNameEip712LibCellType:     {Version: "1.7.0"},
}

func (d *DasCore) CheckContractVersion(contractName common.DasContractName) (defaultV, ChainV string, isSame bool, err error) {
	var defaultContractStatus common.ContractStatus
	var ok bool

	switch d.net {
	case common.DasNetTypeMainNet, common.DasNetTypeTestnet3:
		defaultContractStatus, ok = ContractStatusMapMainNet[contractName]
	case common.DasNetTypeTestnet2:
		defaultContractStatus, ok = ContractStatusMapTestNet[contractName]
	default:
		err = fmt.Errorf("unknow net[%d]", d.net)
		return
	}
	if !ok {
		err = fmt.Errorf("unkonw contract name[%s]", contractName)
		return
	}
	defaultV = defaultContractStatus.Version

	res, err := d.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsSystemStatus)
	if err != nil {
		err = fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
		return
	}
	contractStatus, err := res.GetContractStatus(contractName)
	if err != nil {
		err = fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
		return
	}
	ChainV = contractStatus.Version

	if defaultV != ChainV {
		defaultX, defaultY, _ := defaultContractStatus.VersionInfo()
		x, y, _ := contractStatus.VersionInfo()
		if x != defaultX || defaultY != y {
			err = ErrContractMajorVersionDiff
			return
		}
	}
	isSame = true

	return
}

var ErrContractMajorVersionDiff = errors.New("the major version of the contract is different")

func (d *DasCore) CheckContractStatusOK(contractName common.DasContractName) (bool, error) {
	res, err := d.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsSystemStatus)
	if err != nil {
		return false, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	contractStatus, err := res.GetContractStatus(contractName)
	if err != nil {
		return false, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	return contractStatus.Status == 1, nil
}
