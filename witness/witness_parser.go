package witness

import (
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"strings"
)

func ParserWitnessData(witnessByte []byte) interface{} {
	if len(witnessByte) <= common.WitnessDasTableTypeEndIndex+1 {
		return parserDefaultWitness(witnessByte)
	}
	if string(witnessByte[0:common.WitnessDasCharLen]) != common.WitnessDas {
		return parserDefaultWitness(witnessByte)
	}
	actionDataType := common.Bytes2Hex(witnessByte[common.WitnessDasCharLen:common.WitnessDasTableTypeEndIndex])

	switch actionDataType {
	case common.ActionDataTypeActionData:
		return parserActionData(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ActionDataTypeAccountCell:
		return parserAccountCell(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ActionDataTypeAccountSaleCell:
		return parserAccountSaleCell(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ActionDataTypeAccountAuctionCell:
		return parserAccountAuctionCell(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ActionDataTypeProposalCell:
		return parserProposalCell(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ActionDataTypePreAccountCell:
		return parserPreAccountCell(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ActionDataTypeIncomeCell:
		return parserIncomeCell(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ActionDataTypeOfferCell:
		return parserOfferCell(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])

	case common.ConfigCellTypeArgsAccount:
		return parserConfigCellAccount(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsApply:
		return parserConfigCellApply(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsIncome:
		return parserConfigCellIncome(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsMain:
		return parserConfigCellMain(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsPrice:
		return parserConfigCellPrice(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsProposal:
		return parserConfigCellProposal(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsProfitRate:
		return parserConfigCellProfitRate(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsRecordNamespace:
		return parserConfigCellRecordNamespace(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsRelease:
		return parserConfigCellRelease(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsUnavailable:
		return parserConfigCellUnavailable(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsSecondaryMarket:
		return parserConfigCellSecondaryMarket(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])
	case common.ConfigCellTypeArgsReverseRecord:
		return parserConfigCellReverseRecord(witnessByte, witnessByte[common.WitnessDasTableTypeEndIndex:])

	default:
		return parserDefaultWitness(witnessByte)
	}
}

func parserDefaultWitness(witnessByte []byte) interface{} {
	return map[string]interface{}{
		"unknown": common.Bytes2Hex(witnessByte),
	}
}

func parserActionData(witnessByte, slice []byte) interface{} {
	var builder ActionDataBuilder
	actionData, _ := molecule.ActionDataFromSlice(slice, false)
	if actionData == nil {
		return parserDefaultWitness(witnessByte)
	}
	builder.ActionData = actionData
	builder.Action = string(actionData.Action().RawData())
	if builder.Action == common.DasActionBuyAccount {
		raw := actionData.Params().RawData()

		lenRaw := len(raw)
		inviterLockBytesLen, err := molecule.Bytes2GoU32(raw[:4])
		if err != nil {
			return parserDefaultWitness(witnessByte)
		}
		inviterLockRaw := raw[:inviterLockBytesLen]
		channelLockRaw := raw[inviterLockBytesLen : lenRaw-1]

		builder.Params = append(builder.Params, inviterLockRaw)
		builder.Params = append(builder.Params, channelLockRaw)
		builder.Params = append(builder.Params, raw[lenRaw-1:lenRaw])
		builder.ParamsStr = common.GetMaxHashLenParams(common.Bytes2Hex(inviterLockRaw)) + "," + common.GetMaxHashLenParams(common.Bytes2Hex(channelLockRaw)) + "," + common.Bytes2Hex(raw[lenRaw-1:lenRaw])
	} else {
		builder.Params = append(builder.Params, actionData.Params().RawData())
		builder.ParamsStr = common.Bytes2Hex(actionData.Params().RawData())
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(builder.ActionData.AsSlice())),
		"ActionData": map[string]interface{}{
			"action":      builder.Action,
			"action_hash": common.Bytes2Hex(builder.ActionData.Action().RawData()),
			"params":      builder.ParamsStr,
		},
	}
}

func parserAccountCell(witnessByte, slice []byte) interface{} {
	accountCell, _ := molecule.AccountCellDataFromSlice(slice, false)
	if accountCell != nil {
		return map[string]interface{}{
			"witness":      common.Bytes2Hex(witnessByte),
			"witness_hash": common.Bytes2Hex(common.Blake2b(accountCell.AsSlice())),
			"AccountCell":  map[string]interface{}{},
		}
	}

	accountCellV1, _ := molecule.AccountCellDataV1FromSlice(slice, false)
	if accountCellV1 != nil {
		return map[string]interface{}{
			"witness":      common.Bytes2Hex(witnessByte),
			"witness_hash": common.Bytes2Hex(common.Blake2b(accountCellV1.AsSlice())),
			"AccountCell":  map[string]interface{}{},
		}
	}
	return parserDefaultWitness(witnessByte)
}

func parserAccountSaleCell(witnessByte, slice []byte) interface{} {
	accountSaleCell, _ := molecule.AccountSaleCellDataFromSlice(slice, false)
	if accountSaleCell != nil {
		return map[string]interface{}{
			"witness":         common.Bytes2Hex(witnessByte),
			"witness_hash":    common.Bytes2Hex(common.Blake2b(accountSaleCell.AsSlice())),
			"AccountSaleCell": map[string]interface{}{},
		}
	}

	accountSaleCellV1, _ := molecule.AccountSaleCellDataV1FromSlice(slice, false)
	if accountSaleCellV1 != nil {
		return map[string]interface{}{
			"witness":         common.Bytes2Hex(witnessByte),
			"witness_hash":    common.Bytes2Hex(common.Blake2b(accountSaleCellV1.AsSlice())),
			"AccountSaleCell": map[string]interface{}{},
		}
	}
	return parserDefaultWitness(witnessByte)
}

func parserAccountAuctionCell(witnessByte, slice []byte) interface{} {
	accountAuctionCell, _ := molecule.AccountAuctionCellDataFromSlice(slice, false)
	if accountAuctionCell == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":            common.Bytes2Hex(witnessByte),
		"witness_hash":       common.Bytes2Hex(common.Blake2b(accountAuctionCell.AsSlice())),
		"AccountAuctionCell": map[string]interface{}{},
	}
}

func parserProposalCell(witnessByte, slice []byte) interface{} {
	proposalCell, _ := molecule.ProposalCellDataFromSlice(slice, false)
	if proposalCell == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(proposalCell.AsSlice())),
		"ProposalCell": map[string]interface{}{},
	}
}

func parserPreAccountCell(witnessByte, slice []byte) interface{} {
	preAccountCell, _ := molecule.PreAccountCellDataFromSlice(slice, false)
	if preAccountCell == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":        common.Bytes2Hex(witnessByte),
		"witness_hash":   common.Bytes2Hex(common.Blake2b(preAccountCell.AsSlice())),
		"PreAccountCell": map[string]interface{}{},
	}
}

func parserIncomeCell(witnessByte, slice []byte) interface{} {
	data, _ := molecule.DataFromSlice(slice, false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	var dataEntityOpts []map[string]interface{}
	var incomeCellData []interface{}
	if data.New() != nil && !data.New().IsNone() {
		dataEntityOpts = append(dataEntityOpts, map[string]interface{}{
			"type":   "new",
			"entity": data.New(),
		})
	}
	if data.Old() != nil && !data.Old().IsNone() {
		dataEntityOpts = append(dataEntityOpts, map[string]interface{}{
			"type":   "old",
			"entity": data.Old(),
		})
	}
	if data.Dep() != nil && !data.Dep().IsNone() {
		dataEntityOpts = append(dataEntityOpts, map[string]interface{}{
			"type":   "dep",
			"entity": data.Dep(),
		})
	}

	for _, v := range dataEntityOpts {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		incomeCell, _ := molecule.IncomeCellDataFromSlice(dataEntity.Entity().RawData(), false)
		if incomeCell == nil {
			return parserDefaultWitness(witnessByte)
		}

		var recordsMaps []map[string]interface{}
		for i := uint(0); i < incomeCell.Records().Len(); i++ {
			record := incomeCell.Records().Get(i)
			capacity, _ := molecule.Bytes2GoU64(record.Capacity().RawData())
			recordsMaps = append(recordsMaps, map[string]interface{}{
				"belong_to": map[string]interface{}{
					"code_hash": common.Bytes2Hex(record.BelongTo().CodeHash().RawData()),
					"hash_type": common.Bytes2Hex(record.BelongTo().HashType().AsSlice()),
					"args":      common.Bytes2Hex(record.BelongTo().Args().RawData()),
				},
				"capacity": capacity,
			})
		}

		incomeCellData = append(incomeCellData, map[string]interface{}{
			v["type"].(string): map[string]interface{}{
				"version": version,
				"index":   index,
				"entity": map[string]interface{}{
					"witness_hash": common.Bytes2Hex(common.Blake2b(incomeCell.AsSlice())),
					"creator": map[string]interface{}{
						"code_hash": common.Bytes2Hex(incomeCell.Creator().CodeHash().RawData()),
						"hash_type": common.Bytes2Hex(incomeCell.Creator().HashType().AsSlice()),
					},
					"records": recordsMaps,
				},
			},
		})
	}

	return map[string]interface{}{
		"witness":    common.Bytes2Hex(witnessByte),
		"IncomeCell": incomeCellData,
	}
}

func parserOfferCell(witnessByte, slice []byte) interface{} {
	offerCell, _ := molecule.OfferCellDataFromSlice(slice, false)
	if offerCell == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(offerCell.AsSlice())),
		"OfferCell":    map[string]interface{}{},
	}
}

func parserConfigCellAccount(witnessByte, configCellByte []byte) interface{} {
	configCellAccount, _ := molecule.ConfigCellAccountFromSlice(configCellByte, false)
	if configCellAccount == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellAccount.AsSlice())),
		"ConfigCellAccount": map[string]interface{}{
			"max_length":                common.Bytes2Hex(configCellAccount.MaxLength().RawData()),
			"basic_capacity":            common.Bytes2Hex(configCellAccount.BasicCapacity().RawData()),
			"prepared_fee_capacity":     common.Bytes2Hex(configCellAccount.PreparedFeeCapacity().RawData()),
			"expiration_grace_period":   common.Bytes2Hex(configCellAccount.ExpirationGracePeriod().RawData()),
			"record_min_ttl":            common.Bytes2Hex(configCellAccount.RecordMinTtl().RawData()),
			"record_size_limit":         common.Bytes2Hex(configCellAccount.RecordSizeLimit().RawData()),
			"transfer_account_fee":      common.Bytes2Hex(configCellAccount.TransferAccountFee().RawData()),
			"edit_manager_fee":          common.Bytes2Hex(configCellAccount.EditManagerFee().RawData()),
			"edit_records_fee":          common.Bytes2Hex(configCellAccount.EditRecordsFee().RawData()),
			"common_fee":                common.Bytes2Hex(configCellAccount.CommonFee().RawData()),
			"transfer_account_throttle": common.Bytes2Hex(configCellAccount.TransferAccountThrottle().RawData()),
			"edit_manager_throttle":     common.Bytes2Hex(configCellAccount.EditManagerThrottle().RawData()),
			"edit_records_throttle":     common.Bytes2Hex(configCellAccount.EditRecordsThrottle().RawData()),
			"common_throttle":           common.Bytes2Hex(configCellAccount.CommonThrottle().RawData()),
		},
	}
}

func parserConfigCellApply(witnessByte, configCellByte []byte) interface{} {
	configCellApply, _ := molecule.ConfigCellApplyFromSlice(configCellByte, false)
	if configCellApply == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellApply.AsSlice())),
		"ConfigCellApply": map[string]interface{}{
			"apply_min_waiting_block_number": common.Bytes2Hex(configCellApply.ApplyMinWaitingBlockNumber().RawData()),
			"apply_max_waiting_block_number": common.Bytes2Hex(configCellApply.ApplyMaxWaitingBlockNumber().RawData()),
		},
	}
}

func parserConfigCellIncome(witnessByte, configCellByte []byte) interface{} {
	configCellIncome, _ := molecule.ConfigCellIncomeFromSlice(configCellByte, false)
	if configCellIncome == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellIncome.AsSlice())),
		"ConfigCellIncome": map[string]interface{}{
			"basic_capacity":        common.Bytes2Hex(configCellIncome.BasicCapacity().RawData()),
			"max_records":           common.Bytes2Hex(configCellIncome.MaxRecords().RawData()),
			"min_transfer_capacity": common.Bytes2Hex(configCellIncome.MinTransferCapacity().RawData()),
		},
	}
}

