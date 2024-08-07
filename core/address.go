package core

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
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
	DasAlgorithmId    common.DasAlgorithmId
	DasSubAlgorithmId common.DasSubAlgorithmId
	AddressHex        string
	AddressPayload    []byte
	IsMulti           bool
	ChainType         common.ChainType // format normal address ckb chain type
	ParsedAddress     *address.ParsedAddress
}

func (d *DasAddressHex) Payload() string {
	return hex.EncodeToString(d.AddressPayload)
}

func (d *DasAddressHex) FormatAnyLock() (*DasAddressHex, error) {
	if d == nil {
		return nil, fmt.Errorf("address is nil")
	}
	if d.DasAlgorithmId != common.DasAlgorithmIdAnyLock || d.ParsedAddress == nil {
		return nil, fmt.Errorf("address invalid")
	}
	argsLen := len(d.ParsedAddress.Script.Args)
	argsEnd := d.ParsedAddress.Script.Args[argsLen-1]
	log.Info("FormatAnyLock:", common.Bytes2Hex(d.ParsedAddress.Script.Args), argsLen, argsEnd)
	if argsLen != 22 || argsEnd != 0 {
		return nil, fmt.Errorf("address invalid")
	}
	switch d.ParsedAddress.Script.CodeHash.String() {
	case common.AnyLockCodeHashOfMainnetOmniLock, common.AnyLockCodeHashOfTestnetOmniLock:
		var res DasAddressHex
		args0 := d.ParsedAddress.Script.Args[0]
		log.Info("FormatAnyLock:", args0)
		switch args0 {
		case byte(18), byte(1):
			res.DasAlgorithmId = common.DasAlgorithmIdEth
			res.AddressPayload = d.ParsedAddress.Script.Args[1:21]
			res.AddressHex = common.Bytes2Hex(res.AddressPayload)
			res.ChainType = common.ChainTypeEth
		case byte(4):
			res.DasAlgorithmId = common.DasAlgorithmIdBitcoin
			res.DasSubAlgorithmId = common.DasSubAlgorithmIdBitcoinP2WPKH
			res.AddressPayload = d.ParsedAddress.Script.Args[1:21]
			res.AddressHex = hex.EncodeToString(res.AddressPayload)
			res.ChainType = common.ChainTypeBitcoin
		default:
			return nil, fmt.Errorf("unsupport")
		}
		return &res, nil
	}
	return nil, fmt.Errorf("unsupport")
}

type DasAddressFormat struct {
	DasNetType common.DasNetType
}

