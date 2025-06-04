package job

import (
	"strconv"
	"strings"
)

/*

{
  "rcode": 23,
  "host": "{host}",
  "mx": [
    "mx01.mail.icloud.com",
    "mx02.mail.icloud.com"
  ],
  "spf": [
    "apple-domain=Ocg7wboGhDipjpNm",
    "v=spf1 include:icloud.com ~all"
  ]
}

{
  "rcode": 7,
  "host": "{host}",
  "spf": [
    "v=spf1 -all"
  ],
  "dmarc": [
    "v=DMARC1; p=reject; sp=reject; adkim=s; aspf=s;"
  ],
  "dkim": [
    "v=DKIM1; p="
  ]
}

{
  "rcode": 23,
  "host": "{mail}",
  "mx": [
    "aspmx.l.google.com",
    "alt2.aspmx.l.google.com",
    "alt1.aspmx.l.google.com",
    "aspmx2.googlemail.com",
    "aspmx3.googlemail.com",
    "aspmx5.googlemail.com",
    "aspmx4.googlemail.com"
  ],
  "spf": [
    "MS=ms85975312",
    "v=spf1 include:_spf1.netstar-inc.com include:_spf2.netstar-inc.com include:_spf3.netstar-inc.com include:_spf4.netstar-inc.com include:_spf.google.com include:spf.mptx.jp include:spf.protection.outlook.com ~all",
    "apple-domain-verification=qNTIP4lfsu9oSSI6"
  ],
  "dmarc": [
    "v=DMARC1; p=none; adkim=r; aspf=r; rua=mailto:dmarc@netstar-inc.com,mailto:rua-mptx@mpub.ne.jp"
  ]
}

*/

const (
	// rcode flags
	SPF = 1 << iota
	DMARC
	DKIM
	BIMI

	// MX = 16 // from DNS

)

// NewMail is the Mail job configurator
func NewMail(host string) *Mail { return &Mail{Host: host} }

// Mail struct for detecting the host mail configuration
type Mail struct {
	UUID   uint64   `json:"uuid,omitempty"`   // unique job tracking id
	Status int      `json:"status,omitempty"` // status of request
	RCode  int      `json:"rcode,omitempty"`  // resolution request/response flags
	Host   string   `json:"host,omitempty"`   // host request
	MX     []string `json:"mx,omitempty"`     // MX records
	Spf    []string `json:"spf,omitempty"`    // TXT host
	Dmarc  []string `json:"dmarc,omitempty"`  // TXT _dmarc.
	Bimi   []string `json:"bimi,omitempty"`   // TXT <selector>.bimi.
	Dkim   []string `json:"dkim,omitempty"`   // TXT <selector>._domainkey.
}

func (j *Mail) Okay() bool      { return j.Status == 0 }
func (j *Mail) Request() string { return j.Host }
func (j *Mail) Unpack() any     { return *j }

// MailDecode returns a textual represenation of the rCode record types
//
//	0b 0 0 0 0 0 0 0 0
//	   | | | | 1 = SPF
//	   | | | 2   = DMARC
//	   | | 4     = DKIM
//	   | 8       = BIMI
//	   16        = MX
func MailDecode(rcode *int) (text []string) {

	if *rcode&MX > 0 {
		text = append(text, "MX")
	}
	if *rcode&SPF > 0 {
		text = append(text, "SPF")
	}
	if *rcode&DMARC > 0 {
		text = append(text, "DMARC")
	}
	if *rcode&BIMI > 0 {
		text = append(text, "BIMI")
	}
	if *rcode&DKIM > 0 {
		text = append(text, "DKIM")
	}
	return
}

// SPFResult parse response type
type SPFResult struct {
	Valid   bool     // valid for protection against potential fraud or impersional
	Version string   // spf version
	IP4     []string // ipv4 address
	IP6     []string // ipv6 address
	A       []string // a record
	MX      []string // mx record
	Include []string // include directive

	//	spf all flag
	//	1: 00b0001 Pass   +all; default when no policy stated; not protective
	//	2:0b00010 Neutral ?all; no authority mechanism stated; not protective
	//	4:0b00100 Soft    ~all; authorization mechanism suspicious; protective
	//	8:0b01000 Hard    -all; authorization mechanism reject; protective
	All int
}

// AllDecode returns the text version of the All flag
func (s *SPFResult) AllDecode() (all string) {
	switch s.All {
	case 1:
		all = "+all"
	case 2:
		all = "?all"
	case 4:
		all = "~all"
	case 8:
		all = "-all"
	}
	return
}

