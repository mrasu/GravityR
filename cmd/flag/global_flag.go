package flag

var GlobalFlag = &globalFlag{}

type globalFlag struct {
	Verbose bool
	UseMock bool
}

func (gf *globalFlag) GetAwsEndpoint() string {
	if gf.UseMock {
		return "http://localhost:8080"
	} else {
		return ""
	}
}
