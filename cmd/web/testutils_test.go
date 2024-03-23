package main

import (
	"bytes"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"snippetbox/internal/models/mocks"
	"testing"
	"time"
)

// 사용자 가입 페이지의 HTML에서 CSRF 토큰 값을 캡처하는 정규식을 정의합니다.
var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func extractCSRFToken(t *testing.T, body string) string {
	// FindStringSubmatch 메서드를 사용하여 HTML 본문에서 토큰을 추출합니다.
	// 이는 첫 번째 위치에 일치하는 전체 패턴이 포함된 배열을 반환하고 후속 위치에 캡처된 데이터의 값을 반환합니다.
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}

// 모의 종속성을 포함하는 애플리케이션 구조체의 인스턴스를 반환하는 newTestApplication 도우미를 만듭니다.
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// httptest.Server 인스턴스를 포함하는 사용자 정의 testServer 유형을 정의합니다.
type testServer struct {
	*httptest.Server
}

// 사용자 정의 testServer 유형의 새 인스턴스를 초기화하고 반환하는 newTestServer 도우미를 만듭니다.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// 새 쿠키 jar를 초기화합니다.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// 테스트 서버 클라이언트에 쿠키 jar를 추가합니다.
	// 이제 이 클라이언트를 사용할 때 모든 응답 쿠키가 저장되고 후속 요청과 함께 전송됩니다.
	ts.Client().Jar = jar

	// 사용자 정의 CheckRedirect 기능을 설정하여 테스트 서버 클라이언트에 대한 리디렉션 따르기를 비활성화합니다.
	// 이 함수는 클라이언트가 3xx 응답을 수신할 때마다 호출되며
	// 항상 http.ErrUseLastResponse 오류를 반환하여 클라이언트가 수신된 응답을 즉시 반환하도록 합니다.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

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

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
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
