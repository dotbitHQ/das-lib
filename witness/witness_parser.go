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
		return ParserActionData(witnessByte)
	case common.ActionDataTypeAccountCell:
		return ParserAccountCell(witnessByte)
	case common.ActionDataTypeAccountSaleCell:
		return ParserAccountSaleCell(witnessByte)
	case common.ActionDataTypeAccountAuctionCell:
		return ParserAccountAuctionCell(witnessByte)
	case common.ActionDataTypeProposalCell:
		return ParserProposalCell(witnessByte)
	case common.ActionDataTypePreAccountCell:
		return ParserPreAccountCell(witnessByte)
	case common.ActionDataTypeIncomeCell:
		return ParserIncomeCell(witnessByte)
	case common.ActionDataTypeOfferCell:
		return ParserOfferCell(witnessByte)

	case common.ConfigCellTypeArgsAccount:
		return ParserConfigCellAccount(witnessByte)
	case common.ConfigCellTypeArgsApply:
		return ParserConfigCellApply(witnessByte)
	case common.ConfigCellTypeArgsIncome:
		return ParserConfigCellIncome(witnessByte)
	case common.ConfigCellTypeArgsMain:
		return ParserConfigCellMain(witnessByte)
	case common.ConfigCellTypeArgsPrice:
		return ParserConfigCellPrice(witnessByte)
	case common.ConfigCellTypeArgsProposal:
		return ParserConfigCellProposal(witnessByte)
	case common.ConfigCellTypeArgsProfitRate:
		return ParserConfigCellProfitRate(witnessByte)
	case common.ConfigCellTypeArgsRecordNamespace:
		return ParserConfigCellRecordNamespace(witnessByte)
	case common.ConfigCellTypeArgsRelease:
		return ParserConfigCellRelease(witnessByte)
	case common.ConfigCellTypeArgsUnavailable:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellUnavailable")
	case common.ConfigCellTypeArgsSecondaryMarket:
		return ParserConfigCellSecondaryMarket(witnessByte)
	case common.ConfigCellTypeArgsReverseRecord:
		return ParserConfigCellReverseRecord(witnessByte)

	case common.ConfigCellTypeArgsPreservedAccount00:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount00")
	case common.ConfigCellTypeArgsPreservedAccount01:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount01")
	case common.ConfigCellTypeArgsPreservedAccount02:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount02")
	case common.ConfigCellTypeArgsPreservedAccount03:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount03")
	case common.ConfigCellTypeArgsPreservedAccount04:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount04")
	case common.ConfigCellTypeArgsPreservedAccount05:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount05")
	case common.ConfigCellTypeArgsPreservedAccount06:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount06")
	case common.ConfigCellTypeArgsPreservedAccount07:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount07")
	case common.ConfigCellTypeArgsPreservedAccount08:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount08")
	case common.ConfigCellTypeArgsPreservedAccount09:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount09")
	case common.ConfigCellTypeArgsPreservedAccount10:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount10")
	case common.ConfigCellTypeArgsPreservedAccount11:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount11")
	case common.ConfigCellTypeArgsPreservedAccount12:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount12")
	case common.ConfigCellTypeArgsPreservedAccount13:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount13")
	case common.ConfigCellTypeArgsPreservedAccount14:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount14")
	case common.ConfigCellTypeArgsPreservedAccount15:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount15")
	case common.ConfigCellTypeArgsPreservedAccount16:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount16")
	case common.ConfigCellTypeArgsPreservedAccount17:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount17")
	case common.ConfigCellTypeArgsPreservedAccount18:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount18")
	case common.ConfigCellTypeArgsPreservedAccount19:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount19")

	case common.ConfigCellTypeArgsCharSetEmoji:
		return ParserConfigCellTypeArgsCharSetEmoji(witnessByte)
	case common.ConfigCellTypeArgsCharSetDigit:
		return ParserConfigCellTypeArgsCharSetDigit(witnessByte)
	case common.ConfigCellTypeArgsCharSetEn:
		return ParserConfigCellTypeArgsCharSetEn(witnessByte)
	case common.ConfigCellTypeArgsCharSetHanS:
		return ParserConfigCellTypeArgsCharSetHanS(witnessByte)
	case common.ConfigCellTypeArgsCharSetHanT:
		return ParserConfigCellTypeArgsCharSetHanT(witnessByte)

	default:
		return parserDefaultWitness(witnessByte)
	}
}

func parserDefaultWitness(witnessByte []byte) interface{} {
	return map[string]interface{}{
		"unknown": common.Bytes2Hex(witnessByte),
	}
}

