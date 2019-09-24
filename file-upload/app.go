package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var MaxFileSize int64
var UploadDir string

type Message struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func loadEnv() {
	_ = godotenv.Load()
	var err error
	MaxFileSize, err = strconv.ParseInt(os.Getenv("MAX_FILE_SIZE_MB"), 10, 64)
	if err != nil {
		MaxFileSize = 0
	}
	UploadDir = os.Getenv("UploadDir")
	fmt.Print(UploadDir)
	if UploadDir == "" {
		UploadDir = "upload"
	}
}

func respondMessage(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	messageStruct := Message{
		Code:    code,
		Message: message,
	}

	messageBytes, _ := json.Marshal(messageStruct)
	w.Write(messageBytes)
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(MaxFileSize << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		respondMessage(w, 422, fmt.Sprintf("Max file size is %dMB", MaxFileSize))
		return
	}
	defer file.Close()

	filenameParts := strings.Split(handler.Filename, ".")
	extension := filenameParts[len(filenameParts)-1]
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Print(err)
		respondMessage(w, 422,"Can't read file")
		return
	}

	newFilename := fmt.Sprintf("%s.%s", generateRandomString(8), extension)
	newFilepath := fmt.Sprintf("%s/%s", UploadDir, newFilename)
	err = ioutil.WriteFile(newFilepath, fileBytes, 0644)
	if err != nil {
		fmt.Print("500")
		fmt.Print(err)
		respondMessage(w, 500, "Can't write file")
		return
	}

	respondMessage(w, 200, newFilename)
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query()["file"][0]
	fileData, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", UploadDir, filename))
	if err != nil {
		fmt.Print(err)
		respondMessage(w, 404, "File not found")
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write(fileData)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	http.ListenAndServe(":8000", nil)
}

func main() {
	loadEnv()
	setupRoutes()
}