package ossindex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/A1bemuth/go-oss-index/types"
)

const default_uri = "https://ossindex.sonatype.org/api/v3/component-report"

var (
	ErrMissingCoordinatesVersion = errors.New("Component coordinates should always specify a version")
	ErrTooManyRequests           = errors.New("Too many requests")
)

type Client struct {
	Uri    string
	client http.Client
}

type ossIndexRequest struct {
	Coordinates []string `json:"coordinates"`
}

func (c *Client) Get(purls []string) ([]types.ComponentReport, error) {
	request, err := makeRequest(purls)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Post(c.getUri(), "application/json", request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return readResponse(resp.Body)
	case 400:
		return nil, ErrMissingCoordinatesVersion
	case 429:
		return nil, ErrTooManyRequests
	default:
		return nil, fmt.Errorf("Unexpected response code: %d %s", resp.StatusCode, resp.Status)
	}
}

func (c *Client) getUri() string {
	if c.Uri != "" {
		return c.Uri
	}
	return default_uri
}

func makeRequest(purls []string) (*bytes.Buffer, error) {
	request := ossIndexRequest{
		Coordinates: purls,
	}
	serialized, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(serialized)
	return buffer, nil
}

func readResponse(reader io.Reader) ([]types.ComponentReport, error) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var result []types.ComponentReport
	err = json.Unmarshal(body, &result)

	return result, err
}
