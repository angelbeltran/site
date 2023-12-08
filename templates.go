package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"strings"
)

type htmlTemplateServer struct {
	root string
}

func (s *htmlTemplateServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	const htmlTemplateFileExtension = ".html.tmpl"

	p := r.URL.Path
	var name string
	if isDir := p == "" || strings.HasSuffix(p, "/"); !isDir {
		if path.Ext(p) == "tmpl" {
			p = p[:len(p)-5]
		}
		if path.Ext(p) == "html" {
			p = p[:len(p)-5]
		}
		name = p
		p += htmlTemplateFileExtension
		r.URL.Path = p
	}

	bw := newBufferedResponseWriter()
	http.FileServer(http.Dir(s.root)).ServeHTTP(bw, r)
	if c := bw.statusCode; c != 0 && c != http.StatusOK {
		w.WriteHeader(c)
		if bw.buf.Len() > 0 {
			w.Write(bw.buf.Bytes())
		}
		return
	}

	b := bw.buf.Bytes()

	t, err := template.New(name).Parse(string(b))
	if err != nil {
		s.serveErrorf(
			w,
			http.StatusInternalServerError,
			"Failed to parse html template %s: %v",
			name,
			err,
		)
		return
	}

	out := bytes.NewBuffer(nil)
	if err = t.Execute(out, map[string]any{}); err != nil {
		s.serveErrorf(
			w,
			http.StatusInternalServerError,
			"Failed to execute html template %s: %v",
			name,
			err,
		)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(out.Len()))
	io.Copy(w, out)
}

func (s *htmlTemplateServer) serveErrorf(w http.ResponseWriter, statusCode int, msg string, args ...any) {
	http.Error(w, fmt.Sprintf(msg, args...), statusCode)
}

type bufferedResponseWriter struct {
	header     http.Header
	buf        *bytes.Buffer
	statusCode int
}

func newBufferedResponseWriter() *bufferedResponseWriter {
	return &bufferedResponseWriter{
		header: http.Header{},
		buf:    bytes.NewBuffer(nil),
	}
}

func (w *bufferedResponseWriter) Header() http.Header {
	return w.header
}

func (w *bufferedResponseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *bufferedResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
	}
}
