package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	//"golang.org/x/text/encoding/charmap"
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
