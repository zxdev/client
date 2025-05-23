package job

//
// job method
//  support client.Job interface

// NewMethod is the job.Method configurator
func NewMethod(a string) *Method { return &Method{Url: a} }

type Method struct {
	UUID    uint64 `json:"uuid,omitempty"`    // unique job tracking id
	Status  int    `json:"status,omitempty"`  // request status; !=0 fail
	Url     string `json:"url,omitempty"`     // url or host target
	Options bool   `json:"options,omitempty"` // signals non-interrogation of methods
	Flag    int    `json:"flag,omitempty"`    // method flags and groups
}

func (j *Method) Okay() bool      { return j.Status == 0 }
func (j *Method) Request() string { return j.Url }
func (j *Method) Unpack() any     { return *j }
