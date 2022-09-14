package thelper

import (
	"fmt"
)

var DBName string
var MySQLDsn string

func SetUp() {
	DBName = "gravityr"
	MySQLDsn = fmt.Sprintf("%s@/%s", "root", DBName)
}
