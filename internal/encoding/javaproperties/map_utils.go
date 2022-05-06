package javaproperties

import (
	"strings"

	"github.com/spf13/cast"
)

// THIS CODE IS COPIED HERE: IT SHOULD NOT BE MODIFIED
// AT SOME POINT IT WILL BE MOVED TO A COMMON PLACE
// deepSearch scans deep maps, following the key indexes listed in the
// sequence "path".
// The last value is expected to be another map, and is returned.
//
// In case intermediate keys do not exist, or map to a non-map value,
// a new map is created and inserted, and the search continues from there:
// the initial map "m" may be modified!
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}

// flattenAndMergeMap recursively flattens the given map into a new map
// Code is based on the function with the same name in tha main package.
// TODO: move it to a common place
func flattenAndMergeMap(shadow map[string]interface{}, m map[string]interface{}, prefix string, delimiter string) map[string]interface{} {
	if shadow != nil && prefix != "" && shadow[prefix] != nil {
		// prefix is shadowed => nothing more to flatten
		return shadow
	}
	if shadow == nil {
		shadow = make(map[string]interface{})
	}

	var m2 map[string]interface{}
	if prefix != "" {
		prefix += delimiter
	}
	for k, val := range m {
		fullKey := prefix + k
		switch val.(type) {
		case map[string]interface{}:
			m2 = val.(map[string]interface{})
		case map[interface{}]interface{}:
			m2 = cast.ToStringMap(val)
		default:
			// immediate value
			shadow[strings.ToLower(fullKey)] = val
			continue
		}
		// recursively merge to shadow map
		shadow = flattenAndMergeMap(shadow, m2, fullKey, delimiter)
	}
	return shadow
}

// processlist
func processlist(resMap *map[string]interface{}) {
	count := len(*resMap)
	for k, v := range *resMap {
		// 判断是否key包含了数组的符号
		if strings.Contains(k, "[") && strings.Contains(k, "]") && (strings.LastIndex(k, "]")+1 == len(k)) {
			x, ok := (*resMap)[string([]rune(k)[:strings.LastIndex(k, "[")])]
			// 已经存在
			if ok {
				// 如果是v是map,处理它的孩子并把它加到数组里面
				if xxx, ok := v.(map[string]interface{}); ok {
					processlist(&xxx)
				}
				x = append(x.([]interface{}), v)
				// 创建一个数组子项，给后面生成yaml使用的数组
				(*resMap)[string([]rune(k)[:strings.LastIndex(k, "[")])] = x
				delete(*resMap, k)
				count--
				// key已经全部遍历完成
				if count == 0 {
					return
				}
				// 该项处理完成 继续处理下一项
				continue
			} else {
				if xxx, ok := v.(map[string]interface{}); ok {
					processlist(&xxx)
				}
				// 首次
				x = append([]interface{}{}, v)
				(*resMap)[string([]rune(k)[:strings.LastIndex(k, "[")])] = x
				delete(*resMap, k)
				count--
				// 遍历完成
				if count == 0 {
					return
				}
				continue
			}
		}
		// 如果是字符串说明已经到达最后直接下次循环
		if _, ok := v.(string); ok {
			count--
			if count <= 0 {
				return
			}
			continue
		}
		// 因为加入了数组所以这里遍历到需要把它跳过
		if _, ok := v.([]interface{}); ok {
			continue
		}
		// 正常的情况不带数组标记的
		x := v.(map[string]interface{})
		processlist(&x)
		count--
		if count == 0 {
			return
		}
	}
}
