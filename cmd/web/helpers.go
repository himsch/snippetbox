package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
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

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// http.ResponseWriter에 직접 작성하는 대신 버퍼에 템플릿을 작성합니다.
	// 오류가 있으면 serverError() 도우미를 호출한 다음 반환하세요.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 템플릿이 오류 없이 버퍼에 기록되면 계속 진행하여 HTTP 상태 코드를 http.ResponseWriter에 작성해도 안전합니다.
	w.WriteHeader(status)

	// 버퍼의 내용을 http.ResponseWriter에 씁니다.
	// 참고: 이번에는 io.Writer를 사용하는 함수에 http.ResponseWriter를 전달하는 또 다른 시간입니다.
	buf.WriteTo(w)
}

// 포인터 templateData를 반환하는 newTemplateData() 도우미를 만듭니다.
// 현재 연도로 초기화된 구조체입니다.
// *http.Request 매개변수는 현재 여기에 있지만 나중에 수행하겠습니다.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// 잘못된 대상 대상을 사용하려고 하면 Decode() 메서드는
		// *form.InvalidDecoderError 유형의 오류를 반환합니다.
		// 우리는 오류를 반환하는 대신 오류를 확인하고 패닉을 발생시키기 위해 errors.As()를 사용합니다.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
