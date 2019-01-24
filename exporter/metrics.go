package exporter

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var repoLabels = []string{"repo", "user", "private", "fork", "archived", "license", "language"}

// AddMetrics - Add's all of the metrics to a map of strings, returns the map.
func AddMetrics() map[string]*prometheus.Desc {

	APIMetrics := make(map[string]*prometheus.Desc)

	APIMetrics["Stars"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "stars"),
		"Total number of Stars for given repository",
		repoLabels, nil,
	)
	APIMetrics["OpenIssues"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "open_issues"),
		"Total number of open issues for given repository",
		repoLabels, nil,
	)
	APIMetrics["Watchers"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "watchers"),
		"Total number of watchers/subscribers for given repository",
		repoLabels, nil,
	)
	APIMetrics["Forks"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "forks"),
		"Total number of forks for given repository",
		repoLabels, nil,
	)
	APIMetrics["Size"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "size_kb"),
		"Size in KB for given repository",
		repoLabels, nil,
	)
	APIMetrics["Limit"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "rate", "limit"),
		"Number of API queries allowed in a 60 minute window",
		[]string{}, nil,
	)
	APIMetrics["Remaining"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "rate", "remaining"),
		"Number of API queries remaining in the current window",
		[]string{}, nil,
	)
	APIMetrics["Reset"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "rate", "reset"),
		"The time at which the current rate limit window resets in UTC epoch seconds",
		[]string{}, nil,
	)

	return APIMetrics
}

// processMetrics - processes the response data and sets the metrics using it as a source
func (e *Exporter) processMetrics(data *githubData, ch chan<- prometheus.Metric) error {

	// TODO: Dedupe repositories

	// APIMetrics - range through the data
	for _, o := range data.Organizations {
		e.processRepositories(o.Repositories.Nodes, ch)
	}
	for _, u := range data.Users {
		e.processRepositories(u.Repositories.Nodes, ch)
	}
	e.processRepositories(data.Repositories, ch)

	// Set Rate limit stats
	ch <- prometheus.MustNewConstMetric(e.APIMetrics["Limit"], prometheus.GaugeValue, data.RateLimit.Limit)
	// ch <- prometheus.MustNewConstMetric(e.APIMetrics["Limit"], prometheus.GaugeValue, data.RateLimit.Limit)
	ch <- prometheus.MustNewConstMetric(e.APIMetrics["Remaining"], prometheus.GaugeValue, data.RateLimit.Remaining)
	// ch <- prometheus.MustNewConstMetric(e.APIMetrics["Reset"], prometheus.GaugeValue, data.RateLimit.ResetAt)

	return nil
}

func (e *Exporter) processRepositories(repos []Repository, ch chan<- prometheus.Metric) {
	for _, r := range repos {
		labelValues := r.labelValues()
		metric, err := prometheus.NewConstMetric(e.APIMetrics["Stars"], prometheus.GaugeValue, r.Stargazers.TotalCount, labelValues...)
		if err == nil {
			ch <- metric
		}
	}
}

func (r *Repository) labelValues() []string {
	// repoLabels = []string{"repo", "user", "private", "fork", "archived", "license", "language"}
	return []string{
		r.Name,
		r.Owner.Login,
		strconv.FormatBool(r.IsPrivate),
		strconv.FormatBool(r.IsFork),
		strconv.FormatBool(r.IsArchived),
		r.LicenseInfo.Key,
		r.PrimaryLanguage.Name,
	}
}
