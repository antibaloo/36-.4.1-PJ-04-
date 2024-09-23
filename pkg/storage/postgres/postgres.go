package postgres

import (
	"context"
	"goNews/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	Pool *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func NewStorage(conStr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), conStr)
	if err != nil {
		return nil, err
	}
	s := Store{
		Pool: db,
	}
	return &s, nil
}

// Создает таблицу и заполняет ее тестовыми данными
func (s *Store) Init(news []storage.News) error {
	initRequest := `
		DROP TABLE IF EXISTS news;
		CREATE TABLE news (
    		id SERIAL PRIMARY KEY,
    		title TEXT NOT NULL CHECK (title <> ''),
    		content TEXT NOT NULL CHECK (content <> ''),
    		pub_time INTEGER DEFAULT 0,
    		link TEXT NOT NULL UNIQUE CHECK (link <> '')
	);`
	// Пересоздаем таблицу
	_, err := s.Pool.Exec(context.Background(), initRequest)
	if err != nil {
		return err
	}
	// Итерируем по массиву новостей и добавляем их в таблицу
	for _, n := range news {
		if err := s.AddNews(n); err != nil {
			return err
		}
	}
	return nil
}

// Добавляет в БД запись с новостью из переданной структуры
func (s *Store) AddNews(n storage.News) error {
	_, err := s.Pool.Exec(
		context.Background(),
		`INSERT INTO news (title, content, pub_time, link) VALUES ($1, $2, $3, $4)`,
		n.Title,
		n.Content,
		n.PubTime,
		n.Link,
	)
	if err != nil {
		return err
	}
	return nil
}

// Возвращает последние q (quantity) новостей. quantity - количество
func (s *Store) LastNews(q int) ([]storage.News, error) {
	var news []storage.News
	rows, err := s.Pool.Query(
		context.Background(),
		`SELECT id, title, content, pub_time, link FROM news ORDER BY id DESC LIMIT $1`,
		q,
	)
	if err != nil {
		return news, err
	}
	for rows.Next() {
		var n storage.News
		err = rows.Scan(
			&n.Id,
			&n.Title,
			&n.Content,
			&n.PubTime,
			&n.Link,
		)
		if err != nil {
			return news, err
		}
		news = append(news, n)
	}
	return news, nil
}
