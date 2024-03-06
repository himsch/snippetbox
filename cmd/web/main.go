package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	// 중요한 점은 flag.Parse() 함수를 사용하여 명령줄 플래그를 구문 분석한다는 것입니다.
	// 이는 명령줄 플래그 값을 읽고 이를 addr에 할당합니다.
	// 변수. addr 변수를 사용하기 *전에* 이것을 호출해야 합니다
	// 그렇지 않으면 항상 기본값 ":4000"이 포함됩니다. 오류가 있는 경우
	// 구문 분석 중에 발생하면 애플리케이션이 종료됩니다.
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// mux.Handle() 함수를 사용하여 파일 서버를 처리기로 등록합니다.
	// "/static/"으로 시작하는 모든 URL 경로. 일치하는 경로의 경우 요청이 파일 서버에 도달하기 전에 "/static" 접두사를 제거합니다.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// flag.String() 함수에서 반환된 값은 플래그에 대한 포인터입니다.
	// 값, 값 자체가 아닙니다. 따라서 포인터를 역참조해야 합니다. (즉, 앞에 * 기호를 붙입니다).
	// 로그 메시지에 주소를 삽입하기 위해 log.Printf() 함수를 사용하고 있다는 점에 유의하세요.
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
