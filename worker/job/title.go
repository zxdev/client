package job

//
// job title
//  supports client.job interface

// NewTitle is the job.Title configurator
func NewTitle(a string) *Title { return &Title{Url: a} }

type Title struct {
	UUID   uint64 `json:"uuid,omitempty"`   // unique job tracking id
	Status int    `json:"status,omitempty"` // request status
	Url    string `json:"url,omitempty"`    // target url or host
	Title  string `json:"title,omitempty"`  // title extraction
}

func (j *Title) Okay() bool      { return j.Status == 0 }
func (j *Title) Request() string { return j.Url }
func (j *Title) Unpack() any     { return *j }
