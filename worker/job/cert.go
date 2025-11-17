package job

//
// cert
//  support client.Job
//
// The cert worker performs TLS handshake and comprehensive certificate analysis,
// returning metadata (not scores) that students can use to build their own scoring systems.
//
// Request format (POST /cert):
//   - Simple: "example.com" or {"host": "example.com"} - returns basic checks only
//   - With options: {"host": "example.com", "port": 443, "options": {"include_ocsp": true, "include_disallowed_root": true}}
//     - include_ocsp: perform OCSP revocation check (network call, slower)
//     - include_disallowed_root: check root CA against disallowed list (file load, slower)
//
// Example output:
//
// {
//   "host": "example.com",
//   "port": 443,
//   "timestamp": "2025-01-15T10:30:00Z",
//   "connection": {
//     "tls_version": "TLS 1.3",
//     "cipher_suite": "TLS_AES_128_GCM_SHA256"
//   },
//   "verification": {
//     "hostname_checked": "example.com",
//     "hostname_matches": true,
//     "chain_verified": true,
//     "chain_verify_error": "",
//     "chain_length": 3,
//     "subject_equals_issuer": false,
//     "self_signature_verifies": false,
//     "weak_crypto": false,
//     "validity_risk": false,
//     "incomplete_chain": false,
//     "ocsp_checked": true,
//     "ocsp_status": "good",
//     "root_in_disallowed_list": false
//   },
//   "revocation": {
//     "checked_at": "2025-01-15T10:30:05Z",
//     "revoked": false
//   },
//   "certs": [
//     {
//       "position": 0,
//       "role": "leaf",
//       "subject_cn": "example.com",
//       "subject_dns_names": ["example.com", "www.example.com"],
//       "issuer_cn": "R3",
//       "issuer_dn": "CN=R3, O=Let's Encrypt, C=US",
//       "issuer_friendly_name": "Let's Encrypt",
//       "not_before": "2025-10-01T00:00:00Z",
//       "not_after": "2026-01-01T00:00:00Z",
//       "public_key_algorithm": "RSA",
//       "public_key_size_bits": 2048,
//       "signature_algorithm": "SHA256-RSA",
//       "is_ca": false,
//       "serial_number": "...",
//       "cert_sha256": "abc123...",
//       "pubkey_sha256": "def456...",
//       "tbs_sha256": "ghi789...",
//       "issuer_pubkey_sha256": "jkl012..."
//     },
//     {
//       "position": 1,
//       "role": "intermediate",
//       "subject_cn": "R3",
//       "issuer_cn": "ISRG Root X1",
//       "issuer_dn": "CN=ISRG Root X1, O=Internet Security Research Group, C=US",
//       "issuer_friendly_name": "Internet Security Research Group",
//       "is_ca": true,
//       "cert_sha256": "...",
//       "pubkey_sha256": "..."
//     }
//   ]
// }
//

// NewCert is the job.Cert configurator
func NewCert(a string) *Cert { return &Cert{Host: a} }

// CertJob is the request struct for POST /cert with optional expensive checks
// Example request:
//
//	{
//	  "host": "example.com",
//	  "port": 443,
//	  "options": {
//	    "include_ocsp": true,
//	    "include_disallowed_root": true
//	  }
//	}
type CertJob struct {
	Host    string       `json:"host"`              // hostname or hostname:port
	Port    int          `json:"port,omitempty"`    // port number (defaults to 443)
	Options *CertOptions `json:"options,omitempty"` // optional expensive checks
}

// CertOptions controls which expensive/optional checks are performed
// All basic checks (connection info, chain, hostname, weak crypto, etc.) are always included.
// Only OCSP, disallowed root, and CRLite checks are optional and require explicit opt-in.
type CertOptions struct {
	IncludeOCSP           bool `json:"include_ocsp,omitempty"`            // perform OCSP revocation check (network call)
	IncludeDisallowedRoot bool `json:"include_disallowed_root,omitempty"` // check root CA against disallowed list (file load)
	IncludeCRLite         bool `json:"include_crlite,omitempty"`          // check certificate against CRLite revocation filter (fast, offline)
}

// Cert is the worker cert struct for TLS handshake and certificate analysis
type Cert struct {
	UUID      uint64 `json:"uuid,omitempty"`      // unique job tracking id
	Status    int    `json:"status,omitempty"`    // HTTP-like status: 0 = ok, nonzero = error (connection failed, invalid host, etc.)
	Host      string `json:"host,omitempty"`      // hostname or hostname:port
	Port      int    `json:"port,omitempty"`      // port number (defaults to 443)
	Timestamp string `json:"timestamp,omitempty"` // RFC3339 timestamp of when this measurement was taken

	Connection   *ConnectionInfo   `json:"connection,omitempty"`   // TLS connection details
	Verification *VerificationInfo `json:"verification,omitempty"` // certificate verification metadata
	Certs        []CertificateInfo `json:"certs,omitempty"`        // certificate chain (leaf + intermediates)
	Revocation   *RevocationInfo   `json:"revocation,omitempty"`   // detailed revocation information (only when include_ocsp option is enabled)
}

func (j *Cert) Okay() bool      { return j.Status == 0 }
func (j *Cert) Request() string { return j.Host }
func (j *Cert) Unpack() any     { return *j }

// ConnectionInfo contains TLS connection metadata
type ConnectionInfo struct {
	TLSVersion  string `json:"tls_version,omitempty"`  // e.g., "TLS 1.3"
	CipherSuite string `json:"cipher_suite,omitempty"` // e.g., "TLS_AES_128_GCM_SHA256"
}

