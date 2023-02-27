package unione

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

const host = "https://go2.unisender.ru"
const methodSend = "/ru/transactional/api/v1/email/send.json"

type Client struct {
	host   string
	client http.Client
}

type transport struct {
	http.RoundTripper
	apikey string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-API-KEY", t.apikey)
	return t.RoundTripper.RoundTrip(req)
}

func NewClient(apikey string) Client {
	return Client{
		host: host,
		client: http.Client{
			Transport: &transport{
				RoundTripper: http.DefaultTransport,
				apikey:       apikey,
			},
		},
	}
}

func (c Client) Send(m Message) error {
	path, err := url.JoinPath(c.host, methodSend)
	if err != nil {
		return err
	}

	body, err := json.Marshal(request{Message: m})
	if err != nil {
		return err
	}

	resp, err := c.client.Post(path, "application/json", bytes.NewReader(body))
	{
		b, _ := io.ReadAll(resp.Body)
		log.Print("unione response: ", resp.StatusCode, resp.Status, string(b))
	}

	return err
}
