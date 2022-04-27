package core

import (
	"encoding/hex"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"regexp"
	"strings"
)

type DasAddressNormal struct {
	ChainType     common.ChainType
	AddressNormal string
	Is712         bool
}

type DasAddressHex struct {
	DasAlgorithmId common.DasAlgorithmId
	AddressHex     string
	IsMulti        bool
	ChainType      common.ChainType
}

type DasAddressFormat struct {
	DasNetType common.DasNetType
}

func (d *DasAddressFormat) NormalToHex(p DasAddressNormal) (r DasAddressHex, e error) {
	r.ChainType = p.ChainType
	switch p.ChainType {
	case common.ChainTypeCkb:
		r.DasAlgorithmId = common.DasAlgorithmIdCkb
		if parseAddr, err := address.Parse(p.AddressNormal); err != nil {
			e = fmt.Errorf("address.Parse err: %s", err.Error())
		} else {
			r.AddressHex = common.Bytes2Hex(parseAddr.Script.Args)
			if parseAddr.Script.CodeHash.Hex() == transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH {
				r.IsMulti = true
			}
		}
	case common.ChainTypeCkbDas:
		if parseAddr, err := address.Parse(p.AddressNormal); err != nil {
			e = fmt.Errorf("address.Parse err: %s", err.Error())
		} else {
			r.AddressHex = common.Bytes2Hex(parseAddr.Script.Args)
			if parseAddr.Script.CodeHash.Hex() == transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH {
				r.IsMulti = true
				r.DasAlgorithmId = common.DasAlgorithmIdCkbMulti
			} else {
				e = fmt.Errorf("not support single sign")
			}
		}
	case common.ChainTypeEth:
		r.DasAlgorithmId = common.DasAlgorithmIdEth
		if p.Is712 {
			r.DasAlgorithmId = common.DasAlgorithmIdEth712
		}
		if ok, err := regexp.MatchString("^0x[0-9a-fA-F]{40}$", p.AddressNormal); err != nil {
			e = fmt.Errorf("regexp.MatchString err: %s", err.Error())
		} else if ok {
			r.AddressHex = p.AddressNormal
		} else {
			e = fmt.Errorf("regexp.MatchString fail")
		}
	case common.ChainTypeMixin:
		r.DasAlgorithmId = common.DasAlgorithmIdEd25519
		if ok, err := regexp.MatchString("^0x[0-9a-fA-F]{64}$", p.AddressNormal); err != nil {
			e = fmt.Errorf("regexp.MatchString err: %s", err.Error())
		} else if ok {
			r.AddressHex = p.AddressNormal
		} else {
			e = fmt.Errorf("regexp.MatchString fail")
		}
	case common.ChainTypeTron:
		r.DasAlgorithmId = common.DasAlgorithmIdTron
		if strings.HasPrefix(p.AddressNormal, common.TronBase58PreFix) {
			if addrHex, err := common.TronBase58ToHex(p.AddressNormal); err != nil {
				e = fmt.Errorf("TronBase58ToHex err: %s", err.Error())
			} else {
				r.AddressHex = addrHex
			}
		} else if strings.HasPrefix(p.AddressNormal, common.TronPreFix) {
			if _, err := common.TronHexToBase58(p.AddressNormal); err != nil {
				e = fmt.Errorf("TronHexToBase58 err: %s", err.Error())
			} else {
				r.AddressHex = p.AddressNormal
			}
		}
	default:
		e = fmt.Errorf("not support chain type [%d]", p.ChainType)
	}
	return
}
func (d *DasAddressFormat) NormalToScript(p DasAddressNormal) (lockScript, typeScript *types.Script, e error) {
	if p.ChainType == common.ChainTypeCkb {
		if parseAddr, err := address.Parse(p.AddressNormal); err != nil {
			e = fmt.Errorf("address.Parse err: %s", err.Error())
			return
		} else {
			lockScript = parseAddr.Script
			return
		}
	}

	addrHex, err := d.NormalToHex(p)
	if err != nil {
		e = fmt.Errorf("NormalToHex err: %s", err.Error())
		return
	}
	if p.Is712 && addrHex.DasAlgorithmId == common.DasAlgorithmIdEth712 {
		contractBalance, err := GetDasContractInfo(common.DasContractNameBalanceCellType)
		if err != nil {
			e = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
			return
		}
		typeScript = contractBalance.ToScript(nil)
	}

	halfArgs, err := d.hexToHalfArgs(addrHex)
	if err != nil {
		e = fmt.Errorf("hexToHalfArgs err: %s", err.Error())
		return
	}
	args := append(halfArgs, halfArgs...)

	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		e = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}
	lockScript = contractDispatch.ToScript(args)

	return
}

