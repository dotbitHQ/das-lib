package core

import (
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"time"
)

type CacheConfigCellCharSet struct {
	ConfigCellEmojis       []string `json:"config_cell_emojis"`
	ConfigCellCharSetDigit []string `json:"config_cell_char_set_digit"`
	ConfigCellCharSetEn    []string `json:"config_cell_char_set_en"`
	ConfigCellCharSetHanS  []string `json:"config_cell_char_set_han_s"`
	ConfigCellCharSetHanT  []string `json:"config_cell_char_set_han_t"`
	ConfigCellCharSetJa    []string `json:"config_cell_char_set_ja"`
	ConfigCellCharSetKo    []string `json:"config_cell_char_set_ko"`
	ConfigCellCharSetRu    []string `json:"config_cell_char_set_ru"`
	ConfigCellCharSetTr    []string `json:"config_cell_char_set_tr"`
	ConfigCellCharSetTh    []string `json:"config_cell_char_set_th"`
	ConfigCellCharSetVi    []string `json:"config_cell_char_set_vi"`
}

type CacheConfigCellReservedAccounts struct {
	MapReservedAccounts    map[string]struct{} `json:"map_reserved_accounts"`
	MapUnAvailableAccounts map[string]struct{} `json:"map_un_available_accounts"`
}

type CacheConfigCellBase struct {
	LuckyNumber uint32 `json:"lucky_number"`

	RecordBasicCapacity       uint64 `json:"record_basic_capacity"`
	RecordPreparedFeeCapacity uint64 `json:"record_prepared_fee_capacity"`

	SaleCellBasicCapacity       uint64 `json:"sale_cell_basic_capacity"`
	SaleCellPreparedFeeCapacity uint64 `json:"sale_cell_prepared_fee_capacity"`
	SaleMinPrice                uint64 `json:"sale_min_price"`

	ProfitRateInviter uint32 `json:"profit_rate_inviter"`

	IncomeMinTransferCapacity uint64 `json:"income_min_transfer_capacity"`

	PriceInvitedDiscount uint32 `json:"price_invited_discount"`

	TransferAccountThrottle uint32 `json:"transfer_account_throttle"`
	EditManagerThrottle     uint32 `json:"edit_manager_throttle"`
	EditRecordsThrottle     uint32 `json:"edit_records_throttle"`
	MaxLength               uint32 `json:"max_length"`
	RecordMinTtl            uint32 `json:"record_min_ttl"`
	ExpirationGracePeriod   uint32 `json:"expiration_grace_period"`
}

type CacheConfigCellKey = string

const (
	CacheConfigCellKeyBase             CacheConfigCellKey = "CacheConfigCellKeyBase"
	CacheConfigCellKeyCharSet          CacheConfigCellKey = "CacheConfigCellKeyCharSet"
	CacheConfigCellKeyReservedAccounts CacheConfigCellKey = "CacheConfigCellKeyReservedAccounts"
)

