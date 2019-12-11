package bundle_types

import "time"

type Metadata struct {
	App         string    `json:"app"`
	Timestamp   time.Time `json:"timestamp"`
	Annotations []string  `json:"annotations"`
}