// ParseSPF record and reutrns Valid:true in the reponse to indicate
// spf policy provides protection against fraud or impersional
//
//	minimal threshold that provide valid:true spf protective
//	v=spf1 ~all
//	v=spf1 -all
//
//	spf all flag
//	00b0001 1:Pass    +all; no policy stated, default; not protective
//	0b00010 2:Neutral ?all; no authority stated; not protective
//	0b00100 4:Soft    ~all; soft fail, mark suspicious; protective
//	0b01000 8:Hard    -all; hard fail, rejcet; protective
func ParseSPF(m *Mail) (parse SPFResult) {

	// common
	//  v=spf1 ~all
	for i := range m.Spf {

		// while other records may exist, we only
		// want to parse the spf policy record which
		// should not be larger than 512 bytes and
		// we are looking for normal configurations
		m.Spf[i] = strings.ToLower(m.Spf[i])
		if strings.HasPrefix(m.Spf[i], "v=spf1") {

			// the spf policy record must start with
			// the version reference; version 1
			parse.Version = "spf1"

			// policy says all must be the last reference and
			// everything else beond should be ignored, but
			// when it is not the last reference we fail
			switch {
			case strings.HasSuffix(m.Spf[i], "~all"):
				// ~all is a soft fail of authentication and means that it is unlikely
				// that the ip/domain is authorized send; mark as suspicious
				parse.All, parse.Valid = 0b0100, true // 4, protects against fraud or impersonation

			case strings.HasSuffix(m.Spf[i], "-all"):
				// -all is a hard fail of authentication and means that the ip/domain
				// is not authorized to send mail; hard failure, do not accept
				parse.All, parse.Valid = 0b1000, true // 8, protects againts fraud or impersonation

			case strings.HasSuffix(m.Spf[i], "+all"):
				// +all is the default pass authentication; and assume that ip/domin is authorized
				// to send mail; allow, but this fails to profect against fraud or impersonation
				parse.All, parse.Valid = 0b0001, false // 1, no protection against fraud or impersoniation

			case strings.HasSuffix(m.Spf[i], "?all"):
				// ?all is neutral so it means that ip/host neither passes or fails
				// since the record does not explicity state; allow, but this fails to profect against fraud or impersonation
				parse.All, parse.Valid = 0b0010, false // 2,  no protection against fraud or impersoniation

			default:
				return
			}

			// v=spf1 ~all; this is all that is
			// requried for a minimally valid spf record
			parse.Valid = true

			// parse the ip/host sections
			var idx int
			for pair := range strings.SplitSeq(m.Spf[i], " ") {

				pair = strings.TrimSpace(pair)
				idx = strings.Index(pair, ":")
				if idx < 0 {
					continue
				}

				switch strings.ToLower(pair[:idx]) {
				case "ip4":
					parse.IP4 = append(parse.IP4, pair[idx+1:])
				case "ip6":
					parse.IP6 = append(parse.IP6, pair[idx+1:])
				case "a":
					parse.A = append(parse.A, pair[idx+1:])
				case "mx":
					parse.MX = append(parse.MX, pair[idx+1:])
				case "include":
					parse.Include = append(parse.Include, pair[idx+1:])
				}

			}

			return
		}
	}

	return
}

// DMARCResult response record type
type DMARCResult struct {
	Valid   bool
	Version string   // DMARC1
	P       string   // none, quarantine, reject
	Pct     int      // 1 to 100; with BIMI must be 100
	Rua     []string // mailto:report@example.com,mailto:report2@example.com
	Ruf     []string // mailto:fail@example.com
	SP      string   // optional; none, rua, quarantine, reject
	Adkim   string   // optional; s or r
	Aspf    string   // optional; s or r
}

// ParseDMARC record and return Valid:true when v,p fields are set and pct>0
//
//	"v=DMARC1; p=none; adkim=r; aspf=r; rua=mailto:dmarc@netstar-inc.com,mailto:rua-mptx@mpub.ne.jp"
func ParseDMARC(m *Mail) (result DMARCResult) {

	var idx int
	if len(m.Dmarc) == 1 {
		// you can only have one DMARC record so we do not need to loop over any

		m.Dmarc[0] = strings.ToLower(m.Dmarc[0])
		if strings.HasPrefix(m.Dmarc[0], "v=dmarc1") {
			result.Version = "dmarc1" // required

			for pair := range strings.SplitSeq(m.Dmarc[0], ";") {
				pair = strings.TrimSpace(pair)
				idx = strings.Index(pair, "=")
				if idx < 0 {
					continue
				}

				switch pair[:idx] {
				case "p": // required
					switch pair[idx+1:] {
					case "none": // send report to rua
					case "quarantine": // mark as spam
					case "reject": // reject; send bounce message
					default:
						continue // invalid
					}
					result.P, result.Valid = pair[idx+1:], true

				case "pct": // quarantine percentage
					result.Pct, _ = strconv.Atoi(pair[idx+1:])
					if result.Pct == 0 {
						result.Valid = false
					}

				case "rua": // reporting address; recommended
					for mail := range strings.SplitSeq(pair[idx+1:], ",") {
						result.Rua = append(result.Rua, strings.TrimPrefix(strings.TrimSpace(mail), "mailto:"))
					}
				case "ruf": // failure report address; not recommended
					for mail := range strings.SplitSeq(pair[idx+1:], ",") {
						result.Ruf = append(result.Ruf, strings.TrimPrefix(strings.TrimSpace(mail), "mailto:"))
					}

				case "sp":
					switch pair[idx+1:] {
					case "none": // send report to rua
					case "quarantine": // mark as spam
					case "reject": // reject; send bounce message
					default:
						continue // invalid
					}
					result.SP = pair[idx+1:]
				case "adkim":
					switch pair[idx+1:] {
					case "s": // strict alignment
					case "r": // relaxed alignment
					default:
						continue // invalid
					}
					result.Adkim = pair[idx+1:]
				case "aspf":
					switch pair[idx+1:] {
					case "s": // strict alignment
					case "r": // relaxed alignment
					default:
						continue // invalid
					}
					result.Aspf = pair[idx+1:]
				}
			}
		}
	}
	return
}

