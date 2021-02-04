package cache

// Object represent cache interface
type Object interface {
	Get(args string) (string, error)
	Set(key, val string) error
	IsNil(i interface{}) bool
}
