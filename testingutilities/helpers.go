package testingutilities

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"
)

type MultipartBody struct {
	Reader io.Reader
	Writer *multipart.Writer
}

type FileInfo struct {
	Content  bytes.Buffer
	MimeType string
}

func writeToMultipart(w *multipart.Writer, data map[string]string) error {
	for key, val := range data {
		x, err := w.CreateFormField(key)
		if err != nil {
			return err
		}

		if _, err := x.Write([]byte(val)); err != nil {
			return err
		}
	}

	return nil
}

func GenerateMultipart(t *testing.T, files map[string]FileInfo, content ...map[string]string) (MultipartBody, error) {
	t.Helper()

	buffer := new(bytes.Buffer)
	w := multipart.NewWriter(buffer)
	defer w.Close()

	if len(content) > 0 {
		writeToMultipart(w, content[0])
	}

	for filename, info := range files {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
		header.Set("Content-Type", info.MimeType)

		part, err := w.CreatePart(header)
		if err != nil {
			return MultipartBody{}, err
		}

		if _, err = part.Write(info.Content.Bytes()); err != nil {
			return MultipartBody{}, err
		}
	}

	return MultipartBody{Reader: buffer, Writer: w}, nil
}

func GenerateBody(t *testing.T, content string) (io.Reader, error) {
	t.Helper()

	buffer := new(bytes.Buffer)
	if _, err := buffer.WriteString(content); err != nil {
		return buffer, err
	}

	return buffer, nil
}
