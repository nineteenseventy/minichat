package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ContainsCaseInsensitive(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

func main() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	allowedBucketNames := strings.Split(os.Getenv("MINIO_ALLOWED_BUCKET_NAMES"), ",")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	handle := func(w http.ResponseWriter, request *http.Request) {
		bucket := request.PathValue("bucket")
		object := request.PathValue("object")

		if !ContainsCaseInsensitive(allowedBucketNames, bucket) {
			http.Error(w, "Bucket not allowed", http.StatusForbidden)
			return
		}

		if object == "" {
			http.Error(w, "Object not specified", http.StatusBadRequest)
			return
		}

		bucketExists, err := minioClient.BucketExists(request.Context(), bucket)

		if err != nil {
			http.Error(w, "Error checking bucket", http.StatusInternalServerError)
			return
		}

		if !bucketExists {
			http.Error(w, "Bucket not found", http.StatusNotFound)
			return
		}

		objectInfo, err := minioClient.StatObject(request.Context(), bucket, object, minio.StatObjectOptions{})
		if err != nil {
			http.Error(w, "Object not found", http.StatusNotFound)
			return
		}

		objectReader, err := minioClient.GetObject(request.Context(), bucket, object, minio.GetObjectOptions{})
		if err != nil {
			http.Error(w, "Error getting object", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(objectInfo.Size, 10))
		w.Header().Set("Content-Type", objectInfo.ContentType)
		w.Header().Set("Last-Modified", objectInfo.LastModified.Format(http.TimeFormat))

		_, err = io.Copy(w, objectReader)
	}

	http.HandleFunc("/{bucket}/{object...}", handle)
	http.ListenAndServe(":8080", nil)
}
