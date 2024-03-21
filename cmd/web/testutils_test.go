package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 모의 종속성을 포함하는 애플리케이션 구조체의 인스턴스를 반환하는 newTestApplication 도우미를 만듭니다.
func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

// httptest.Server 인스턴스를 포함하는 사용자 정의 testServer 유형을 정의합니다.
type testServer struct {
	*httptest.Server
}

// 사용자 정의 testServer 유형의 새 인스턴스를 초기화하고 반환하는 newTestServer 도우미를 만듭니다.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

// 사용자 정의 testServer 유형에 get() 메소드를 구현하십시오.
// 테스트 서버 클라이언트를 사용하여 지정된 URL 경로에 대해
// GET 요청을 수행하고 응답 상태 코드, 헤더 및 본문을 반환합니다.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
