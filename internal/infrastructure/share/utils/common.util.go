package utils

import (
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func ToUUID(id string) (pgtype.UUID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgtype.UUID{
		Bytes: uuid,
		Valid: true,
	}, nil
}

func StringPtrToPgUUID(s *string) (pgtype.UUID, error) {
	if s == nil || *s == "" {
		return pgtype.UUID{Valid: false}, nil // SQL NULL
	}

	u, err := uuid.Parse(*s)
	if err != nil {
		return pgtype.UUID{}, err
	}

	var bytes [16]byte
	copy(bytes[:], u[:])

	return pgtype.UUID{
		Bytes: bytes,
		Valid: true,
	}, nil
}

func RoundToTwoDecimal(val float64) float64 {
	return math.Round(val*100) / 100
}

func FromPgTypeTimeToUnix(t pgtype.Timestamp) *int64 {
	if !t.Valid {
		return nil
	}
	unixTime := t.Time.Unix()
	return &unixTime
}

func FromPgTypeTimeStamptZToUnix(t pgtype.Timestamptz) *int64 {
	if !t.Valid {
		return nil
	}
	unixTime := t.Time.Unix()
	return &unixTime
}

func Difference[T comparable](a, b []T) []T {
	m := make(map[T]struct{}, len(b))
	for _, item := range b {
		m[item] = struct{}{}
	}

	var diff []T
	for _, item := range a {
		if _, found := m[item]; !found {
			diff = append(diff, item)
		}
	}
	return diff
}

func ToJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func ToBoolPointer(b bool) *bool {
	return &b
}

func Ptr[T any](v T) *T {
	return &v
}

func FromTimeStampToTimePtr(timestamp *int64) *time.Time {
	if timestamp == nil {
		return nil
	}
	// 13 digits timestamp
	t := time.UnixMilli(*timestamp).UTC()
	return &t
}

func FromTimeStampToTime(timestamp int64) time.Time {
	// 13 digits timestamp
	return time.UnixMilli(timestamp).UTC()
}

func FromTimeToTimeStamp(t time.Time) int64 {
	return t.UnixMilli()
}

func FromTimePtrToTimeStamp(t *time.Time) *int64 {
	if t == nil {
		return nil
	}

	timestamp := t.UnixMilli()
	return &timestamp
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func SafeStringWithDefault(ptr *string, def string) string {
	if ptr == nil || *ptr == "" {
		return def
	}
	return *ptr
}

func SafeInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func SafeFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func FromInt64ToInt32(i *int64) *int32 {
	if i == nil {
		return nil
	}
	i32 := int32(*i)
	return &i32
}
