package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var markdownTemplate, _ = template.New("markdown").Parse(`<!doctype html>
<html>
<head>
	<link rel="shortcut icon"type="image/x-icon" href="data:image/x-icon;,">
	<title>{{.Name}}</title>
	{{- if .CSS}}
	<link rel="stylesheet" href="{{.CSS}}">
	{{- else}}
	<style>
html {
    -ms-text-size-adjust: 100%;
    -webkit-text-size-adjust: 100%;
    font-family: sans-serif;
}

body {
    margin: 0;
    padding: 30px;
    /*min-width: 1020px;*/
    background-color: #fff;
    color: #333333;
    font: 13px Helvetica, arial, freesans, clean, sans-serif;
    line-height: 1.4;
}

    body > *:first-child {
        margin-top: 0 !important;
    }

    body > *:last-child {
        margin-bottom: 0 !important;
    }

a {
    color: #4183c4;
    text-decoration: none;
}

    a:hover {
        outline: 0;
        text-decoration: underline;
    }

    a:focus {
        outline: thin dotted;
        text-decoration: underline;
    }

    a:active {
        outline: 0;
        text-decoration: underline;
    }

h1, h2, h3, h4, h5, h6 {
    position: relative;
    margin: 1em 0 15px;
    padding: 0;
    font-weight: bold;
    line-height: 1.7;
    cursor: text;
}

    /*h1:hover a.anchor, h2:hover a.anchor, h3:hover a.anchor, h4:hover a.anchor, h5:hover a.anchor, h6:hover a.anchor {
        top: 15%;
        margin-left: -30px;
        padding-left: 8px;
        text-decoration: none;
        line-height: 1;
    }*/

    h1 code, h2 code, h3 code, h4 code, h5 code, h6 code {
        font-size: inherit;
    }

h1 {
    border-bottom: 1px solid #ddd;
    font-size: 2.5em;
}

h2 {
    border-bottom: 1px solid #eee;
    font-size: 2em;
}

h3 {
    font-size: 1.5em;
}

h4 {
    font-size: 1.2em;
}

h5 {
    font-size: 1em;
}

h6 {
    color: #777;
    font-size: 1em;
}

b, strong {
    font-weight: bold;
}

hr:before, hr:after {
    display: table;
    content: " ";
}

hr:after {
    clear: both;
}

sub, sup {
    position: relative;
    vertical-align: baseline;
    font-size: 75%;
    line-height: 0;
}

sup {
    top: -0.5em;
}

sub {
    bottom: -0.25em;
}

img {
    -moz-box-sizing: border-box;
    box-sizing: border-box;
    max-width: 100%;
    border: 0;
}

code, pre {
    font-size: 12px;
    font-family: Consolas, "Liberation Mono", Courier, monospace;
}

pre {
    margin-top: 0;
    margin-bottom: 0;
}

a.anchor:focus {
    outline: none;
}

p, blockquote, ul, ol, dl, table, pre {
    margin: 15px 0;
}

ul, ol {
    margin-top: 0;
    margin-bottom: 0;
    padding: 0;
    padding-left: 30px;
}

    ul.no-list, ol.no-list {
        padding: 0;
        list-style-type: none;
    }

    ul ul, ul ol, ol ol, ol ul {
        margin-top: 0;
        margin-bottom: 0;
    }

dl {
    padding: 0;
}

blockquote {
    padding: 0 15px;
    border-left: 4px solid #DDD;
    color: #777;
}

table {
    display: block;
    overflow: auto;
    width: 100%;
	border-collapse: collapse;
}

    table th, table td {
        padding: 6px 13px;
        border: 1px solid #ddd;
    }

    table tr:nth-child(2n) {
        background-color: #f8f8f8;
    }

/* TODO: Something in here is stripping line-breaks */
code {
    display: inline-block;
    overflow: auto;
    margin: 0;
    padding: 0;
    max-width: 100%;
    border: 1px solid #ddd;
    border-radius: 3px;
    background-color: #f8f8f8;
    vertical-align: middle;
    /*white-space: nowrap;*/
    line-height: 1.3;
}

    code:before, code:after {
        content: "\00a0";
        letter-spacing: -0.2em;
    }

pre {
    overflow: auto;
    /*padding: 6px 10px;*/
    padding: 2px 3px;
    border: 1px solid #ddd;
    border-radius: 3px;
    background-color: #f8f8f8;
    word-wrap: normal;
    font-size: 13px;
	line-height: 1.2em;
}

    pre code {
        display: inline;
        overflow: initial;
        margin: 0;
        padding: 0;
        max-width: initial;
        border: none;
        background-color: transparent;
        word-wrap: normal;
        line-height: inherit;
    }

        pre code:before, pre code:after {
            content: normal;
        }

/* breadcrumbs */
nav {
    position: fixed;
    width: 100%;
    top: 0;
    left: 0;
    padding: 0 15px;
    border-bottom: 1px solid rgb(220, 220, 220);
    background: rgb(235, 235, 235);
	color: rgb(120, 120, 120);
	z-index: 9999;
}
nav ol {
    list-style-type: none;
    padding: 0;
}
nav ol li {
    display: inline;
}
nav ol li+li::before {
    content: ' | ';
    margin: 5px;
    
}
nav a:visited {
    color: #4183c4;
}
	</style>{{end}}

	<!--<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.10/styles/default.min.css">
	<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.10/highlight.min.js"></script>-->
</head>
<body>

{{.Content}}

	<script>

~function (d) {
	'use strict';

	//var pres = d.querySelectorAll('div.highlight > pre');
	var codes = d.querySelectorAll('pre > code');
	console.log(codes.length)
	if (codes.length > 0) {
		var css = d.createElement('link'),
			js = d.createElement('script');

		//css.setAttribute('rel', 'stylesheet');
		//css.setAttribute('href', '//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.10/styles/default.min.css');
		//js.setAttribute('src', '//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.10/highlight.min.js');

		css.async = true;
		css.rel = 'stylesheet';
		//css.href = '//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.1/styles/default.min.css';
		css.href = '/_/highlightjs/styles/github-gist.css';
		js.async = true;
		//js.src = '//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.1/highlight.min.js';
		js.src = '/_/highlightjs/highlight.pack.js';
		js.onreadystatechange = js.onload = function () {
			if (!js.readyState || /loaded|complete/.test(js.readyState)) {
				Array.prototype.forEach.call(codes, function (block) { hljs.highlightBlock(block); });
			}
		}

		d.head.appendChild(css);
		d.head.appendChild(js);
	}
}(document);

// breadcrumbs
~function (d) {
	'use strict';

	function pathElem(path, name, isCurrent) {
		var elem = document.createElement('li'),
			text;
		
		if (isCurrent) {
			text = document.createElement('span');
		} else {
			text = document.createElement('a');
			//text.href = path == '/' ? '/index.md' : path + '/index.md';
			text.href = path
		}
		text.textContent = name;
		elem.appendChild(text);

		return elem;
	}

	function breadcrumbElem(paths) {
		var elem = d.createElement('ol'),
			path;

		elem.appendChild(pathElem('/', 'home'));
		for (var i=0; i<paths.length; i++) {
			path = [path, paths[i]].join('/');
			elem.appendChild(pathElem(path, paths[i], i == paths.length-1));
		}

		return elem;
	}

	var paths = d.location.pathname.split('/'),
		breadcrumb = breadcrumbElem(paths.slice(1)),
		nav = d.createElement('nav');
	
	nav.appendChild(breadcrumb);
	d.body.insertBefore(nav, d.body.firstChild);
	//hljs.initHighlightingOnLoad();
	
}(document);

	</script>
</body>
</html>
`)

