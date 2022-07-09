package thelper

import (
	"fmt"
)

var DbName string
var MySQLDsn string

func SetUp() {
	DbName = "gravityr"
	MySQLDsn = fmt.Sprintf("%s@/%s", "root", DbName)
}
