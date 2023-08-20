package fetch

import (
	"io"
	"net/http"
	"strings"
)

type Options struct {
	cli *http.Client
	m   string //method
	h   http.Header
	b   io.Reader
}

var DefaultOptions = &Options{
	cli: http.DefaultClient,
	h: http.Header{
		"Content-Type": []string{"application/json"},
	},
	m: http.MethodGet,
}

func (o *Options) Client(cli *http.Client) *Options {
	o.cli = cli
	return o
}

func (o *Options) ContentType(contentTyp string) *Options {
	o.h.Set("Content-Type", contentTyp)
	return o
}

func (o *Options) Header(h http.Header) *Options {
	for k, v := range h {
		o.h.Set(k, strings.Join(v, "; "))
	}
	return o
}

func Fetch(url string, callback func(resp *http.Response) error, o *Options) error {
	req, err := http.NewRequest(o.m, url, o.b)
	if err != nil {
		return err
	}
	req.Header = o.h

	resp, err := o.cli.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	return callback(resp)
}
