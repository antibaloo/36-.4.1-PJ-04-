package postgres

import (
	"goNews/pkg/storage"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestStore_AddNews(t *testing.T) {
	type args struct {
		n storage.News
	}
	pwd := os.Getenv("postgrespass")
	conString := "postgres://postgres:" + pwd + "@localhost:5432/gonews?sslmode=disable"
	s, err := NewStorage(conString)
	if err != nil {
		t.Errorf("Error while connecting database: %s", err.Error())
	}
	if err = s.Init([]storage.News{}); err != nil {
		t.Errorf("Error while initialization database: %s", err.Error())
	}

	tests := []struct {
		name    string
		s       *Store
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"No error", s, args{storage.News{Title: "first title", Content: "first content", PubTime: 0, Link: "www.firstlink.ru"}}, false},
		{"No pub_time", s, args{storage.News{Title: "second title", Content: "second content", Link: "www.secondlink.ru"}}, false},
		{"No title", s, args{storage.News{Content: "third content", PubTime: 0, Link: "www.thirdlink.ru"}}, true},
		{"No content", s, args{storage.News{Title: "fourth title", PubTime: 0, Link: "www.fourthlink.ru"}}, true},
		{"No link", s, args{storage.News{Title: "fifth title", Content: "fifth content", PubTime: 0}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddNews(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Store.AddNews() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_LastNews(t *testing.T) {
	type args struct {
		q int
	}

	pwd := os.Getenv("postgrespass")
	conString := "postgres://postgres:" + pwd + "@localhost:5432/gonews?sslmode=disable"
	// Инициализируем пустое хранилище
	noContent := []storage.News{}
	s, err := NewStorage(conString)
	if err != nil {
		t.Errorf("при подключении к ДБ произошла ошибка: %s", err.Error())
	}
	if err = s.Init(noContent); err != nil {
		t.Errorf("при инициализации БД произошла ошибка: %s", err.Error())
	}
	var emptyResult []storage.News
	tests := []struct {
		name    string
		s       *Store
		args    args
		want    []storage.News
		wantErr bool
	}{
		// TODO: Add test cases.
		{"empty result", s, args{1}, emptyResult, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.LastNews(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.LastNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.LastNews() = %v, want %v", got, tt.want)
			}
		})
	}
	content := []storage.News{
		{Id: 1, Title: "first title", Content: "first content", PubTime: time.Now().Unix(), Link: "www.firstlink.ru"},
		{Id: 2, Title: "second title", Content: "second content", PubTime: time.Now().Unix(), Link: "www.secondtlink.ru"},
		{Id: 3, Title: "fourth title", Content: "fourth content", PubTime: time.Now().Unix(), Link: "www.thirdlink.ru"},
		{Id: 4, Title: "fourth title", Content: "fourth content", PubTime: time.Now().Unix(), Link: "www.fourthlink.ru"},
		{Id: 5, Title: "fifth title", Content: "fifth content", PubTime: time.Now().Unix(), Link: "www.fifthlink.ru"},
		{Id: 6, Title: "sixth title", Content: "sixth content", PubTime: time.Now().Unix(), Link: "www.sixthlink.ru"},
		{Id: 7, Title: "seventh title", Content: "seventh content", PubTime: time.Now().Unix(), Link: "www.seventhlink.ru"},
	}
	if err = s.Init(content); err != nil {
		t.Errorf("при инициализации БД произошла ошибка: %s", err.Error())
	}

	moreTests := []struct {
		name    string
		s       *Store
		args    args
		want    []storage.News
		wantErr bool
	}{
		// TODO: Add test cases.
		{"last 1", s, args{1}, []storage.News{content[6]}, false},
		{"last 2", s, args{2}, []storage.News{content[6], content[5]}, false},
		{"last 3", s, args{3}, []storage.News{content[6], content[5], content[4]}, false},
		{"last 4", s, args{4}, []storage.News{content[6], content[5], content[4], content[3]}, false},
		{"last 5", s, args{5}, []storage.News{content[6], content[5], content[4], content[3], content[2]}, false},
		{"last 6", s, args{6}, []storage.News{content[6], content[5], content[4], content[3], content[2], content[1]}, false},
		{"last 7", s, args{7}, []storage.News{content[6], content[5], content[4], content[3], content[2], content[1], content[0]}, false},
	}
	for _, tt := range moreTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.LastNews(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.LastNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.LastNews() = %v, want %v", got, tt.want)
			}
		})
	}
}
