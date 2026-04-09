package r2

import (
	"time"

	"github.com/thanvuc/go-core-lib/storage"
)

func PresignURLs(contentType string, fileName string) storage.PresignOptions {
	return storage.PresignOptions{
		KeyPrefix:   "ai-sprint-generation/",
		ContentType: contentType,
		Expiry:      5 * time.Minute,
		FileName:    fileName,
	}
}
