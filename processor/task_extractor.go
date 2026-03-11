package processor

import (
	"bytes"
	"io"

	"github.com/ledongthuc/pdf"
)

func ExtractText(pdfReader io.Reader) (string, int, error) {
	data, err := io.ReadAll(pdfReader)
	if err != nil {
		return "", 0, err
	}

	pdfFile, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", 0, err
	}

	pageCount := pdfFile.NumPage()

	textReader, err := pdfFile.GetPlainText()
	if err != nil {
		return "", 0, err
	}

	textBytes, err := io.ReadAll(textReader)
	if err != nil {
		return "", 0, err
	}

	return string(textBytes), pageCount, nil
}