func parserConfigCellMain(witnessByte, configCellByte []byte) interface{} {
	configCellMain, _ := molecule.ConfigCellMainFromSlice(configCellByte, false)
	if configCellMain == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellMain.AsSlice())),
		"ConfigCellMain": map[string]interface{}{
			"status": common.Bytes2Hex(configCellMain.Status().RawData()),
			"type_id_table": map[string]interface{}{
				"account_cell":         common.Bytes2Hex(configCellMain.TypeIdTable().AccountCell().RawData()),
				"apply_register_cell":  common.Bytes2Hex(configCellMain.TypeIdTable().ApplyRegisterCell().RawData()),
				"balance_cell":         common.Bytes2Hex(configCellMain.TypeIdTable().BalanceCell().RawData()),
				"income_cell":          common.Bytes2Hex(configCellMain.TypeIdTable().IncomeCell().RawData()),
				"pre_account_cell":     common.Bytes2Hex(configCellMain.TypeIdTable().PreAccountCell().RawData()),
				"proposal_cell":        common.Bytes2Hex(configCellMain.TypeIdTable().ProposalCell().RawData()),
				"account_sale_cell":    common.Bytes2Hex(configCellMain.TypeIdTable().AccountSaleCell().RawData()),
				"account_auction_cell": common.Bytes2Hex(configCellMain.TypeIdTable().AccountAuctionCell().RawData()),
				"offer_cell":           common.Bytes2Hex(configCellMain.TypeIdTable().OfferCell().RawData()),
				"reverse_record_cell":  common.Bytes2Hex(configCellMain.TypeIdTable().ReverseRecordCell().RawData()),
			},
			"das_lock_out_point_table": map[string]interface{}{
				"ckb_signall": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbSignall().TxHash().RawData()),
					"index":   common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbSignall().Index().RawData()),
				},
				"ckb_multisign": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbMultisign().TxHash().RawData()),
					"index":   common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbMultisign().Index().RawData()),
				},
				"ckb_anyone_can_pay": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbAnyoneCanPay().TxHash().RawData()),
					"index":   common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbAnyoneCanPay().Index().RawData()),
				},
				"eth": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().Eth().TxHash().RawData()),
					"index":   common.Bytes2Hex(configCellMain.DasLockOutPointTable().Eth().Index().RawData()),
				},
				"tron": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().Tron().TxHash().RawData()),
					"index":   common.Bytes2Hex(configCellMain.DasLockOutPointTable().Tron().Index().RawData()),
				},
			},
		},
	}
}

