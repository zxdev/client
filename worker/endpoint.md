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

/* curl "https://w4.netstar.dev/cert/zxdev.com"
{"host":"zxdev.com","data":[{"id":21535723622,"issuer_ca_id":286236,"issuer_name":"C=US, O=Google Trust Services, CN=WE1","common_name":"zxdev.com","name_value":"*.zxdev.com\nzxdev.com","entry_timestamp":"2025-10-07T02:42:38.34","not_before":"2025-10-07T01:42:38","not_after":"2026-01-05T02:40:17","serial_number":"302b013802f0cb5513ad36506b2e5959","result_count":3},{"id":21155961688,"issuer_ca_id":295813,"issuer_name":"C=US, O=Let's Encrypt, CN=E7","common_name":"zxdev.com","name_value":"*.zxdev.com\nzxdev.com","entry_timestamp":"2025-09-20T13:00:44.354","not_before":"2025-09-20T12:02:14","not_after":"2025-12-19T12:02:13","serial_number":"0676e3f05c270fa5af31a3a416ccc12fd074","result_count":3},{"id":21155962885,"issuer_ca_id":295813,"issuer_name":"C=US, O=Let's Encrypt, CN=E7","common_name":"zxdev.com","name_value":"*.zxdev.com\nzxdev.com","entry_timestamp":"2025-09-20T13:00:44.149","not_before":"2025-09-20T12:02:14","not_after":"2025-12-19T12:02:13","serial_number":"0676e3f05c270fa5af31a3a416ccc12fd074","result_count":3},{"id":20600816078,"issuer_ca_id":295817,"issuer_name":"C=US, O=Let's Encrypt, CN=R13","common_name":"zxdev.com","name_value":"www.zxdev.com\nzxdev.com","entry_timestamp":"2025-08-27T04:33:09.947","not_before":"2025-08-27T03:34:39","not_after":"2025-11-25T03:34:38","serial_number":"0603cfef517089fefaa0fca914e8a9d38a59","result_count":3},{"id":20600817466,"issuer_ca_id":295817,"issuer_name":"C=US, O=Let's Encrypt, CN=R13","common_name":"zxdev.com","name_value":"www.zxdev.com\nzxdev.com","entry_timestamp":"2025-08-27T04:33:09.55","not_before":"2025-08-27T03:34:39","not_after":"2025-11-25T03:34:38","serial_number":"0603cfef517089fefaa0fca914e8a9d38a59","result_count":3},{"id":20213881401,"issuer_ca_id":286236,"issuer_name":"C=US, O=Google Trust Services, CN=WE1","common_name":"zxdev.com","name_value":"*.zxdev.com\nzxdev.com","entry_timestamp":"2025-08-09T02:20:22.603","not_before":"2025-08-09T01:20:22","not_after":"2025-11-07T02:18:52","serial_number":"0086a28e859ea6fd09130742534fc2b8e5","result_count":3},{"id":19779099779,"issuer_ca_id":295810,"issuer_name":"C=US, O=Let's Encrypt, CN=E5","common_name":"zxdev.com","name_value":"*.zxdev.com\nzxdev.com","entry_timestamp":"2025-07-19T17:04:47.544","not_before":"2025-07-19T16:06:16","not_after":"2025-10-17T16:06:15","serial_number":"0655c8bb655cfec7b70ba9eeb9059be3f9cd","result_count":3},{"id":19779100583,"issuer_ca_id":295810,"issuer_name":"C=US, O=Let's Encrypt, CN=E5","common_name":"zxdev.com","name_value":"*.zxdev.com\nzxdev.com","entry_timestamp":"2025-07-19T17:04:46.51","not_before":"2025-07-19T16:06:16","not_after":"2025-10-17T16:06:15","serial_number":"0655c8bb655cfec7b70ba9eeb9059be3f9cd","result_count":3}]}
*/

        // drives rdap.net
		rdaper := rdap.NewRDAP(time.Second * 15)
		rx.Get("/rdap/{host}", rdaper.GetHandler())
		rx.Post("/rdap", rdaper.PostHandler(grace.Context()))

/*
curl "https://w4.netstar.dev/rdap/zxdev.com"  
{"host":"zxdev.com","nameserver":["aron.ns.cloudflare.com","marek.ns.cloudflare.com"]}
*/

		titler := title.NewTitle(time.Second * 10) 
		rx.Get("/title/{host}", titler.GetHandler())
		rx.Post("/title", titler.PostHandler(grace.Context()))
