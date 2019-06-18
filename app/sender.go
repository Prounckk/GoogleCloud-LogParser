package app

import (
	"bufio"
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	ftp "github.com/martinr92/goftp"
	"log"
	"os"
	"strings"
)

type GCSEvent struct {
	Bucket         string `json:"bucket"`
	Name           string `json:"name"`
	Metageneration string `json:"metageneration"`
	ResourceState  string `json:"resourceState"`
}

//  consumes a GCS event.
func GCSwatcher(ctx context.Context, e GCSEvent) error {
	if e.ResourceState == "not_exists" {
		log.Printf("File %v deleted.", e.Name)
		return nil
	}
	if e.Metageneration == "1" {
		// The metageneration attribute is updated on metadata changes.
		// The on create value is 1.
		log.Printf("File %v created.", e.Name)
		file := e.Name
		bucket := e.Bucket
		if strings.Contains(file, "json") {
			Reader(file, bucket)
		}
		return nil
	}
	log.Printf("File %v metadata updated.", e.Name)
	return nil
}

func Reader(file, bucket string) {

	ctx := context.Background()
	// Creates a client to connect to the bucket.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to open: %v", err)
	}

	scanner := bufio.NewScanner(rc)

	var j map[string]interface{}


	ff := strings.Replace(strings.Replace (file,"/", "_", -1), "json", "txt", 1)
	filename := "/tmp/" + ff

	for scanner.Scan() {
		t := scanner.Text()
		err = json.NewDecoder(bytes.NewReader([]byte(t))).Decode(&j)
		if textPayload, ok := j["textPayload"]; ok {

			fmt.Println()
			file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}

			datawriter := bufio.NewWriter(file)

			datawriter.WriteString(textPayload.(string) + "\n")

			datawriter.Flush()
			file.Close()
		}
	}

	fmt.Println("updated object:", filename)

	SenderToFTP(filename, ff)
}
func SenderToFTP(filename, ff string) {

	ftphost := os.Getenv("FTPHOST")
	ftplogin := os.Getenv("FTPLOGIN")
	ftppass := os.Getenv("FTPPASS")
	ftpfolder := os.Getenv("FTPFOLDER")

	ftpClient, err := ftp.NewFtp(ftphost)

	//uncomment if you need it
	//ftpClient.ActiveMode = true
	if err != nil {
		fmt.Printf("failed connect to ftp: %s", err)
	}
	if err = ftpClient.Login(ftplogin, ftppass); err != nil {
		fmt.Printf("failed to get autorisation: %s", err)
	}
	if err = ftpClient.OpenDirectory(ftpfolder); err != nil {
		fmt.Printf("failed to open the folder: %s", err)
	}
	if err = ftpClient.Upload(filename, ff); err != nil {
		fmt.Printf("failed to upload the file: %s", err)
	}
	defer ftpClient.Close()

	df := os.Remove(filename)
	fmt.Printf("%s File Deleted ", df)
}
