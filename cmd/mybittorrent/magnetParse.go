package main

import (
	"fmt"
	"net/url"
	"strings"
)

type MagnetLink struct {
	infoHash   string
	fileName   string
	trackerURL string
}

func parseMagentLink(magnet string) (MagnetLink, error) {
	var magnetLink MagnetLink
	pos := strings.Index(magnet, "xt=urn:btih:")
	if pos == -1 {
		return MagnetLink{}, fmt.Errorf("invalid magnet link : xt=urn:btih: field required")
	}
	magnetLink.infoHash = magnet[pos+12 : pos+52]

	pos = strings.Index(magnet, "tr")
	if pos == -1 {
		return MagnetLink{}, fmt.Errorf("invalid magnet link : tr= field required")
	}
	magnetLink.trackerURL, _ = url.QueryUnescape(magnet[pos+3:])

	return magnetLink, nil
}
