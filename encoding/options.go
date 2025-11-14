package encoding

type prettyOpt struct{}

func (p prettyOpt) applyHumanEncodeOption(opt *humanEncodeConfig) {
	opt.compact = true
}

type hideParamName struct{}

func (p hideParamName) applyHumanEncodeOption(opt *humanEncodeConfig) {
	opt.showParamName = false
}

func Compact() *prettyOpt {
	return &prettyOpt{}
}

func Unnamed() *hideParamName {
	return &hideParamName{}
}
