package smtp

import (
	"errors"
	"fmt"
	"net/url"
)

var errInvalidScheme = errors.New("invalid smtp server scheme")

func ErrInvalidScheme(scheme string) error {
	return fmt.Errorf("%w `%s`, must be one of `plain`, `tls` or `starttls`", errInvalidScheme, scheme)
}

type URI struct {
	Scheme  string
	Address string
}

func NewURI(uri string) (*URI, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("error parsing `%s` uri: %w", uri, err)
	}

	if url.Scheme != "plain" && url.Scheme != "tls" && url.Scheme != "starttls" {
		return nil, ErrInvalidScheme(url.Scheme)
	}

	return &URI{
		Scheme:  url.Scheme,
		Address: url.Host,
	}, nil
}

func (u *URI) String() string {
	return fmt.Sprintf("%s://%s", u.Scheme, u.Address)
}
