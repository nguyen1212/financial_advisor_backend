package entity

type Job struct {
	Payload []byte
}

type WebDomain string

const (
	WebDomainVnExpress WebDomain = "vnexpress.net"
)

type WebScrapperJob struct {
	Domain WebDomain `json:"domain"`
	URL    string    `json:"url"`
	NewsID uint64    `json:"id"`
}
