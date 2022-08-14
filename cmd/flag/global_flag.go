package flag

var GlobalFlag = &globalFlag{}

type globalFlag struct {
	Verbose bool
	UseMock bool
}