type variables struct {
	Name, CSS string
	Content   template.HTML
}

func markdownToHTML(file string) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(path.Join(cwd, file))
	if err != nil {
		return nil, err
	}
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Linkify,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithBlockParsers(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	err = md.Convert(data, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func writeMarkDown(w http.ResponseWriter, name, css string, md []byte) error {
	content := template.HTML(md)
	return markdownTemplate.Execute(w, variables{name, css, content})
}

func handleFileError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	status := 500
	if err.Error() == "no such file or directory" {
		status = 404
	}
	http.Error(w, err.Error(), status)
	return true
}

func cleanRequestPath(r *http.Request) string {
	p := r.URL.Path
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
		r.URL.Path = p
	}
	return path.Clean(p)
}

func markdownHandler(root string) http.Handler {
	svr := http.FileServer(http.Dir(root))
	fn := func(w http.ResponseWriter, r *http.Request) {
		p := cleanRequestPath(r)
		_, isRaw := r.URL.Query()["raw"]
		if filepath.Ext(p) != ".md" || isRaw {
			svr.ServeHTTP(w, r)
			return
		}
		md, err := markdownToHTML(p)
		if handleFileError(w, err) {
			return
		}
		css := r.URL.Query().Get("css")
		writeMarkDown(w, p, css, md)
	}
	return http.HandlerFunc(fn)
}

func singlePageHandler(root string) http.Handler {
	indexPath := path.Join(root, "index.html")
	svr := http.FileServer(http.Dir(root))
	fn := func(w http.ResponseWriter, r *http.Request) {
		p := cleanRequestPath(r)
		if filepath.Ext(p) != "" {
			svr.ServeHTTP(w, r)
			return
		}
		b, err := ioutil.ReadFile(indexPath)
		if handleFileError(w, err) {
			return
		}
		w.Header().Set("Content-type", "text/html")
		w.Write(b)
	}
	return http.HandlerFunc(fn)
}

func main() {
	d, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	p := flag.Int("port", 9000, "Port number")
	path := flag.String("path", d, "Content path")
	isSinglePage := flag.Bool("single", false, "Whether this is a single page site or not (eg. react-router)")
	flag.Parse()

	port := fmt.Sprintf(":%d", *p)
	fmt.Printf("static listener on \"%s\" at path \"%s\"\n", port, *path)
	var handler http.Handler
	if *isSinglePage {
		handler = singlePageHandler(*path)
	} else {
		http.Handle("/", markdownHandler(*path))
	}
	err = http.ListenAndServe(port, handler)
	log.Fatalln(err)
}
