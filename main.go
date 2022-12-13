package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
)

type Employee struct {
	Name string `json:"name"`
}
type Handler struct {
}
type UploadHandler struct {
	UploadDir string
	HostAddr  string
}

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fileNameLength = 16
)

func generateFileName() string {
	b := make([]byte, fileNameLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		name := r.FormValue("name")

		fmt.Fprintf(w, "GET \"name\": %s", name)
	case http.MethodPost:
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable s", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var employee Employee
		err = json.NewDecoder(r.Body).Decode(&employee)
		if err != nil {
			http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "got a new employee!\nName: %s\n",
			employee.Name,
		)
	}
}
func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "no read file1", http.StatusBadRequest)
			return
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "no read file2", http.StatusBadRequest)
			return
		}
		filePath := h.UploadDir + "\\" + header.Filename
		err = ioutil.WriteFile(filePath, data, 0777)
		if err != nil {
			log.Println(err)
			http.Error(w, "no save file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "successfully load %s", header.Filename)

		fileLink := h.HostAddr + "/" + header.Filename
		fmt.Fprintln(w, fileLink)

	case http.MethodGet:
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "no read file555", http.StatusBadRequest)
			return
		}
		defer file.Close()
		namefile := header.Filename
		size := header.Size
		buffer := make([]byte, size)
		file.Read(buffer)
		//type := http.DetectContentType(buffer)
		files, _ := filepath.Glob("*.txt")
		fmt.Printf("%q\n", files)

		fmt.Fprintf(w, "name: %s\n , size: %d", namefile, size)

	}

}
func main() {
	handler := &Handler{}
	http.Handle("/", handler)
	uploadHandler := &UploadHandler{
		UploadDir: "upload",
	}
	http.Handle("/upload", uploadHandler)
	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	srv.ListenAndServe()
	dirToServe := http.Dir(uploadHandler.UploadDir)
	fs := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fs.ListenAndServe()
}
