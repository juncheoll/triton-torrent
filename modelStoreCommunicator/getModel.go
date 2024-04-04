package modelStoreCommunicator

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/juncheoll/triton-torrent/setting"
)

/*
type modelStoreResponse struct {
	File []byte `json:"file"`
}
*/

// The ProgressReader is a Reader that displays progress status.
type ProgressReader struct {
	io.Reader
	totalRead int64
}

// Read is an implementation of the io.Reader interface for ProgressReader.
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.totalRead += int64(n)

	log.Printf("Downloaded %.2f MB\n", float64(pr.totalRead)/1024/1024)

	return n, err
}

/* downloading the model from the Model Store. */
func GetModel(provider string, model string, version string) ([]byte, error) {
	log.Println("Starting model download.")
	url := fmt.Sprintf("http://%s/getModel?id=%s&model=%s&version=%s", setting.ModelStoreUrl, provider, model, version)

	// Sending a download request to the model store.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Use progressReader.
	progressReader := &ProgressReader{Reader: resp.Body}
	body, err := io.ReadAll(progressReader)
	if err != nil {
		return nil, err
	}

	log.Printf("Download completed. Total size: %.2f MB\n", float64(len(body))/1024/1024)

	return body, nil
}
