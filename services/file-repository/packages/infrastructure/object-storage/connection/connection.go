package storageconnection

type Status uint8

const (
	Disconnected Status = iota
	Connected
)

// Better to turn this into interface, but it will be kinda hard due to how each
// storage driver handles connections. So for now better to leave it as it is.
// Even if this would require some refactored in the future it won't cause much problems.
// The same goes for DI.
type Config struct {
	URL 		string
	Login 		string
	Password 	string
	Token		string
	Secure		bool
}

type Manager interface {
	Status() Status
	Connect(cfg *Config) error
	Disconnect() error
}

