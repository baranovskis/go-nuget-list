package nuget

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	endpoints []string
	client    *http.Client
}

// NewNugetClient returns a new instance of Parser.
func NewNugetClient() *Client {
	return &Client{client: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) Search(sourceUrl string, id string) (*ResponseQuery, error) {
	r, err := c.client.Get(sourceUrl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	if mediaType == "application/xml" {
		return c.QueryApiV2(sourceUrl, id)
	}

	if mediaType == "application/json" {
		response := &ResponseResources{}
		err = json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			return nil, err
		}

		for _, resource := range response.Resources {
			if resource.Type == "SearchQueryService" {
				return c.QueryApiV3(resource.Id, id)
			}
		}

		return nil, errors.New("search query service not found")
	}

	return nil, errors.New(fmt.Sprintf("unknown media type: %s", mediaType))
}

func (c *Client) QueryApiV2(sourceUrl string, id string) (*ResponseQuery, error) {
	u, err := url.Parse(sourceUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("id", "'"+id+"'")

	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, "FindPackagesById()")

	r, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	response := &ResponseQuery{}
	err = xml.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) QueryApiV3(sourceUrl string, id string) (*ResponseQuery, error) {
	u, err := url.Parse(sourceUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", id)
	q.Set("prerelease", "false")
	q.Set("semVerLevel", "2.0.0")
	u.RawQuery = q.Encode()

	r, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	response := &ResponseQuery{}
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
