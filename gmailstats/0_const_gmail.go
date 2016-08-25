// Token cache file and scope information for GMail API calls
package main

import gmail "google.golang.org/api/gmail/v1"

const (
	defaultGmailTokenFile = "gmail_token.json"
	defaultGmailScope     = gmail.GmailReadonlyScope
	defaultGmailUser      = "me"
)