func (d *DasCore) RunSetConfigCellByCache(keyList []CacheConfigCellKey) {
	ticUpdate := time.NewTicker(time.Minute * 10)
	d.wg.Add(1)
	go func() {
		for {
			select {
			case <-ticUpdate.C:
				log.Info("RunSetConfigCellByCache start")
				for _, v := range keyList {
					cacheStr := ""
					switch v {
					case CacheConfigCellKeyBase:
						builder, err := d.ConfigCellDataBuilderByTypeArgsList(
							common.ConfigCellTypeArgsAccount,
							common.ConfigCellTypeArgsPrice,
							common.ConfigCellTypeArgsIncome,
							common.ConfigCellTypeArgsProfitRate,
							common.ConfigCellTypeArgsSecondaryMarket,
							common.ConfigCellTypeArgsReverseRecord,
							common.ConfigCellTypeArgsRelease,
						)
						if err != nil {
							log.Error("ConfigCellDataBuilderByTypeArgsList err: ", err.Error(), v)
						} else {
							var cacheBuilder CacheConfigCellBase
							//
							cacheBuilder.LuckyNumber, _ = builder.LuckyNumber()

							cacheBuilder.RecordPreparedFeeCapacity, _ = builder.RecordPreparedFeeCapacity()
							cacheBuilder.RecordBasicCapacity, _ = builder.RecordBasicCapacity()

							cacheBuilder.SaleCellBasicCapacity, _ = builder.SaleCellBasicCapacity()
							cacheBuilder.SaleCellPreparedFeeCapacity, _ = builder.SaleCellPreparedFeeCapacity()
							cacheBuilder.SaleMinPrice, _ = builder.SaleMinPrice()

							cacheBuilder.ProfitRateInviter, _ = builder.ProfitRateInviter()

							cacheBuilder.IncomeMinTransferCapacity, _ = builder.IncomeMinTransferCapacity()

							cacheBuilder.PriceInvitedDiscount, _ = builder.PriceInvitedDiscount()

							cacheBuilder.TransferAccountThrottle, _ = builder.TransferAccountThrottle()
							cacheBuilder.EditManagerThrottle, _ = builder.EditManagerThrottle()
							cacheBuilder.EditRecordsThrottle, _ = builder.EditRecordsThrottle()
							cacheBuilder.MaxLength, _ = builder.MaxLength()
							cacheBuilder.RecordMinTtl, _ = builder.RecordMinTtl()
							cacheBuilder.ExpirationGracePeriod, _ = builder.ExpirationGracePeriod()

							cacheStrBys, _ := json.Marshal(&cacheBuilder)
							cacheStr = string(cacheStrBys)
						}
					case CacheConfigCellKeyReservedAccounts:
						builderConfigCell, err := d.ConfigCellDataBuilderByTypeArgsList(
							common.ConfigCellTypeArgsPreservedAccount00,
							common.ConfigCellTypeArgsPreservedAccount01,
							common.ConfigCellTypeArgsPreservedAccount02,
							common.ConfigCellTypeArgsPreservedAccount03,
							common.ConfigCellTypeArgsPreservedAccount04,
							common.ConfigCellTypeArgsPreservedAccount05,
							common.ConfigCellTypeArgsPreservedAccount06,
							common.ConfigCellTypeArgsPreservedAccount07,
							common.ConfigCellTypeArgsPreservedAccount08,
							common.ConfigCellTypeArgsPreservedAccount09,
							common.ConfigCellTypeArgsPreservedAccount10,
							common.ConfigCellTypeArgsPreservedAccount11,
							common.ConfigCellTypeArgsPreservedAccount12,
							common.ConfigCellTypeArgsPreservedAccount13,
							common.ConfigCellTypeArgsPreservedAccount14,
							common.ConfigCellTypeArgsPreservedAccount15,
							common.ConfigCellTypeArgsPreservedAccount16,
							common.ConfigCellTypeArgsPreservedAccount17,
							common.ConfigCellTypeArgsPreservedAccount18,
							common.ConfigCellTypeArgsPreservedAccount19,
							common.ConfigCellTypeArgsUnavailable,
						)
						if err != nil {
							log.Error("ConfigCellDataBuilderByTypeArgsList err: ", err.Error(), v)
						} else {
							var cacheBuilder CacheConfigCellReservedAccounts
							cacheBuilder.MapReservedAccounts = builderConfigCell.ConfigCellPreservedAccountMap
							cacheBuilder.MapUnAvailableAccounts = builderConfigCell.ConfigCellUnavailableAccountMap
							cacheStrBys, _ := json.Marshal(&cacheBuilder)
							cacheStr = string(cacheStrBys)
						}
					case CacheConfigCellKeyCharSet:
						builder, err := d.ConfigCellDataBuilderByTypeArgsList(
							common.ConfigCellTypeArgsCharSetEmoji,
							common.ConfigCellTypeArgsCharSetDigit,
							common.ConfigCellTypeArgsCharSetEn,
							common.ConfigCellTypeArgsCharSetHanS,
							common.ConfigCellTypeArgsCharSetHanT,
							common.ConfigCellTypeArgsCharSetJa,
							common.ConfigCellTypeArgsCharSetKo,
							common.ConfigCellTypeArgsCharSetRu,
							common.ConfigCellTypeArgsCharSetTr,
							common.ConfigCellTypeArgsCharSetTh,
							common.ConfigCellTypeArgsCharSetVi,
						)
						if err != nil {
							log.Error("ConfigCellDataBuilderByTypeArgsList err: ", err.Error(), v)
						} else {
							var cacheBuilder CacheConfigCellCharSet
							cacheBuilder.ConfigCellEmojis = builder.ConfigCellEmojis
							cacheBuilder.ConfigCellCharSetDigit = builder.ConfigCellCharSetDigit
							cacheBuilder.ConfigCellCharSetEn = builder.ConfigCellCharSetEn
							cacheBuilder.ConfigCellCharSetHanS = builder.ConfigCellCharSetHanS
							cacheBuilder.ConfigCellCharSetHanT = builder.ConfigCellCharSetHanT
							cacheBuilder.ConfigCellCharSetJa = builder.ConfigCellCharSetJa
							cacheBuilder.ConfigCellCharSetKo = builder.ConfigCellCharSetKo
							cacheBuilder.ConfigCellCharSetRu = builder.ConfigCellCharSetRu
							cacheBuilder.ConfigCellCharSetTr = builder.ConfigCellCharSetTr
							cacheBuilder.ConfigCellCharSetTh = builder.ConfigCellCharSetTh
							cacheBuilder.ConfigCellCharSetVi = builder.ConfigCellCharSetVi
							cacheStrBys, _ := json.Marshal(&cacheBuilder)
							cacheStr = string(cacheStrBys)
						}
					}
					if err := d.setConfigCellByCache(v, cacheStr); err != nil {
						log.Error("setConfigCellByCache err:", err.Error(), v)
					}
				}
				log.Info("RunSetConfigCellByCache ok")
			case <-d.ctx.Done():
				log.Info("RunSetConfigCellByCache Done")
				d.wg.Done()
				return
			}
		}
	}()
}

func (d *DasCore) setConfigCellByCache(key CacheConfigCellKey, value string) error {
	if d.red == nil {
		return fmt.Errorf("d.red is nil")
	}
	if value == "" {
		return nil
	}
	if err := d.red.Set(key, value, 0).Err(); err != nil {
		return fmt.Errorf("d.red.Set err: %s [%s]", err.Error(), key)
	}
	return nil
}

func (d *DasCore) GetConfigCellByCache(key CacheConfigCellKey) (string, error) {
	if d.red == nil {
		return "", fmt.Errorf("d.red is nil")
	}
	res, err := d.red.Get(key).Result()
	if err != nil {
		return "", fmt.Errorf("d.red.Get err: %s [%s]", err.Error(), key)
	}
	return res, nil
}
