package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) 
			VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires 
			FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	// 빈 값인 새로운 Snippet 구조체에 대한 포인터를 초기화합니다.
	s := &Snippet{}

	// row.Scan()을 사용하여 sql.Row의 각 필드 값을
	// Snippet 구조체의 해당 필드입니다. 인수는 다음과 같습니다.
	// row.Scan은 데이터를 복사하려는 위치에 대한 *포인터*입니다.
	// 인수 개수는 인수 개수와 정확히 같아야 합니다.
	// 명령문에서 반환된 열입니다.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
			FROM snippets 
			WHERE expires > UTC_TIMESTAMP() 
			ORDER BY id DESC
			LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// sql.Rows 결과 집합이 Latest() 메서드가 반환되기 전에 항상 적절하게 닫히도록
	// rows.Close()를 defer합니다.
	// 이 defer 문은 Query() 메서드에서 오류를 확인한 *후에* 와야 합니다.
	// 그렇지 않고 Query()가 오류를 반환하면 nil 결과 집합을 닫으려고 할 때 패닉이 발생합니다.
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// rows.Next() 루프가 끝나면 rows.Err()을 호출하여 검색합니다.
	// *반복하는 동안 발생한 오류*입니다.
	// 성공적인 반복이 완료되었다고 가정하지 마세요.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
