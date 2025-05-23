package job

const (
	// method verb flags
	HEAD = 1 << iota
	GET
	POST
	PUT
	PATCH
	DELETE
	TRACE
	CONNECT

	// OPTIONS is used directly and only for setting supported
// methods when it is supported
)

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

// MethodStandard reports head,get,post and their combinations as valid for the Standard group
func MethodStandard(flag *int) bool {
	// head,get,post combined 1|2|3 =7 will always be less than the put flag 8
	return *flag < PUT
}

// MethodDecoder returns textual representation of the method flags;
// pass nil option will always ingore reporting the OPTION verb
func MethodDecoder(flag *int, option *bool) (text []string) {
	if *flag&HEAD > 0 {
		text = append(text, "HEAD")
	}
	if *flag&GET > 0 {
		text = append(text, "GET")
	}
	if *flag&POST > 0 {
		text = append(text, "POST")
	}
	if *flag&PUT > 0 {
		text = append(text, "PUT")
	}
	if *flag&PATCH > 0 {
		text = append(text, "PATCH")
	}
	if *flag&DELETE > 0 {
		text = append(text, "DELETE")
	}
	if *flag&TRACE > 0 {
		text = append(text, "TRACE")
	}

	if option != nil && *option {
		text = append(text, "OPTION")
	}

	if *flag&CONNECT > 0 {
		text = append(text, "CONNECT")
	}

	return
}
