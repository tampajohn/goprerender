Prerender Go
===========================

Bots are constantly hitting your site, and a lot of times they're unable to render
javascript.  Prerender.io is awesome, and allows a headless browser to render you
page.  

This middleware allows you to intercept requests from crawlers and route them
to an external Prerender Service to retrieve the static HTML for the requested page.

Prerender adheres to google's `_escaped_fragment_` proposal, which we recommend you use. It's easy:
- Just add &lt;meta name="fragment" content="!"> to the &lt;head> of all of your pages
- If you use hash urls (#), change them to the hash-bang (#!)
- That's it! Perfect SEO on javascript pages.

## Features
I tried to replicate the features found in the [Prerender-node](https://github.com/prerender/prerender-node/)
middleware.

## Using it in [negroni](https://github.com/codegangsta/negroni)
``` go
package main

import (
  "net/http"

  "github.com/codegangsta/negroni"
  "github.com/tampajohn/prerender"
  )

  func main() {
    n := negroni.New()
    n.Use(negroni.NewLogger())
    n.Use(prerender.NewOptions().NewPrerender())
    n.Use(negroni.NewStatic(http.Dir(".")))
    n.Run(":80")
  }


```
... or if you want to use a custom prerender server

``` go
package main

import (
  "net/http"
  "net/url"

  "github.com/codegangsta/negroni"
  "github.com/tampajohn/prerender"
  )

  func main() {
    n := negroni.New()
    n.Use(negroni.NewLogger())
    o := prerender.NewOptions()
    o.PrerenderURL, _ = url.Parse("http://prerender.powerchord.io/")
    n.Use(o.NewPrerender())
    n.Use(negroni.NewStatic(http.Dir(".")))
    n.Run(":80")
  }


  ```
  ... or if you want to use it without negroni
  ``` go
  package main

  import (
    "net/http"

    "github.com/tampajohn/prerender"
    )

    func main() {
      m := http.NewServeMux()
      m.HandleFunc("/", prerender.NewOptions().NewPrerender().PreRenderHandler)
      http.ListenAndServe(":80", m)
    }


```

... or if you want to use it on a single page application
  ``` go
  package main
  
  import (
  	"net/http"
  	"net/url"
  	"os"
  	"path/filepath"
  
  	"github.com/codegangsta/negroni"
  	"github.com/tampajohn/prerender"
  	"github.com/gorilla/mux"
  )
  
  type spaHandler struct {
  	staticPath string
  	indexPath  string
  }
  
  func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  	path, err := filepath.Abs(r.URL.Path)
  	if err != nil {
  		http.Error(w, err.Error(), http.StatusBadRequest)
  		return
  	}
  
  	path = filepath.Join(h.staticPath, path)
  
  	_, err = os.Stat(path)
  	if os.IsNotExist(err) {
  		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
  		return
  	} else if err != nil {
  		http.Error(w, err.Error(), http.StatusInternalServerError)
  		return
  	}
  
  	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
  }
  func main() {
  	n := negroni.New()
  	n.Use(negroni.NewLogger())
  	n.Use(prerender.NewOptions().NewPrerender())
  	r := mux.NewRouter()
  	spa := spaHandler{staticPath: "src", indexPath: "index.html"}
  	r.PathPrefix("/").Handler(spa)
  	n.UseHandler(r)
  	n.Run(":8099")
  }



```

### Special Thanks
I stole almost all of the logic from prerender-node (thanks prerender guys :))

I also want to thank [CodeGangsta](https://github.com/codegangsta) for creating
Negroni and making it so freaking awesome to use.
