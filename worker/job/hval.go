package job

import (
	"net"
	"net/http"
)

//
// job hval
//  supports client.Job interface

const (
	// security flags

	HSTS = 1 << iota // Strict-Transport-Security       ; basic
	CSP              // Content-Security-Policy         ; basic
	XCTO             // X-Content-Type-Options: nosniff
	ACAO             // Access-Control-Allow-Origin     ; CORS
	COOP             // Cross-Origin-Opener-Policy      ; CORS
	CORP             // Cross-Origin-Resource-Policy    ; CORS
	COEP             // Cross-Origin-Embedder-Policy    ; CORS
)

// NewHval is the job.Hval configurator
func NewHval(a string) *Hval { return &Hval{Item: a} }

// HHeader type is the hval payload response
type HHeader struct {
	Status int         `json:"status,omitempty"` // status code response; !=0 fail
	Url    string      `json:"url,omitempty"`    // targeted url or host
	Header http.Header `json:"header,omitempty"` // header payload, full
	IP     []net.IP    `json:"ip,omitempty"`     // ip address resolution of target server
	TLS    string      `json:"tls,omitempty"`    // tlsCipher in use
}

// Response struct
type Hval struct {
	UUID     uint64    `json:"uuid,omitempty"`     // unique job tracking id
	Status   int       `json:"status,omitempty"`   // response status; !=0 fail
	Item     string    `json:"item,omitempty"`     // target url or host
	Head     []HHeader `json:"head,omitempty"`     // headers
	N        int       `json:"n,omitempty"`        // n hop|redirect counter; same as len(Head) when > 0 and no method/scheme prefix
	Security int       `json:"security,omitempty"` // security flags HSTS|CPS|XCTO|ACAO|COOP|CORP|COEP
}

func (j *Hval) Okay() bool      { return j.Status == 0 }
func (j *Hval) Request() string { return j.Item }
func (j *Hval) Unpack() any     { return *j }

// SecurityBasic reports true on the minimal valid security combinations of HSTS,CSP
func SecurityBasic(security *int) bool {
	// HSTS|CSP combined 1|2 = 3 will always be less than XCTO flag 4
	return *security < XCTO
}

// SecurityDeocder returns encoded *secuirty flags as http text header key slice
func SecurityDecoder(security *int) (text []string) {

	if *security&HSTS > 0 {
		text = append(text, "Strict-Transport-Security")
	}
	if *security&CSP > 0 {
		text = append(text, "Content-Security-Policy")
	}
	if *security&XCTO > 0 {
		text = append(text, "X-Content-Type-Options")
	}
	if *security&ACAO > 0 {
		text = append(text, "Access-Control-Allow-Origin")
	}
	if *security&COOP > 0 {
		text = append(text, "Cross-Origin-Opener-Policy")
	}
	if *security&CORP > 0 {
		text = append(text, "Cross-Origin-Resource-Policy")
	}
	if *security&COEP > 0 {
		text = append(text, "Cross-Origin-Embedder-Policy")
	}
	return
}
