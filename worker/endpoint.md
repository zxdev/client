# worker - endpoint

* GET /{method}/{host} 
    * supports single hostname level submissions

* POST /{method} 
    * supports bulk and full url submission and accepts \n delimited as well as json:Host or json:Url


```golang
        // drives crt.sh
        certer := cert.NewCRTsh(time.Second * 20)
		rx.Get("/cert/{host}", certer.GetHandler())
		rx.Post("/cert", certer.PostHandler(grace.Context()))

        // drives rdap.net
		rdaper := rdap.NewRDAP(time.Second * 15)
		rx.Get("/rdap/{host}", rdaper.GetHandler())
		rx.Post("/rdap", rdaper.PostHandler(grace.Context()))

		titler := title.NewTitle(time.Second * 10) 
		rx.Get("/title/{host}", titler.GetHandler())
		rx.Post("/title", titler.PostHandler(grace.Context()))

        // add ?header for full header responses
		hvaler := hval.NewHVAL(time.Second * 10) 
		rx.Get("/hval/{host}", hvaler.GetHandler())
		rx.Post("/hval", hvaler.PostHandler(grace.Context()))

        // add ?{shortCode || A & AAAA & CNAME & NS & MX & TXT & DOMAIN}
        //	  0b 0 0 0 0 0 0 0 0
        //	     |   | | | | | 1 = A
        //	     |   | | | | 2   = AAAA
        //	     |   | | | 4     = CNAME
        //	     |   | | 8       = NS
        //	     |   | 16        = MX
        //	     |   32          = TXT
        //		 128             = DOMAIN (ip only)
        // ex. ?15 or ?A&AAAA&CNAME&NS for A,AAAA,CNAME,NS
        // ex. ?128 or ?DOMAIN for reverse DNS
		dnser := dns.NewDNS(time.Millisecond) 
		rx.Get("/dns/{host}", dnser.GetHandler())
		rx.Post("/dns", dnser.PostHandler())

        // add ?{shortCode || spf &dmarc &dkim &bimi &mx}
        //	0b 0 0 0 0 0 0 0 0
        //	   | | | | 1 = SPF
        //	   | | | 2   = DMARC
        //	   | | 4     = DKIM={selector}
        //	   | 8       = BIMI={selector}
        //	   16        = MX
        // ex ?19 or ?SPF&DMARC&MX for SPF,DMARC,MX
		mailer := mail.NewMail(time.Second * 10)
		rx.Get("/mail/{host}", mailer.GetHandler())
		rx.Post("/mail/", mailer.PostHandler())

		methoder := method.NewMethod(time.Second * 7)
		rx.Get("/method/{host}", methoder.GetHandler())
		rx.Post("/method", methoder.PostHandler(grace.Context()))
```
