package config

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"context"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	cfg "github.com/infinityworks/go-common/config"
)

// Config struct holds all of the runtime confgiguration for the application
type Config struct {
	*cfg.BaseConfig
	Repositories  []Repository
	Organisations []string
	Users         []string
	APIToken      string
	Client        *githubv4.Client
}

type Repository struct {
	Owner string
	Name  string
}

const defaultAPIURL = "https://api.github.com"

// Init populates the Config struct based on environmental runtime configuration
func Init() Config {

	ac := cfg.Init()
	url := cfg.GetEnv("API_URL", defaultAPIURL)
	repoStrings := splitEnvVar("REPOS")
	orgs := splitEnvVar("ORGS")
	users := splitEnvVar("USERS")
	tokenEnv := os.Getenv("GITHUB_TOKEN")
	tokenFile := os.Getenv("GITHUB_TOKEN_FILE")
	token, err := getAuth(tokenEnv, tokenFile)

	if err != nil {
		log.Errorf("Error initialising Configuration, Error: %v", err)
	}

	repos := []Repository{}

	for _, repo := range repoStrings {
		parts := strings.Split(repo, "/")
		repos = append(repos, Repository{
			Owner: parts[0],
			Name:  parts[1],
		})
	}

	var (
		httpClient *http.Client
		client     *githubv4.Client
	)
	if token != "" {
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient = oauth2.NewClient(context.Background(), tokenSource)
	} else {
		httpClient = &http.Client{}
	}
	httpClient.Timeout = time.Second * 15
	if url != defaultAPIURL {
		client = githubv4.NewEnterpriseClient(url, httpClient)
	} else {
		client = githubv4.NewClient(httpClient)
	}

	appConfig := Config{
		&ac,
		repos,
		orgs,
		users,
		token,
		client,
	}

	return appConfig
}

func splitEnvVar(envVar string) []string {
	raw := os.Getenv(envVar)
	var values []string
	for _, val := range strings.Split(raw, ",") {
		if val != "" {
			values = append(values, strings.TrimSpace(val))
		}
	}

	return values
}

// getAuth returns oauth2 token as string for usage in http.request
func getAuth(token string, tokenFile string) (string, error) {

	if token != "" {
		return token, nil
	} else if tokenFile != "" {
		b, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), err

	}

	return "", nil
}
