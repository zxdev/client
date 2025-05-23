# worker | remote server driver

All of the worker ```client``` methods and supporting types are made available under the client folder and the server configuration is also easily configured via ```env.Conf``` and should be a json object with ```json:Host``` and ```json:secret``` configured for default machine-to-machine passkey authentication.

```golang
// worker server configuration
var work worker.Worker
env.Conf(&work, "path/to/dev.worker.json")
```


The ```client.Worker``` is a generic driver suite that can be used to drive the all remote server worker processes, the ```client_test.go``` demonstrates how to do this.

The client supports ```GET``` or bulk lookup ```POST``` methods against the cluster which are automatically determined based on the ```worker.Size``` setting or the```worker.FullURL``` setting. Basically, GET mdethod only support hostnamr or IP address configurations while POST can support the same as well as include ports and paths.

Since this is build around using generic interfaces for simplicity, it requires simple type casting for handling the server responses.


```/etc/dev.worker.json```

```json
{
    "host":"127.0.0.1:1455",
    "secret":"OOIPTYIG6NZ4BMCRLKNPC54XYQ4UZ3W4"
}
```

example job.Title process configuration and logic flow


```golang

	// configure the generic client.Worker and populate the Host 
	// and Secret fields in client.Worker struct using a basic 
	// default configuration resource disk file and then connect
	// to the remote server cluster
	var work = client.Worker{
		Size:       5, 
		Path:       "title",
	}
	env.Conf(&work,"/etc/dev.worker.json") 
	work.Connect(ctx)
	if work == nil {
		// failed; Host was not configured
		return
	}

	// submit job requests via the spawner worker.Inbox channel
	// and loop until finished; work.Done() signals no more 
	// items to spawn for processing

	items := []string{"one.com", "two.com", "three.com", "zxdev.com", "example.com"}

	// spawner; items for processing
	go func() {
		defer work.Done()
		for i := range items {
			work.Inbox <- &job.Title{Url:items[i]}
			// work.Inbox <- job.NewTitle(items[i])
		}
	}()

	// listen for server job responses on the work.Outbox receiver channel
	// and loop until all jobs complete; valid responses require a simple
	// type cast to degenericize the server interface responses

	for j := range work.Outbox {
		r := j.Response().(job.Title)
		if r.Ok {
			t.Log(r.Url, r.Title)
		}
	}
	

```

