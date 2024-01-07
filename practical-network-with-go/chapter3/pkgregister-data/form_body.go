package pkgregister

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

type pkgData struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader
}

func createMultiPartMessage(data pkgData) (dataBytes []byte, contentType string, err error) {
	var (
		b  bytes.Buffer
		fw io.Writer
	)

	mw := multipart.NewWriter(&b)
	fw, err = mw.CreateFormField("name")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Name)

	fw, err = mw.CreateFormField("version")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Version)

	fw, err = mw.CreateFormFile("filedata", data.Filename)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(fw, data.Bytes)
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	return b.Bytes(), mw.FormDataContentType(), nil
}
