package job

// job rdap
//  support client.Job interface

// NewRdap is the job.Rdap configurator
func NewRdap(a string) *Rdap { return &Rdap{Host: a} }

// RdapJob is the request struct for POST /rdap with optional full response
// Example request:
//
//	{
//	  "host": "example.com",
//	  "full": true  // Request full RDAP domain object
//	}
type RdapJob struct {
	Host string `json:"host"`           // hostname/domain
	Full bool   `json:"full,omitempty"` // if true, include full RDAP domain object
}

// Rdap is your job wrapper.
// It now includes the full RDAP Domain response under Domain when requested.
// For backward compatibility, existing fields (NRD, NameServer) are always populated.
type Rdap struct {
	UUID       uint64   `json:"uuid,omitempty"`       // unique job tracking id
	Status     int      `json:"status,omitempty"`     // status of request; !=0 fail
	Host       string   `json:"host,omitempty"`       // zxdev.com; request
	NRD        bool     `json:"nrd,omitempty"`        // NRD flag (backward compatible)
	NameServer []string `json:"nameserver,omitempty"` // related NameServers (backward compatible)

	// Full RDAP payload (only populated when full=true in request)
	Domain *RdapDomain `json:"domain,omitempty"`
}

func (j *Rdap) Okay() bool      { return j.Status == 0 }
func (j *Rdap) Request() string { return j.Host }
func (j *Rdap) Unpack() any     { return *j }

// Full RDAP Domain object
type RdapDomain struct {
	ObjectClassName  string           `json:"objectClassName,omitempty"`
	Handle           string           `json:"handle,omitempty"`
	LDHName          string           `json:"ldhName,omitempty"`
	UnicodeName      string           `json:"unicodeName,omitempty"`
	Status           []string         `json:"status,omitempty"`
	Entities         []RdapEntity     `json:"entities,omitempty"`
	Events           []RdapEvent      `json:"events,omitempty"`
	Nameservers      []RdapNameserver `json:"nameservers,omitempty"`
	Links            []RdapLink       `json:"links,omitempty"`
	PublicIDs        []RdapPublicID   `json:"publicIds,omitempty"`
	Remarks          []RdapRemark     `json:"remarks,omitempty"`
	Notices          []RdapNotice     `json:"notices,omitempty"`
	SecureDNS        *RdapSecureDNS   `json:"secureDNS,omitempty"`
	Variants         []RdapVariant    `json:"variants,omitempty"`
	RdapConformance  []string         `json:"rdapConformance,omitempty"`
	Port43           string           `json:"port43,omitempty"`
	DelegationStatus []string         `json:"delegationStatus,omitempty"`
}

// =========================
// RDAP Substructures
// =========================

// Entities: Registrar, Registrant, Admin, Tech
type RdapEntity struct {
	ObjectClassName string         `json:"objectClassName,omitempty"`
	Handle          string         `json:"handle,omitempty"`
	Roles           []string       `json:"roles,omitempty"`
	PublicIDs       []RdapPublicID `json:"publicIds,omitempty"`
	VCardArray      interface{}    `json:"vcardArray,omitempty"` // raw vCard array
	Links           []RdapLink     `json:"links,omitempty"`
	Events          []RdapEvent    `json:"events,omitempty"`
	Remarks         []RdapRemark   `json:"remarks,omitempty"`
	Notices         []RdapNotice   `json:"notices,omitempty"`
	Entities        []RdapEntity   `json:"entities,omitempty"` // nested entities (e.g., abuse contacts)
}

// Event structure
type RdapEvent struct {
	EventAction string     `json:"eventAction,omitempty"`
	EventDate   string     `json:"eventDate,omitempty"`
	Links       []RdapLink `json:"links,omitempty"`
}

// Nameserver object
type RdapNameserver struct {
	ObjectClassName string       `json:"objectClassName,omitempty"`
	LDHName         string       `json:"ldhName,omitempty"`
	UnicodeName     string       `json:"unicodeName,omitempty"`
	IPAddresses     *RdapIPAddrs `json:"ipAddresses,omitempty"`
	Links           []RdapLink   `json:"links,omitempty"`
}

// IPv4 + IPv6
type RdapIPAddrs struct {
	IPv4 []string `json:"v4,omitempty"`
	IPv6 []string `json:"v6,omitempty"`
}

// RDAP Link object
type RdapLink struct {
	Value    string   `json:"value,omitempty"`
	Rel      string   `json:"rel,omitempty"`
	Href     string   `json:"href,omitempty"`
	Type     string   `json:"type,omitempty"`
	HRefLang []string `json:"hreflang,omitempty"`
	Title    string   `json:"title,omitempty"`
	Media    string   `json:"media,omitempty"`
}

// Public ID object
type RdapPublicID struct {
	Type       string `json:"type,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

// Remarks
type RdapRemark struct {
	Title       string     `json:"title,omitempty"`
	Type        string     `json:"type,omitempty"`
	Description []string   `json:"description,omitempty"`
	Links       []RdapLink `json:"links,omitempty"`
}

// Notices
type RdapNotice struct {
	Title       string     `json:"title,omitempty"`
	Type        string     `json:"type,omitempty"`
	Description []string   `json:"description,omitempty"`
	Links       []RdapLink `json:"links,omitempty"`
}

// DNSSEC
type RdapSecureDNS struct {
	DelegationSigned bool          `json:"delegationSigned,omitempty"`
	DSData           []RdapDSData  `json:"dsData,omitempty"`
	KeyData          []RdapKeyData `json:"keyData,omitempty"`
	MaxSigLife       int           `json:"maxSigLife,omitempty"`
}

type RdapDSData struct {
	KeyTag     int        `json:"keyTag,omitempty"`
	Algorithm  int        `json:"algorithm,omitempty"`
	Digest     string     `json:"digest,omitempty"`
	DigestType int        `json:"digestType,omitempty"`
	Links      []RdapLink `json:"links,omitempty"`
}

type RdapKeyData struct {
	Flags     int        `json:"flags,omitempty"`
	Protocol  int        `json:"protocol,omitempty"`
	Algorithm int        `json:"algorithm,omitempty"`
	PublicKey string     `json:"publicKey,omitempty"`
	Links     []RdapLink `json:"links,omitempty"`
}

// IDN domain variants
type RdapVariant struct {
	Relation []string         `json:"relation,omitempty"`
	IDNs     []RdapIDNVariant `json:"idns,omitempty"`
}

type RdapIDNVariant struct {
	LDHName     string `json:"ldhName,omitempty"`
	UnicodeName string `json:"unicodeName,omitempty"`
	VariantType string `json:"variantType,omitempty"`
}
