package core

import "net/url"

// User supplied config
type Config struct {
	URL string
}

type Media struct {
	URL  *url.URL
	Type string
}
