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
