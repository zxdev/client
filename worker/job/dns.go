package job

const (
	// rcode flags
	A = 1 << iota
	AAAA
	CNAME
	NS
	MX
	TXT
	PTR
	DOMAIN
)

//
// dns job
//  support client.Job interface

// NewDNS is the job.DNS configurator
func NewDNS(a string) *DNS { return &DNS{Host: a} }

// numeric RCode type request combinations (short codes)
//
//	  0b 0 0 0 0 0 0 0 0
//	     |   | | | | | 1 = A
//	     |   | | | | 2   = AAAA
//	     |   | | | 4     = CNAME
//	     |   | | 8       = NS
//	     |   | 16        = MX
//	     |   32          = TXT
//		 128             = DOMAIN
//
//	1 A
//	2 AAAA
//	3 A,AAAA
//	4 CNAME
//	5 A,CNAME
//	6 AAAA,CNAME
//	7 A,AAAA,CNAME
//	8 NS
//	9 A,NS
//	10 AAAA,NS
//	11 A,AAAA,NS
//	12 CNAME,NS
//	13 A,CNAME,NS
//	14 AAAA,CNAME,NS
//	15 A,AAAA,CNAME,NS
//	16 MX
//	32 TXT
//	48 MX,TXT
//	128 DOMAIN
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

// check Rcode response flag
func HasA(rcode *int) bool      { return *rcode&A != 0 }
func HasAAAA(rcode *int) bool   { return *rcode&AAAA != 0 }
func HasCNAME(rcode *int) bool  { return *rcode&CNAME != 0 }
func HasNS(rcode *int) bool     { return *rcode&NS != 0 }
func HasMX(rcode *int) bool     { return *rcode&MX != 0 }
func HasTXT(rcode *int) bool    { return *rcode&TXT != 0 }
func HasDOMAIN(rcode *int) bool { return *rcode&DOMAIN != 0 }

// DNSDecode returns a textual represenation of the record types
func DNSDecode(rcode *int) (text []string) {

	if *rcode&A > 0 {
		text = append(text, "A")
	}
	if *rcode&AAAA > 0 {
		text = append(text, "AAAA")
	}
	if *rcode&CNAME > 0 {
		text = append(text, "CNAME")
	}
	if *rcode&NS > 0 {
		text = append(text, "NS")
	}
	if *rcode&MX > 0 {
		text = append(text, "MX")
	}
	if *rcode&TXT > 0 {
		text = append(text, "TXT")
	}

	if *rcode&DOMAIN > 0 {
		text = append(text, "DOMAIN")
	}
	return
}
