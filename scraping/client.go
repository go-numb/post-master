package scraping

import (
	"context"
	"fmt"
	"net/url"
	"time"

	playwright "github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

// AddFunc Scrap & Save to Target Body, next raw uri and err for return values
type AddFunc func(playwright.Page) (string, error)
type Client struct {
	TermMinits int
	TargetURL  *url.URL

	pw *playwright.Playwright

	f AddFunc
}

func New(termMinites int, baseURI string, setModelsToDatabase ...any) *Client {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatal().Msgf("%v", fmt.Errorf("playwright client has not launch"))
	}

	u, err := url.Parse(baseURI)
	if err != nil || baseURI == "" || u == nil {
		log.Fatal().Msg("nothing base URI")
		return nil
	}

	return &Client{
		TermMinits: termMinites,
		TargetURL:  u,
		pw:         pw,
	}
}

func (p *Client) ChangeTargetURI(uri string) error {
	u, err := url.Parse(uri)
	if err != nil || uri == "" || u == nil {
		return err
	}

	p.TargetURL = u

	return nil
}

func (p *Client) Start(ctx context.Context, f func(playwright.Page) (string, error)) {
	t := time.NewTicker(time.Duration(p.TermMinits) * time.Minute)
	defer t.Stop()

	// スクレイピング実行関数を登録する
	p.AddFunc(f)

ENDED:
	for {
		select {
		case <-t.C:
			log.Info().Msg("goto")
			if err := p.do(); err != nil {
				log.Err(err)
				break ENDED
			}

		case <-ctx.Done():
			log.Err(ctx.Err())
			break ENDED
		}
	}

	log.Info().Msg("ended goroutine")
}

func (p *Client) Close() error {
	return p.pw.Stop()
}

func (p *Client) AddFunc(f func(playwright.Page) (string, error)) {
	p.f = f
}
