package witness

type ReverseSmtRecordAction string
type ReverseSmtRecordVersion uint32

const (
	ReverseSmtRecordVersion1 ReverseSmtRecordVersion = 1

	ReverseSmtRecordActionUpdate ReverseSmtRecordAction = "update"
	ReverseSmtRecordActionRemove ReverseSmtRecordAction = "remove"
)

type ReverseSmtRecord struct {
	Version     ReverseSmtRecordVersion
	Action      ReverseSmtRecordAction
	Signature   []byte
	SignType    uint8
	Address     []byte
	Proof       []byte
	PrevNonce   uint32 `witness:",omitempty"`
	PrevAccount string
	NextRoot    []byte
	NextAccount string
}
