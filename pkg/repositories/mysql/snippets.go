package mysql

import (
	"criozone.net/snippetbox/pkg/domain"
	"database/sql"
)

type SnippetMysqlRep struct {
	DB *sql.DB
}

//goland:noinspection SqlNoDataSourceInspection
func (sr *SnippetMysqlRep) Insert(title, content, expires string) (int, error) {
	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"
	result, err := sr.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sr *SnippetMysqlRep) Get(id int) (*domain.Snippet, error) {
	return nil, nil
}

func (sr *SnippetMysqlRep) Latest() ([]*domain.Snippet, error) {
	return nil, nil
}