func parserConfigCellPrice(witnessByte, configCellByte []byte) interface{} {
	configCellPrice, _ := molecule.ConfigCellPriceFromSlice(configCellByte, false)
	if configCellPrice != nil {
		return parserDefaultWitness(witnessByte)
	}

	var prices []interface{}
	for i := uint(0); i < configCellPrice.Prices().Len(); i++ {
		price, _ := molecule.PriceConfigFromSlice(configCellPrice.Prices().Get(i).AsSlice(), false)
		if price == nil {
			return parserDefaultWitness(witnessByte)
		}
		prices = append(prices, map[string]interface{}{
			"length": common.Bytes2Hex(price.Length().RawData()),
			"new":    common.Bytes2Hex(price.New().RawData()),
			"renew":  common.Bytes2Hex(price.Renew().RawData()),
		})
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellPrice.AsSlice())),
		"ConfigCellPrice": map[string]interface{}{
			"discount": map[string]interface{}{
				"invited_discount": common.Bytes2Hex(configCellPrice.Discount().InvitedDiscount().RawData()),
			},
			"prices": prices,
		},
	}
}

func parserConfigCellProposal(witnessByte, configCellByte []byte) interface{} {
	configCellProposal, _ := molecule.ConfigCellProposalFromSlice(configCellByte, false)
	if configCellProposal == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellProposal.AsSlice())),
		"ConfigCellProposal": map[string]interface{}{
			"proposal_min_confirm_interval":    common.Bytes2Hex(configCellProposal.ProposalMinConfirmInterval().RawData()),
			"proposal_min_recycle_interval":    common.Bytes2Hex(configCellProposal.ProposalMinRecycleInterval().RawData()),
			"proposal_min_extend_interval":     common.Bytes2Hex(configCellProposal.ProposalMinExtendInterval().RawData()),
			"proposal_max_account_affect":      common.Bytes2Hex(configCellProposal.ProposalMaxAccountAffect().RawData()),
			"proposal_max_pre_account_contain": common.Bytes2Hex(configCellProposal.ProposalMaxPreAccountContain().RawData()),
		},
	}
}

