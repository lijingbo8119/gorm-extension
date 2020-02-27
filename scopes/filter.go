package scopes

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

/*
	filter参数类型
	1、{"field":"some value"} field='some value'
	2、{"field >=":"some value"} field >='some value' 注意key中必须有空格
	3、{"some operation = ?":"some value"} some operation = 'some value' 问号为gorm原生替代符号
 */

func Filter(filter map[string]interface{}, model interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}

		operate := func(k string, v interface{}, db *gorm.DB) *gorm.DB {
			kSlice := strings.Split(k, " ")
			query := ""
			if len(kSlice) == 1 {
				query = k + " = ?"
			} else {
				query = kSlice[0] + " " + kSlice[1] + " ?"
			}
			return db.Where(query, gconv.String(v))
		}

		t := reflect.TypeOf(model)
		for k, v := range filter {
			hasReplaceMark := strings.Contains(k, "?")
			if hasReplaceMark {
				db = db.Where(k, v)
				continue
			}

			kSlice := strings.Split(k, " ")
			K := ""
			switch len(kSlice) {
			case 1:
				K = k
			case 2:
				K = kSlice[0]
			default:
				K = kSlice[0]
			}

			switch K {
			case "id", "create_time", "update_time":
				db = operate(k, v, db)
			default:
				for i := 0; i < t.NumField(); i++ {
					jsonTag := t.Field(i).Tag.Get("json")
					if jsonTag == k {
						db = operate(k, v, db)
					}
				}
			}
		}
		return db
	}
}
