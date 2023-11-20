package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	//"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"

	"net/url"
	"path"

	"github.com/schollz/progressbar/v3"
)

// Print - print value
func Print(s interface{}) {
	os.Setenv("LC_ALL", "cs_CZ.UTF-8")
	switch value := s.(type) {
	case string:
		//windows1250Encoder := charmap.Windows1250.NewEncoder()
		//res, _ := windows1250Encoder.String(value)
		fmt.Println(value)

	case int:
		// Parase int to string
		strValue := strconv.Itoa(value)
		fmt.Println(strValue)
	default:
		fmt.Println("Nerozpoznaný typ:", value)
	}
}

// Read user input, return string
func Read() string {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSuffix(s, "\n")
	return s
}

func Contains(s string, array []string) bool {
	for _, value := range array {
		if s == value {
			return true
		}
	}
	return false
}

// from ChatGPT
func Download(destinationPath, downloadUrl string) error {
	tempDestinationPath := destinationPath + ".tmp"
	req, _ := http.NewRequest("GET", downloadUrl, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile(tempDestinationPath, os.O_CREATE|os.O_WRONLY, 0644)

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	os.Rename(tempDestinationPath, destinationPath)
	return nil
}

// from ChatGPT
func ExtractFileName(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	fileName := path.Base(parsedURL.Path)
	return fileName, nil
}
