package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	//"golang.org/x/text/encoding/charmap"

	"net/url"
	"path"
)

// Remove trailing zeros from float64
func removeTrailingZeros(str string) string {

	// Remove trailing zeros
	str = strings.TrimRight(str, "0")

	// Remove trailing dot
	str = strings.TrimRight(str, ".")

	return str
}
func Print(s interface{}) {
	_, err := strconv.ParseFloat(s.(string), 64)

	if err == nil {
		fmt.Println(removeTrailingZeros(s.(string)))
	} else {
		fmt.Println(s)
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
func ExtractFileName(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	fileName := path.Base(parsedURL.Path)
	return fileName, nil
}

func ParseFloat(s string) float64 {
	if strings.Contains(s, ".") {
		// Pokud obsahuje desetinnou tečku, převedeme na float
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Println("Chyba při převodu na float:", err)

		}
		return f
	} else {
		// Pokud neobsahuje desetinnou tečku, převedeme na int a poté na float
		i, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Chyba při převodu na int:", err)

		}
		f := float64(i)
		return f
	}
}
