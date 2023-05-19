package coverage

import (
	"fmt"
	"strconv"
	"strings"
)

func MarkCovered(root interface{}, path string, lookupTable map[string]bool, discriminatorTable map[string]string) {
	if root == nil {
		return
	}

	if _, ok := lookupTable[path]; ok {
		lookupTable[path] = true
	}

	// https://pkg.go.dev/encoding/json#Unmarshal
	switch value := root.(type) {
	case string:
		if _, exist := lookupTable[path+"()"]; exist {
			lookupTable[path+"()"] = true
			lookupTable[path+"("+value+")"] = true
		}

	case bool:
		lookupTable[path+"()"] = true
		lookupTable[path+"("+strconv.FormatBool(value)+")"] = true

	case float64:

	case []interface{}:
		path += "[]"
		lookupTable[path] = true
		for _, item := range value {
			MarkCovered(item, path, lookupTable, discriminatorTable)
		}

	case map[string]interface{}:
		discriminator, ok := discriminatorTable[path]
		if ok {
			step := ""
			for k, v := range value {
				if k == discriminator {
					step = "{" + discriminator + "(" + v.(string) + ")}"
					break
				}
			}
			if step == "" {
				panic(fmt.Errorf("block %s has no discriminator %s", value, discriminator))
			}
			path += step
			lookupTable[path] = true
		}

		for k, v := range value {
			if strings.Contains(k, ".") {
				k = "\"" + k + "\""
			}
			MarkCovered(v, strings.TrimLeft(path+"."+k, "."), lookupTable, discriminatorTable)
		}

	default:
		panic(fmt.Errorf("unexpect type %T for json unmarshaled value", value))
	}
}
