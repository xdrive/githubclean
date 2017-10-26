package github_api

import (
	"net/http"
	"testing"
)

func TestNewLinks(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/search/code?q=addClass+user%3Amozilla&page=15>; rel="next", ` +
  				`<https://api.github.com/search/code?q=addClass+user%3Amozilla&page=34>; rel="last", ` +
  				`<https://api.github.com/search/code?q=addClass+user%3Amozilla&page=1>; rel="first", ` +
  				`<https://api.github.com/search/code?q=addClass+user%3Amozilla&page=13>; rel="prev"`,
			},
		},
	}
	l := NewLinks(resp)

	if l.First != 1 {
		t.Error("Expeceted 1, got", l.First)
	}

	if l.Prev != 13 {
		t.Error("Expeceted 13, got", l.Prev)
	}

	if l.Next != 15 {
		t.Error("Expeceted 15, got", l.Next)
	}

	if l.Last != 34 {
		t.Error("Expeceted 34, got", l.Last)
	}
}

func TestNewLinksNoHeader(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	l := NewLinks(resp)

	if l.First != 0 {
		t.Error("Expeceted 0, got", l.First)
	}

	if l.Prev != 0 {
		t.Error("Expeceted 0, got", l.Prev)
	}

	if l.Next != 0 {
		t.Error("Expeceted 0, got", l.Next)
	}

	if l.Last != 0 {
		t.Error("Expeceted 0, got", l.Last)
	}
}

func TestNewLinksNoRelPart(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/search/code?q=addClass+user%3Amozilla&page=15>; rel="next", ` +
					`<https://api.github.com/search/code?q=addClass+user%3Amozilla&page=34>`,
			},
		},
	}
	l := NewLinks(resp)

	if l.First != 0 {
		t.Error("Expeceted 0, got", l.First)
	}

	if l.Prev != 0 {
		t.Error("Expeceted 0, got", l.Prev)
	}

	if l.Next != 15 {
		t.Error("Expeceted 15, got", l.Next)
	}

	if l.Last != 0 {
		t.Error("Expeceted 0, got", l.Last)
	}
}

func TestNewLinksNoPagePArt(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.github.com/search/code?q=addClass+user%3Amozilla&page=15>; rel="next", ` +
					`<https://api.github.com/search/code?q=addClass+user%3Amozilla>; rel="last"`,
			},
		},
	}
	l := NewLinks(resp)

	if l.First != 0 {
		t.Error("Expeceted 0, got", l.First)
	}

	if l.Prev != 0 {
		t.Error("Expeceted 0, got", l.Prev)
	}

	if l.Next != 15 {
		t.Error("Expeceted 15, got", l.Next)
	}

	if l.Last != 0 {
		t.Error("Expeceted 0, got", l.Last)
	}
}