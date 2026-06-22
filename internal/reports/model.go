package reports

import "time"

type ReportRequest struct {
	ProjectID string     `form:"project_id" binding:"required,uuid"`
	DateFrom  *time.Time `form:"date_from" time_format:"2006-01-02"`
	DateTo    *time.Time `form:"date_to" time_format:"2006-01-02"`
}

type ProjectData struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Location    *string `db:"location"`
	Status      string  `db:"status"`
	StartDate   *string `db:"start_date"`
	EndDate     *string `db:"end_date"`
}

type LogData struct {
	ID           string  `db:"id"`
	Date         string  `db:"date"`
	CreatedBy    string  `db:"created_by_name"`
	Weather      *string `db:"weather"`
	Notes        *string `db:"notes"`
	Status       string  `db:"status"`
	GPSLat       *float64 `db:"gps_lat"`
	GPSLng       *float64 `db:"gps_lng"`
}

type IssueData struct {
	ID          string  `db:"id"`
	Title       string  `db:"title"`
	Priority    string  `db:"priority"`
	Status      string  `db:"status"`
	ReportedBy  string  `db:"reported_by_name"`
	AssignedTo  *string `db:"assigned_to_name"`
	CreatedAt   string  `db:"created_at"`
}

type AttendanceData struct {
	WorkerName  string  `db:"worker_name"`
	Trade       *string `db:"trade"`
	CheckIn     string  `db:"check_in"`
	CheckOut    *string `db:"check_out"`
}

type ReportData struct {
	Project    ProjectData
	Logs       []LogData
	Issues     []IssueData
	Attendance []AttendanceData
	GeneratedAt time.Time
	DateFrom   string
	DateTo     string
}
