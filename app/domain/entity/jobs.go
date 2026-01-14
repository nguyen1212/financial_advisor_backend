package entity

import "time"

type WebDomain string

const (
	WebDomainVnExpress WebDomain = "vnexpress.net"
)

type WebScrapperJob struct {
	Domain WebDomain `json:"domain"`
	URL    string    `json:"url"`
	NewsID uint64    `json:"id"`
}

type JobStatus int

const (
	JobStatusNew JobStatus = iota
	JobStatusProcessing
	JobStatusCompleted
	JobStatusFailed
)

type JobType int

const (
	JobTypeUnknown JobType = iota
	JobTypeWebScrapper
)

func ToJobType(jobType string) JobType {
	switch jobType {
	case "web_scrapper":
		return JobTypeWebScrapper
	default:
		return JobTypeWebScrapper
	}
}

type JobResult struct {
	Error string `json:"error,omitempty"`
}

type Job struct {
	UUID string `gorm:"primaryKey"`

	Payload   []byte
	Result    JobResult `gorm:"-"`
	ResultEnc []byte    `gorm:"column:result"`

	Status JobStatus
	Type   JobType

	CreatedAt time.Time
	UpdatedAt time.Time
}
