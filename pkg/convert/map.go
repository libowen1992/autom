package convert

import (
	"encoding/json"
)

func MapStr2Str(m map[string]interface{}) string {
	mj, _ := json.Marshal(m)
	return string(mj)
}
