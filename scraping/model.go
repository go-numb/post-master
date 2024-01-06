package scraping

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (p *Client) do() error {
	if p.pw == nil {
		log.Fatal().Msg("pw is nil")
	}

	browser, err := p.pw.Chromium.Launch()
	if err != nil {
		log.Fatal().Err(fmt.Errorf("chromium browser has not launch"))
		return nil
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatal().Msgf("could not create page: %v", err)
	}
	_, err = page.Goto(p.TargetURL.String())
	if err != nil {
		log.Fatal().Msgf("could not goto: %v", err)
	}

	// s, err := res.Text()
	// if err != nil {
	// 	log.Fatal().Err(err)
	// }
	// log.Debug().Msgf("%#v\n", s)

	// Scrap & Save
	nextRawURI, err := p.f(page)
	if err != nil {
		return err
	}
	if err := p.ChangeTargetURI(nextRawURI); err != nil {
		return err
	}

	return nil
}