func parserData(data *molecule.Data) (dataEntityOpts []map[string]interface{}) {
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

	return dataEntityOpts
}

func parserScript(script *molecule.Script) map[string]interface{} {
	if script == nil {
		return nil
	}

	return map[string]interface{}{
		"code_hash": common.Bytes2Hex(script.CodeHash().RawData()),
		"hash_type": common.Bytes2Hex(script.HashType().AsSlice()),
		"args":      common.Bytes2Hex(script.Args().RawData()),
	}
}

func parserConfig(priceConfig *molecule.PriceConfig) map[string]interface{} {
	if priceConfig == nil {
		return nil
	}

	length, _ := molecule.Bytes2GoU8(priceConfig.Length().RawData())
	newP, _ := molecule.Bytes2GoU64(priceConfig.New().RawData())
	renew, _ := molecule.Bytes2GoU64(priceConfig.Renew().RawData())

	return map[string]interface{}{
		"length": length,
		"new":    newP,
		"renew":  renew,
	}
}

func ParserActionData(witnessByte []byte) interface{} {
	builder, err := ActionDataBuilderFromWitness(witnessByte)
	if err != nil {
		return parserDefaultWitness(witnessByte)
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

func ParserAccountCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	accountCellsData := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		var accountCellData map[string]interface{}
		switch version {
		case common.GoDataEntityVersion1:
			accountCellData = parserAccountCellDataV1(dataEntity)
		case common.GoDataEntityVersion2:
			accountCellData = parserAccountCellData(dataEntity)
		}
		if accountCellData == nil {
			return parserDefaultWitness(witnessByte)
		}
		accountCellsData[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": accountCellData["witness_hash"],
			"entity":       accountCellData["entity"],
		}
	}

	return map[string]interface{}{
		"witness":     common.Bytes2Hex(witnessByte),
		"AccountCell": accountCellsData,
	}
}

func parserAccountCellDataV1(dataEntity *molecule.DataEntity) map[string]interface{} {
	accountCellV1, _ := molecule.AccountCellDataV1FromSlice(dataEntity.Entity().RawData(), false)
	if accountCellV1 == nil {
		return nil
	}

	registeredAt, _ := molecule.Bytes2GoU64(accountCellV1.RegisteredAt().RawData())
	updatedAt, _ := molecule.Bytes2GoU64(accountCellV1.UpdatedAt().RawData())
	status, _ := molecule.Bytes2GoU64(accountCellV1.Status().RawData())
	var recordsMaps []map[string]interface{}
	for i := uint(0); i < accountCellV1.Records().Len(); i++ {
		record := accountCellV1.Records().Get(i)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		recordsMaps = append(recordsMaps, map[string]interface{}{
			"key":   string(record.RecordKey().RawData()),
			"type":  string(record.RecordType().RawData()),
			"label": string(record.RecordLabel().RawData()),
			"value": string(record.RecordValue().RawData()),
			"ttl":   ttl,
		})
	}

	return map[string]interface{}{
		"witness_hash": common.Bytes2Hex(common.Blake2b(accountCellV1.AsSlice())),
		"entity": map[string]interface{}{
			"id":            common.Bytes2Hex(accountCellV1.Id().RawData()),
			"account":       common.AccountCharsToAccount(accountCellV1.Account()),
			"registered_at": registeredAt,
			"updated_at":    updatedAt,
			"status":        status,
			"records":       recordsMaps,
		},
	}
}

func parserAccountCellData(dataEntity *molecule.DataEntity) map[string]interface{} {
	accountCell, _ := molecule.AccountCellDataFromSlice(dataEntity.Entity().RawData(), false)
	if accountCell == nil {
		return nil
	}

	registeredAt, _ := molecule.Bytes2GoU64(accountCell.RegisteredAt().RawData())
	lastTransferAccountAt, _ := molecule.Bytes2GoU64(accountCell.LastTransferAccountAt().RawData())
	lastEditManagerAt, _ := molecule.Bytes2GoU64(accountCell.LastEditManagerAt().RawData())
	lastEditRecordsAt, _ := molecule.Bytes2GoU64(accountCell.LastEditRecordsAt().RawData())
	status, _ := molecule.Bytes2GoU64(accountCell.Status().RawData())
	var recordsMaps []map[string]interface{}
	for i := uint(0); i < accountCell.Records().Len(); i++ {
		record := accountCell.Records().Get(i)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		recordsMaps = append(recordsMaps, map[string]interface{}{
			"key":   string(record.RecordKey().RawData()),
			"type":  string(record.RecordType().RawData()),
			"label": string(record.RecordLabel().RawData()),
			"value": string(record.RecordValue().RawData()),
			"ttl":   ttl,
		})
	}

	return map[string]interface{}{
		"witness_hash": common.Bytes2Hex(common.Blake2b(accountCell.AsSlice())),
		"entity": map[string]interface{}{
			"id":                       common.Bytes2Hex(accountCell.Id().RawData()),
			"account":                  common.AccountCharsToAccount(accountCell.Account()),
			"registered_at":            registeredAt,
			"last_transfer_account_at": lastTransferAccountAt,
			"last_edit_manager_at":     lastEditManagerAt,
			"last_edit_records_at":     lastEditRecordsAt,
			"status":                   status,
			"records":                  recordsMaps,
		},
	}
}

