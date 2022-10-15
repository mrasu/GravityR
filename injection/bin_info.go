package injection

var BinInfo = binInfo{}

type binInfo struct {
	Version string
	Commit  string
}

func SetBinInfo(version, commit string) {
	BinInfo.Version = version
	BinInfo.Commit = commit
}
