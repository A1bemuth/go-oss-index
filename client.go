package ossindex

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/A1bemuth/go-oss-index/types"
)

const DEFAULT_URI = "https://ossindex.sonatype.org/api/v3/component-report"

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

	reports, err := readResponse(resp.Body)

	return reports, err
}

func (c *Client) getUri() string {
	if c.Uri != "" {
		return c.Uri
	}
	return DEFAULT_URI
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
