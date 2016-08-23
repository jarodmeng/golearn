package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	// get cache file's path
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	// reads the cached token, if it exists
	tok, err := tokenFromFile(cacheFile)
	// if there's no existing cached token, get one and save it
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	// Client returns an HTTP client using the provided token
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// returns a URL to OAuth 2.0 provider's consent page that asks for
	// permissions for the required scopes explicitly.
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	// scan the standard input to populate code
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	// Exchange converts an authorization code into a token
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	// get the current user. It returns a pointer to a User instance
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	// create a .credentials directory in current user's home directory, if
	// necessary
	os.MkdirAll(tokenCacheDir, 0700)
	// return the credential file's path
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	// open the file for reading. It returns a pointer to a File instance
	f, err := os.Open(file)
	if err != nil {
		// if the file cannot be opened, return nil and the error so the receiving
		// side can act accordingly
		return nil, err
	}
	// create a pointer to an empty Token instance
	t := &oauth2.Token{}
	// NewDecoder returns a new decoder that reads from f.
	// Decode reads the JSON-encoded value and stores it in t.
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	// encode the token and save it
	json.NewEncoder(f).Encode(token)
}

// create gmail client
func createGmailClient(file string) *http.Client {
	// create an empty context
	ctx := context.Background()

	// read the client secret json file which contains the client id and secret
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json
	// ConfigFromJSON uses a Google Developers Console client_credentials.json
	// file to construct a config. It returns a pointer to an oauth2.Config
	// instance
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	// get an HTTP client using the context and config/token
	client := getClient(ctx, config)

	return client
}

func main() {
	client := createGmailClient("client_secret.json")

	// create a pointer to a Service instance using the client
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	user := "me"
	// execute the api call to get results in ListLabelsResponse format
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels. %v", err)
	}
	if len(r.Labels) > 0 {
		fmt.Print("Labels:\n")
		for _, l := range r.Labels {
			fmt.Printf("- %s\n", l.Name)
		}
	} else {
		fmt.Print("No labels found.")
	}
}
