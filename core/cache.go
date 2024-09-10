package core

import "fmt"

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

type CacheConfigCellKey = string

const (
	CacheConfigCellKeyCharSet CacheConfigCellKey = "CacheConfigCellKeyCharSet"
)

func (d *DasCore) RunSetConfigCellByCache(key []CacheConfigCellKey) {
	go func() {
		// todo
	}()
}

func (d *DasCore) setConfigCellByCache(key CacheConfigCellKey, value string) error {
	if d.red == nil {
		return fmt.Errorf("d.red is nil")
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
