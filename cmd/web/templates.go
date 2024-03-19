package main

import (
	"html/template"
	"path/filepath"
	"snippetbox/internal/models"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// 형식이 좋은 문자열을 반환하는 humanDate 함수를 만듭니다.
// time.Time 객체의 표현.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// template.FuncMap 객체를 초기화하고 전역 변수에 저장합니다.
// 이는 본질적으로 사용자 정의 템플릿 함수의 이름과 함수 자체 사이를 조회하는 역할을 하는 문자열 키 맵입니다.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// 전체 파일 경로에서 파일 이름(예: 'home.tmpl')을 추출합니다.
		name := filepath.Base(page)

		// ParseFiles() 메서드를 호출하기 전에 template.FuncMap을 템플릿 세트에 등록해야 합니다.
		// 즉, template.New()를 사용하여 빈 템플릿 세트를 만들고,
		// Funcs() 메서드를 사용하여 template.FuncMap을 등록한 다음 파일을 정상적으로 구문 분석해야 합니다.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