func parserConfigCellProfitRate(witnessByte, configCellByte []byte) interface{} {
	configCellProfitRate, _ := molecule.ConfigCellProfitRateFromSlice(configCellByte, false)
	if configCellProfitRate == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellProfitRate.AsSlice())),
		"ConfigCellProfitRate": map[string]interface{}{
			"inviter":                common.Bytes2Hex(configCellProfitRate.Inviter().RawData()),
			"channel":                common.Bytes2Hex(configCellProfitRate.Channel().RawData()),
			"proposal_create":        common.Bytes2Hex(configCellProfitRate.ProposalCreate().RawData()),
			"proposal_confirm":       common.Bytes2Hex(configCellProfitRate.ProposalConfirm().RawData()),
			"income_consolidate":     common.Bytes2Hex(configCellProfitRate.IncomeConsolidate().RawData()),
			"sale_buyer_inviter":     common.Bytes2Hex(configCellProfitRate.SaleBuyerInviter().RawData()),
			"sale_buyer_channel":     common.Bytes2Hex(configCellProfitRate.SaleBuyerChannel().RawData()),
			"sale_das":               common.Bytes2Hex(configCellProfitRate.SaleDas().RawData()),
			"auction_bidder_inviter": common.Bytes2Hex(configCellProfitRate.AuctionBidderInviter().RawData()),
			"auction_bidder_channel": common.Bytes2Hex(configCellProfitRate.AuctionBidderChannel().RawData()),
			"auction_das":            common.Bytes2Hex(configCellProfitRate.AuctionDas().RawData()),
			"auction_prev_bidder":    common.Bytes2Hex(configCellProfitRate.AuctionPrevBidder().RawData()),
		},
	}
}

