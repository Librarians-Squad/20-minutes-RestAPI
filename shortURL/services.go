package main

import (
	"fmt"
	"time"

	"github.com/speps/go-hashids"
)

type model struct {
	Link string      `json:"link"`
	ShortLink string `json:"short_link"`
	Ttl int64        `json:"ttl"`
}

type Datastore interface {
	NewURL(string) (model, error)
	Get(string) (model, error)
	SetTtl(int64) error
	GetAll() (int, []model, error)
}

func (db *DB)CheckURL(URL string) string {
	var obj model
	row := db.QueryRow("SELECT * FROM link_table WHERE link = $1", URL)
	err := row.Scan(&obj.Link, &obj.ShortLink, &obj.Ttl)
	if err != nil {
		return ""
	}
	return URL
}

func GenerateShortURL() (shortURL string, err error) {
	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	now := time.Now()
	shortURL, err = h.Encode([]int{int(now.Unix())})
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (db *DB)NewURL(URL string) (model, error) {
	shortURL, err := GenerateShortURL()
	if err != nil {
		return model{}, err
	}
	obj := model{}
	obj.ShortLink = shortURL
	obj.Link = db.CheckURL(URL)
	// если url нет в бд
	if obj.Link == "" {
		_, err := db.Query("INSERT INTO link_table (link, short_link, ttl) values ($1, $2, $3)", URL, shortURL, time.Now().Add(2*time.Minute).Unix())
		if err != nil {
			return model{}, err
		}
		return obj, nil
	}
	// если есть
	_, err = db.Exec("UPDATE link_table SET short_link = $1, ttl = $2 where link = $3", shortURL, time.Now().Add(2*time.Minute).Unix(), URL)
	if err != nil {
		return model{}, err
	}
	return obj, nil
}

func (db *DB)Get(shortURL string) (model, error) {
	obj := model{}
	row := db.QueryRow("SELECT * FROM link_table WHERE short_link = $1", shortURL)
	err := row.Scan(&obj.Link, &obj.ShortLink, &obj.Ttl)
	if err != nil {
		return model{}, err
	}
	if obj.Ttl < time.Now().Unix() {
		_, err = db.Exec("DELETE FROM link_table WHERE short_link = $1", shortURL)
		if err != nil {
			return model{}, err
		}
		return model{}, fmt.Errorf("Time is up ")
	}
	_, err = db.Exec("UPDATE link_table SET ttl = $1 where short_link = $2", time.Now().Add(2*time.Minute).Unix(), shortURL)
	if err != nil {
		return model{}, err
	}
	return obj, nil
}

func (db *DB)SetTtl(ttl int64) (err error) {
	_, err = db.Exec("UPDATE link_table SET ttl = $1", time.Now().Unix() + ttl)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB)GetAll() (total int, links []model, err error) {
	total = 0
	rows, err  := db.Query("SELECT link, short_link, ttl FROM link_table")
	if err != nil {
		return total, []model{}, err
	}

	links = make([]model, 0)
	for rows.Next() {
		link := model{}
		err = rows.Scan(&link.Link, &link.ShortLink, &link.Ttl)
		if err != nil {
			return total, []model{}, err
		}
		links = append(links, link)
		total++
	}
	return total, links, nil
}
