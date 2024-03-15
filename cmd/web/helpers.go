package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError 도우미는 오류 메시지와 스택 추적을 errorLog에 기록한 다음
// 일반 500 내부 서버 오류 응답을 사용자에게 보냅니다.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError 도우미는 특정 상태 코드와 해당 설명을 사용자에게 보냅니다.
// 사용자가 보낸 요청에 문제가 있을 때 이 책의 뒷부분에서 400"Bad Request"와 같은 응답을 보내는 데 사용할 것입니다.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 일관성을 위해 notFound 도우미도 구현합니다.
// 이는 사용자에게 404 Not Found 응답을 보내는 clientError에 대한 편리한 래퍼입니다.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}