func parserConfigCellRecordNamespace(witnessByte, configCellByte []byte) interface{} {
	dataLength, err := molecule.Bytes2GoU32(configCellByte[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(configCellByte),
		"ConfigCellRecordNamespace": map[string]interface{}{
			"data_length": dataLength,
			"record_keys": strings.Split(string(configCellByte[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func parserConfigCellRelease(witnessByte, configCellByte []byte) interface{} {
	configCellRelease, _ := molecule.ConfigCellReleaseFromSlice(configCellByte, false)
	if configCellRelease != nil {
		return parserDefaultWitness(witnessByte)
	}

	var releaseRules []interface{}
	for i := uint(0); i < configCellRelease.ReleaseRules().Len(); i++ {
		releaseRule := configCellRelease.ReleaseRules().Get(i)
		releaseRules = append(releaseRules, map[string]interface{}{
			"length":        common.Bytes2Hex(releaseRule.Length().RawData()),
			"release_start": common.Bytes2Hex(releaseRule.ReleaseStart().RawData()),
			"release_end":   common.Bytes2Hex(releaseRule.ReleaseEnd().RawData()),
		})
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellRelease.AsSlice())),
		"ConfigCellRelease": map[string]interface{}{
			"release_rules": releaseRules,
		},
	}
}

func parserConfigCellUnavailable(witnessByte, configCellByte []byte) interface{} {
	dataLength, err := molecule.Bytes2GoU32(configCellByte[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(configCellByte),
		"ConfigCellUnavailable": map[string]interface{}{
			"data_length": dataLength,
		},
	}
}

func parserConfigCellSecondaryMarket(witnessByte, configCellByte []byte) interface{} {
	configCellSecondaryMarket, _ := molecule.ConfigCellSecondaryMarketFromSlice(configCellByte, false)
	if configCellSecondaryMarket == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellSecondaryMarket.AsSlice())),
		"ConfigCellSecondaryMarket": map[string]interface{}{
			"common_fee":                          common.Bytes2Hex(configCellSecondaryMarket.CommonFee().RawData()),
			"sale_min_price":                      common.Bytes2Hex(configCellSecondaryMarket.SaleMinPrice().RawData()),
			"sale_expiration_limit":               common.Bytes2Hex(configCellSecondaryMarket.SaleExpirationLimit().RawData()),
			"sale_description_bytes_limit":        common.Bytes2Hex(configCellSecondaryMarket.SaleDescriptionBytesLimit().RawData()),
			"sale_cell_basic_capacity":            common.Bytes2Hex(configCellSecondaryMarket.SaleCellBasicCapacity().RawData()),
			"sale_cell_prepared_fee_capacity":     common.Bytes2Hex(configCellSecondaryMarket.SaleCellPreparedFeeCapacity().RawData()),
			"auction_max_extendable_duration":     common.Bytes2Hex(configCellSecondaryMarket.AuctionMaxExtendableDuration().RawData()),
			"auction_duration_increment_each_bid": common.Bytes2Hex(configCellSecondaryMarket.AuctionDurationIncrementEachBid().RawData()),
			"auction_min_opening_price":           common.Bytes2Hex(configCellSecondaryMarket.AuctionMinOpeningPrice().RawData()),
			"auction_min_increment_rate_each_bid": common.Bytes2Hex(configCellSecondaryMarket.AuctionMinIncrementRateEachBid().RawData()),
			"auction_description_bytes_limit":     common.Bytes2Hex(configCellSecondaryMarket.AuctionDescriptionBytesLimit().RawData()),
			"auction_cell_basic_capacity":         common.Bytes2Hex(configCellSecondaryMarket.AuctionCellBasicCapacity().RawData()),
			"auction_cell_prepared_fee_capacity":  common.Bytes2Hex(configCellSecondaryMarket.AuctionCellPreparedFeeCapacity().RawData()),
			"offer_min_price":                     common.Bytes2Hex(configCellSecondaryMarket.OfferMinPrice().RawData()),
			"offer_cell_basic_capacity":           common.Bytes2Hex(configCellSecondaryMarket.OfferCellBasicCapacity().RawData()),
			"offer_cell_prepared_fee_capacity":    common.Bytes2Hex(configCellSecondaryMarket.OfferCellPreparedFeeCapacity().RawData()),
			"offer_message_bytes_limit":           common.Bytes2Hex(configCellSecondaryMarket.OfferMessageBytesLimit().RawData()),
		},
	}
}

func parserConfigCellReverseRecord(witnessByte, configCellByte []byte) interface{} {
	configCellReverseRecord, _ := molecule.ConfigCellReverseResolutionFromSlice(configCellByte, false)
	if configCellReverseRecord == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellReverseRecord.AsSlice())),
		"ConfigCellReverseRecord": map[string]interface{}{
			"common_fee":                   common.Bytes2Hex(configCellReverseRecord.CommonFee().RawData()),
			"record_prepared_fee_capacity": common.Bytes2Hex(configCellReverseRecord.RecordPreparedFeeCapacity().RawData()),
			"record_basic_capacity":        common.Bytes2Hex(configCellReverseRecord.RecordBasicCapacity().RawData()),
		},
	}
}
