package bitcoin

func (t *TxTool) GetUnspentOutputsBtc(addr, privateKey string, value int64) (int64, []UnspentOutputs, error) {
	var uos []UnspentOutputs
	total := int64(0)
	return total, uos, nil
}
