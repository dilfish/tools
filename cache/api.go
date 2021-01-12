package cache

type CacheObject interface {
	Get(args string) (string, error)
	Set(key, val string) error
	IsNil(i interface{}) bool
}