// VerificationInfo contains certificate verification metadata (no scoring, just flags)
//
// Data tiers:
//   - Always included (cheap): HostnameChecked, HostnameMatches, ChainVerified, ChainVerifyError,
//     ChainLength, SubjectEqualsIssuer, SelfSignatureVerifies, WeakCrypto, ValidityRisk, IncompleteChain
//   - Optional/expensive (requires options): OCSPChecked, OCSPStatus, RootInDisallowedList
type VerificationInfo struct {
	// Basic checks (always included)
	HostnameChecked       string `json:"hostname_checked,omitempty"`        // hostname used for verification
	HostnameMatches       bool   `json:"hostname_matches,omitempty"`        // true if certificate matches hostname
	ChainVerified         bool   `json:"chain_verified,omitempty"`          // true if chain verified successfully
	ChainVerifyError      string `json:"chain_verify_error,omitempty"`      // error string or "" if ChainVerified == true
	ChainLength           int    `json:"chain_length,omitempty"`            // number of certificates in chain
	SubjectEqualsIssuer   bool   `json:"subject_equals_issuer,omitempty"`   // true if subject == issuer (self-signed indicator)
	SelfSignatureVerifies bool   `json:"self_signature_verifies,omitempty"` // true if cert is cryptographically self-signed
	WeakCrypto            bool   `json:"weak_crypto,omitempty"`             // true if uses weak algorithms (SHA1, MD5, etc.)
	ValidityRisk          bool   `json:"validity_risk,omitempty"`           // true if unusual validity period (<7d or >825d)
	IncompleteChain       bool   `json:"incomplete_chain,omitempty"`        // true if chain is incomplete

	// Optional/expensive checks (only populated if requested via CertOptions)
	OCSPChecked          bool   `json:"ocsp_checked,omitempty"`            // true if OCSP check was attempted (requires include_ocsp option)
	OCSPStatus           string `json:"ocsp_status,omitempty"`             // "good" | "revoked" | "unknown" | "error" | "" (requires include_ocsp option)
	RootInDisallowedList *bool  `json:"root_in_disallowed_list,omitempty"` // nil = check not run, true = root is disallowed, false = root is not disallowed (requires include_disallowed_root option)
	CRLiteChecked        bool   `json:"crlite_checked,omitempty"`          // true if CRLite check was performed (requires include_crlite option)
	CRLiteRevoked        bool   `json:"crlite_revoked,omitempty"`          // true if certificate is revoked according to CRLite
	CRLiteSource         string `json:"crlite_source,omitempty"`           // source of CRLite data (e.g., "mozilla-crlite")
	CRLiteVersion        string `json:"crlite_version,omitempty"`          // version/date of CRLite filter (e.g., "2025-11-15")
}

// RevocationInfo contains detailed revocation information from OCSP
type RevocationInfo struct {
	CheckedAt string `json:"checked_at,omitempty"` // RFC3339 timestamp when revocation was checked
	Revoked   bool   `json:"revoked,omitempty"`    // true if certificate is revoked
	Reason    string `json:"reason,omitempty"`     // revocation reason code (e.g., "keyCompromise", "CACompromise", "affiliationChanged")
	Time      string `json:"time,omitempty"`       // RFC3339 timestamp when certificate was revoked (only if revoked)
}

// CertificateInfo contains details for a single certificate in the chain
type CertificateInfo struct {
	Position           int      `json:"position,omitempty"`             // position in chain (0 = leaf)
	Role               string   `json:"role,omitempty"`                 // "leaf" | "intermediate" | "root"
	SubjectCN          string   `json:"subject_cn,omitempty"`           // subject common name
	SubjectDNSNames    []string `json:"subject_dns_names,omitempty"`    // subject alternative names
	IssuerCN           string   `json:"issuer_cn,omitempty"`            // issuer common name
	IssuerDN           string   `json:"issuer_dn,omitempty"`            // full issuer Distinguished Name (DN)
	IssuerFriendlyName string   `json:"issuer_friendly_name,omitempty"` // human-readable CA name (typically from Organization field)
	NotBefore          string   `json:"not_before,omitempty"`           // RFC3339
	NotAfter           string   `json:"not_after,omitempty"`            // RFC3339
	PublicKeyAlgorithm string   `json:"public_key_algorithm,omitempty"` // e.g., "RSA", "ECDSA"
	PublicKeySizeBits  int      `json:"public_key_size_bits,omitempty"` // key size in bits
	CurveName          string   `json:"curve_name,omitempty"`           // elliptic curve name for EC-based keys (e.g., "secp256r1", "secp384r1", "x25519")
	SignatureAlgorithm string   `json:"signature_algorithm,omitempty"`  // e.g., "SHA256-RSA"
	IsCA               bool     `json:"is_ca,omitempty"`                // true if this is a CA certificate (leaf should be false)
	SerialNumber       string   `json:"serial_number,omitempty"`        // certificate serial number
	CertSHA256         string   `json:"cert_sha256,omitempty"`          // SHA-256 hash of entire certificate (DER)
	PubkeySHA256       string   `json:"pubkey_sha256,omitempty"`        // SHA-256 hash of public key
	TbsSHA256          string   `json:"tbs_sha256,omitempty"`           // SHA-256 hash of "to be signed" portion
	IssuerPubkeySHA256 string   `json:"issuer_pubkey_sha256,omitempty"` // SHA-256 hash of issuer's public key (only for non-root certs)
}
