package shortened_url

import "errors"

var ErrSlugAlreadyExists = errors.New("Provided slug already exists")
var ErrInvalidSlug = errors.New("Slug is not valid")
var ErrInvalidUrl = errors.New("Url is not valid")
