package attendance

import "time"

type Worker struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	ProjectID string    `db:"project_id" json:"project_id"`
	Phone     *string   `db:"phone" json:"phone"`
	Trade     *string   `db:"trade" json:"trade"`
	UserID    *string   `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Attendance struct {
	ID        string     `db:"id" json:"id"`
	ProjectID string     `db:"project_id" json:"project_id"`
	WorkerID  string     `db:"worker_id" json:"worker_id"`
	CheckIn   time.Time  `db:"check_in" json:"check_in"`
	CheckOut  *time.Time `db:"check_out" json:"check_out"`
	Notes     *string    `db:"notes" json:"notes"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type AttendanceWithWorker struct {
	Attendance
	WorkerName  string  `db:"worker_name" json:"worker_name"`
	WorkerTrade *string `db:"worker_trade" json:"worker_trade"`
}

type CreateWorkerRequest struct {
	Name      string  `json:"name" binding:"required,min=2"`
	ProjectID string  `json:"project_id" binding:"required,uuid"`
	Phone     *string `json:"phone"`
	Trade     *string `json:"trade"`
}

type CheckInRequest struct {
	ProjectID string    `json:"project_id" binding:"required,uuid"`
	WorkerID  string    `json:"worker_id" binding:"required,uuid"`
	CheckIn   time.Time `json:"check_in" binding:"required"`
	Notes     *string   `json:"notes"`
}

type CheckOutRequest struct {
	CheckOut time.Time `json:"check_out" binding:"required"`
}

type LinkUserRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

type SelfCheckInRequest struct {
	ProjectID string  `json:"project_id" binding:"required,uuid"`
	Notes     *string `json:"notes"`
}

type AttendanceFilter struct {
	ProjectID string
	WorkerID  string
	DateFrom  string
	DateTo    string
	Page      int
	Limit     int
}
