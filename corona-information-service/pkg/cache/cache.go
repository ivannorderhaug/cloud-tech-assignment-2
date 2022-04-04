package cache

// New cache*/
func New() map[string]interface{} {
	return make(map[string]interface{}, 0)
}

// Get item from cache using key
func Get(cache map[string]interface{}, key string) interface{} {
	val, exist := cache[key]
	if !exist {
		return nil
	}
	return val
}

// Put value in cache by using key
func Put(cache map[string]interface{}, key, value string) {
	cache[key] = value
}
