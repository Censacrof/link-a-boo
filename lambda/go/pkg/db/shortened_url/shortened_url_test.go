package shortened_url

import (
	"fmt"
	"strings"
	"testing"
)

func Test_NewReturnsErrorWhenUrlIsNotValid(t *testing.T) {
	testCases := []string{
		"anInvalidUrl",
		"file:///an/invalid/url.txt",
		"/an/invalid/url.txt",
		"data:text/html,%3Ch1%3EHello%2C%20World%21%3C%2Fh1%3E",
		"http:",
	}

	for _, invalidUrl := range testCases {
		t.Run(fmt.Sprintf("test %s", invalidUrl), func(t *testing.T) {
			_, err := New(invalidUrl, "aSlug")

			if err == nil {
				t.Fatalf("Expected New to return error")
			}
		})
	}
}

func Test_NewDoesntReturnErrorWhenUrlIsValid(t *testing.T) {
	testCases := []string{
		"http://www.mywebsite.com/path/to/stuff#someSection",
		"https://www.mywebsite.com/path/to/stuff#someSection",
	}

	for _, invalidUrl := range testCases {
		t.Run(fmt.Sprintf("test %s", invalidUrl), func(t *testing.T) {
			slug := "aSlug"
			su, err := New(invalidUrl, slug)

			if err != nil {
				t.Fatalf("Expected New not to return error: %v", err)
			}

			if su == nil {
				t.Fatalf("Expected New to return a shortenedUrl but returned nil instead")
			}

			if su.Slug != "aSlug" {
				t.Fatalf("Expected Slug to be '%s', but got '%s' instead", slug, su.Slug)
			}
		})
	}
}

func Test_NewReturnsErrorWhenUrlExceedsMaxLength(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("http://my.website.com/")

	for sb.Len() < UrlMaxLength {
		sb.WriteString("a")
	}

	validUrl := sb.String()

	_, err := New(validUrl, "aSlug")
	if err != nil {
		t.Fatalf("Expected New not to return error when url length is %d characters long. Error returned is: %v", len(validUrl), err)
	}

	urlThatsTooLong := validUrl + "a"

	_, err = New(urlThatsTooLong, "aSlug")
	if err == nil {
		t.Fatalf("Expected New to return error when url length is %d characters long", len(urlThatsTooLong))
	}
}

func Test_NewReturnsErrorWhenSlugExceedsMaxLength(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("http://my.website.com/")

	for sb.Len() < SlugMaxLength {
		sb.WriteString("a")
	}

	validSlug := sb.String()

	_, err := New("http://my.website.com/", validSlug)
	if err != nil {
		t.Fatalf("Expected New not to return error when slug length is %d characters long. Error returned is: %v", len(validSlug), err)
	}

	slugThatsTooLong := validSlug + "a"

	_, err = New("http://my.website.com/", slugThatsTooLong)
	if err == nil {
		t.Fatalf("Expected New to return error when slug length is %d characters long", len(slugThatsTooLong))
	}
}
