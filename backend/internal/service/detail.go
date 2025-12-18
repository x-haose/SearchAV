package service

import (
	"context"
	"fmt"
	"strings"

	"searchav/internal/config"
	"searchav/internal/model"
	"searchav/internal/source"

	"github.com/rs/zerolog"
)

// DetailService handles video detail retrieval
type DetailService struct {
	config *config.Config
	client *source.Client
	logger *zerolog.Logger
}

// NewDetailService creates a new detail service
func NewDetailService(cfg *config.Config, client *source.Client, logger *zerolog.Logger) *DetailService {
	return &DetailService{
		config: cfg,
		client: client,
		logger: logger,
	}
}

// GetDetail retrieves video details from a source
func (s *DetailService) GetDetail(ctx context.Context, sourceCode string, vodID int) (*model.VideoDetail, error) {
	src, ok := s.config.GetSourceByCode(sourceCode)
	if !ok {
		return nil, fmt.Errorf("source not found: %s", sourceCode)
	}

	raw, err := s.client.GetDetail(ctx, *src, vodID)
	if err != nil {
		return nil, err
	}

	// Parse play URLs into episodes
	episodes := s.parseEpisodes(raw.VodPlayURL)

	return &model.VideoDetail{
		VodName:     raw.VodName,
		VodPic:      raw.VodPic,
		VodContent:  raw.VodContent,
		VodYear:     raw.VodYear,
		VodArea:     raw.VodArea,
		VodDirector: raw.VodDirector,
		VodActor:    raw.VodActor,
		Episodes:    episodes,
	}, nil
}

// parseEpisodes parses play URL string into episode URLs
// Format: name1$url1#name2$url2$$$name1$url1#name2$url2
// $$$ separates multiple sources, # separates episodes, $ separates name and URL
func (s *DetailService) parseEpisodes(playURL string) []string {
	if playURL == "" {
		return nil
	}

	// Split by multiple sources
	sources := strings.Split(playURL, "$$$")
	if len(sources) == 0 {
		return nil
	}

	// Prefer sources containing m3u8
	selectedSource := sources[0]
	for _, src := range sources {
		if strings.Contains(src, ".m3u8") || strings.Contains(src, "m3u8") {
			selectedSource = src
			break
		}
	}

	// Split episodes
	epList := strings.Split(selectedSource, "#")

	episodes := make([]string, 0, len(epList))
	for _, ep := range epList {
		parts := strings.Split(ep, "$")
		if len(parts) >= 2 {
			url := parts[1]
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				episodes = append(episodes, url)
			}
		}
	}

	return episodes
}
