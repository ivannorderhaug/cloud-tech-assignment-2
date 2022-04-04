package cache

// New returns a new empty map
func New() map[string]interface{} {
	return make(map[string]interface{}, 0)
}

// NewNestedMap returns a new empty nested map
func NewNestedMap() map[string]map[string]interface{} {
	return make(map[string]map[string]interface{}, 0)
}

// Get returns item from cache using single key
func Get(cache map[string]interface{}, key string) interface{} {
	val, exist := cache[key]
	if !exist {
		return nil
	}
	return val
}

// Put inserts value in cache by using key
func Put(cache map[string]interface{}, key string, value interface{}) {
	cache[key] = value
}

// GetNestedMap returns item from cache using multiple keys
func GetNestedMap(cache map[string]map[string]interface{}, key1 string, key2 string) interface{} {
	val, exist := cache[key1][key2]
	if !exist {
		return nil
	}
	return val
}

// PutNestedMap inserts value into cache using multiple keys
func PutNestedMap(cache map[string]map[string]interface{}, key1 string, key2 string, value interface{}) {
	inner, ok := cache[key1]
	if !ok {
		inner = make(map[string]interface{})
		cache[key1] = inner
	}
	inner[key2] = value
}
