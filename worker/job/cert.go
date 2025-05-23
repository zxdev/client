package job

import (
	"strings"
	"time"
)

//
// cert
//  support client.Job
//

/*

[
    {
      "issuer_ca_id": 286236,
      "issuer_name": "C=US, O=Google Trust Services, CN=WE1",
      "common_name": "zxdev.com",
      "name_value": "*.zxdev.com\nzxdev.com",
      "id": 13489789012,
      "entry_timestamp": "2024-06-23T13:12:25.328",
      "not_before": "2024-06-23T12:12:24",
      "not_after": "2024-09-21T12:12:23",
      "serial_number": "474dec9ce6330ac613470f3474b728e0",
      "result_count": 3
	},
	...
]
*/

// NewCert is the job.Cert configurator
func NewCert(a string) *Cert { return &Cert{Host: a} }

// worker cert struct
type Cert struct {
	UUID        uint64        `json:"uuid,omitempty"`   // unique job tracking id
	Status      int           `json:"status,omitempty"` // status of request
	Host        string        `json:"host,omitempty"`   // zxdev.com; request
	Certificate []certificate `json:"data,omitempty"`
}

func (j *Cert) Okay() bool      { return j.Status == 0 }
func (j *Cert) Request() string { return j.Host }
func (j *Cert) Unpack() any     { return *j }

// worker response certificate record
type certificate struct {
	ID             uint64 `json:"id,omitempty"`              // 13489789012; crt.sh database ID
	IssuerCaId     int    `json:"issuer_ca_id,omitempty"`    // 286236
	IssuerName     string `json:"issuer_name,omitempty"`     // "C=US, O=Google Trust Services, CN=WE1"
	CommonName     string `json:"common_name,omitempty"`     // zxdev.com
	NameValue      string `json:"name_value,omitempty"`      // zxdev.com\n*zxdev.com; conver to []string
	EntryTimestamp string `json:"entry_timestamp,omitempty"` // RFC3339; convert to unix timestamp
	NotBefore      string `json:"not_before,omitempty"`      // RFC3339; convert to unix timestamp
	NotAfter       string `json:"not_after,omitempty"`       // RFC3339; convert to unix timestamp
	SerialNumber   string `json:"serial_number,omitempty"`   // 474dec9ce6330ac613470f3474b728e0
	ResultCount    int    `json:"result_count,omitempty"`    // 3; = CommonName + NameValue
}

// NameValues parses c.NameValue into a slice
func (c *certificate) NameValues(n string) []string { return strings.Split(c.NameValue, "\n") }

// EntryUnix parses c.EntryTimestamp into a unix timestamp
func (c *certificate) EntryUnix() int64 { return c.unix(c.EntryTimestamp) }

// NotBeforeUnix parses c.NotBefore into a unix timestamp
func (c *certificate) NotBeforeUnix() int64 { return c.unix(c.NotBefore) }

// NotAfterUnix parse c.NotAfter into a unix timestamp
func (c *certificate) NotAfterUnix() int64 { return c.unix(c.NotAfter) }

func (c *certificate) unix(a string) int64 {
	ts, _ := time.Parse("2006-01-02T15:04:05", a)
	return ts.UTC().Unix()
}
