package app

import (
	"bufio"
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	ftp "github.com/martinr92/goftp"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"

	"reflect"
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
		//file := e.Name
		//bucket := e.Bucket
		//Reader(file, bucket)
		return nil
	}
	log.Printf("File %v metadata updated.", e.Name)
	return nil
}

//func Reader(file, bucket string) {
func Reader() {
	file := "nginx/2019/06/10/nginx_2019_06_09_01_00_00_01_59_59_S0.json"
	bucket := "corbeil"

	ctx := context.Background()
	// Creates a client to connect to the bucket.
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/Users/sergei.eremeev/.config/lightspeed/mkt-hq-website/admins-sa2.json"))
	//client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to open: %v", err)
	}

	scanner := bufio.NewScanner(rc)

	var j map[string]interface{}

	filename := strings.Replace(file, "json", "txt", 1)
	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	_ = wc
	wc.ContentType = "text/plain"
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	for scanner.Scan() {
		t := scanner.Text()
		err = json.NewDecoder(bytes.NewReader([]byte(t))).Decode(&j)

		if textPayload, ok := j["textPayload"]; ok {

			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}

			if _, err := wc.Write([]byte(textPayload.(string))); err != nil {
				fmt.Println(err)
			}

		}
	}

	if err := wc.Close(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("updated object:", wc.Attrs())

	SenderToFTP(filename)
}
func SenderToFTP(filename string) {
	ctx := context.Background()
	//client, err := storage.NewClient(ctx)
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/Users/sergei.eremeev/.config/lightspeed/mkt-hq-website/admins-sa2.json"))

	bucket := "corbeil"
	rcc, err := client.Bucket(bucket).Object(filename).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to open: %v", err)
	}

	slurp, err := ioutil.ReadAll(rcc)
	rcc.Close()

	fmt.Printf("File contents: %s", slurp)

	fmt.Println(reflect.TypeOf(rcc))
	fmt.Println(rcc)
	//ftphost := os.Getenv("FTPHOST")
	//ftplogin := os.Getenv("FTPLOGIN")
	//ftppass := os.Getenv("FTPPASS")
	//ftpfolder := os.Getenv("FTPFOLDER")


	ftphost := "ftp.oncrawl.com:21"
	ftplogin := "hqmarketing"
	ftppass := "hqdeepsthgiL8*"
	ftpfolder := "/HQ - COM/nginx/"

	ftpClient, err := ftp.NewFtp(ftphost)
	fileserv := "https://storage.googleapis.com/corbeil/nginx/2019/06/10/nginx_2019_06_09_01_00_00_01_59_59_S0.txt"
	filenamenew := "foo.txt"
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
	if err = ftpClient.Upload(fileserv, filenamenew); err != nil {
		fmt.Printf("failed to upload the file: %s", err)
	}
	defer ftpClient.Close()
	//once all is done, let's delete the file
	//df := os.Remove(fn)
	//fmt.Printf("%s File Deleted ", df)
}
