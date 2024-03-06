package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	// 정보 메시지 작성을 위한 로거를 생성하려면 log.New()를 사용합니다.
	// 여기에는 로그를 쓸 대상(os.Stdout), 메시지의 문자열 접두사(INFO 다음에 탭이 옴),
	// 포함할 추가 정보를 나타내는 플래그(현지 날짜 및 시간) 등 세 가지 매개변수가 사용됩니다.
	// 플래그는 비트 OR 연산자 |를 사용하여 결합됩니다.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// 오류 메시지 작성을 위한 로거도 같은 방식으로 생성하되 stderr을 대상으로 사용하고
	// log.Lshortfile 플래그를 사용하여 해당 파일명과 라인 번호를 포함시킨다.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// 표준 로거 대신 두 개의 새로운 로거를 사용하여 메시지를 작성합니다.
	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
