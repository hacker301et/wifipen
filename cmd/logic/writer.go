package logic

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type myNewWriter struct {
	buf     bytes.Buffer
	counter int
}

func (m *myNewWriter) Write(p []byte) (n int, err error) {
	if len(p)+m.buf.Len() > m.buf.Cap() {
		m.Flush()
	}
	return m.buf.WriteString(string(p))
}

func (m *myNewWriter) Flush() {
	file, err := os.OpenFile("store/wifi.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)

	}
	if m.counter == 0 {
		m.buf.Reset()
		m.counter = 1000
		return
	}
	fmt.Print(m.buf.String())
	file.Write([]byte(removeNonPrintableCharacters(m.buf.String())))
	m.buf.Reset()
	m.counter = m.counter - 1
}

func removeNonPrintableCharacters(s string) string {
	regex := regexp.MustCompile("[^[:print:]\n\t]")

	cleanedString := regex.ReplaceAllString(s, "")
	cleanedString = strings.ReplaceAll(cleanedString, "[0K", "")
	cleanedString = strings.ReplaceAll(cleanedString, "[1B", "")
	cleanedString = strings.TrimSpace(cleanedString)
	pattern := `(?mi)^.*Elapsed.*$`
	regex = regexp.MustCompile(pattern)
	cleanedString = regex.ReplaceAllString(cleanedString, "")
	return cleanedString
}
