package job

// job rdap
//  support client.Job interface

// NewRdap is the job.Rdap configurator
func NewRdap(a string) *Rdap { return &Rdap{Host: a} }

type Rdap struct {
	UUID       uint64   `json:"uuid,omitempty"`       // unique job tracking id
	Status     int      `json:"status,omitempty"`     // status of request; !=0 fail
	Host       string   `json:"host,omitempty"`       // zxdev.com; request
	NRD        bool     `json:"nrd,omitempty"`        // NRD flag
	NameServer []string `json:"nameserver,omitempty"` // realted NameServers
}

func (j *Rdap) Okay() bool      { return j.Status == 0 }
func (j *Rdap) Request() string { return j.Host }
func (j *Rdap) Unpack() any     { return *j }
