package main

import (
	"bufio"
	"bytes"
	"cloud.google.com/go/functions/metadata"
	"context"
	"fmt"
	"os"

	"encoding/json"

	"time"

	"log"

	"cloud.google.com/go/storage"
	ftp "github.com/martinr92/goftp"
	"google.golang.org/api/option"
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

HelloGCSInfo prints information about a GCS event.
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

func Reader() {

	ctx := context.Background()
	// Creates a client to connect to the bucket.
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/Users/sergei.eremeev/.config/lightspeed/mkt-hq-website/admins-sa2.json"))
	//client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	file := "nginx/2019/06/10/nginx_2019_06_09_01_00_00_01_59_59_S0.json"
	//file := "nginx/2019/06/09/01:00:00_01:59:59_S0.txt"

	rc, err := client.Bucket("corbeil").Object(file).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to open: %v", err)
	}

	scanner := bufio.NewScanner(rc)

	var j map[string]interface{}
	filename := "test.txt"
	//filename := strings.Replace(file, "json", "txt", 1))
	ff, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	d := bufio.NewWriter(ff)

	for scanner.Scan() {
		t := scanner.Text()
		err = json.NewDecoder(bytes.NewReader([]byte(t))).Decode(&j)

		if textPayload, ok := j["textPayload"]; ok {

			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}

			d.WriteString(textPayload.(string))
			d.Flush()

		}
	}
	ff.Close()
	SenderToFTP(ff)

}
func SenderToFTP(*os.File) {

	ftphost :=os.Getenv("FTPHOST")
	ftplogin :=os.Getenv("FTPLOGIN")
	ftppass :=os.Getenv("FTPPASS")
	ftpfolder := os.Getenv("FTPFOLDER")

	ftpClient, err := ftp.NewFtp(ftphost)
	//ftpClient.ActiveMode = true
	if err != nil {
		log.Fatalf("failed connect to ftp: %s", err)
	}
	if err = ftpClient.Login(ftplogin, ftppass); err != nil {
		log.Fatalf("failed to get autorisation: %s", err)
	}
	if err = ftpClient.OpenDirectory(ftpfolder); err != nil {
		log.Fatalf("failed to open the folder: %s", err)
	}
	if err = ftpClient.Upload("test.txt", "file.txt"); err != nil {
		log.Fatalf("failed to upload the file: %s", err)
	}
	defer ftpClient.Close()
}
