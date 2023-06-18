package tools

func Flatten(m map[string]interface{}) map[string]interface{} {
	f := make(map[string]interface{})
	for k, v := range m {
		switch v.(type) {
		case map[string]interface{}:
			for k2, v2 := range Flatten(v.(map[string]interface{})) {
				f[k+"."+k2] = v2
			}
		default:
			f[k] = v
		}
	}
	return f
}
