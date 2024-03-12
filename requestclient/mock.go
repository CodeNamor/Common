package requestclient

import (
	"go/build"
	"io"
	"net/http"
	"os"

	"github.com/CodeNamor/Common/logging"
	"github.com/CodeNamor/Common/path"
)

// MockRequestFileClient holds keys that identify a locally cached file
type MockRequestFileClient struct {
	FilePath string // relative to GOPATH
}

// Do retrieves local content to mocks an http request
func (s *MockRequestFileClient) Do(req *http.Request) (*http.Response, error) {

	readCloser, err := FileReadCloser(s.FilePath)
	if err != nil {
		return nil, err
	}
	tempResponse := http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       readCloser,
	}
	return &tempResponse, nil
}

// FileReadCloser is a helper function that will search for the fileName.json file in the testdata subdirectory
func FileReadCloser(filePath string) (io.ReadCloser, error) {
	absPath := path.Resolve(build.Default.GOPATH, filePath)
	logging.Trace("MockRequestFileClient", absPath)
	fileHandler, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	return fileHandler, nil
}
