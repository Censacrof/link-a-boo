package shortened_url

import (
	"fmt"
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