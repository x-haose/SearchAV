package source

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"searchav/internal/config"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

const (
	apiPath = "/provide/vod/at/json"
)

// Client is the video source API client
type Client struct {
	http   *resty.Client
	logger *zerolog.Logger
}

// NewClient creates a new source client
func NewClient(cfg *config.Config, logger *zerolog.Logger) *Client {
	client := resty.New().
		SetTimeout(cfg.Source.Timeout).
		SetRetryCount(cfg.Source.Retry).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	return &Client{
		http:   client,
		logger: logger,
	}
}

// Search searches videos from a source
func (c *Client) Search(ctx context.Context, src config.SourceItem, keyword string) ([]RawVideo, error) {
	reqURL := fmt.Sprintf("%s%s?ac=videolist&wd=%s", src.URL, apiPath, url.QueryEscape(keyword))

	c.logger.Info().Str("url", reqURL).Str("source", src.Code).Msg("search request")

	var resp MacCMSResponse
	httpResp, err := c.http.R().
		SetContext(ctx).
		ForceContentType("application/json").
		SetResult(&resp).
		Get(reqURL)

	if err != nil {
		c.logger.Error().Err(err).Str("source", src.Code).Msg("search request failed")
		return nil, err
	}

	c.logger.Info().
		Str("source", src.Code).
		Int("status", httpResp.StatusCode()).
		Int("code", resp.Code).
		Str("msg", resp.Msg).
		Int("total", resp.Total).
		Int("count", len(resp.List)).
		Msg("search response")

	// Inject source info
	for i := range resp.List {
		resp.List[i].SourceCode = src.Code
		resp.List[i].SourceName = src.Name
	}

	return resp.List, nil
}

// GetDetail gets video detail from a source
func (c *Client) GetDetail(ctx context.Context, src config.SourceItem, vodID int) (*RawVideo, error) {
	reqURL := fmt.Sprintf("%s%s?ac=videolist&ids=%d", src.URL, apiPath, vodID)

	c.logger.Debug().Str("url", reqURL).Str("source", src.Code).Msg("detail request")

	var resp MacCMSResponse
	_, err := c.http.R().
		SetContext(ctx).
		ForceContentType("application/json").
		SetResult(&resp).
		Get(reqURL)

	if err != nil {
		return nil, err
	}

	if len(resp.List) == 0 {
		return nil, fmt.Errorf("video not found")
	}

	return &resp.List[0], nil
}

// SearchWithTimeout searches with a timeout
func (c *Client) SearchWithTimeout(ctx context.Context, src config.SourceItem, keyword string, timeout time.Duration) ([]RawVideo, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Search(ctx, src, keyword)
}
