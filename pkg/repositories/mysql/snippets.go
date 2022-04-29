package mysql

import (
	"criozone.net/snippetbox/pkg/domain"
	"database/sql"
	"errors"
	"log"
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

//goland:noinspection SqlNoDataSourceInspection
func (sr *SnippetMysqlRep) Get(id int) (*domain.Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"
	row := sr.DB.QueryRow(stmt, id)
	s := &domain.Snippet{}
	err := row.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

//goland:noinspection SqlNoDataSourceInspection
func (sr *SnippetMysqlRep) Latest() ([]*domain.Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"
	rows, err := sr.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	snippets := make([]*domain.Snippet, 0, 10)
	for rows.Next() {
		s := &domain.Snippet{}
		err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
