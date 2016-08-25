package main

import (
	"github.com/jarodmeng/googleauth"
	gmail "google.golang.org/api/gmail/v1"
)

func createGmailService(secretFile string, tokenFile string, scope string) (*gmail.Service, error) {
	gmailClient, err := googleauth.CreateClientFromFile(secretFile, tokenFile, scope)
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.New(gmailClient)
	if err != nil {
		return nil, err
	}

	return gmailService, nil
}
