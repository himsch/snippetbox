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

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// 새로운 http.Server 구조체를 초기화합니다. Addr 및 Handler 필드를 다음과 같이 설정했습니다.
	// 서버는 이전과 동일한 네트워크 주소와 경로를 사용하고,
	// 이제 문제가 발생할 경우 서버가 사용자 정의 errorLog 로거를 사용하도록 ErrorLog 필드를 설정합니다.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// 새 http.Server 구조체에서 ListenAndServe() 메서드를 호출합니다.
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
