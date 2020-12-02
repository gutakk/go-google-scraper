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
)

func createFormFile(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", "text/csv")
	return w.CreatePart(h)
}

func CreateMultipartPayload(filename string) (http.Header, *bytes.Buffer) {
	path := filename
	file, _ := os.Open(path)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := createFormFile(writer, "file", filepath.Base(path))
	_, _ = io.Copy(part, file)
	writer.Close()

	headers := http.Header{}
	headers.Set("Content-Type", writer.FormDataContentType())

	return headers, body
}
