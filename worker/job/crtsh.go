package job

//
// crtsh
//  support client.Job
//
// The crtsh worker queries crt.sh for historical certificate data.
// This is separate from the cert worker which performs live TLS handshakes.
//
// Request format (POST /crtsh):
//   - Simple: {"host": "example.com"} or {"domain": "example.com"}
//
// Example output:
//
// {
//   "host": "example.com",
//   "count": 5,
//   "certs": [
//     {
//       "id": 16593729128,
//       "sha256": "abc123...",
//       "sha1": "def456...",
//       "common_name": "*.example.com",
//       "alt_names": ["*.example.com", "example.com"],
//       "issuer_name": "C=US, O=Let's Encrypt, CN=R3",
//       "not_before": "2025-02-07T00:00:00Z",
//       "not_after": "2026-03-10T23:59:59Z",
//       "logged_at": "2025-02-07T19:16:30Z",
//       "is_precert": false
//     }
//   ]
// }
//

// NewCRTSH is the job.CRTSH configurator
func NewCRTSH(a string) *CRTSH { return &CRTSH{Host: a} }

// CRTSH is the worker crtsh struct for historical certificate data
type CRTSH struct {
	UUID   uint64        `json:"uuid,omitempty"`   // unique job tracking id
	Status int           `json:"status,omitempty"` // HTTP-like status: 0 = ok, nonzero = error
	Host   string        `json:"host,omitempty"`   // hostname/domain queried
	Count  int           `json:"count,omitempty"`  // number of historical certificates found
	Certs  []CRTSHCert   `json:"certs,omitempty"`  // array of historical certificates (sorted by not_before desc)
}

func (j *CRTSH) Okay() bool      { return j.Status == 0 }
func (j *CRTSH) Request() string { return j.Host }
func (j *CRTSH) Unpack() any     { return *j }

// CRTSHCert contains historical certificate data from crt.sh
type CRTSHCert struct {
	ID         int64    `json:"id,omitempty"`          // crt.sh internal ID
	Sha256     string   `json:"sha256,omitempty"`       // SHA-256 fingerprint
	Sha1       string   `json:"sha1,omitempty"`         // SHA-1 fingerprint
	CommonName string   `json:"common_name,omitempty"`  // certificate common name
	AltNames   []string `json:"alt_names,omitempty"`    // subject alternative names
	IssuerName string   `json:"issuer_name,omitempty"`  // issuer distinguished name
	NotBefore  string   `json:"not_before,omitempty"`   // RFC3339 validity start
	NotAfter   string   `json:"not_after,omitempty"`   // RFC3339 validity end
	LoggedAt   string   `json:"logged_at,omitempty"`   // RFC3339 when logged in CT
	IsPrecert  bool     `json:"is_precert,omitempty"`  // true if precertificate
}