// only for .bit normal address
func (d *DasAddressFormat) NormalToHex(p DasAddressNormal) (r DasAddressHex, e error) {
	r.ChainType = p.ChainType
	switch p.ChainType {
	case common.ChainTypeCkb:
		if parseAddr, err := address.Parse(p.AddressNormal); err != nil {
			e = fmt.Errorf("address.Parse err: %s", err.Error())
		} else {
			if d.DasNetType == common.DasNetTypeMainNet && parseAddr.Mode != address.Mainnet {
				e = fmt.Errorf("not support testnet address[%s]", p.AddressNormal)
				return
			}
			if d.DasNetType != common.DasNetTypeMainNet && parseAddr.Mode == address.Mainnet {
				e = fmt.Errorf("not support mainnet address[%s]", p.AddressNormal)
				return
			}

			r.AddressHex = common.Bytes2Hex(parseAddr.Script.Args)
			r.AddressPayload = common.Hex2Bytes(r.AddressHex)

			var envNet Env
			switch d.DasNetType {
			case common.DasNetTypeMainNet:
				envNet = EnvMainNet
			case common.DasNetTypeTestnet2:
				envNet = EnvTestnet2
			case common.DasNetTypeTestnet3:
				envNet = EnvTestnet3
			default:
				e = fmt.Errorf("not support DasNetType[%d]", d.DasNetType)
				return
			}

			dasLockCodeHash := common.ScriptToTypeId(&types.Script{
				CodeHash: types.HexToHash(envNet.ContractCodeHash),
				HashType: types.HashTypeType,
				Args:     common.Hex2Bytes(envNet.MapContract[common.DasContractNameDispatchCellType]),
			}).Hex()

			switch parseAddr.Script.CodeHash.Hex() {
			case transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH:
				r.IsMulti = true
				r.DasAlgorithmId = common.DasAlgorithmIdCkbMulti
				r.ChainType = common.ChainTypeCkbMulti
			case dasLockCodeHash:
				switch common.DasAlgorithmId(parseAddr.Script.Args[0]) {
				case common.DasAlgorithmIdWebauthn:
					r.ChainType = common.ChainTypeWebauthn
					r.DasAlgorithmId = common.DasAlgorithmId(parseAddr.Script.Args[0])
					r.DasSubAlgorithmId = common.DasSubAlgorithmId(parseAddr.Script.Args[1])
					r.AddressHex = hex.EncodeToString(parseAddr.Script.Args[2:22])
					r.AddressPayload = parseAddr.Script.Args[2:22]
					return
				default:
					e = fmt.Errorf("not support DasAlgorithmId[%d]", parseAddr.Script.Args[0])
					return
				}
			case common.AnyLockCodeHashOfMainnetNoStrLock,
				common.AnyLockCodeHashOfTestnetNoStrLock:
				argsLen := len(parseAddr.Script.Args)
				args0 := parseAddr.Script.Args[0]
				if argsLen != 21 {
					e = fmt.Errorf("not support nostr-lock args[%s]", common.Bytes2Hex(parseAddr.Script.Args))
					return
				}
				if args0 != byte(0) {
					e = fmt.Errorf("not support nostr-lock args[%s]", common.Bytes2Hex(parseAddr.Script.Args))
					return
				}
				r.ChainType = common.ChainTypeAnyLock
				r.DasAlgorithmId = common.DasAlgorithmIdAnyLock
				r.AddressHex = p.AddressNormal
				r.ParsedAddress = parseAddr
			case common.AnyLockCodeHashOfMainnetOmniLock,
				common.AnyLockCodeHashOfTestnetOmniLock:
				argsLen := len(parseAddr.Script.Args)
				argsEnd := parseAddr.Script.Args[argsLen-1]
				args0 := parseAddr.Script.Args[0]
				if argsLen != 22 || argsEnd != 0 {
					e = fmt.Errorf("not support omni-lock args[%s]", common.Bytes2Hex(parseAddr.Script.Args))
					return
				}
				if args0 != byte(18) && args0 != byte(4) && args0 != byte(1) {
					e = fmt.Errorf("not support omni-lock args[%s]", common.Bytes2Hex(parseAddr.Script.Args))
					return
				}
				r.ChainType = common.ChainTypeAnyLock
				r.DasAlgorithmId = common.DasAlgorithmIdAnyLock
				r.AddressHex = p.AddressNormal
				r.ParsedAddress = parseAddr
			case common.AnyLockCodeHashOfMainnetJoyIDLock,
				common.AnyLockCodeHashOfTestnetJoyIDLock:
				argsLen := len(parseAddr.Script.Args)
				args0 := parseAddr.Script.Args[0]
				args1 := parseAddr.Script.Args[1]
				if argsLen != 22 {
					e = fmt.Errorf("not support joyid-lock args[%s]", common.Bytes2Hex(parseAddr.Script.Args))
					return
				}
				if args0 != byte(0) || args1 != byte(1) {
					e = fmt.Errorf("not support joyid-lock args[%s]", common.Bytes2Hex(parseAddr.Script.Args))
					return
				}
				r.ChainType = common.ChainTypeAnyLock
				r.DasAlgorithmId = common.DasAlgorithmIdAnyLock
				r.AddressHex = p.AddressNormal
				r.ParsedAddress = parseAddr
			case transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH:
				//r.ChainType = common.ChainTypeCkbSingle
				//r.DasAlgorithmId = common.DasAlgorithmIdCkbSingle
				r.ChainType = common.ChainTypeAnyLock
				r.DasAlgorithmId = common.DasAlgorithmIdAnyLock
				r.AddressHex = p.AddressNormal
				r.ParsedAddress = parseAddr
			default:
				e = fmt.Errorf("not support CodeHash, address invalid")
				return
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
			addr := p.AddressNormal[2:]
			if strings.ToLower(addr) != addr && strings.ToUpper(addr) != addr && gethcommon.HexToAddress(p.AddressNormal).Hex() != p.AddressNormal {
				e = fmt.Errorf("eth checkSum fail")
			} else {
				r.AddressHex = p.AddressNormal
				r.AddressPayload = common.Hex2Bytes(r.AddressHex)
			}
		} else {
			e = fmt.Errorf("regexp.MatchString fail")
		}
	case common.ChainTypeMixin:
		r.DasAlgorithmId = common.DasAlgorithmIdEd25519
		if ok, err := regexp.MatchString("^0x[0-9a-fA-F]{64}$", p.AddressNormal); err != nil {
			e = fmt.Errorf("regexp.MatchString err: %s", err.Error())
		} else if ok {
			r.AddressHex = p.AddressNormal
			r.AddressPayload = common.Hex2Bytes(r.AddressHex)
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
		} else {
			e = fmt.Errorf("invalid chain-type address: [%d][%s]", p.ChainType, p.AddressNormal)
		}
		if e == nil {
			r.AddressPayload = common.Hex2Bytes(strings.TrimPrefix(r.AddressHex, common.TronPreFix))
		}
	case common.ChainTypeDogeCoin:
		r.DasAlgorithmId = common.DasAlgorithmIdDogeChain
		addr, err := common.Base58CheckDecode(p.AddressNormal, common.DogeCoinBase58Version)
		if err != nil {
			e = fmt.Errorf("Base58CheckDecode err: %s", err.Error())
		} else {
			r.AddressHex = addr
			r.AddressPayload = common.Hex2Bytes(addr)
		}
	case common.ChainTypeBitcoin:
		r.DasAlgorithmId = common.DasAlgorithmIdBitcoin
		if strings.HasPrefix(p.AddressNormal, "bc1q") || strings.HasPrefix(p.AddressNormal, "tb1q") {
			if len(p.AddressNormal) != 42 {
				e = fmt.Errorf("invalid address [%s]", p.AddressNormal)
				return
			}
			r.DasSubAlgorithmId = common.DasSubAlgorithmIdBitcoinP2WPKH
		} else if strings.HasPrefix(p.AddressNormal, "1") || strings.HasPrefix(p.AddressNormal, "m") {
			r.DasSubAlgorithmId = common.DasSubAlgorithmIdBitcoinP2PKH
		} else if strings.HasPrefix(p.AddressNormal, "3") || strings.HasPrefix(p.AddressNormal, "2") {
			r.DasSubAlgorithmId = common.DasSubAlgorithmIdBitcoinP2SHP2WPKH
		} else if strings.HasPrefix(p.AddressNormal, "bc1p") || strings.HasPrefix(p.AddressNormal, "tb1p") {
			r.DasSubAlgorithmId = common.DasSubAlgorithmIdBitcoinP2TR
		} else {
			e = fmt.Errorf("invalid address [%s]", p.AddressNormal)
			return
		}
		netParams := bitcoin.GetBTCMainNetParams()
		if d.DasNetType != common.DasNetTypeMainNet {
			netParams = bitcoin.GetBTCTestNetParams()
		}
		addr, err := btcutil.DecodeAddress(p.AddressNormal, &netParams)
		if err != nil {
			e = fmt.Errorf("btcutil.DecodeAddress [%s] err: %s", p.AddressNormal, err.Error())
			return
		}
		r.AddressHex = hex.EncodeToString(addr.ScriptAddress())
		r.AddressPayload = addr.ScriptAddress()
	case common.ChainTypeWebauthn:
		if parseAddr, err := address.Parse(p.AddressNormal); err != nil {
			e = fmt.Errorf("address.Parse err: %s", err.Error())
		} else {
			r.DasAlgorithmId = common.DasAlgorithmId(parseAddr.Script.Args[0])
			r.DasSubAlgorithmId = common.DasSubAlgorithmId(parseAddr.Script.Args[1])
			r.AddressHex = hex.EncodeToString(parseAddr.Script.Args[2:22])
			r.AddressPayload = parseAddr.Script.Args[2:22]
		}
	default:
		e = fmt.Errorf("not support chain type [%d]", p.ChainType)
		return
	}
	return
}

// only for .bit hex address
func (d *DasAddressFormat) HexToNormal(p DasAddressHex) (r DasAddressNormal, e error) {
	switch p.DasAlgorithmId {
	case common.DasAlgorithmIdCkbMulti, common.DasAlgorithmIdCkbSingle:
		script := common.GetNormalLockScript(p.AddressHex)
		r.ChainType = common.ChainTypeCkbSingle
		if p.DasAlgorithmId == common.DasAlgorithmIdCkbMulti {
			r.ChainType = common.ChainTypeCkbMulti
			script = common.GetNormalLockScriptByMultiSig(p.AddressHex)
		}

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
	case common.DasAlgorithmIdDogeChain:
		r.ChainType = common.ChainTypeDogeCoin
		addr, err := common.Base58CheckEncode(p.AddressHex, common.DogeCoinBase58Version)
		if err != nil {
			e = fmt.Errorf("doge coin DecodeString err: %s", err.Error())
		} else {
			r.AddressNormal = addr
		}
	case common.DasAlgorithmIdBitcoin:
		r.ChainType = common.ChainTypeBitcoin
		netParams := bitcoin.GetBTCMainNetParams()
		if d.DasNetType != common.DasNetTypeMainNet {
			netParams = bitcoin.GetBTCTestNetParams()
		}
		switch p.DasSubAlgorithmId {
		case common.DasSubAlgorithmIdBitcoinP2PKH:
			r.AddressNormal = base58.CheckEncode(common.Hex2Bytes(p.AddressHex), netParams.PubKeyHashAddrID)
		case common.DasSubAlgorithmIdBitcoinP2WPKH:
			converted, err := bech32.ConvertBits(common.Hex2Bytes(p.AddressHex), 8, 5, true)
			if err != nil {
				e = fmt.Errorf("bech32.ConvertBits err: %s", err.Error())
				return
			}
			bech32Addr, err := bech32.Encode(netParams.Bech32HRPSegwit, append([]byte{0x00}, converted...))
			if err != nil {
				e = fmt.Errorf("bech32.Encode err: %s", err.Error())
				return
			}
			r.AddressNormal = bech32Addr
		case common.DasSubAlgorithmIdBitcoinP2SHP2WPKH:
			r.AddressNormal = base58.CheckEncode(common.Hex2Bytes(p.AddressHex), netParams.ScriptHashAddrID)
		case common.DasSubAlgorithmIdBitcoinP2TR:
			converted, err := bech32.ConvertBits(common.Hex2Bytes(p.AddressHex), 8, 5, true)
			if err != nil {
				e = fmt.Errorf("bech32.ConvertBits err: %s", err.Error())
				return
			}
			bech32Addr, err := bech32.EncodeM(netParams.Bech32HRPSegwit, append([]byte{0x01}, converted...))
			if err != nil {
				e = fmt.Errorf("bech32.Encode err: %s", err.Error())
				return
			}
			r.AddressNormal = bech32Addr
		default:
			e = fmt.Errorf("unknow sub algorith id [%d]", p.DasSubAlgorithmId)
		}
	case common.DasAlgorithmIdWebauthn:
		r.ChainType = common.ChainTypeWebauthn
		lock, _, err := d.HexToScript(p)
		if err != nil {
			e = err
			return
		}
		mode := address.Mainnet
		if d.DasNetType != common.DasNetTypeMainNet {
			mode = address.Testnet
		}
		if addr, err := common.ConvertScriptToAddress(mode, lock); err != nil {
			e = fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
		} else {
			r.AddressNormal = addr
		}
	default:
		e = fmt.Errorf("not support DasAlgorithmId [%d]", p.DasAlgorithmId)
	}

	return
}

// only for .bit hex address
func (d *DasAddressFormat) HexToScript(p DasAddressHex) (lockScript, typeScript *types.Script, e error) {
	if p.DasAlgorithmId == common.DasAlgorithmIdEth712 {
		contractBalance, err := GetDasContractInfo(common.DasContractNameBalanceCellType)
		if err != nil {
			e = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
			return
		}
		typeScript = contractBalance.ToScript(nil)
	}

	args, err := d.HexToArgs(p, p)
	if err != nil {
		e = fmt.Errorf("HexToArgs err: %s", err.Error())
		return
	}

	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		e = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}
	lockScript = contractDispatch.ToScript(args)
	return
}

