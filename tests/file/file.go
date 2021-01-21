package file

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func createFormFile(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", "text/csv")
	return w.CreatePart(h)
}

func CreateMultipartPayload(filename string) (http.Header, *bytes.Buffer) {
	path := filename
	file, openFileErr := os.Open(path)
	if openFileErr != nil {
		log.Errorf("Cannot open file: %s", openFileErr)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, createFormFileErr := createFormFile(writer, "file", filepath.Base(path))
	if createFormFileErr != nil {
		log.Errorf("Cannot create form file: %s", createFormFileErr)
	}

	_, copyErr := io.Copy(part, file)
	if copyErr != nil {
		log.Errorf("Cannot copy file part: %s", copyErr)
	}
	writer.Close()

	headers := http.Header{}
	headers.Set("Content-Type", writer.FormDataContentType())

	return headers, body
}
