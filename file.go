package main

import (
	"cloud.google.com/go/functions/metadata"
	"fmt"
	"time"
)

func main() {
	Reader()
}

type GCSEvent struct {
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}

//HelloGCSInfo prints information about a GCS event.
func HelloGCSInfo(ctx context.Context, e GCSEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}

	BucketName := e.Bucket
	FileName := e.Name
	fmt.Println(BucketName + FileName)
	return nil
}
