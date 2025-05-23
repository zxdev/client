package job

//
// dns job
//  support client.Job interface

// NewDNS is the job.DNS configurator
func NewDNS(a string) *DNS { return &DNS{Host: a} }

// DNS job numeric RCode type request combinations and short codes
//
//	  0b 0 0 0 0 0 0 0 0
//	     |   | | | | | 1 = A
//	     |   | | | | 2   = AAAA
//	     |   | | | 4     = CNAME
//	     |   | | 8       = NS
//	     |   | 16        = MX
//	     |   32          = TXT
//		 128             = DOMAIN
type DNS struct {
	UUID   uint64   `json:"uuid,omitempty"`   // unique job tracking id
	RCode  int      `json:"rcode,omitempty"`  // resolution request/response flags
	Host   string   `json:"host,omitempty"`   // host request
	A      []string `json:"a,omitempty"`      // A records
	AAAA   []string `json:"aaaa,omitempty"`   // AAAA records
	CNAME  []string `json:"cname,omitempty"`  // CNAME records
	NS     []string `json:"ns,omitempty"`     // NS records
	MX     []string `json:"mx,omitempty"`     // MX records
	TXT    []string `json:"txt,omitempty"`    // TXT records
	Domain []string `json:"domain,omitempty"` // rDNS resolution target
}

func (j *DNS) Okay() bool      { return j.RCode != 0 }
func (j *DNS) Request() string { return j.Host }
func (j *DNS) Unpack() any     { return *j }

// Check Rcode response flag
func HasA(rcode *int) bool      { return *rcode&0b00000001 != 0 }
func HasAAAA(rcode *int) bool   { return *rcode&0b00000010 != 0 }
func HasCNAME(rcode *int) bool  { return *rcode&0b00000100 != 0 }
func HasNS(rcode *int) bool     { return *rcode&0b00001000 != 0 }
func HasMS(rcode *int) bool     { return *rcode&0b00010000 != 0 }
func HasTXT(rcode *int) bool    { return *rcode&0b00100000 != 0 }
func HasDOMAIN(rcode *int) bool { return *rcode&0b10000000 != 0 }
