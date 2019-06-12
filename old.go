package p

import (
	"context"
	"log"
)

type GCSEvent struct {
	Bucket         string `json:"bucket"`
	Name           string `json:"name"`
	Metageneration string `json:"metageneration"`
	ResourceState  string `json:"resourceState"`
}

// GCSwatcher consumes a GCS event.
func GCSwatcher(ctx context.Context, e GCSEvent) error {
	if e.ResourceState == "not_exists" {
		log.Printf("File %v deleted.", e.Name)
		return nil
	}
	if e.Metageneration == "1" {
		// The metageneration attribute is updated on metadata changes.
		// The on create value is 1.
		log.Printf("File %v created.", e.Name)

		return nil
	}
	log.Printf("File %v metadata updated.", e.Name)
	return nil
}
