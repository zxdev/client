package client_test

import (
	"testing"

	"github.com/zxdev/client/worker/client"
	"github.com/zxdev/client/worker/job"
	"github.com/zxdev/passkey"
)

var (

	// note: for testing we have to manually set the
	// HOST and SECRET value; these are normally loaded
	// via env.Conf(&work,path) from the resouce file

	host   = "http://localhost:1455"            // requires localhost server instance runnning
	secret = "OOIPTYIG6NZ4BMCRLKNPC54XYQ4UZ3W4" // required but is ignored by localhost server
)

/*

	// PRODUCTION EXAMPLES

	// GET, 1x per reqeust; only Host supported
	var work = client.Worker{
			Path:   "title",
		}
	env.Conf(&work,"/etc/dev.worker.json")
	work.Connect(t.Context())

	// POST, 5x per request; full Url supported
	var work = client.Worker{
			Size:	5,
			Path:   "title",
		}
	env.Conf(&work,"/etc/dev.worker.json")
	work.Connect(t.Context())

	// POST, 1x per request; full Url supported
	var work = client.Worker{
			POST:	true,
			Path:   "title",
		}
	env.Conf(&work,"/etc/dev.worker.json")
	work.Connect(t.Context())

*/

// go test -v client/client_test.go --run=TITLE
func TestTITLE(t *testing.T) {

	// note: for testing, we have to set the Host and AuthHeader

	items := []string{"one.com", "two.com", "three.com", "zxdev.com", "example.com"}

	// test GET and POST via a loop
	for size := range 2 {

		t.Log("== Title", size+1, "==")
		var work = client.Worker{
			Host:       host,                                             // testing
			AuthHeader: passkey.NewClient(t.Context(), secret).SetHeader, // testing
			Size:       size + 1,
			Path:       "title",
		}
		work.Connect(t.Context())

		// submit job requests on the worker.Inbox channel
		// and loop until finished; work.Done() signals no
		// more items to consume

		go func() {
			defer work.Done()
			for i := range items {
				//work.Inbox <- title.Job(items[i]) //title.Job(items[i])
				work.Inbox <- job.NewTitle(items[i]) //&job.Title{Url: items[i]}
			}
		}()

		// listen for result jobs on the work.Outbox channel
		// and loop until jobs complete; valid responses must
		// be type cast for access

		for j := range work.Outbox {
			if j.Okay() { // checks status == 0
				r := j.Unpack().(job.Title) // type case resposne
				t.Log(r.Url, r.Title)       // use results
			}
		}

	}
}

// go test -v client/client_test.go --run=TestRDAP
func TestRDAP(t *testing.T) {

	items := []string{"one.com", "two.com", "three.com", "zxdev.com", "example.com"}

	// test GET and POST via a loop
	for size := range 2 {

		t.Log("== Rdap", size+1, "==")
		var work = client.Worker{
			Host:       host,                                             // testing
			AuthHeader: passkey.NewClient(t.Context(), secret).SetHeader, // testing
			Size:       size + 1,
			Path:       "rdap",
		}
		work.Connect(t.Context())

		// submit job requests on the worker.Inbox channel
		// and loop until finished; work.Done() signals no
		// more items to consume

		go func() {
			defer work.Done()
			for i := range items {
				work.Inbox <- job.NewRdap(items[i])
			}
		}()

		// listen for result jobs on the work.Outbox channel
		// and loop until jobs complete; valid responses
		// are type cast and processed

		for j := range work.Outbox {
			if j.Okay() {
				r := j.Unpack().(job.Rdap)
				t.Log(r.Host, r.NRD, r.NameServer)
			}

		}
	}
}

// go test -v client/client_test.go --run=METHOD
func TestMETHOD(t *testing.T) {

	items := []string{"one.com", "two.com", "three.com", "zxdev.com", "example.com"}

	// test GET and POST via a loop
	for size := range 2 {

		t.Log("== Method", size+1, "==")
		var work = client.Worker{
			Host:       host,                                             // testing
			AuthHeader: passkey.NewClient(t.Context(), secret).SetHeader, // testing
			Size:       size + 1,
			Path:       "method",
		}
		work.Connect(t.Context())

		// submit job requests on the worker.Inbox channel
		// and loop until finished; work.Done() signals no
		// more items to consume

		go func() {
			defer work.Done()
			for i := range items {
				work.Inbox <- job.NewMethod(items[i])
			}
		}()

		// listen for result jobs on the work.Outbox channel
		// and loop until jobs complete; valid responses
		// are type cast and processed

		for j := range work.Outbox {
			if j.Okay() {
				r := j.Unpack().(job.Method)
				t.Log(r.Url, r.Flag, r.Options)
			}
		}
	}
}

