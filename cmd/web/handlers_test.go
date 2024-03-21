package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"snippetbox/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	// 애플리케이션 구조체의 새 인스턴스를 만듭니다.
	// 현재 여기에는 두 개의 모의 로거만 포함되어 있습니다(기록된 모든 항목을 삭제함).
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	// 그런 다음 httptest.NewTLSServer() 함수를 사용하여 새 테스트 서버를 생성하고
	// app.routes() 메서드에서 반환된 값을 서버의 핸들러로 전달합니다.
	// 테스트 기간 동안 로컬 시스템에서 무작위로 선택된 포트를 수신하는 HTTPS 서버가 시작됩니다.
	// 테스트가 완료되면 서버가 종료되도록 ts.Close() 호출을 연기합니다.
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// 테스트 서버가 수신 대기 중인 네트워크 주소는 ts.URL 필드에 포함되어 있습니다.
	// 이를 ts.Client().Get()과 함께 사용할 수 있습니다.
	// 테스트 서버에 대해 GET /ping 요청을 보내는 방법입니다.
	// 그러면 응답이 포함된 http.Response 구조체가 반환됩니다.
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
