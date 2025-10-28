package minioconnection

import (
	StorageConnection "vega/packages/infrastructure/object-storage/connection"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Init() StorageConnection.Manager {
	return &defaultConnectionManager{
		// It would be 0 (StatusDisconnected) even if left it uninitialized,
		// but anyway better to specify this explicitly
		status: StorageConnection.Disconnected,
	}
}

type defaultConnectionManager struct {
	client *minio.Client
	status StorageConnection.Status
}

func (m *defaultConnectionManager) Status() StorageConnection.Status {
	return m.status
}

func (m *defaultConnectionManager) Connect(cfg *StorageConnection.Config) error {
	client, err := minio.New(cfg.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Login, cfg.Password, cfg.Token),
		Secure: cfg.Secure,
	})
	if err != nil {
		return err
	}

	m.client = client
	m.status = StorageConnection.Connected

	return nil
}

// The Go MinIO client manages network connections automatically using Go's standard http.Client,
// which handles connection pooling and cleanup automatically, so there no need in this method.
// P.S. Furthermore MinIO client has no methods that even makes manual disconnection posible, so there no choice.
func (m *defaultConnectionManager) Disconnect() error {
	return nil
}

