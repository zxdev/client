package job

// NewFirewall is the job.Firewall configurator
func NewFirewall(a string) *Firewall { return &Firewall{Host: a} }

// Firewall is the worker struct
type Firewall struct {
	UUID    uint64   `json:"uuid,omitempty"`    // unique job tracking id
	Status  int      `json:"status,omitempty"`  // HTTP-like status: 0 = ok, nonzero = error
	Host    string   `json:"host,omitempty"`    // query; domain or ip,ip, ... ip list
	IP      []string `json:"ip,omitempty"`      // IPv4 queried or resolved from host
	Block   bool     `json:"block,omitempty"`   // block flag
	Version int64    `json:"version,omitempty"` // unix timestamp of last update
	//Version uint64   `json:"version,omitempty"` // version hash of the current pulled object
}

func (j *Firewall) Okay() bool      { return j.Status == 0 }
func (j *Firewall) Request() string { return j.Host }
func (j *Firewall) Unpack() any     { return *j }
