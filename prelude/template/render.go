package template

import (
	"bytes"
	"html/template"
	"net/http"
	"sync"
)

type RenderFunc func(rw http.ResponseWriter, r *http.Request, data interface{})

func Render(filenames ...string) (RenderFunc, error) {
	var (
		init sync.Once

		tpl *template.Template
		buf *bytes.Buffer
		err error
	)

	init.Do(func() { tpl, err = template.ParseFiles(filenames...) })

	var f = func(rw http.ResponseWriter, r *http.Request, data interface{}) {
		buf = &bytes.Buffer{}
		err = tpl.Execute(buf, data)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(rw)
	}

	return f, err
}
