// Generated by Molecule 0.7.5
// Generated by Moleculec-Go 0.1.11

package molecule

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type AccountCellDataV1Builder struct {
	id            AccountId
	account       AccountChars
	registered_at Uint64
	updated_at    Uint64
	status        Uint8
	records       Records
}

func (s *AccountCellDataV1Builder) Build() AccountCellDataV1 {
	b := new(bytes.Buffer)

	totalSize := HeaderSizeUint * (6 + 1)
	offsets := make([]uint32, 0, 6)

	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.id.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.account.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.registered_at.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.updated_at.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.status.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.records.AsSlice()))

	b.Write(packNumber(Number(totalSize)))

	for i := 0; i < len(offsets); i++ {
		b.Write(packNumber(Number(offsets[i])))
	}

	b.Write(s.id.AsSlice())
	b.Write(s.account.AsSlice())
	b.Write(s.registered_at.AsSlice())
	b.Write(s.updated_at.AsSlice())
	b.Write(s.status.AsSlice())
	b.Write(s.records.AsSlice())
	return AccountCellDataV1{inner: b.Bytes()}
}

func (s *AccountCellDataV1Builder) Id(v AccountId) *AccountCellDataV1Builder {
	s.id = v
	return s
}

func (s *AccountCellDataV1Builder) Account(v AccountChars) *AccountCellDataV1Builder {
	s.account = v
	return s
}

func (s *AccountCellDataV1Builder) RegisteredAt(v Uint64) *AccountCellDataV1Builder {
	s.registered_at = v
	return s
}

func (s *AccountCellDataV1Builder) UpdatedAt(v Uint64) *AccountCellDataV1Builder {
	s.updated_at = v
	return s
}

func (s *AccountCellDataV1Builder) Status(v Uint8) *AccountCellDataV1Builder {
	s.status = v
	return s
}

func (s *AccountCellDataV1Builder) Records(v Records) *AccountCellDataV1Builder {
	s.records = v
	return s
}

func NewAccountCellDataV1Builder() *AccountCellDataV1Builder {
	return &AccountCellDataV1Builder{id: AccountIdDefault(), account: AccountCharsDefault(), registered_at: Uint64Default(), updated_at: Uint64Default(), status: Uint8Default(), records: RecordsDefault()}
}

type AccountCellDataV1 struct {
	inner []byte
}

func AccountCellDataV1FromSliceUnchecked(slice []byte) *AccountCellDataV1 {
	return &AccountCellDataV1{inner: slice}
}
func (s *AccountCellDataV1) AsSlice() []byte {
	return s.inner
}

func AccountCellDataV1Default() AccountCellDataV1 {
	return *AccountCellDataV1FromSliceUnchecked([]byte{73, 0, 0, 0, 28, 0, 0, 0, 48, 0, 0, 0, 52, 0, 0, 0, 60, 0, 0, 0, 68, 0, 0, 0, 69, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0})
}

