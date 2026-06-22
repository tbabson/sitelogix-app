package search

type ResultType string

const (
	TypeProject  ResultType = "project"
	TypeLog      ResultType = "log"
	TypeIssue    ResultType = "issue"
	TypeWorker   ResultType = "worker"
	TypeMedia    ResultType = "media"
)

type Result struct {
	Type    ResultType  `json:"type"`
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	Excerpt string      `json:"excerpt"`
	Meta    interface{} `json:"meta"`
}

type SearchResponse struct {
	Query   string   `json:"query"`
	Total   int      `json:"total"`
	Results []Result `json:"results"`
}

type projectRow struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Status      string  `db:"status"`
}

type logRow struct {
	ID          string  `db:"id"`
	ProjectName string  `db:"project_name"`
	Date        string  `db:"date"`
	Notes       *string `db:"notes"`
	Status      string  `db:"status"`
}

type issueRow struct {
	ID          string  `db:"id"`
	ProjectName string  `db:"project_name"`
	Title       string  `db:"title"`
	Description *string `db:"description"`
	Priority    string  `db:"priority"`
	Status      string  `db:"status"`
}

type workerRow struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	ProjectName string  `db:"project_name"`
	Trade       *string `db:"trade"`
}