// go test -v client/client_test.go -run=HVAL
func TestHVAL(t *testing.T) {

	items := []string{"one.com", "two.com", "walmart.com", "zxdev.com", "example.com"}

	// test GET and POST via a loop
	for size := range 2 {

		t.Log("== HVAL", size+1, "==")
		var work = client.Worker{
			Host:       host,                                             // testing
			AuthHeader: passkey.NewClient(t.Context(), secret).SetHeader, // testing
			Size:       size + 1,
			Path:       "hval",
		}
		work.Connect(t.Context())

		// submit job requests on the worker.Inbox channel
		// and loop until finished; work.Done() signals no
		// more items to consume

		go func() {
			defer work.Done()
			for i := range items {
				work.Inbox <- job.NewHval(items[i])
			}
		}()

		// listen for result jobs on the work.Outbox channel
		// and loop until jobs complete; valid responses
		// are type cast and processed

		// listen for result jobs on the wg.Outbox channel
		for j := range work.Outbox {
			if j.Okay() {
				r := j.Unpack().(job.Hval)
				t.Log(r.Item, r.N, r.Security, r.Head)

			}

		}
	}
}

// go test -v client/client_test.go --run=DNS
func TestDNS(t *testing.T) {

	items := []string{"one.com", "two.com", "three.com", "zxdev.com", "example.com"}

	// test GET and POST via a loop
	for size := range 2 {

		t.Log("== DNS", size+1, "==")
		var work = client.Worker{
			Host:       host,                                             // testing
			AuthHeader: passkey.NewClient(t.Context(), secret).SetHeader, // testing
			Size:       size + 1,
			Path:       "dns",
			Params:     "15", // A,AAAA,CNAME,NS; quick code
		}
		work.Connect(t.Context())

		// submit job requests on the worker.Inbox channel
		// and loop until finished; work.Done() signals no
		// more items to consume

		go func() {
			defer work.Done()
			for i := range items {
				work.Inbox <- job.NewDNS(items[i])
			}
		}()

		// listen for result jobs on the work.Outbox channel
		// and loop until jobs complete; valid responses
		// are type cast and processed

		for j := range work.Outbox {
			if j.Okay() {
				r := j.Unpack().(job.DNS)
				t.Log(r.Host, r.A, r.AAAA, r.CNAME, r.NS, r.Domain)
			}
		}

	}

}

// go test -v client/client_test.go --run=RDNS
func TestRDNS(t *testing.T) {

	items := []string{"1.1.1.1", "8.8.8.8", "1.164.99.136", "185.199.111.153", " 185.199.109.153"}

	// test GET and POST via a loop
	for size := range 2 {

		t.Log("== RDNS", size+1, "==")
		var work = client.Worker{
			Host:       host,                                             // testing
			AuthHeader: passkey.NewClient(t.Context(), secret).SetHeader, // testing
			Size:       size + 1,
			Path:       "dns",
			Params:     "128", // reverse DNS quick code
		}
		work.Connect(t.Context())

		// submit job requests on the worker.Inbox channel
		// and loop until finished; work.Done() signals no
		// more items to consume

		go func() {
			defer work.Done()
			for i := range items {
				work.Inbox <- job.NewDNS(items[i])
			}
		}()

		// listen for result jobs on the work.Outbox channel
		// and loop until jobs complete; valid responses
		// are type cast and processed

		for j := range work.Outbox {
			if j.Okay() {
				r := j.Unpack().(job.DNS)
				t.Log(r.Host, r.Domain)
			}
		}
	}

}

// go test -v client/client_test.go --run=CERT
func TestCERT(t *testing.T) {

	items := []string{"one.com", "two.com", "three.com", "zxdev.com", "example.com"}

	// test GET and POST via a loop
	for size := range 2 {

		t.Log("== Cert", size+1, "==")
		var work = client.Worker{
			Host:       host,                                             // testing
			AuthHeader: passkey.NewClient(t.Context(), secret).SetHeader, // testing
			Size:       size + 1,
			Path:       "cert",
		}
		work.Connect(t.Context())

		// submit job requests on the worker.Inbox channel
		// and loop until finished; work.Done() signals no
		// more items to consume

		go func() {
			defer work.Done()
			for i := range items {
				work.Inbox <- job.NewCert(items[i])
			}
		}()

		// listen for result jobs on the work.Outbox channel
		// and loop until jobs complete; valid responses
		// are type cast and processed

		for j := range work.Outbox {
			if j.Okay() {
				r := j.Unpack().(job.Cert)
				if len(r.Certificate) == 0 {
					t.Log(r.Host, len(r.Certificate))
				} else {
					entry, notBefore, notAfter := job.CertificateTimestamps(&r.Certificate[0])
					t.Log(r.Host, entry, notBefore, notAfter)
				}
			}
		}
	}
}
