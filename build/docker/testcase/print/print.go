package print

import (
	"log"
	"os"
)

func fileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	filename := "/tmp/tst.sql"

	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}

	log.SetOutput(logFile)
	log.SetPrefix("[fuzz]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func PrintSQL(sql string) {
	log.Println(sql)
}
