package shortened_url

import (
	"net/url"
)

const UrlMaxLength int = 2048
const SlugMaxLength int = 50
const SlugMinLength int = 4

type shortenedUrl struct {
	Slug string  `dynamodbav:"slug"`
	Url  url.URL `dynamodbav:"url"`
}

func New(rawUrl string, slug string) (*shortenedUrl, error) {
	if len(slug) > SlugMaxLength {
		return nil, ErrInvalidSlug
	}

	if len(slug) < SlugMinLength {
		return nil, ErrInvalidSlug
	}

	if len(rawUrl) > UrlMaxLength {
		return nil, ErrInvalidUrl
	}

	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, ErrInvalidUrl
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, ErrInvalidUrl
	}

	if url.Host == "" {
		return nil, ErrInvalidUrl
	}

	return &shortenedUrl{
		Slug: slug,
		Url:  *url,
	}, nil
}
