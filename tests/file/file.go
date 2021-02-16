package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"

	"github.com/gutakk/go-google-scraper/helpers/log"
)

func createFormFile(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", "text/csv")
	return w.CreatePart(h)
}

func CreateMultipartPayload(filename string) (http.Header, *bytes.Buffer) {
	path := filename
	file, err := os.Open(path)
	if err != nil {
		log.Error("Failed to open file: ", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := createFormFile(writer, "file", filepath.Base(path))
	if err != nil {
		log.Error("Failed to create part from file: ", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Error("Failed to copy file: ", err)
	}
	writer.Close()

	headers := http.Header{}
	headers.Set("Content-Type", writer.FormDataContentType())

	return headers, body
}