// BIMI record result
type BIMIResult struct {
	Valid   bool   `json:"valid,omitempty"`
	Version string `json:"version,omitempty"`
	L       string `json:"l,omitempty"`
	A       string `json:"a,omitempty"`
}

// ParseBIMI record sets the Valid:true bool when the required elements
// are present, howeever BIMI also requires pct=100 in the DMARC record so
// this must also be part of the request set
//
//	v=BIMI1;l=https://images.solarmora.com/brand/bimi-logo.svg
//	v=BIMI1;l=;a=https://images.solarmora.com/brand/certificate.pem
func ParseBIMI(m *Mail) (result BIMIResult) {

	var idx int

	if len(m.Bimi) == 1 {
		// you should only have one bimi record so no reason to loop

		m.Bimi[0] = strings.ToLower(m.Bimi[0])
		if strings.HasPrefix(m.Bimi[0], "v=bimi1") {
			result.Version = "bimi1"

			for pair := range strings.SplitSeq(m.Bimi[0], ";") {
				pair = strings.TrimSpace(pair)
				idx = strings.Index(pair, "=")
				if idx < 0 {
					continue
				}

				switch pair[:idx] {
				case "l": // required
					result.L = pair[idx+1:]
				case "a": // required
					result.A = pair[idx+1:]
				}
			}

			if len(m.Dmarc) == 1 {
				// can only have one dmarc record
				result.Valid = len(result.A) > 0 || len(result.L) > 0
				result.Valid = strings.Contains(m.Dmarc[0], "pct=100")
			}

		}
	}
	return
}

type DKIMResult struct {
	Valid   bool   `json:"valid,omitempty"`
	Version string `json:"version,omitempty"`
	P       string `json:"p,omitempty"` // public key
	K       string `json:"k,omitempty"` // key type
	T       string `json:"t,omitempty"` // sender is testing dkim
}

// ParseDKIM will parse the DKIM record and extract the public key, however
// this will only valid:true indicating that the record has the required
// elements present and that is all
func ParseDKIM(m *Mail) (result DKIMResult) {

	// this is only the DNS DKIM key process, the other part of DKIM
	// comes from processing the email to verify the dkim key
	//
	// email header (v,b,bd,d,s,a,h are required)
	//
	// v = the version
	// b = the actual digital signature of the contents (headers and body) of the mail message
	// bh = the body hash
	// d = the signing domain
	// s = the selector
	// a = the signing algorithm
	// c = the canonicalization algorithm(s) for header and body
	// q = the default query method
	// l = the length of the canonicalized part of the body that has been signed
	// t = the signature timestamp
	// x = the expire time
	// h = the list of signed header fields, repeated for fields that occur multiple times

	var idx int

	if len(m.Dkim) == 1 {
		// you should only have one dkim record so no reason to loop

		m.Dkim[0] = strings.ToLower(m.Dkim[0])
		if strings.HasPrefix(m.Dkim[0], "v=dkim1") {
			result.Version = "dkim1" // required

			for pair := range strings.SplitSeq(m.Dkim[0], ";") {
				pair = strings.TrimSpace(pair)
				idx = strings.Index(pair, "=")
				if idx < 0 {
					continue
				}
				switch pair[:idx] {
				case "p": // required
					result.P = pair[idx+1:]
				case "k": // required
					result.K = pair[idx+1:]
				case "t": // optional
					result.T = pair[idx+1:]
				}
			}

			if len(m.Dkim) > 0 {
				result.Valid = len(result.P) > 0 && len(result.K) > 0
			}
		}
	}
	return

}
