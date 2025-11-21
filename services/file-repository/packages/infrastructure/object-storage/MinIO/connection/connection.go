package minioconnection

import (
	"context"
	"time"
	StorageConnection "vega/packages/infrastructure/object-storage/connection"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Manager = &defaultConnectionManager{
	// It would be 0 (StatusDisconnected) even if left it uninitialized,
	// but anyway better to specify this explicitly
	status: StorageConnection.Disconnected,
}

type defaultConnectionManager struct {
	Client *minio.Client
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

	m.Client = client
	m.status = StorageConnection.Connected

	return nil
}

// The Go MinIO client manages network connections automatically using Go's standard http.Client,
// which handles connection pooling and cleanup automatically, so there no need in this method.
// P.S. Furthermore MinIO client has no methods that even makes manual disconnection posible, so there no choice.
func (m *defaultConnectionManager) Disconnect() error {
	m.status = StorageConnection.Disconnected
	return nil
}

func (m *defaultConnectionManager) Ping(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// There are no built-in Ping function
	_, err := m.Client.BucketExists(ctx, "vega--health-check")
	if err != nil {
		return err
	}
	return nil
}
