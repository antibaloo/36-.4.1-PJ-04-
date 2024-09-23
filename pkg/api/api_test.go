package api

import (
	"encoding/json"
	"goNews/pkg/storage"
	"goNews/pkg/storage/postgres"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestAPI_lastNews(t *testing.T) {
	pwd := os.Getenv("postgrespass")
	conString := "postgres://postgres:" + pwd + "@localhost:5432/gonews?sslmode=disable"
	db, err := postgres.NewStorage(conString)
	if err != nil {
		t.Errorf("произошла ошибка при создание хранилища: %s", err.Error())
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
	if err := db.Init(content); err != nil {
		t.Errorf("произошла ошибка при инициализации хранилища: %s", err.Error())
	}
	api := NewAPI(db)
	// Создаем запрос для проверки обработчика
	req := httptest.NewRequest(http.MethodGet, "/news/7", nil)
	// Создаем объект для записи ответа обработчика
	rr := httptest.NewRecorder()

	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать ответ сервера: %v", err)
	}
	// Раскодируем тело ответа
	var data []storage.News
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Errorf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что в массиве ровно один элемент.
	const wantLen = 7
	if len(data) != wantLen {
		t.Errorf("получено %d записей, ожидалось %d", len(data), wantLen)
	}
	reversedContent := []storage.News{}
	for i := len(content) - 1; i >= 0; i-- {
		reversedContent = append(reversedContent, content[i])
	}
	if !reflect.DeepEqual(data, reversedContent) {
		t.Errorf("lastNews(7) = %v, want %v", data, reversedContent)
	}
}
