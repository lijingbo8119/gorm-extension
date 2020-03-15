package scopes

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/jinzhu/gorm"

	"reflect"
	"strings"
)

func Sorter(sorter map[string]interface{}, model interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if sorter == nil {
			return db
		}

		t := reflect.TypeOf(model)
		for k, v := range sorter {
			switch strings.ToLower(gconv.String(v)) {
			case "", "asc", "desc":
			default:
				g.Throw("Sorter value error:" + gconv.String(v))
			}

			switch k {
			case "id", "create_time", "update_time":
				db = db.Order(k + " " + gconv.String(v))
			default:
				tagKey := k
				if strings.Contains(k, "_formatted") {
					tagKey = strings.Split(k, "_formatted")[0]
				}

				for i := 0; i < t.NumField(); i++ {
					jsonTag := t.Field(i).Tag.Get("json")
					if jsonTag == tagKey {
						db = db.Order(tagKey + " " + gconv.String(v))
					}
				}
			}

		}
		return db
	}
}