// only for .bit hex address
func (d *DasAddressFormat) HexToArgs(owner, manager DasAddressHex) (args []byte, e error) {
	ownerArgs, err := d.HexToHalfArgs(owner)
	if err != nil {
		e = fmt.Errorf("HexToHalfArgs err: %s", err.Error())
		return
	}
	managerArgs, err := d.HexToHalfArgs(manager)
	if err != nil {
		e = fmt.Errorf("HexToHalfArgs err: %s", err.Error())
		return
	}
	args = append(ownerArgs, managerArgs...)
	return
}

// only for .bit hex address
func (d *DasAddressFormat) HexToHalfArgs(p DasAddressHex) (args []byte, e error) {
	argsStr := ""
	switch p.DasAlgorithmId {
	case common.DasAlgorithmIdCkbMulti:
		argsStr = common.DasLockCkbMultiPreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdCkbSingle:
		argsStr = common.DasLockCkbSinglePreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdEth:
		argsStr = common.DasLockEthPreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdTron:
		argsStr = common.DasLockTronPreFix + strings.TrimPrefix(p.AddressHex, common.TronPreFix)
	case common.DasAlgorithmIdEth712:
		argsStr = common.DasLockEth712PreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdEd25519:
		argsStr = common.DasLockEd25519PreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdCkb:
		argsStr = common.DasLockCkbPreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	case common.DasAlgorithmIdDogeChain:
		argsStr = common.DasLockDogePreFix + p.AddressHex
	case common.DasAlgorithmIdBitcoin:
		switch p.DasSubAlgorithmId {
		case common.DasSubAlgorithmIdBitcoinP2PKH:
			argsStr = common.DasLockBitcoinPreFix + common.DasLockBitcoinSubPreFixP2PKH + p.AddressHex
		case common.DasSubAlgorithmIdBitcoinP2WPKH:
			argsStr = common.DasLockBitcoinPreFix + common.DasLockBitcoinSubPreFixP2WPKH + p.AddressHex
		default:
			e = fmt.Errorf("unknow sub algorithm id[%d]", p.DasSubAlgorithmId)
		}
	case common.DasAlgorithmIdWebauthn:
		// TODO Temporarily written as a fixed sub-algorithm id
		argsStr = common.DasLockWebauthnPreFix + common.DasLockWebauthnSubPreFix + strings.TrimPrefix(p.AddressHex, common.HexPreFix)
	default:
		e = fmt.Errorf("not support DasAlgorithmId[%d]", p.DasAlgorithmId)
	}
	if argsStr != "" {
		args = common.Hex2Bytes(argsStr)
	}
	return
}

