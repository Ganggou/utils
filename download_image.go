package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	jsonFile, _ := os.Open("data.json")
	// [[url, name], [url, name]]
	byteValue, _ := ioutil.ReadAll(jsonFile)
	list := make([][]string, 70)
	json.Unmarshal(byteValue, &list)
	for _, i := range list {
		putFile("images/"+i[1], i[0])
	}
}

func putFile(fileName, url string) {
	client := httpClient()
	file := createFile(fileName)
	resp, err := client.Get(url)

	checkError(err)

	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	checkError(err)

	fmt.Println(fmt.Sprintf("Just Downloaded a file %s with size %v", fileName, size))
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func createFile(fileName string) *os.File {
	file, err := os.Create(fileName)

	checkError(err)
	return file
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
