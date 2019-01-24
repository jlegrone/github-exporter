package exporter

import (
	"time"

	"github.com/infinityworks/github-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shurcooL/githubv4"
)

// Exporter is used to store Metrics data and embeds the config struct.
// This is done so that the relevant functions have easy access to the
// user defined runtime configuration when the Collect method is called.
type Exporter struct {
	APIMetrics map[string]*prometheus.Desc
	config.Config
}

type githubData struct {
	RateLimit     RateLimit
	Organizations []Organization
	Users         []User
	Repositories  []Repository
}

type TotalCount struct {
	TotalCount float64
}

type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

type Ref struct {
	Name   string
	Target struct {
		CommitURL      string
		OID            string
		AbbreviatedOID string
		Commit         struct {
			History TotalCount
		} `graphql:"... on Commit"`
	}
}

type LanguageConnection struct {
	TotalCount int
	TotalSize  int
	Nodes      []struct {
		Name string
	}
	Edges []struct {
		Size int
	}
}

type Repository struct {
	Owner struct {
		Login string
	}
	Name                  string
	NameWithOwner         string
	ID                    string
	Description           string
	URL                   string
	HomepageURL           string
	DefaultBranchRef      Ref
	DiskUsage             int
	UpdatedAt             githubv4.DateTime
	ForkCount             int
	Stargazers            TotalCount
	Watchers              TotalCount
	AssignableUsers       TotalCount
	MentionableUsers      TotalCount
	Releases              TotalCount
	Deployments           TotalCount
	BranchProtectionRules TotalCount
	CommitComments        TotalCount
	Labels                TotalCount
	RepositoryTopics      TotalCount
	OpenIssues            TotalCount `graphql:"openIssues: issues(states: OPEN)"`
	ClosedIssues          TotalCount `graphql:"closedIssues: issues(states: CLOSED)"`
	OpenPullRequests      TotalCount `graphql:"openPullRequests: pullRequests(states: OPEN)"`
	ClosedPullRequests    TotalCount `graphql:"closedPullRequests: pullRequests(states: CLOSED)"`
	MergedPullRequests    TotalCount `graphql:"mergedPullRequests: pullRequests(states: MERGED)"`
	OpenMilestones        TotalCount `graphql:"openMilestones: milestones(states: OPEN)"`
	ClosedMilestones      TotalCount `graphql:"closedMilestones: milestones(states: CLOSED)"`
	IsArchived            bool
	IsPrivate             bool
	IsFork                bool
	IsLocked              bool
	IsMirror              bool
	HasIssuesEnabled      bool
	HasWikiEnabled        bool
	LicenseInfo           struct {
		Name string
		Key  string
	}
	PrimaryLanguage struct {
		Name string
	}
	Languages LanguageConnection `graphql:"languages(first: $languageCount, orderBy: {field: SIZE, direction: DESC})"`
}

type RepositoryConnection struct {
	TotalCount int
	PageInfo   PageInfo
	Nodes      []Repository
}

type User struct {
	Name                    string
	Login                   string
	ID                      string
	ContributionsCollection struct {
		TotalIssueContributions                      int
		TotalCommitContributions                     int
		TotalRepositoryContributions                 int
		TotalPullRequestContributions                int
		TotalPullRequestReviewContributions          int
		TotalRepositoriesWithContributedIssues       int
		TotalRepositoriesWithContributedCommits      int
		TotalRepositoriesWithContributedPullRequests int
	}
	CommitComments      TotalCount
	IssueComments       TotalCount
	GistComments        TotalCount
	Gists               TotalCount
	OpenIssues          TotalCount `graphql:"openIssues: issues(states: OPEN)"`
	ClosedIssues        TotalCount `graphql:"closedIssues: issues(states: CLOSED)"`
	OpenPullRequests    TotalCount `graphql:"openPullRequests: pullRequests(states: OPEN)"`
	ClosedPullRequests  TotalCount `graphql:"closedPullRequests: pullRequests(states: CLOSED)"`
	MergedPullRequests  TotalCount `graphql:"mergedPullRequests: pullRequests(states: MERGED)"`
	Followers           TotalCount
	Following           TotalCount
	StarredRepositories TotalCount
	Repositories        RepositoryConnection `graphql:"repositories(first: $repoCount, orderBy: {field: UPDATED_AT, direction: DESC})"`
}

type Organization struct {
	Name            string
	Login           string
	ID              string
	MembersWithRole TotalCount
	Teams           TotalCount
	Projects        TotalCount
	Repositories    RepositoryConnection `graphql:"repositories(first: $repoCount)"`
}

type RateLimit struct {
	Cost      float64
	Limit     float64
	NodeCount float64
	Remaining float64 `graphql:"remaining"`
	ResetAt   time.Time
}

type UserQuery struct {
	RateLimit RateLimit
	User      User `graphql:"user(login: $login)"`
}

type OrganizationQuery struct {
	RateLimit    RateLimit
	Organization Organization `graphql:"organization(login: $login)"`
}

type RepositoryQuery struct {
	RateLimit  RateLimit
	Repository Repository `graphql:"repository(owner: $owner, name: $name)"`
}
