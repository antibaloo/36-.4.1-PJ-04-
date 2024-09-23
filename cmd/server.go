package main

import (
	"encoding/json"
	"fmt"
	"goNews/pkg/api"
	"goNews/pkg/rss"
	"goNews/pkg/storage"
	"goNews/pkg/storage/postgres"
	"net/http"
	"os"
	"time"
)

type config struct {
	URLs          []string      `json:"rss"`
	RequestPeriod time.Duration `json:"request_period"`
}

// Метод чтения фида новостей, заданного ссылкой. Добполнительно переданы каналы для передачи полученных новостей и
// ошибок в работе метода, также в качестве параметра передается период получения данных из фида
func readNews(url string, newsCh chan<- []storage.News, errorsCh chan<- error, period time.Duration) {
	fmt.Printf("%v: чтение новостей из канала %s начато\n", time.Now().Format("02.01.2005 15:04:05 GMT"), url)

	for {
		news, err := rss.ParseFeed(url)
		if err != nil {
			errorsCh <- fmt.Errorf("адрес фида: %s, описание: %s", url, err.Error())
			continue
		}
		newsCh <- news
		fmt.Printf("%v: получено %d новостей из фида с адресом: %s \n", time.Now().Format("02.01.2005 15:04:05 GMT"), len(news), url)
		time.Sleep(time.Minute * period)
	}
}

func main() {
	// чтение и раскодирование файла конфигурации
	b, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("%v: при чтении файла конфигурации произошла ошибка: %s", time.Now().Format("02.01.2005 15:04:05 GMT"), err.Error())
		return
	}

	fmt.Printf("%v файл конфигурации прочитан\n", time.Now().Format("02.01.2005 15:04:05 GMT"))

	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		fmt.Printf("%v: при раскодировании файла конфигурации произошла ошибка: %s", time.Now().Format("02.01.2005 15:04:05 GMT"), err.Error())
		return
	}

	fmt.Printf("%v: файл конфигурации раскодирован\n", time.Now().Format("02.01.2005 15:04:05 GMT"))

	pwd := os.Getenv("postgrespass")
	conString := "postgres://postgres:" + pwd + "@localhost:5432/gonews?sslmode=disable"
	// Подключаемся к БД
	db, err := postgres.NewStorage(conString)
	if err != nil {
		fmt.Printf("%v: при подключении к БД произошла ошибка: %s", time.Now().Format("02.01.2005 15:04:05 GMT"), err.Error())
		return
	}

	fmt.Printf("%v: подключение к БД создано\n", time.Now().Format("02.01.2005 15:04:05 GMT"))

	// Канал для передачи в обработчик списка полученных новостей
	newsCh := make(chan []storage.News)
	// Канал для передачи в обработчик ошибок, при чтении фидов
	errorsCh := make(chan error)
	// Каждый фид, заданный конфигорационным файлом читается в отдельной горутине
	for _, url := range config.URLs {
		go readNews(url, newsCh, errorsCh, config.RequestPeriod)
	}

	// Обработчик канала новостей
	go func() {
		for news := range newsCh {
			for _, n := range news {
				err := db.AddNews(n)
				if err != nil {
					fmt.Printf("%v: при записи полученной новости в БД произошла ошибка: %s\n", time.Now().Format("02.01.2005 15:04:05 GMT"), err.Error())
				}
			}
		}
	}()

	// Обработчик канала ошибок
	go func() {
		for err := range errorsCh {
			fmt.Printf("%v: при чтении фида новостей произошла ошибка: %s\n", time.Now().Format("02.01.2005 15:04:05 GMT"), err.Error())
		}
	}()

	// Создаем экземпляр API и запускаем веб-сервера с API и приложением
	api := api.NewAPI(db)

	fmt.Printf("%v: запускаем веб-сервер\n", time.Now().Format("02.01.2005 15:04:05 GMT"))

	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		fmt.Printf("%v: при запуске веб-сервера произошла ошибка: %s", time.Now().Format("02.01.2005 15:04:05 GMT"), err.Error())
		return
	}
}
