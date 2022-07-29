package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

/* -------------------- Oauth2 Token -------------------- */

type tokenSource struct {
	AccessToken string
}

// Token creates and returns an Oauth2 token
func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

/* -------------------- Client -------------------- */

// NewClient creates and returns a new DigitalOcean API client
func NewClient(apiKey string) *godo.Client {
	tokenSource := &tokenSource{
		AccessToken: apiKey,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	return godo.NewClient(oauthClient)
}
