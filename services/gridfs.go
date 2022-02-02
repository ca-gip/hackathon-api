package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"hackathon-api/configs"
)

var conn = configs.ConnectDB()
var db = conn.Database("donationFiles")

func UploadFile(file []byte, filename string) (int, error) {

	bucket, err := gridfs.NewBucket(
		db,
	)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(file)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	log.Printf("Write file to DB was successful. File size: %d \n", fileSize)
	return fileSize, nil
}

func DownloadFile(fileName string) ([]byte, error) {

	// CRUD operation
	fsFiles := db.Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// you can print out the results
	fmt.Println(results)

	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Printf("File size to download: %v\n", dStream)
	return buf.Bytes(), nil
}
