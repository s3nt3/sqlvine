package print

import (
	"log"
	"strings"
)

func PrintSQL(sql string) {
	if len(sql) < 128 {
		switch {
		case strings.Contains(sql, "COUNT("):
			if strings.Contains(sql, "COUNT((SELECT ") {
				log.Panicln(sql)
			} else {
				log.Println(sql)
			}
		case strings.Contains(sql, "SUM("):
			if strings.Contains(sql, "SUM((SELECT ") {
				log.Panicln(sql)
			} else {
				log.Println(sql)
			}
		case strings.Contains(sql, "MAX("):
			if strings.Contains(sql, "MAX((SELECT ") {
				log.Panicln(sql)
			} else {
				log.Println(sql)
			}
		case strings.Contains(sql, "GROUP_CONCAT("):
			if strings.Contains(sql, "GROUP_CONCAT((SELECT ") {
				log.Panicln(sql)
			} else {
				log.Println(sql)
			}
		case strings.Contains(sql, "JSON_OBJECTAGG("):
			if strings.Contains(sql, "JSON_OBJECTAGG((SELECT ") {
				log.Panicln(sql)
			} else {
				log.Println(sql)
			}
		}
	}
}