// only for .bit args
func (d *DasAddressFormat) ArgsToNormal(args []byte) (ownerNormal, managerNormal DasAddressNormal, e error) {
	log.Info("ArgsToNormal:", common.Bytes2Hex(args))
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

// only for .bit args
func (d *DasAddressFormat) ArgsToHex(args []byte) (ownerHex, managerHex DasAddressHex, e error) {
	//log.Info("ArgsToHex:", common.Bytes2Hex(args))
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
	case common.DasAlgorithmIdCkbMulti, common.DasAlgorithmIdCkbSingle, common.DasAlgorithmIdTron,
		common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712, common.DasAlgorithmIdDogeChain:
		splitLen = common.DasLockArgsLen / 2
	case common.DasAlgorithmIdEd25519:
		splitLen = common.DasLockArgsLenMax / 2
	case common.DasAlgorithmIdCkb:
		splitLen = common.DasLockArgsLen / 2
	case common.DasAlgorithmIdBitcoin:
		splitLen = common.DasLockArgsLenBitcoin / 2
	case common.DasAlgorithmIdWebauthn:
		splitLen = common.DasLockArgsLenWebAuthn / 2
		if d.DasNetType != common.DasNetTypeMainNet && len(args) == 48 {
			splitLen = 24
			owner = args[2:splitLen]
			manager = args[splitLen+2:]
			return
		}
	default:
		e = fmt.Errorf("unknow DasAlgorithmId[%d]", oID)
		return
	}
	owner = args[:splitLen]
	manager = args[splitLen:]
	return
}
func (d *DasAddressFormat) halfArgsToHex(args []byte) (r DasAddressHex, e error) {
	r.DasAlgorithmId = common.DasAlgorithmId(args[0])
	switch r.DasAlgorithmId {
	case common.DasAlgorithmIdCkbMulti:
		r.ChainType = common.ChainTypeCkbMulti
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
		r.AddressPayload = args[1:]
		r.IsMulti = true
	case common.DasAlgorithmIdCkbSingle:
		r.ChainType = common.ChainTypeCkbSingle
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
		r.AddressPayload = args[1:]
	case common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712:
		r.ChainType = common.ChainTypeEth
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
		r.AddressPayload = args[1:]
	case common.DasAlgorithmIdTron:
		r.ChainType = common.ChainTypeTron
		r.AddressHex = common.TronPreFix + hex.EncodeToString(args[1:])
		r.AddressPayload = args[1:]
	case common.DasAlgorithmIdEd25519:
		r.ChainType = common.ChainTypeMixin
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
		r.AddressPayload = args[1:]
	case common.DasAlgorithmIdCkb:
		r.ChainType = common.ChainTypeCkb
		r.AddressHex = common.HexPreFix + hex.EncodeToString(args[1:])
		r.AddressPayload = args[1:]
	case common.DasAlgorithmIdDogeChain:
		r.ChainType = common.ChainTypeDogeCoin
		r.AddressHex = hex.EncodeToString(args[1:])
		r.AddressPayload = args[1:]
	case common.DasAlgorithmIdBitcoin:
		r.ChainType = common.ChainTypeBitcoin
		r.AddressHex = hex.EncodeToString(args[2:])
		r.AddressPayload = args[2:]
		r.DasSubAlgorithmId = common.DasSubAlgorithmId(args[1])
	case common.DasAlgorithmIdWebauthn:
		r.ChainType = common.ChainTypeWebauthn
		r.AddressHex = hex.EncodeToString(args[2:])
		r.AddressPayload = args[2:]
		r.DasSubAlgorithmId = common.DasSubAlgorithmId(args[1])
	default:
		e = fmt.Errorf("not support DasAlgorithmId [%d]", r.DasAlgorithmId)
	}
	return
}

// for .bit or normal ckb script
func (d *DasAddressFormat) ScriptToHex(s *types.Script) (ownerHex, managerHex DasAddressHex, e error) {
	if s == nil {
		e = fmt.Errorf("script is nil")
		return
	}
	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		e = fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		return
	}
	if contractDispatch.IsSameTypeId(s.CodeHash) {
		return d.ArgsToHex(s.Args)
	} else {
		ownerHex.ChainType = common.ChainTypeCkb
		ownerHex.DasAlgorithmId = common.DasAlgorithmIdCkb
		ownerHex.AddressPayload = s.Args
		ownerHex.AddressHex = common.Bytes2Hex(s.Args)

		mode := address.Mainnet
		if d.DasNetType != common.DasNetTypeMainNet {
			mode = address.Testnet
		}
		addr, err := address.ConvertScriptToAddress(mode, s)
		if err == nil {
			ownerHex.ChainType = common.ChainTypeAnyLock
			ownerHex.DasAlgorithmId = common.DasAlgorithmIdAnyLock
			ownerHex.AddressHex = addr
		}
	}
	return
}
