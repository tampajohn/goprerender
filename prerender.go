package prerender

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	e "github.com/jqatampa/gadget-arm/errors"
)

type Options struct {
	PrerenderURL *url.URL
	Token        string
}

type Prerender struct {
	Options Options
}

func NewPrerender(options ...Options) *Prerender {
	var o Options
	defaultServiceURL, _ := url.Parse("https://service.prerender.io/")

	if len(options) == 0 {
		o = Options{
			PrerenderURL: defaultServiceURL,
			Token:        os.Getenv("PRERENDER_TOKEN"),
		}
	} else {
		o = options[0]
	}
	fmt.Println(o.Token)
	return &Prerender{Options: o}
}

func (p *Prerender) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if p.shouldPrerender(r) {
		p.getRenderedPage(rw, r)
		// p.writeRenderedPage(rw, r)
	}

	if next != nil {
		next(rw, r)
	}
}

func (p *Prerender) shouldPrerender(r *http.Request) bool {
	return true
}

func (p *Prerender) buildURL(or *http.Request) string {
	url := p.Options.PrerenderURL

	if !strings.HasSuffix(url.String(), "/") {
		url.Path = url.Path + "/"
	}

	var protocol = or.URL.Scheme

	if cf := or.Header.Get("CF-Visitor"); cf != "" {
		match := cfSchemeRegex.FindStringSubmatch(cf)
		if len(match) > 1 {
			protocol = match[1]
		}
	}

	if fp := or.Header.Get("X-Forwarded-Proto"); fp != "" {
		protocol = strings.Split(fp, ",")[0]
	}

	apiUrl := url.String() + protocol + "://" + or.URL.Host + or.URL.Path + "?" + or.URL.RawQuery

	return apiUrl
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (p *Prerender) getRenderedPage(rw http.ResponseWriter, or *http.Request) {
	//Figure out whether the client accepts gzip responses
	doGzip := strings.Contains(or.Header.Get("Accept-Encoding"), "gzip")

	client := &http.Client{}
	req, err := http.NewRequest("GET", p.buildURL(or), nil)
	e.Check(err)

	req.Header.Set("Accept-Encoding", "gzip")

	res, err := client.Do(req)
	e.Check(err)

	defer res.Body.Close()

	if doGzip && res.Header.Get("Content-Encoding") == "gzip" {
		fmt.Println("Accept and Content gzip")
		io.Copy(rw, res.Body)
	} else if doGzip {
		fmt.Println("Accept gzip content raw")
		gz := gzip.NewWriter(rw)
		io.Copy(gz, res.Body)
		gz.Flush()
	} else if res.Header.Get("Content-Encoding") == "gzip" {
		fmt.Println("Don't accept gzip, content gzip")
		gz, err := gzip.NewReader(res.Body)
		e.Check(err)
		defer gz.Close()
		io.Copy(rw, gz)
	} else {
		// Pass through, gzip/gzip or raw/raw
		fmt.Println("Don't accept gzip, content not gzip")
		io.Copy(rw, res.Body)
	}
}

func (p *Prerender) writeRenderedPage(rw http.ResponseWriter, or *http.Request) {
	//Figure out whether the client accepts gzip responses
	doGzip := strings.Contains(or.Header.Get("Accept-Encoding"), "gzip")

	or.Header.Set("Accept-Encoding", "gzip")
	res, err := http.Get(p.buildURL(or))
	e.Check(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if doGzip {
		rw.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(rw)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: rw}
		gzr.Write(body)
	} else {
		rw.Write(body)
	}

	fmt.Println(string(body))
}
