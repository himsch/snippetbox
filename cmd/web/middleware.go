package main

import (
	"context"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Go가 스택을 풀면서 패닉이 발생할 경우 항상 실행되는 지연된 함수를 만듭니다.
		defer func() {
			// 내장된 recover 기능을 사용하여 패닉이 발생했는지 여부를 확인하세요. 만약 있다면...
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// 그렇지 않으면 인증이 필요한 페이지가 사용자 브라우저 캐시(다른 중간 캐시)에 저장되지 않도록
		// "Cache-Control: no-store" 헤더를 설정합니다.
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GetInt() 메서드를 사용하여 세션에서 authenticatedUserID 값을 검색합니다.
		// 세션에 "authenticatedUserID" 값이 없으면 int(0)에 대해 0 값을 반환합니다.
		// 이 경우 체인의 다음 핸들러를 정상적으로 호출하고 반환합니다.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// 일치하는 사용자가 발견되면 해당 요청이 데이터베이스에 존재하는 인증된 사용자로부터 온 것임을 알 수 있습니다.
		// 요청의 새 복사본을 생성하고(요청 컨텍스트에서 isAuthenticatedContextKey 값이 true임) 이를 r에 할당합니다.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
