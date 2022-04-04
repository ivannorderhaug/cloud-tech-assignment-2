package cache

// New cache*/
func New() map[string]interface{} {
	return make(map[string]interface{}, 0)
}

// Get item from cache using key
func Get(cache map[string]interface{}, key string) interface{} {
	return cache[key]
}

// Put value in cache by using key
func Put(cache map[string]interface{}, key, value string) {
	cache[key] = value
}

// Contains will search the cache for the specified key
func Contains(cache map[string]interface{}, key string) bool {
	for k, _ := range cache {
		if k == key {
			return true
		}
	}
	return false
}
