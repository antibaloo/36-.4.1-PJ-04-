package rss

import (
	"testing"
)

func TestParseRSS(t *testing.T) {
	news, err := ParseFeed("https://habr.com/ru/rss/hubs/go/articles/all/?fl=ru")
	if err != nil {
		t.Errorf("при получения массива новостей произошла ошибка: %s", err.Error())
	}
	t.Logf("получено %d новостей\n%+v", len(news), news)
}

func Test_stripHtmlTags(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"simple", args{"<div>content</div>"}, "content"},
		{"with errors", args{"</div>content"}, "content"},
		{"table", args{"<table><tr><td>1</td><td>2</td></tr><tr><td>3</td><td>4</td></tr></table>"}, "1234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripHtmlTags(tt.args.s); got != tt.want {
				t.Errorf("stripHtmlTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
