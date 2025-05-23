package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zxdev/passkey"
)

// Job interface for GET Inbox/Outbox expectation and response
type Job interface {
	Request() string // request url|host key
	Unpack() any     // unpack the job
	Okay() bool      // job response result status
}

// Jobs interface for POST Inbox/Outbox expectation and response
type Jobs []Job

func (j *Jobs) Reset() { *j = []Job{} }

// Worker configuration to specify how to interact with the worker cluster
// using generic types with auto selection of GET vs POST methods based on
// the worker.Bulk value setting
type Worker struct {
	Host   string `json:"host,omitempty"`   // host: scheme://host:port
	Secret string `json:"secret,omitempty"` // worker: passKey secret

	Workers       int                 `json:"-"` // number of workers
	Path          string              `json:"-"` // host: endpoint path segment
	Params        string              `json:"-"` // host: endpoint ?param segment
	AuthHeader    func(*http.Request) `json:"_"` // set the auth header
	Client        *http.Client        `json:"-"` // client; default timeout 10-second
	Pacer         time.Duration       `json:"-"` // pacer time delay
	Size          int                 `json:"-"` // GET=0|1 (default); POST>1
	FullURL       bool                `json:"-"` // FullURL flag; set POST:true, Size=1/+ for full url processing
	Inbox, Outbox chan Job            // worker communication channels

	pacer *time.Ticker   // pace control signaler
	jobs  sync.WaitGroup // job state control monitor

}

// Connect configures the Worker and starts listening for jobs on worker.Inbox; the endpoint method
// selection of GET or POST is auto selected based on the worker.Size setting (GET:default) or
// the FullURL:true (POST:default,1) setting. GET methods only support single requests and hostnames
// while POST methods support muliple reqest processing and hostname or full url processing.
//
// For for ease of use all requests via worker.Inbox and
// worker.Outbox are handled individually but not necessarily in a true FIFO ordering
func (w *Worker) Connect(ctx context.Context) *Worker {

	// size assurance
	if w.FullURL && w.Size == 0 || w.Size == 0 {
		//  w.FullURL && w.Size == 1; POST
		//  !w.POST && w.Size == 1; GET
		w.Size++
	}

	// must use POST method for bulk processing
	if w.Size > 1 {
		w.FullURL = true
	}

	// workers assurance and channel configuration
	if w.Workers == 0 {
		w.Workers = 10
	}
	w.Inbox = make(chan Job, w.Workers*3/2)
	w.Outbox = make(chan Job, w.Workers*w.Size*3/2)

	// client with default timeout
	if w.Client == nil {
		w.Client = &http.Client{Timeout: time.Second * 10}
	}

	// configure authentication
	if w.AuthHeader == nil {
		w.AuthHeader = passkey.NewClient(ctx, w.Secret).SetHeader
	}

	// configure pacer
	if w.Pacer == 0 {
		w.Pacer = time.Millisecond * 10 // 100 rps
	}
	w.pacer = time.NewTicker(w.Pacer)

	// configure host/method and ?param assurance
	//  GET  .../method/{host}?{param}
	//  POST .../method?param
	if len(w.Host) == 0 {
		w.Host = "http://localhost:1455"
	}
	if !strings.HasPrefix(w.Host, "http://") && !strings.HasPrefix(w.Host, "https://") {
		w.Host = "http://" + w.Host
	}
	if len(w.Path) > 0 {
		w.Host += "/" + strings.TrimPrefix(w.Path, "/")
	}
	if len(w.Params) > 0 && !strings.HasPrefix(w.Params, "?") {
		w.Params = "?" + w.Params
	}
	w.Host = strings.TrimSuffix(w.Host, "/")

	if !w.FullURL { // GET

		// uses a single item Job object
		w.jobs.Add(w.Workers)
		for range w.Workers {
			go func() {
				for job := range w.Inbox {
					w.get(ctx, job)
				}
				w.jobs.Done()
			}()
		}

	} else { // POST

		// uses a multi item Jobs object
		w.jobs.Add(w.Workers)
		for range w.Workers {
			go func() {
				var jobs Jobs
				for job := range w.Inbox {
					jobs = append(jobs, job)
					if len(jobs) == w.Size {
						w.post(ctx, jobs)
						jobs.Reset()
					}
				}
				if len(jobs) > 0 {
					w.post(ctx, jobs)
				}
				w.jobs.Done()
			}()
		}

	}

	return w
}

// Done shuts down the channels and cleanly exitis
func (w *Worker) Done() {
	close(w.Inbox)
	w.jobs.Wait()
	close(w.Outbox)
	w.pacer.Stop()
}

// GET .../method/{host}?{param}
func (w *Worker) get(ctx context.Context, job Job) {

	w.jobs.Add(1)

	req, _ := http.NewRequest("GET", w.Host+"/"+job.Request()+w.Params, nil)

	w.AuthHeader(req)
	resp, err := w.Client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		json.NewDecoder(resp.Body).Decode(&job)
		<-w.pacer.C
	}

	select {
	case w.Outbox <- job:
	case <-ctx.Done():
	}
	w.jobs.Done()

}

// POST .../method?{param}
func (w *Worker) post(ctx context.Context, jobs Jobs) {

	w.jobs.Add(len(jobs))

	var buf bytes.Buffer
	for i := range jobs {
		buf.WriteString(jobs[i].Request())
		buf.WriteByte(10) // \n
	}

	req, _ := http.NewRequest("POST", w.Host+w.Params, &buf)
	w.AuthHeader(req)
	resp, err := w.Client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		json.NewDecoder(resp.Body).Decode(&jobs)
		<-w.pacer.C
	}

	for i := range jobs {
		select {
		case w.Outbox <- jobs[i]:
		case <-ctx.Done():
		}
		w.jobs.Done()
	}

}
