package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Censacrof/link-a-boo/lambda/shorten/cmd/pkg"
	"github.com/aws/aws-lambda-go/lambda"
)

func getHandler() (interface{}, error) {
	handlerName := os.Getenv("_HANDLER")

	switch handlerName {
	case "shorten":
		return shorten.HandleShortenRequest, nil

	default:
		return nil, errors.New(fmt.Sprintf("Invalid handler %s", handlerName))
	}
}

func main() {
	handler, err := getHandler()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	lambda.Start(handler)
}
