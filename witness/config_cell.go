package witness

//
//// ConfigCellDataBuilderRefByTypeArgs Deprecated, GetConfigCellDataBuilderRefByTxCon
//func ConfigCellDataBuilderRefByTypeArgs(builder *ConfigCellDataBuilder, tx *types.Transaction, configCellTypeArgs common.ConfigCellTypeArgs) error {
//	var configCellDataBys []byte
//	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte, index int) (bool, error) {
//		if actionDataType == configCellTypeArgs {
//			configCellDataBys = dataBys
//			return false, nil
//		}
//		return true, nil
//	})
//	if err != nil {
//		return fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
//	}
//
//	switch configCellTypeArgs {
//	case common.ConfigCellTypeArgsDPoint:
//		configCellTypeArgsDPoint, err := molecule.ConfigCellDPointFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellDPointFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellDPoint = configCellTypeArgsDPoint
//	case common.ConfigCellTypeArgsAccount:
//		ConfigCellAccount, err := molecule.ConfigCellAccountFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellAccountFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellAccount = ConfigCellAccount
//	case common.ConfigCellTypeArgsPrice:
//		ConfigCellPrice, err := molecule.ConfigCellPriceFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellPriceFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellPrice = ConfigCellPrice
//		builder.PriceConfigMap = make(map[uint8]*molecule.PriceConfig)
//		prices := builder.ConfigCellPrice.Prices()
//		for i, count := uint(0), prices.Len(); i < count; i++ {
//			price, err := molecule.PriceConfigFromSlice(prices.Get(i).AsSlice(), true)
//			if err != nil {
//				return fmt.Errorf("PriceConfigFromSlice err: %s", err.Error())
//			}
//			length, err := molecule.Bytes2GoU8(price.Length().RawData())
//			if err != nil {
//				return fmt.Errorf("price.Length() err: %s", err.Error())
//			}
//			if builder.PriceMaxLength < length {
//				builder.PriceMaxLength = length
//			}
//			builder.PriceConfigMap[length] = price
//		}
//	case common.ConfigCellTypeArgsApply:
//		ConfigCellApply, err := molecule.ConfigCellApplyFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellApply = ConfigCellApply
//	case common.ConfigCellTypeArgsRelease:
//		ConfigCellRelease, err := molecule.ConfigCellReleaseFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellRelease = ConfigCellRelease
//	case common.ConfigCellTypeArgsSecondaryMarket:
//		ConfigCellSecondaryMarket, err := molecule.ConfigCellSecondaryMarketFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellSecondaryMarketFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellSecondaryMarket = ConfigCellSecondaryMarket
//	case common.ConfigCellTypeArgsIncome:
//		ConfigCellIncome, err := molecule.ConfigCellIncomeFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellIncomeFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellIncome = ConfigCellIncome
//	case common.ConfigCellTypeArgsProfitRate:
//		ConfigCellProfitRate, err := molecule.ConfigCellProfitRateFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellProfitRate = ConfigCellProfitRate
//	case common.ConfigCellTypeArgsMain:
//		ConfigCellMain, err := molecule.ConfigCellMainFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellMainFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellMain = ConfigCellMain
//	case common.ConfigCellTypeArgsReverseRecord:
//		ConfigCellReverseResolution, err := molecule.ConfigCellReverseResolutionFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellReverseResolutionFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellReverseResolution = ConfigCellReverseResolution
//	case common.ConfigCellTypeArgsSubAccount:
//		ConfigCellSubAccount, err := molecule.ConfigCellSubAccountFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellSubAccountFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellSubAccount = ConfigCellSubAccount
//	case common.ConfigCellTypeArgsProposal:
//		ConfigCellProposal, err := molecule.ConfigCellProposalFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellProposalFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellProposal = ConfigCellProposal
//	case common.ConfigCellTypeArgsRecordNamespace:
//		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//		if err != nil {
//			return fmt.Errorf("key name space len err: %s", err.Error())
//		}
//		builder.ConfigCellRecordKeys = strings.Split(string(configCellDataBys[4:dataLength]), string([]byte{0x00}))
//
//	case common.ConfigCellTypeArgsCharSetEmoji:
//		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//		if err != nil {
//			return fmt.Errorf("char set emoji err: %s", err.Error())
//		}
//		builder.ConfigCellEmojis = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//	case common.ConfigCellTypeArgsCharSetDigit:
//		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//		if err != nil {
//			return fmt.Errorf("char set digit err: %s", err.Error())
//		}
//		builder.ConfigCellCharSetDigit = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//	case common.ConfigCellTypeArgsCharSetEn:
//		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//		if err != nil {
//			return fmt.Errorf("char set en err: %s", err.Error())
//		}
//		builder.ConfigCellCharSetEn = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//	case common.ConfigCellTypeArgsCharSetHanS:
//		if len(configCellDataBys) != 0 {
//			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//			if err != nil {
//				return fmt.Errorf("char set hans err: %s", err.Error())
//			}
//			builder.ConfigCellCharSetHanS = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//		}
//	case common.ConfigCellTypeArgsCharSetHanT:
//		if len(configCellDataBys) != 0 {
//			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//			if err != nil {
//				return fmt.Errorf("char set hant err: %s", err.Error())
//			}
//			builder.ConfigCellCharSetHanT = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//		}
//	case common.ConfigCellTypeArgsCharSetJa:
//		if len(configCellDataBys) != 0 {
//			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//			if err != nil {
//				return fmt.Errorf("char set jp err: %s", err.Error())
//			}
//			builder.ConfigCellCharSetJa = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//		}
//	case common.ConfigCellTypeArgsCharSetKo:
//		if len(configCellDataBys) != 0 {
//			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//			if err != nil {
//				return fmt.Errorf("char set kr err: %s", err.Error())
//			}
//			builder.ConfigCellCharSetKo = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//		}
//	case common.ConfigCellTypeArgsCharSetRu:
//		if len(configCellDataBys) != 0 {
//			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//			if err != nil {
//				return fmt.Errorf("char set ru err: %s", err.Error())
//			}
//			builder.ConfigCellCharSetRu = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//		}
//	case common.ConfigCellTypeArgsCharSetTr:
//		if len(configCellDataBys) != 0 {
//			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//			if err != nil {
//				return fmt.Errorf("char set tr err: %s", err.Error())
//			}
//			builder.ConfigCellCharSetTr = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//		}
//	case common.ConfigCellTypeArgsCharSetTh:
//		if len(configCellDataBys) != 0 {
//			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//			if err != nil {
//				return fmt.Errorf("char set th err: %s", err.Error())
//			}
//			builder.ConfigCellCharSetTh = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//		}
//	case common.ConfigCellTypeArgsCharSetVi:
//		if len(configCellDataBys) != 0 {
//			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//			if err != nil {
//				return fmt.Errorf("char set vn err: %s", err.Error())
//			}
//			builder.ConfigCellCharSetVi = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
//		}
//	case common.ConfigCellTypeArgsUnavailable:
//		if builder.ConfigCellUnavailableAccountMap == nil {
//			builder.ConfigCellUnavailableAccountMap = make(map[string]struct{})
//		}
//		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//		if err != nil {
//			return fmt.Errorf("unavailable account err: %s", err.Error())
//		}
//		for i := 20; i <= len(configCellDataBys[4:dataLength]); i += 20 {
//			tmp := common.Bytes2Hex(configCellDataBys[4:dataLength][i-20 : i])
//			if _, ok := builder.ConfigCellUnavailableAccountMap[tmp]; ok {
//				fmt.Println(tmp, "ok")
//			}
//			builder.ConfigCellUnavailableAccountMap[tmp] = struct{}{}
//		}
//	case common.ConfigCellTypeArgsSubAccountWhiteList:
//		if builder.ConfigCellSubAccountWhiteListMap == nil {
//			builder.ConfigCellSubAccountWhiteListMap = make(map[string]struct{})
//		}
//		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//		if err != nil {
//			return fmt.Errorf("preserved account err: %s", err.Error())
//		}
//		for i := 20; i <= len(configCellDataBys[4:dataLength]); i += 20 {
//			tmp := common.Bytes2Hex(configCellDataBys[4:dataLength][i-20 : i])
//			builder.ConfigCellSubAccountWhiteListMap[tmp] = struct{}{}
//		}
//	case common.ConfigCellTypeArgsSystemStatus:
//		configCellSystemStatus, err := molecule.ConfigCellSystemStatusFromSlice(configCellDataBys, true)
//		if err != nil {
//			return fmt.Errorf("ConfigCellSystemStatusFromSlice err: %s", err.Error())
//		}
//		builder.ConfigCellSystemStatus = configCellSystemStatus
//	case common.ConfigCellTypeArgsPreservedAccount00,
//		common.ConfigCellTypeArgsPreservedAccount01,
//		common.ConfigCellTypeArgsPreservedAccount02,
//		common.ConfigCellTypeArgsPreservedAccount03,
//		common.ConfigCellTypeArgsPreservedAccount04,
//		common.ConfigCellTypeArgsPreservedAccount05,
//		common.ConfigCellTypeArgsPreservedAccount06,
//		common.ConfigCellTypeArgsPreservedAccount07,
//		common.ConfigCellTypeArgsPreservedAccount08,
//		common.ConfigCellTypeArgsPreservedAccount09,
//		common.ConfigCellTypeArgsPreservedAccount10,
//		common.ConfigCellTypeArgsPreservedAccount11,
//		common.ConfigCellTypeArgsPreservedAccount12,
//		common.ConfigCellTypeArgsPreservedAccount13,
//		common.ConfigCellTypeArgsPreservedAccount14,
//		common.ConfigCellTypeArgsPreservedAccount15,
//		common.ConfigCellTypeArgsPreservedAccount16,
//		common.ConfigCellTypeArgsPreservedAccount17,
//		common.ConfigCellTypeArgsPreservedAccount18,
//		common.ConfigCellTypeArgsPreservedAccount19:
//		if builder.ConfigCellPreservedAccountMap == nil {
//			builder.ConfigCellPreservedAccountMap = make(map[string]struct{})
//		}
//		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
//		if err != nil {
//			return fmt.Errorf("preserved account err: %s", err.Error())
//		}
//		for i := 20; i <= len(configCellDataBys[4:dataLength]); i += 20 {
//			tmp := common.Bytes2Hex(configCellDataBys[4:dataLength][i-20 : i])
//			builder.ConfigCellPreservedAccountMap[tmp] = struct{}{}
//		}
//	}
//	return nil
//}
//
//// ConfigCellDataBuilderByTypeArgs Deprecated, GetConfigCellDataBuilderByTx
//func ConfigCellDataBuilderByTypeArgs(tx *types.Transaction, configCellTypeArgs common.ConfigCellTypeArgs) (*ConfigCellDataBuilder, error) {
//	var resp ConfigCellDataBuilder
//
//	err := ConfigCellDataBuilderRefByTypeArgs(&resp, tx, configCellTypeArgs)
//	if err != nil {
//		return nil, fmt.Errorf("ConfigCellDataBuilderRefByTypeArgs err: %s", err.Error())
//	}
//
//	return &resp, nil
//}
//
