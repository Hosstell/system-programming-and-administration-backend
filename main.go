package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name   string
	Path   string
	IsFile bool
}

func getListFile(path string) []File {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	fileList := []File{}
	for _, f := range files {
		fileList = append(fileList, File{
			f.Name(),
			filepath.Join(path, f.Name()),
			!f.IsDir(),
		})
	}
	return fileList
}

func moveFile(filePath string, newPath string) bool {
	err := os.Rename(filePath, newPath)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func makeNewBasket(newName string) string {
	pathToBaskets := "/home/andrey/baskets/"
	pathToNewDir := pathToBaskets + newName
	os.MkdirAll(pathToNewDir, os.ModePerm)
	return pathToNewDir
}

func deleteFile(pathFile string) {
	os.RemoveAll(pathFile + "/")
	var err = os.Remove(pathFile)
	if isError(err) {
		return
	}

	fmt.Println("File Deleted")
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func main() {
	pathToBaskets := "/home/andrey/baskets/"

	http.HandleFunc("/file_list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		value := strings.Split(r.URL.RawQuery, "=")
		newValue, _ := url.QueryUnescape(value[1])

		files := getListFile(newValue)

		jsonFiles, _ := json.Marshal(files)
		fmt.Fprintf(w, string(jsonFiles))
	})

	http.HandleFunc("/new_basket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		value := strings.Split(r.URL.RawQuery, "=")
		fmt.Println(value[1])
		newValue, _ := url.QueryUnescape(value[1])
		makeNewBasket(newValue)
	})

	http.HandleFunc("/move_file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		params := strings.Split(r.URL.RawQuery, "&")
		oldFile := strings.Split(params[0], "=")
		newFile := strings.Split(params[1], "=")

		newOldFile, _ := url.QueryUnescape(oldFile[1])
		newNewFile, _ := url.QueryUnescape(newFile[1])

		fileName := strings.Split(newOldFile, "/")

		moveFile(newOldFile, pathToBaskets+newNewFile+"/"+fileName[len(fileName)-1])
		fmt.Println(newOldFile, pathToBaskets+newNewFile+"/"+fileName[len(fileName)-1])
	})

	http.HandleFunc("/restore_file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		params := strings.Split(r.URL.RawQuery, "&")
		oldFile := strings.Split(params[0], "=")
		newFile := strings.Split(params[1], "=")

		newOldFile, _ := url.QueryUnescape(oldFile[1])
		newNewFile, _ := url.QueryUnescape(newFile[1])

		fileName := strings.Split(newOldFile, "/")

		moveFile(newOldFile, newNewFile+"/"+fileName[len(fileName)-1])
		fmt.Println(newOldFile, newNewFile+"/"+fileName[len(fileName)-1])
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		newFile := strings.Split(r.URL.RawQuery, "=")
		newNewFile, _ := url.QueryUnescape(newFile[1])
		fmt.Println(newNewFile)
		deleteFile(newNewFile)
	})

	http.ListenAndServe(":8888", nil)
}
