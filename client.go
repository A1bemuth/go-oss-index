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
	ErrUnauthorized              = errors.New("Unauthorized")
	ErrTooManyRequests           = errors.New("Too many requests")
)

type Client struct {
	Uri      string
	User     string
	Password string
	client   http.Client
}

type ossIndexRequest struct {
	Coordinates []string `json:"coordinates"`
}

func (c *Client) Get(purls []string) ([]types.ComponentReport, error) {
	request, err := c.makeRequest(purls)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return readResponse(resp.Body)
	case 400:
		return nil, ErrMissingCoordinatesVersion
	case 401:
		return nil, ErrUnauthorized
	case 429:
		return nil, ErrTooManyRequests
	default:
		return nil, fmt.Errorf("Unexpected response code: %s", resp.Status)
	}
}

func (c *Client) getUri() string {
	if c.Uri != "" {
		return c.Uri
	}
	return default_uri
}

func (c *Client) makeRequest(purls []string) (*http.Request, error) {
	body, err := makeRequestBody(purls)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.getUri(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.ossindex.component-report-request.v1+json")
	if c.User != "" && c.Password != "" {
		req.SetBasicAuth(c.User, c.Password)
	}

	return req, nil
}

func makeRequestBody(purls []string) (*bytes.Buffer, error) {
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
