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

	"github.com/cheggaaa/pb/v3"
)

func Print(s interface{}) {
	os.Setenv("LC_ALL", "cs_CZ.UTF-8")
	switch value := s.(type) {
	case string:
		//windows1250Encoder := charmap.Windows1250.NewEncoder()
		//res, _ := windows1250Encoder.String(value)
		fmt.Println(value)

	case int:
		// Převedení int na string
		strValue := strconv.Itoa(value)
		// Zde zpracuj string
		fmt.Println(strValue)
	default:
		fmt.Println("Nerozpoznaný typ:", value)
	}
}

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

type progressBarWriter struct {
	bar *pb.ProgressBar
}

func (pw *progressBarWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.bar.Add(n)
	return n, nil
}

// from ChatGPT
func DownloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create progress bar
	fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	bar := pb.Full.Start(fileSize)
	bar.SetWidth(80)

	// Create multi writer
	writer := io.MultiWriter(out, &progressBarWriter{bar: bar})

	// Write the body to file with progress
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	// Finish progress bar
	bar.Finish()

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
