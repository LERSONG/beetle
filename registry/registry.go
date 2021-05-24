package registry

type Registry interface {
	Register(ops ...RegisterOption) error
	Deregister() error
	String() string
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}
