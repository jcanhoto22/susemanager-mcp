package models

type CveAuditSystem struct {
	SystemID     int      `xmlrpc:"system_id"`
	PatchStatus  string   `xmlrpc:"patch_status"`
	ChannelLabels []string `xmlrpc:"channel_labels"`
	ErrataAdvisories []string `xmlrpc:"errata_advisories"`
}

type Errata struct {
	ID          int    `xmlrpc:"id"`
	AdvisoryName string `xmlrpc:"advisory_name"`
	AdvisoryType string `xmlrpc:"advisory_type"`
	AdvisoryRel  string `xmlrpc:"advisory_rel"`
	Product      string `xmlrpc:"product"`
	Description  string `xmlrpc:"description"`
	Synopsis     string `xmlrpc:"synopsis"`
	Topic        string `xmlrpc:"topic"`
	Severity     string `xmlrpc:"severity"`
	IssueDate    string `xmlrpc:"issue_date"`
	UpdateDate   string `xmlrpc:"update_date"`
	CVEs         string `xmlrpc:"cves"`
}
