package utils

func Merge(base map[string]string, overlay ...map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range base {
		result[k] = v
	}

	i := len(overlay)
	switch {
	case i == 1:
		for k, v := range overlay[0] {
			result[k] = v
		}
	case i > 1:
		result = Merge(Merge(base, overlay[1]), overlay[2:]...)
	}

	return result
}
func GetValueFrom(maps map[string]string, key string) string {
	for k, v := range maps {
		if k == key {
			return v
		}
	}
	return ""
}

func HasLabel(maps map[string]string, key string) bool {
	for k := range maps {
		if k == key {
			return true
		}
	}
	return false
}
