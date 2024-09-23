package rss

import (
	"encoding/xml"
	"goNews/pkg/storage"
	"io"
	"net/http"
	"regexp"
	"time"
)

// Набор вложенныъ структур для раскодировки xml rss фида
type Feed struct {
	RSS     string  `xml:"rss"`
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"description"`
	Link    string `xml:"link"`
	PubTime string `xml:"pubDate"`
}

func ParseFeed(url string) ([]storage.News, error) {
	// Сохраняем ответ на запрос по адресу url
	response, err := http.Get(url)
	if err != nil {
		return []storage.News{}, err
	}

	// Читаем тело ответа в массив байт
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return []storage.News{}, err
	}

	var feed Feed
	// Раскодиреум xml в структуру
	err = xml.Unmarshal(b, &feed)
	if err != nil {
		return []storage.News{}, err
	}

	// Регулярное выражение для удаления html тэгов
	const regex = `<.*?>`
	r := regexp.MustCompile(regex)

	var news []storage.News
	// Итерируем по массиву новостей
	for _, item := range feed.Channel.Items {
		var n storage.News
		n.Title = item.Title
		// Удаляем html тэги с помощью регулярного выражения
		n.Content = r.ReplaceAllString(item.Content, "")
		n.Link = item.Link
		// Парсим время публикации по одному формату
		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubTime)
		// Если получаем ошибку, то парсим другой формат
		if err != nil {
			t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 GMT", item.PubTime)
		}
		if err != nil {
			return []storage.News{}, err
		}
		n.PubTime = t.Unix()
		news = append(news, n)
	}
	return news, nil
}
