package mysql

import (
	"database/sql"
	"webpage/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ?
DAY))`

result, err := m.DB.Exec(stmt,title,content,expires)
if err != nil {
	return 0, err
}

id ,err := result.LastInsertId()
if err != nil{
	return 0,err
}

return int(id),nil
	
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `select id, title,content,created,expires from snippets where expires > UTC_TIMESTAMP() and id = ?`

	row := m.DB.QueryRow(stmt,id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID,&s.Title,&s.Content,&s.Created,&s.Expires)
	if err == sql.ErrNoRows{
		return nil,models.ErrNoRecord
	}else if err != nil{
		return nil , err
	}

	return s ,nil


}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `select id , title, content ,created,expires from snippets where expires > UTC_TIMESTAMP() order by created desc limit 10`

	rows,err := m.DB.Query(stmt)
	if err != nil{
		return nil,err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next(){
		s := &models.Snippet{}

		err = rows.Scan(&s.ID,&s.Title,&s.Content,&s.Created,&s.Expires)
		if err != nil{
			return nil,err
		}
		snippets = append(snippets,s)
		}

		if err = rows.Err();err != nil{
			return nil,err
		}
		return snippets,nil
}

