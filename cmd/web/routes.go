package main

import (
	"net/http"
	"snippetbox/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// ui.Files 내장 파일 시스템을 가져와서 http.FS 유형으로 변환합니다.
	// http.FileSystem 인터페이스를 만족시킵니다. 그런 다음 이를 전달합니다.
	// http.FileServer() 함수를 사용하여 파일 서버 핸들러를 생성합니다.
	fileServer := http.FileServer(http.FS(ui.Files))

	// 정적 파일은 ui.Files 내장 파일 시스템의 "static" 폴더에 포함되어 있습니다.
	// 예를 들어 CSS 스타일시트는 "static/css/main.css"에 있습니다.
	// 이는 이제 요청 URL(/static/으로 시작하는 모든 요청)에서 접두사를 더 이상 제거해야 함을 의미합니다.
	// 파일 서버로 직접 전달될 수 있으며 해당 정적 파일이 존재하는 한 제공됩니다.
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
