package r2

import (
	"team_service/internal/infrastructure/share/settings"

	"github.com/thanvuc/go-core-lib/log"
	"github.com/thanvuc/go-core-lib/storage"
)

func NewR2Client(cfg settings.R2, logger log.Logger) (*storage.R2Client, error) {
	r2Clients, err := storage.NewClient(storage.Config{
		AccountID: cfg.AccountID,
		Endpoint:  cfg.Endpoint,
		AccessKey: cfg.AccessKeyID,
		SecretKey: cfg.SecrecAccessKey,
		Bucket:    cfg.BucketName,
		UseSSL:    cfg.UseSSL,
		PublicURL: cfg.PublicURL,
	})

	if err != nil {
		return nil, err
	}

	logger.Info("R2 client created successfully", "")

	return r2Clients, nil
}
