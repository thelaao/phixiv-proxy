package pixiv

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type PixivClient struct {
	HttpClient *http.Client
	UserAgent  string
	Cookie     string
	UserId     string
}

func NewPixivClient(httpClient *http.Client, userAgent string, cookie string) *PixivClient {
	client := &PixivClient{
		HttpClient: httpClient,
		UserAgent:  userAgent,
		Cookie:     cookie,
	}
	if len(cookie) > 0 {
		parts := strings.SplitN(cookie, "_", 1)
		if len(parts) == 2 {
			client.UserId = parts[0]
		}
	}
	return client
}

func (client *PixivClient) CallApi(ctx context.Context, url string, output any) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", client.UserAgent)
	if len(client.Cookie) > 0 {
		req.Header.Set("Cookie", "PHPSESSID="+client.Cookie)
	}
	if len(client.UserId) > 0 {
		req.Header.Set("X-User-Id", client.UserId)
	}
	res, err := client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return &PixivRequestError{StatusCode: res.StatusCode}
	}
	err = json.NewDecoder(res.Body).Decode(output)
	return err
}

func (client *PixivClient) Download(ctx context.Context, url string) (contentType string, body []byte, err error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", client.UserAgent)
	req.Header.Set("Referer", "https://www.pixiv.net/")
	res, err := client.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = &PixivRequestError{StatusCode: res.StatusCode}
		return
	}
	contentType = res.Header.Get("Content-Type")
	body, err = io.ReadAll(res.Body)
	return
}
