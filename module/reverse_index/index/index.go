package index

var cache map[string]string = make(map[string]string)

func Set(key, value string) {
	cache[key] = value
}

func Get(key string) (value string, ok bool) {
	value, ok = cache[key]
	return
}

func Search(prefix string) (values string) {
	return ""
}