func ParserAccountSaleCell(witnessByte []byte) interface{} {
	accountSaleCell, _ := molecule.AccountSaleCellDataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if accountSaleCell != nil {
		return map[string]interface{}{
			"witness":         common.Bytes2Hex(witnessByte),
			"witness_hash":    common.Bytes2Hex(common.Blake2b(accountSaleCell.AsSlice())),
			"AccountSaleCell": map[string]interface{}{},
		}
	}

	accountSaleCellV1, _ := molecule.AccountSaleCellDataV1FromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if accountSaleCellV1 != nil {
		return map[string]interface{}{
			"witness":         common.Bytes2Hex(witnessByte),
			"witness_hash":    common.Bytes2Hex(common.Blake2b(accountSaleCellV1.AsSlice())),
			"AccountSaleCell": map[string]interface{}{},
		}
	}
	return parserDefaultWitness(witnessByte)
}

func ParserAccountAuctionCell(witnessByte []byte) interface{} {
	accountAuctionCell, _ := molecule.AccountAuctionCellDataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if accountAuctionCell == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":            common.Bytes2Hex(witnessByte),
		"witness_hash":       common.Bytes2Hex(common.Blake2b(accountAuctionCell.AsSlice())),
		"AccountAuctionCell": map[string]interface{}{},
	}
}

func ParserProposalCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	proposalCellsData := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		proposalCell, _ := molecule.ProposalCellDataFromSlice(dataEntity.Entity().RawData(), false)
		if proposalCell == nil {
			return parserDefaultWitness(witnessByte)
		}

		proposalLock, _ := molecule.ScriptFromSlice(proposalCell.ProposerLock().AsSlice(), false)
		createdAtHeight, _ := molecule.Bytes2GoU64(proposalCell.CreatedAtHeight().RawData())
		var slices []interface{}
		for i := uint(0); i < proposalCell.Slices().Len(); i++ {
			slice := proposalCell.Slices().Get(i)
			var proposalItems []interface{}
			for k := uint(0); k < slice.Len(); k++ {
				proposalItem := slice.Get(k)
				id := common.Bytes2Hex(proposalItem.AccountId().RawData())
				itemType, _ := molecule.Bytes2GoU8(proposalItem.ItemType().RawData())
				next := common.Bytes2Hex(proposalItem.Next().RawData())
				proposalItems = append(proposalItems, map[string]interface{}{
					"id":        id,
					"item_type": itemType,
					"next":      next,
				})
			}
			slices = append(slices, proposalItems)
		}
		proposalCellsData[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": common.Bytes2Hex(common.Blake2b(proposalCell.AsSlice())),
			"entity": map[string]interface{}{
				"proposal_lock":     parserScript(proposalLock),
				"created_at_height": createdAtHeight,
				"slices":            slices,
			},
		}
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"ProposalCell": proposalCellsData,
	}
}

func ParserPreAccountCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	var preAccountCellData []interface{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		preAccountCell, _ := molecule.PreAccountCellDataFromSlice(dataEntity.Entity().RawData(), false)
		if preAccountCell == nil {
			return parserDefaultWitness(witnessByte)
		}

		refundLock, _ := molecule.ScriptFromSlice(preAccountCell.RefundLock().AsSlice(), false)
		inviterLock, _ := molecule.ScriptFromSlice(preAccountCell.InviterLock().AsSlice(), false)
		channelLock, _ := molecule.ScriptFromSlice(preAccountCell.ChannelLock().AsSlice(), false)
		price, _ := molecule.PriceConfigFromSlice(preAccountCell.Price().AsSlice(), false)
		quote, _ := molecule.Bytes2GoU64(preAccountCell.Quote().RawData())
		invitedDiscount, _ := molecule.Bytes2GoU32(preAccountCell.InvitedDiscount().RawData())
		createdAt, _ := molecule.Bytes2GoU64(preAccountCell.CreatedAt().RawData())

		preAccountCellData = append(preAccountCellData, map[string]interface{}{
			v["type"].(string): map[string]interface{}{
				"version":      version,
				"index":        index,
				"witness_hash": common.Bytes2Hex(common.Blake2b(preAccountCell.AsSlice())),
				"entity": map[string]interface{}{
					"account":          common.AccountCharsToAccount(preAccountCell.Account()),
					"owner_lock_args":  common.Bytes2Hex(preAccountCell.OwnerLockArgs().RawData()),
					"inviter_id":       common.Bytes2Hex(preAccountCell.InviterId().RawData()),
					"refund_lock":      parserScript(refundLock),
					"inviter_lock":     parserScript(inviterLock),
					"channel_lock":     parserScript(channelLock),
					"price":            parserConfig(price),
					"quote":            quote,
					"invited_discount": invitedDiscount,
					"created_at":       createdAt,
				},
			},
		})
	}

	return map[string]interface{}{
		"witness":        common.Bytes2Hex(witnessByte),
		"PreAccountCell": preAccountCellData,
	}
}

func ParserIncomeCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	var incomeCellData []interface{}
	for _, v := range parserData(data) {
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
				"version":      version,
				"index":        index,
				"witness_hash": common.Bytes2Hex(common.Blake2b(incomeCell.AsSlice())),
				"entity": map[string]interface{}{
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

func ParserOfferCell(witnessByte []byte) interface{} {
	offerCell, _ := molecule.OfferCellDataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if offerCell == nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(offerCell.AsSlice())),
		"OfferCell":    map[string]interface{}{},
	}
}

func ParserConfigCellAccount(witnessByte []byte) interface{} {
	configCellAccount, _ := molecule.ConfigCellAccountFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
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

func ParserConfigCellApply(witnessByte []byte) interface{} {
	configCellApply, _ := molecule.ConfigCellApplyFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
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

func ParserConfigCellIncome(witnessByte []byte) interface{} {
	configCellIncome, _ := molecule.ConfigCellIncomeFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
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

func ParserConfigCellMain(witnessByte []byte) interface{} {
	configCellMain, _ := molecule.ConfigCellMainFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
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

func ParserConfigCellPrice(witnessByte []byte) interface{} {
	configCellPrice, _ := molecule.ConfigCellPriceFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellPrice == nil {
		return parserDefaultWitness(witnessByte)
	}

	var prices []interface{}
	for i := uint(0); i < configCellPrice.Prices().Len(); i++ {
		price, _ := molecule.PriceConfigFromSlice(configCellPrice.Prices().Get(i).AsSlice(), false)
		prices = append(prices, parserConfig(price))
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

func ParserConfigCellProposal(witnessByte []byte) interface{} {
	configCellProposal, _ := molecule.ConfigCellProposalFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
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

func ParserConfigCellProfitRate(witnessByte []byte) interface{} {
	configCellProfitRate, _ := molecule.ConfigCellProfitRateFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
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

func ParserConfigCellRecordNamespace(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"ConfigCellRecordNamespace": map[string]interface{}{
			"length":                       dataLength,
			"config_cell_record_namespace": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellRelease(witnessByte []byte) interface{} {
	configCellRelease, _ := molecule.ConfigCellReleaseFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellRelease == nil {
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

func ParserConfigCellUnavailable(witnessByte []byte, action string) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		action: map[string]interface{}{
			"length": dataLength,
		},
	}
}

func ParserConfigCellSecondaryMarket(witnessByte []byte) interface{} {
	configCellSecondaryMarket, _ := molecule.ConfigCellSecondaryMarketFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
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

func ParserConfigCellReverseRecord(witnessByte []byte) interface{} {
	configCellReverseRecord, _ := molecule.ConfigCellReverseResolutionFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
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

func ParserConfigCellTypeArgsCharSetEmoji(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"ConfigCellTypeArgsCharSetEmoji": map[string]interface{}{
			"length":            dataLength,
			"config_cell_emoji": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellTypeArgsCharSetDigit(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"ConfigCellTypeArgsCharSetDigit": map[string]interface{}{
			"length":            dataLength,
			"config_cell_digit": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellTypeArgsCharSetEn(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"ConfigCellTypeArgsCharSetEn": map[string]interface{}{
			"length":         dataLength,
			"config_cell_en": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellTypeArgsCharSetHanS(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"ConfigCellTypeArgsCharSetHanS": map[string]interface{}{
			"length":            dataLength,
			"config_cell_han_s": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellTypeArgsCharSetHanT(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"ConfigCellTypeArgsCharSetHanT": map[string]interface{}{
			"length":            dataLength,
			"config_cell_han_t": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}
