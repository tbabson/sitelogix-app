package media

import "time"

type Media struct {
	ID          string    `db:"id" json:"id"`
	ProjectID   string    `db:"project_id" json:"project_id"`
	UploadedBy  string    `db:"uploaded_by" json:"uploaded_by"`
	Filename    string    `db:"filename" json:"filename"`
	ContentType string    `db:"content_type" json:"content_type"`
	URL         string    `db:"url" json:"url"`
	GitHubSHA   string    `db:"github_sha" json:"github_sha"`
	GitHubKey   string    `db:"github_key" json:"github_key"`
	Size        int64     `db:"size" json:"size"`
	LogID       *string   `db:"log_id" json:"log_id"`
	IssueID     *string   `db:"issue_id" json:"issue_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type MediaFilter struct {
	ProjectID string
	LogID     string
	IssueID   string
	DateFrom  string
	DateTo    string
	Page      int
	Limit     int
}

const (
	MaxUploadSize = 100 * 1024 * 1024 // 100MB
)

var AllowedTypes = map[string]string{
	"image/jpeg":      "images",
	"image/png":       "images",
	"image/gif":       "images",
	"image/webp":      "images",
	"video/mp4":       "videos",
	"video/quicktime": "videos",
	"video/webm":      "videos",
	"audio/mpeg":      "audio",
	"audio/mp4":       "audio",
	"audio/ogg":       "audio",
	"audio/webm":      "audio",
	"application/pdf": "documents",
	"application/msword":                                                        "documents",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "documents",
	"application/vnd.ms-excel":                                                  "documents",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "documents",
}
