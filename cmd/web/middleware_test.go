package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox/internal/assert"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 200 상태 코드와 "OK" 응답 본문을 작성하는 secureHeaders 미들웨어에 전달할 수 있는 모의 HTTP 핸들러를 만듭니다.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// 모의 HTTP 핸들러를 secureHeaders 미들웨어에 전달합니다.
	// secureHeaders는 http.Handler를 *반환*하기 때문에 ServeHTTP() 메서드를 호출하여
	// http.ResponseRecorder 및 더미 http.Request를 전달하여 실행할 수 있습니다.
	secureHeaders(next).ServeHTTP(rr, r)

	// 테스트 결과를 얻으려면 http.ResponseRecorder에서 Result() 메서드를 호출하세요.
	rs := rr.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
