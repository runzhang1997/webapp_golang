// main.go

package main

import (
	"testing"
)

func Test_extractHtmlVersionFromDoctype(t *testing.T) {
	type args struct {
		docType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Test_extractHtmlVersionFromDoctype_1",
			args{"HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd"},
			"HTML 4.01",
		},
		{
			"Test_extractHtmlVersionFromDoctype_2",
			args{"html"},
			"HTML 5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractHtmlVersionFromDoctype(tt.args.docType); got != tt.want {
				t.Errorf("extractHtmlVersionFromDoctype() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isexternallinks(t *testing.T) {
	type args struct {
		link     string
		inputURL string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Test_isexternallinks_1",
			args{"https://godoc.org/golang.org/x/net/html", "https://godoc.org/golang.org/x/"},
			false,
		},
		{
			"Test_isexternallinks_2",
			args{"#section1", "https://godoc.org/golang.org/x/"},
			false,
		},
		{
			"test2",
			args{"https://www.youtube.com/", "https://godoc.org/golang.org/x/"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isexternallinks(tt.args.link, tt.args.inputURL); got != tt.want {
				t.Errorf("isexternallinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findElementID(t *testing.T) {
	type args struct {
		url string
		id  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Test_findElementID_1",
			args{"https://www.w3.org/TR/html401/sgml/dtd.html", "#HTMLsymbol"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findElementID(tt.args.url, tt.args.id); got != tt.want {
				t.Errorf("findElementID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containsLoginForm(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Contains Login Form Test 1:",
			args{"https://www.home24.de/meinkundenkonto/kunde/login"},
			true,
		},
		{
			"Contains Login Form Test 2:",
			args{"https://www.home24.de/?gclid=Cj0KCQjw28T8BRDbARIsAEOMBczxN-PHwUjC515w2roln2-Bgr5Mmybnqpzuc4ldroLnm07knzKuDWQaAm0-EALw_wcB"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsLoginForm(tt.args.url); got != tt.want {
				t.Errorf("containsLoginForm() = %v, want %v", got, tt.want)
			}
		})
	}
}
