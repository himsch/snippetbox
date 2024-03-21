package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	// 새로운 httptest.ResponseRecorder를 초기화합니다.
	rr := httptest.NewRecorder()

	// 새로운 더미 http.Request를 초기화합니다.
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// httptest.ResponseRecorder 및 http.Request를 전달하여 ping 처리기 함수를 호출합니다.
	ping(rr, r)

	// http.ResponseRecorder에서 Result() 메서드를 호출하여 ping 핸들러에서 생성된 http.Response를 가져옵니다.
	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// 그리고 ping 핸들러가 작성한 응답 본문이 "OK"인지 확인할 수 있습니다.
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
