package infrastructure

type Config interface {
	Name() string
	MarshalConfig() ([]byte, error)
}
