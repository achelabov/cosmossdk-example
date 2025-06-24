package types

func NewMsgCreateNote(creator string, text string) *MsgCreateNote {
	return &MsgCreateNote{
		Creator: creator,
		Text:    text,
	}
}
