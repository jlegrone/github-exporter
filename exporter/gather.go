package exporter

import (
	"context"
	"sync"

	conf "github.com/infinityworks/github-exporter/config"
	"github.com/shurcooL/githubv4"
)

func (e *Exporter) gatherData() (*githubData, error) {
	var (
		mux sync.Mutex
		wg  sync.WaitGroup
	)

	result := &githubData{
		RateLimit{},
		[]Organization{},
		[]User{},
		[]Repository{},
	}

	for _, login := range e.Config.Users {
		wg.Add(1)
		go func(login string) {
			defer wg.Done()
			var query UserQuery
			e.Client.Query(context.Background(), &query, map[string]interface{}{
				"login":         githubv4.String(login),
				"repoCount":     githubv4.Int(10),
				"languageCount": githubv4.Int(5),
			})
			if query.User.ID != "" {
				mux.Lock()
				result.Users = append(result.Users, query.User)
				result.RateLimit = query.RateLimit
				mux.Unlock()
			}
		}(login)
	}

	for _, login := range e.Config.Organisations {
		wg.Add(1)
		go func(login string) {
			defer wg.Done()
			var query OrganizationQuery
			e.Client.Query(context.Background(), &query, map[string]interface{}{
				"login":         githubv4.String(login),
				"repoCount":     githubv4.Int(10),
				"languageCount": githubv4.Int(5),
			})
			if query.Organization.ID != "" {
				mux.Lock()
				result.Organizations = append(result.Organizations, query.Organization)
				result.RateLimit = query.RateLimit
				mux.Unlock()
			}
		}(login)
	}

	for _, repo := range e.Config.Repositories {
		wg.Add(1)
		go func(repo conf.Repository) {
			defer wg.Done()
			var query RepositoryQuery
			e.Client.Query(context.Background(), &query, map[string]interface{}{
				"owner":         githubv4.String(repo.Owner),
				"name":          githubv4.String(repo.Name),
				"languageCount": githubv4.Int(5),
			})
			if query.Repository.ID != "" {
				mux.Lock()
				result.Repositories = append(result.Repositories, query.Repository)
				result.RateLimit = query.RateLimit
				mux.Unlock()
			}
		}(repo)
	}

	wg.Wait()

	return result, nil
}