/*
curl https://w4.netstar.dev/title/zxdev.com      
{"url":"zxdev.com","title":"Zx Development","hash":966290809148717815}
*/

        // add ?header for full header responses
		hvaler := hval.NewHVAL(time.Second * 10) 
		rx.Get("/hval/{host}", hvaler.GetHandler())
		rx.Post("/hval", hvaler.PostHandler(grace.Context()))

/*
curl "https://w4.netstar.dev/hval/zxdev.com"  
{"item":"zxdev.com","head":[{"status":200,"url":"http://zxdev.com","ip":["2606:50c0:8001::153","2606:50c0:8002::153","2606:50c0:8003::153","2606:50c0:8000::153","185.199.109.153","185.199.108.153","185.199.111.153","185.199.110.153"]}],"n":1,"security":8}

"https://w4.netstar.dev/hval/zxdev.com?header"
{"item":"zxdev.com","head":[{"status":200,"url":"http://zxdev.com","header":{"Accept-Ranges":["bytes"],"Access-Control-Allow-Origin":["*"],"Age":["344"],"Cache-Control":["max-age=600"],"Content-Length":["208"],"Content-Type":["text/html; charset=utf-8"],"Date":["Mon, 13 Oct 2025 04:07:20 GMT"],"Etag":["\"6258d5fb-d0\""],"Expires":["Mon, 13 Oct 2025 03:05:12 GMT"],"Last-Modified":["Fri, 15 Apr 2022 02:18:35 GMT"],"Server":["GitHub.com"],"Vary":["Accept-Encoding"],"Via":["1.1 varnish"],"X-Cache":["HIT"],"X-Cache-Hits":["1"],"X-Fastly-Request-Id":["8bdf1b515061165d4b05107f09d68757f0570ad6"],"X-Github-Request-Id":["E0FC:87506:867A0:95EC3:68EC6A10"],"X-Proxy-Cache":["MISS"],"X-Served-By":["cache-dfw-kdfw8210121-DFW"],"X-Timer":["S1760328440.274719,VS0,VE1"]},"ip":["2606:50c0:8001::153","2606:50c0:8002::153","2606:50c0:8003::153","2606:50c0:8000::153","185.199.110.153","185.199.109.153","185.199.108.153","185.199.111.153"]}],"n":1,"security":8}
mac@zxdev ~ % curl "https://w4.netstar.dev/mail/zxdev.com"   
{"rcode":17,"host":"zxdev.com","mx":["mx01.mail.icloud.com","mx02.mail.icloud.com"],"spf":["v=spf1 include:icloud.com ~all","apple-domain=Ocg7wboGhDipjpNm"]}
mac@zxdev ~ % curl "https://w4.netstar.dev/rdap/zxdev.com"
{"host":"zxdev.com","nameserver":["aron.ns.cloudflare.com","marek.ns.cloudflare.com"]}
*/

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

/*
curl "https://w4.netstar.dev/dns/zxdev.com"          
{"rcode":1,"host":"zxdev.com","a":["185.199.108.153","185.199.111.153","185.199.110.153","185.199.109.153"],"aaaa":["2606:50c0:8001::153","2606:50c0:8002::153","2606:50c0:8003::153","2606:50c0:8000::153"]}

curl "https://w4.netstar.dev/dns/zxdev.com?15"       
{"rcode":15,"host":"zxdev.com","a":["185.199.109.153","185.199.108.153","185.199.111.153","185.199.110.153"],"aaaa":["2606:50c0:8001::153","2606:50c0:8002::153","2606:50c0:8003::153","2606:50c0:8000::153"],"cname":["zxdev.com"],"ns":["aron.ns.cloudflare.com","marek.ns.cloudflare.com"]}

curl "https://w4.netstar.dev/dns/185.199.108.153?128"
{"rcode":128,"host":"185.199.108.153","domain":["cdn-185-199-108-153.github.com."]}
*/

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

/*
 curl "https://w4.netstar.dev/mail/zxdev.com"  
{"rcode":17,"host":"zxdev.com","mx":["mx02.mail.icloud.com","mx01.mail.icloud.com"],"spf":["apple-domain=Ocg7wboGhDipjpNm","v=spf1 include:icloud.com ~all"]}
*/

		methoder := method.NewMethod(time.Second * 7)
		rx.Get("/method/{host}", methoder.GetHandler())
		rx.Post("/method", methoder.PostHandler(grace.Context()))

/*
curl "https://w4.netstar.dev/method/zxdev.com"
{"url":"zxdev.com","flag":3}
*/
```
