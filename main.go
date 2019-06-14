package main

import (
	"github.com/prounckk/GoogleCloudLogParser/app"
	"os"
)

func main() {
	file :=os.Getenv("FILE")
	bucket := os.Getenv("BUCKET")
	app.Reader(file, bucket)
}