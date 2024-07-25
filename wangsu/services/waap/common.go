package waap

/*
除了 exceptKey 之前的属性是否发生变化，暂不支持嵌套类型
*/
func CheckChangeExcept(oldConfig []interface{}, newConfig []interface{}, exceptKey string) bool {
	if (oldConfig == nil || len(oldConfig) == 0) || (newConfig == nil || len(newConfig) == 0) {
		return false
	}
	// 遍历每个字段
	oldConfigMap := oldConfig[0].(map[string]interface{})
	newConfigMap := newConfig[0].(map[string]interface{})

	for key := range oldConfigMap {
		if key == exceptKey {
			continue
		}
		_, ok := oldConfigMap[key].(string)
		if ok {
			// value 是 string 类型
			if oldConfigMap[key] != newConfigMap[key] {
				// 如果有变化，返回错误
				return true
			}
		}
	}
	return false
}

/*
*
取数组slice1的差集
*/
func Difference(slice1, slice2 []string) []string {
	diff := make([]string, 0)
	hash := make(map[string]bool)

	for _, item := range slice1 {
		hash[item] = true
	}

	for _, item := range slice2 {
		if _, found := hash[item]; found {
			delete(hash, item)
		}
	}

	for item := range hash {
		diff = append(diff, item)
	}

	return diff
}
