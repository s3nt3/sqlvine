package coverage

import (
	"fmt"
	"log"
	"strings"
)

func PrintSQL(sql string) {
	if len(sql) < 256 {
		if strings.Contains(sql, "MAX(") {
			if strings.Contains(sql, " JOIN ") {
				if len(sql) > 64 {
					panic(fmt.Sprintf("bingo: %s", sql))
				}
				log.Println(sql)
			}
			log.Println(sql)
			return
		}
		return
	}
}