func AccountCellDataV1FromSlice(slice []byte, compatible bool) (*AccountCellDataV1, error) {
	sliceLen := len(slice)
	if uint32(sliceLen) < HeaderSizeUint {
		errMsg := strings.Join([]string{"HeaderIsBroken", "AccountCellDataV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}

	totalSize := unpackNumber(slice)
	if Number(sliceLen) != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "AccountCellDataV1", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}

	if uint32(sliceLen) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "AccountCellDataV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	offsetFirst := unpackNumber(slice[HeaderSizeUint:])
	if uint32(offsetFirst)%HeaderSizeUint != 0 || uint32(offsetFirst) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"OffsetsNotMatch", "AccountCellDataV1", strconv.Itoa(int(offsetFirst % 4)), "!= 0", strconv.Itoa(int(offsetFirst)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	if sliceLen < int(offsetFirst) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "AccountCellDataV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(offsetFirst))}, " ")
		return nil, errors.New(errMsg)
	}

	fieldCount := uint32(offsetFirst)/HeaderSizeUint - 1
	if fieldCount < 6 {
		return nil, errors.New("FieldCountNotMatch")
	} else if !compatible && fieldCount > 6 {
		return nil, errors.New("FieldCountNotMatch")
	}

	offsets := make([]uint32, fieldCount)

	for i := 0; i < int(fieldCount); i++ {
		offsets[i] = uint32(unpackNumber(slice[HeaderSizeUint:][int(HeaderSizeUint)*i:]))
	}
	offsets = append(offsets, uint32(totalSize))

	for i := 0; i < len(offsets); i++ {
		if i&1 != 0 && offsets[i-1] > offsets[i] {
			return nil, errors.New("OffsetsNotMatch")
		}
	}

	var err error

	_, err = AccountIdFromSlice(slice[offsets[0]:offsets[1]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = AccountCharsFromSlice(slice[offsets[1]:offsets[2]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint64FromSlice(slice[offsets[2]:offsets[3]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint64FromSlice(slice[offsets[3]:offsets[4]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint8FromSlice(slice[offsets[4]:offsets[5]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = RecordsFromSlice(slice[offsets[5]:offsets[6]], compatible)
	if err != nil {
		return nil, err
	}

	return &AccountCellDataV1{inner: slice}, nil
}

func (s *AccountCellDataV1) TotalSize() uint {
	return uint(unpackNumber(s.inner))
}
func (s *AccountCellDataV1) FieldCount() uint {
	var number uint = 0
	if uint32(s.TotalSize()) == HeaderSizeUint {
		return number
	}
	number = uint(unpackNumber(s.inner[HeaderSizeUint:]))/4 - 1
	return number
}
func (s *AccountCellDataV1) Len() uint {
	return s.FieldCount()
}
func (s *AccountCellDataV1) IsEmpty() bool {
	return s.Len() == 0
}
func (s *AccountCellDataV1) CountExtraFields() uint {
	return s.FieldCount() - 6
}

func (s *AccountCellDataV1) HasExtraFields() bool {
	return 6 != s.FieldCount()
}

func (s *AccountCellDataV1) Id() *AccountId {
	start := unpackNumber(s.inner[4:])
	end := unpackNumber(s.inner[8:])
	return AccountIdFromSliceUnchecked(s.inner[start:end])
}

func (s *AccountCellDataV1) Account() *AccountChars {
	start := unpackNumber(s.inner[8:])
	end := unpackNumber(s.inner[12:])
	return AccountCharsFromSliceUnchecked(s.inner[start:end])
}

func (s *AccountCellDataV1) RegisteredAt() *Uint64 {
	start := unpackNumber(s.inner[12:])
	end := unpackNumber(s.inner[16:])
	return Uint64FromSliceUnchecked(s.inner[start:end])
}

func (s *AccountCellDataV1) UpdatedAt() *Uint64 {
	start := unpackNumber(s.inner[16:])
	end := unpackNumber(s.inner[20:])
	return Uint64FromSliceUnchecked(s.inner[start:end])
}

func (s *AccountCellDataV1) Status() *Uint8 {
	start := unpackNumber(s.inner[20:])
	end := unpackNumber(s.inner[24:])
	return Uint8FromSliceUnchecked(s.inner[start:end])
}

func (s *AccountCellDataV1) Records() *Records {
	var ret *Records
	start := unpackNumber(s.inner[24:])
	if s.HasExtraFields() {
		end := unpackNumber(s.inner[28:])
		ret = RecordsFromSliceUnchecked(s.inner[start:end])
	} else {
		ret = RecordsFromSliceUnchecked(s.inner[start:])
	}
	return ret
}

func (s *AccountCellDataV1) AsBuilder() AccountCellDataV1Builder {
	ret := NewAccountCellDataV1Builder().Id(*s.Id()).Account(*s.Account()).RegisteredAt(*s.RegisteredAt()).UpdatedAt(*s.UpdatedAt()).Status(*s.Status()).Records(*s.Records())
	return *ret
}

type PreAccountCellDataV1Builder struct {
	account          AccountChars
	refund_lock      Script
	owner_lock_args  Bytes
	inviter_id       Bytes
	inviter_lock     ScriptOpt
	channel_lock     ScriptOpt
	price            PriceConfig
	quote            Uint64
	invited_discount Uint32
	created_at       Uint64
}

func (s *PreAccountCellDataV1Builder) Build() PreAccountCellDataV1 {
	b := new(bytes.Buffer)

	totalSize := HeaderSizeUint * (10 + 1)
	offsets := make([]uint32, 0, 10)

	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.account.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.refund_lock.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.owner_lock_args.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.inviter_id.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.inviter_lock.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.channel_lock.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.price.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.quote.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.invited_discount.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.created_at.AsSlice()))

	b.Write(packNumber(Number(totalSize)))

	for i := 0; i < len(offsets); i++ {
		b.Write(packNumber(Number(offsets[i])))
	}

	b.Write(s.account.AsSlice())
	b.Write(s.refund_lock.AsSlice())
	b.Write(s.owner_lock_args.AsSlice())
	b.Write(s.inviter_id.AsSlice())
	b.Write(s.inviter_lock.AsSlice())
	b.Write(s.channel_lock.AsSlice())
	b.Write(s.price.AsSlice())
	b.Write(s.quote.AsSlice())
	b.Write(s.invited_discount.AsSlice())
	b.Write(s.created_at.AsSlice())
	return PreAccountCellDataV1{inner: b.Bytes()}
}

func (s *PreAccountCellDataV1Builder) Account(v AccountChars) *PreAccountCellDataV1Builder {
	s.account = v
	return s
}

func (s *PreAccountCellDataV1Builder) RefundLock(v Script) *PreAccountCellDataV1Builder {
	s.refund_lock = v
	return s
}

func (s *PreAccountCellDataV1Builder) OwnerLockArgs(v Bytes) *PreAccountCellDataV1Builder {
	s.owner_lock_args = v
	return s
}

func (s *PreAccountCellDataV1Builder) InviterId(v Bytes) *PreAccountCellDataV1Builder {
	s.inviter_id = v
	return s
}

func (s *PreAccountCellDataV1Builder) InviterLock(v ScriptOpt) *PreAccountCellDataV1Builder {
	s.inviter_lock = v
	return s
}

func (s *PreAccountCellDataV1Builder) ChannelLock(v ScriptOpt) *PreAccountCellDataV1Builder {
	s.channel_lock = v
	return s
}

func (s *PreAccountCellDataV1Builder) Price(v PriceConfig) *PreAccountCellDataV1Builder {
	s.price = v
	return s
}

func (s *PreAccountCellDataV1Builder) Quote(v Uint64) *PreAccountCellDataV1Builder {
	s.quote = v
	return s
}

func (s *PreAccountCellDataV1Builder) InvitedDiscount(v Uint32) *PreAccountCellDataV1Builder {
	s.invited_discount = v
	return s
}

func (s *PreAccountCellDataV1Builder) CreatedAt(v Uint64) *PreAccountCellDataV1Builder {
	s.created_at = v
	return s
}

func NewPreAccountCellDataV1Builder() *PreAccountCellDataV1Builder {
	return &PreAccountCellDataV1Builder{account: AccountCharsDefault(), refund_lock: ScriptDefault(), owner_lock_args: BytesDefault(), inviter_id: BytesDefault(), inviter_lock: ScriptOptDefault(), channel_lock: ScriptOptDefault(), price: PriceConfigDefault(), quote: Uint64Default(), invited_discount: Uint32Default(), created_at: Uint64Default()}
}

type PreAccountCellDataV1 struct {
	inner []byte
}

func PreAccountCellDataV1FromSliceUnchecked(slice []byte) *PreAccountCellDataV1 {
	return &PreAccountCellDataV1{inner: slice}
}
func (s *PreAccountCellDataV1) AsSlice() []byte {
	return s.inner
}

func PreAccountCellDataV1Default() PreAccountCellDataV1 {
	return *PreAccountCellDataV1FromSliceUnchecked([]byte{162, 0, 0, 0, 44, 0, 0, 0, 48, 0, 0, 0, 101, 0, 0, 0, 105, 0, 0, 0, 109, 0, 0, 0, 109, 0, 0, 0, 109, 0, 0, 0, 142, 0, 0, 0, 150, 0, 0, 0, 154, 0, 0, 0, 4, 0, 0, 0, 53, 0, 0, 0, 16, 0, 0, 0, 48, 0, 0, 0, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 33, 0, 0, 0, 16, 0, 0, 0, 17, 0, 0, 0, 25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func PreAccountCellDataV1FromSlice(slice []byte, compatible bool) (*PreAccountCellDataV1, error) {
	sliceLen := len(slice)
	if uint32(sliceLen) < HeaderSizeUint {
		errMsg := strings.Join([]string{"HeaderIsBroken", "PreAccountCellDataV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}

	totalSize := unpackNumber(slice)
	if Number(sliceLen) != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "PreAccountCellDataV1", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}

	if uint32(sliceLen) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "PreAccountCellDataV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	offsetFirst := unpackNumber(slice[HeaderSizeUint:])
	if uint32(offsetFirst)%HeaderSizeUint != 0 || uint32(offsetFirst) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"OffsetsNotMatch", "PreAccountCellDataV1", strconv.Itoa(int(offsetFirst % 4)), "!= 0", strconv.Itoa(int(offsetFirst)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	if sliceLen < int(offsetFirst) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "PreAccountCellDataV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(offsetFirst))}, " ")
		return nil, errors.New(errMsg)
	}

	fieldCount := uint32(offsetFirst)/HeaderSizeUint - 1
	if fieldCount < 10 {
		return nil, errors.New("FieldCountNotMatch")
	} else if !compatible && fieldCount > 10 {
		return nil, errors.New("FieldCountNotMatch")
	}

	offsets := make([]uint32, fieldCount)

	for i := 0; i < int(fieldCount); i++ {
		offsets[i] = uint32(unpackNumber(slice[HeaderSizeUint:][int(HeaderSizeUint)*i:]))
	}
	offsets = append(offsets, uint32(totalSize))

	for i := 0; i < len(offsets); i++ {
		if i&1 != 0 && offsets[i-1] > offsets[i] {
			return nil, errors.New("OffsetsNotMatch")
		}
	}

	var err error

	_, err = AccountCharsFromSlice(slice[offsets[0]:offsets[1]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = ScriptFromSlice(slice[offsets[1]:offsets[2]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[2]:offsets[3]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[3]:offsets[4]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = ScriptOptFromSlice(slice[offsets[4]:offsets[5]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = ScriptOptFromSlice(slice[offsets[5]:offsets[6]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = PriceConfigFromSlice(slice[offsets[6]:offsets[7]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint64FromSlice(slice[offsets[7]:offsets[8]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint32FromSlice(slice[offsets[8]:offsets[9]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint64FromSlice(slice[offsets[9]:offsets[10]], compatible)
	if err != nil {
		return nil, err
	}

	return &PreAccountCellDataV1{inner: slice}, nil
}

func (s *PreAccountCellDataV1) TotalSize() uint {
	return uint(unpackNumber(s.inner))
}
func (s *PreAccountCellDataV1) FieldCount() uint {
	var number uint = 0
	if uint32(s.TotalSize()) == HeaderSizeUint {
		return number
	}
	number = uint(unpackNumber(s.inner[HeaderSizeUint:]))/4 - 1
	return number
}
func (s *PreAccountCellDataV1) Len() uint {
	return s.FieldCount()
}
func (s *PreAccountCellDataV1) IsEmpty() bool {
	return s.Len() == 0
}
func (s *PreAccountCellDataV1) CountExtraFields() uint {
	return s.FieldCount() - 10
}

func (s *PreAccountCellDataV1) HasExtraFields() bool {
	return 10 != s.FieldCount()
}

func (s *PreAccountCellDataV1) Account() *AccountChars {
	start := unpackNumber(s.inner[4:])
	end := unpackNumber(s.inner[8:])
	return AccountCharsFromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) RefundLock() *Script {
	start := unpackNumber(s.inner[8:])
	end := unpackNumber(s.inner[12:])
	return ScriptFromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) OwnerLockArgs() *Bytes {
	start := unpackNumber(s.inner[12:])
	end := unpackNumber(s.inner[16:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) InviterId() *Bytes {
	start := unpackNumber(s.inner[16:])
	end := unpackNumber(s.inner[20:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) InviterLock() *ScriptOpt {
	start := unpackNumber(s.inner[20:])
	end := unpackNumber(s.inner[24:])
	return ScriptOptFromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) ChannelLock() *ScriptOpt {
	start := unpackNumber(s.inner[24:])
	end := unpackNumber(s.inner[28:])
	return ScriptOptFromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) Price() *PriceConfig {
	start := unpackNumber(s.inner[28:])
	end := unpackNumber(s.inner[32:])
	return PriceConfigFromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) Quote() *Uint64 {
	start := unpackNumber(s.inner[32:])
	end := unpackNumber(s.inner[36:])
	return Uint64FromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) InvitedDiscount() *Uint32 {
	start := unpackNumber(s.inner[36:])
	end := unpackNumber(s.inner[40:])
	return Uint32FromSliceUnchecked(s.inner[start:end])
}

func (s *PreAccountCellDataV1) CreatedAt() *Uint64 {
	var ret *Uint64
	start := unpackNumber(s.inner[40:])
	if s.HasExtraFields() {
		end := unpackNumber(s.inner[44:])
		ret = Uint64FromSliceUnchecked(s.inner[start:end])
	} else {
		ret = Uint64FromSliceUnchecked(s.inner[start:])
	}
	return ret
}

func (s *PreAccountCellDataV1) AsBuilder() PreAccountCellDataV1Builder {
	ret := NewPreAccountCellDataV1Builder().Account(*s.Account()).RefundLock(*s.RefundLock()).OwnerLockArgs(*s.OwnerLockArgs()).InviterId(*s.InviterId()).InviterLock(*s.InviterLock()).ChannelLock(*s.ChannelLock()).Price(*s.Price()).Quote(*s.Quote()).InvitedDiscount(*s.InvitedDiscount()).CreatedAt(*s.CreatedAt())
	return *ret
}

type SubAccountV1Builder struct {
	lock                    Script
	id                      AccountId
	account                 AccountChars
	suffix                  Bytes
	registered_at           Uint64
	expired_at              Uint64
	status                  Uint8
	records                 Records
	nonce                   Uint64
	enable_sub_account      Uint8
	renew_sub_account_price Uint64
}

func (s *SubAccountV1Builder) Build() SubAccountV1 {
	b := new(bytes.Buffer)

	totalSize := HeaderSizeUint * (11 + 1)
	offsets := make([]uint32, 0, 11)

	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.lock.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.id.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.account.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.suffix.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.registered_at.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.expired_at.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.status.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.records.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.nonce.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.enable_sub_account.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.renew_sub_account_price.AsSlice()))

	b.Write(packNumber(Number(totalSize)))

	for i := 0; i < len(offsets); i++ {
		b.Write(packNumber(Number(offsets[i])))
	}

	b.Write(s.lock.AsSlice())
	b.Write(s.id.AsSlice())
	b.Write(s.account.AsSlice())
	b.Write(s.suffix.AsSlice())
	b.Write(s.registered_at.AsSlice())
	b.Write(s.expired_at.AsSlice())
	b.Write(s.status.AsSlice())
	b.Write(s.records.AsSlice())
	b.Write(s.nonce.AsSlice())
	b.Write(s.enable_sub_account.AsSlice())
	b.Write(s.renew_sub_account_price.AsSlice())
	return SubAccountV1{inner: b.Bytes()}
}

func (s *SubAccountV1Builder) Lock(v Script) *SubAccountV1Builder {
	s.lock = v
	return s
}

func (s *SubAccountV1Builder) Id(v AccountId) *SubAccountV1Builder {
	s.id = v
	return s
}

func (s *SubAccountV1Builder) Account(v AccountChars) *SubAccountV1Builder {
	s.account = v
	return s
}

func (s *SubAccountV1Builder) Suffix(v Bytes) *SubAccountV1Builder {
	s.suffix = v
	return s
}

func (s *SubAccountV1Builder) RegisteredAt(v Uint64) *SubAccountV1Builder {
	s.registered_at = v
	return s
}

func (s *SubAccountV1Builder) ExpiredAt(v Uint64) *SubAccountV1Builder {
	s.expired_at = v
	return s
}

func (s *SubAccountV1Builder) Status(v Uint8) *SubAccountV1Builder {
	s.status = v
	return s
}

func (s *SubAccountV1Builder) Records(v Records) *SubAccountV1Builder {
	s.records = v
	return s
}

func (s *SubAccountV1Builder) Nonce(v Uint64) *SubAccountV1Builder {
	s.nonce = v
	return s
}

func (s *SubAccountV1Builder) EnableSubAccount(v Uint8) *SubAccountV1Builder {
	s.enable_sub_account = v
	return s
}

func (s *SubAccountV1Builder) RenewSubAccountPrice(v Uint64) *SubAccountV1Builder {
	s.renew_sub_account_price = v
	return s
}

func NewSubAccountV1Builder() *SubAccountV1Builder {
	return &SubAccountV1Builder{lock: ScriptDefault(), id: AccountIdDefault(), account: AccountCharsDefault(), suffix: BytesDefault(), registered_at: Uint64Default(), expired_at: Uint64Default(), status: Uint8Default(), records: RecordsDefault(), nonce: Uint64Default(), enable_sub_account: Uint8Default(), renew_sub_account_price: Uint64Default()}
}

type SubAccountV1 struct {
	inner []byte
}

func SubAccountV1FromSliceUnchecked(slice []byte) *SubAccountV1 {
	return &SubAccountV1{inner: slice}
}
func (s *SubAccountV1) AsSlice() []byte {
	return s.inner
}

func SubAccountV1Default() SubAccountV1 {
	return *SubAccountV1FromSliceUnchecked([]byte{167, 0, 0, 0, 48, 0, 0, 0, 101, 0, 0, 0, 121, 0, 0, 0, 125, 0, 0, 0, 129, 0, 0, 0, 137, 0, 0, 0, 145, 0, 0, 0, 146, 0, 0, 0, 150, 0, 0, 0, 158, 0, 0, 0, 159, 0, 0, 0, 53, 0, 0, 0, 16, 0, 0, 0, 48, 0, 0, 0, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func SubAccountV1FromSlice(slice []byte, compatible bool) (*SubAccountV1, error) {
	sliceLen := len(slice)
	if uint32(sliceLen) < HeaderSizeUint {
		errMsg := strings.Join([]string{"HeaderIsBroken", "SubAccountV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}

	totalSize := unpackNumber(slice)
	if Number(sliceLen) != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "SubAccountV1", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}

	if uint32(sliceLen) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "SubAccountV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	offsetFirst := unpackNumber(slice[HeaderSizeUint:])
	if uint32(offsetFirst)%HeaderSizeUint != 0 || uint32(offsetFirst) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"OffsetsNotMatch", "SubAccountV1", strconv.Itoa(int(offsetFirst % 4)), "!= 0", strconv.Itoa(int(offsetFirst)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	if sliceLen < int(offsetFirst) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "SubAccountV1", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(offsetFirst))}, " ")
		return nil, errors.New(errMsg)
	}

	fieldCount := uint32(offsetFirst)/HeaderSizeUint - 1
	if fieldCount < 11 {
		return nil, errors.New("FieldCountNotMatch")
	} else if !compatible && fieldCount > 11 {
		return nil, errors.New("FieldCountNotMatch")
	}

	offsets := make([]uint32, fieldCount)

	for i := 0; i < int(fieldCount); i++ {
		offsets[i] = uint32(unpackNumber(slice[HeaderSizeUint:][int(HeaderSizeUint)*i:]))
	}
	offsets = append(offsets, uint32(totalSize))

	for i := 0; i < len(offsets); i++ {
		if i&1 != 0 && offsets[i-1] > offsets[i] {
			return nil, errors.New("OffsetsNotMatch")
		}
	}

	var err error

	_, err = ScriptFromSlice(slice[offsets[0]:offsets[1]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = AccountIdFromSlice(slice[offsets[1]:offsets[2]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = AccountCharsFromSlice(slice[offsets[2]:offsets[3]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[3]:offsets[4]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint64FromSlice(slice[offsets[4]:offsets[5]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint64FromSlice(slice[offsets[5]:offsets[6]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint8FromSlice(slice[offsets[6]:offsets[7]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = RecordsFromSlice(slice[offsets[7]:offsets[8]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint64FromSlice(slice[offsets[8]:offsets[9]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint8FromSlice(slice[offsets[9]:offsets[10]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint64FromSlice(slice[offsets[10]:offsets[11]], compatible)
	if err != nil {
		return nil, err
	}

	return &SubAccountV1{inner: slice}, nil
}

func (s *SubAccountV1) TotalSize() uint {
	return uint(unpackNumber(s.inner))
}
func (s *SubAccountV1) FieldCount() uint {
	var number uint = 0
	if uint32(s.TotalSize()) == HeaderSizeUint {
		return number
	}
	number = uint(unpackNumber(s.inner[HeaderSizeUint:]))/4 - 1
	return number
}
func (s *SubAccountV1) Len() uint {
	return s.FieldCount()
}
func (s *SubAccountV1) IsEmpty() bool {
	return s.Len() == 0
}
func (s *SubAccountV1) CountExtraFields() uint {
	return s.FieldCount() - 11
}

func (s *SubAccountV1) HasExtraFields() bool {
	return 11 != s.FieldCount()
}

func (s *SubAccountV1) Lock() *Script {
	start := unpackNumber(s.inner[4:])
	end := unpackNumber(s.inner[8:])
	return ScriptFromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) Id() *AccountId {
	start := unpackNumber(s.inner[8:])
	end := unpackNumber(s.inner[12:])
	return AccountIdFromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) Account() *AccountChars {
	start := unpackNumber(s.inner[12:])
	end := unpackNumber(s.inner[16:])
	return AccountCharsFromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) Suffix() *Bytes {
	start := unpackNumber(s.inner[16:])
	end := unpackNumber(s.inner[20:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) RegisteredAt() *Uint64 {
	start := unpackNumber(s.inner[20:])
	end := unpackNumber(s.inner[24:])
	return Uint64FromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) ExpiredAt() *Uint64 {
	start := unpackNumber(s.inner[24:])
	end := unpackNumber(s.inner[28:])
	return Uint64FromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) Status() *Uint8 {
	start := unpackNumber(s.inner[28:])
	end := unpackNumber(s.inner[32:])
	return Uint8FromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) Records() *Records {
	start := unpackNumber(s.inner[32:])
	end := unpackNumber(s.inner[36:])
	return RecordsFromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) Nonce() *Uint64 {
	start := unpackNumber(s.inner[36:])
	end := unpackNumber(s.inner[40:])
	return Uint64FromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) EnableSubAccount() *Uint8 {
	start := unpackNumber(s.inner[40:])
	end := unpackNumber(s.inner[44:])
	return Uint8FromSliceUnchecked(s.inner[start:end])
}

func (s *SubAccountV1) RenewSubAccountPrice() *Uint64 {
	var ret *Uint64
	start := unpackNumber(s.inner[44:])
	if s.HasExtraFields() {
		end := unpackNumber(s.inner[48:])
		ret = Uint64FromSliceUnchecked(s.inner[start:end])
	} else {
		ret = Uint64FromSliceUnchecked(s.inner[start:])
	}
	return ret
}

func (s *SubAccountV1) AsBuilder() SubAccountV1Builder {
	ret := NewSubAccountV1Builder().Lock(*s.Lock()).Id(*s.Id()).Account(*s.Account()).Suffix(*s.Suffix()).RegisteredAt(*s.RegisteredAt()).ExpiredAt(*s.ExpiredAt()).Status(*s.Status()).Records(*s.Records()).Nonce(*s.Nonce()).EnableSubAccount(*s.EnableSubAccount()).RenewSubAccountPrice(*s.RenewSubAccountPrice())
	return *ret
}
