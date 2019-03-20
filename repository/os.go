package repository

// IOSRepository is interface for interaction with local files
// include configuration files
type IOSRepository interface {
	Copy(src, dst string) error
	GetConfig(key string) string
	SetConfig(key, value string) error
}