func (d *DasAddressFormat) HexToNormal(p DasAddressHex) (r DasAddressNormal, e error) {
	switch p.DasAlgorithmId {
	case common.DasAlgorithmIdCkb:
		r.ChainType = common.ChainTypeCkb
		script := common.GetNormalLockScript(p.AddressHex)
		mode := address.Mainnet
		if d.DasNetType != common.DasNetTypeMainNet {
			mode = address.Testnet
		}
		if addr, err := common.ConvertScriptToAddress(mode, script); err != nil {
			e = fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
		} else {
			r.AddressNormal = addr
		}
	case common.DasAlgorithmIdCkbMulti:
		r.ChainType = common.ChainTypeCkbDas
		script := common.GetNormalLockScriptByMultiSig(p.AddressHex)
		mode := address.Mainnet
		if d.DasNetType != common.DasNetTypeMainNet {
			mode = address.Testnet
		}
		if addr, err := common.ConvertScriptToAddress(mode, script); err != nil {
			e = fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
		} else {
			r.AddressNormal = addr
		}
	case common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712:
		r.ChainType = common.ChainTypeEth
		r.Is712 = p.DasAlgorithmId == common.DasAlgorithmIdEth712
		if ok, err := regexp.MatchString("^0x[0-9a-fA-F]{40}$", p.AddressHex); err != nil {
			e = fmt.Errorf("regexp.MatchString err: %s", err.Error())
		} else if ok {
			r.AddressNormal = p.AddressHex
		} else {
			e = fmt.Errorf("regexp.MatchString fail")
		}
	case common.DasAlgorithmIdTron:
		r.ChainType = common.ChainTypeTron
		if addr, err := common.TronHexToBase58(p.AddressHex); err != nil {
			e = fmt.Errorf("TronHexToBase58 err: %s", err.Error())
		} else {
			r.AddressNormal = addr
		}
	case common.DasAlgorithmIdEd25519:
		r.ChainType = common.ChainTypeMixin
		if ok, err := regexp.MatchString("^0x[0-9a-fA-F]{64}$", p.AddressHex); err != nil {
			e = fmt.Errorf("regexp.MatchString err: %s", err.Error())
		} else if ok {
			r.AddressNormal = p.AddressHex
		} else {
			e = fmt.Errorf("regexp.MatchString fail")
		}
	default:
		e = fmt.Errorf("not support DasAlgorithmId [%d]", p.DasAlgorithmId)
	}

	return
}
func (d *DasAddressFormat) hexToHalfArgs(p DasAddressHex) (args []byte, e error) {
	argsStr := ""
	switch p.DasAlgorithmId {
	case common.DasAlgorithmIdCkb:
		argsStr = strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdCkbMulti:
		argsStr = common.DasLockCkbMultiPreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdEth:
		argsStr = common.DasLockEthPreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdTron:
		argsStr = common.DasLockTronPreFix + strings.TrimPrefix(p.AddressHex, common.TronPreFix)
	case common.DasAlgorithmIdEth712:
		argsStr = common.DasLockEth712PreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdEd25519:
		argsStr = common.DasLockEd25519PreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	default:
		e = fmt.Errorf("not support DasAlgorithmId[%d]", p.DasAlgorithmId)
	}
	if argsStr != "" {
		args = common.Hex2Bytes(argsStr)
	}
	return
}

func (d *DasAddressFormat) ArgsToNormal(args []byte) (ownerNormal, managerNormal DasAddressNormal, e error) {
	ownerHex, managerHex, err := d.ArgsToHex(args)
	if err != nil {
		e = fmt.Errorf("ArgsToHex err: %s", err.Error())
	} else {
		if ownerNormal, err = d.HexToNormal(ownerHex); err != nil {
			e = fmt.Errorf("owner HexToNormal err: %s", err.Error())
		} else if managerNormal, err = d.HexToNormal(managerHex); err != nil {
			e = fmt.Errorf("manager HexToNormal err: %s", err.Error())
		}
	}
	return
}
func (d *DasAddressFormat) ArgsToHex(args []byte) (ownerHex, managerHex DasAddressHex, e error) {
	owner, manager, err := d.argsToHalfArgs(args)
	if err != nil {
		e = fmt.Errorf("argsToHalfArgs err: %s", err.Error())
	} else {
		if ownerHex, err = d.halfArgsToHex(owner); err != nil {
			e = fmt.Errorf("owner halfArgsToHex err: %s", err.Error())
		} else if managerHex, err = d.halfArgsToHex(manager); err != nil {
			e = fmt.Errorf("manager halfArgsToHex err: %s", err.Error())
		}
	}
	return
}
func (d *DasAddressFormat) argsToHalfArgs(args []byte) (owner, manager []byte, e error) {
	if len(args) < common.DasLockArgsLen || len(args) > common.DasLockArgsLenMax {
		e = fmt.Errorf("len(args) error")
		return
	}
	oID := common.DasAlgorithmId(args[0])
	splitLen := 0
	switch oID {
	case common.DasAlgorithmIdCkb, common.DasAlgorithmIdCkbMulti, common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712, common.DasAlgorithmIdTron:
		splitLen = common.DasLockArgsLen / 2
	case common.DasAlgorithmIdEd25519:
		splitLen = common.DasLockArgsLenMax / 2
	default:
		return
	}
	owner = args[:splitLen]
	manager = args[splitLen:]
	return
}
func (d *DasAddressFormat) halfArgsToHex(args []byte) (r DasAddressHex, e error) {
	r.DasAlgorithmId = common.DasAlgorithmId(args[0])
	switch r.DasAlgorithmId {
	case common.DasAlgorithmIdCkb:
		r.ChainType = common.ChainTypeCkb
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
	case common.DasAlgorithmIdCkbMulti:
		r.ChainType = common.ChainTypeCkbDas
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
		r.IsMulti = true
	case common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712:
		r.ChainType = common.ChainTypeEth
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
	case common.DasAlgorithmIdTron:
		r.ChainType = common.ChainTypeTron
		r.AddressHex = common.TronPreFix + hex.EncodeToString(args[1:])
	case common.DasAlgorithmIdEd25519:
		r.ChainType = common.ChainTypeMixin
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
	default:
		e = fmt.Errorf("not support DasAlgorithmId [%d]", r.DasAlgorithmId)
	}
	return
}
