package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// mux.Handle() 함수를 사용하여 파일 서버를 처리기로 등록합니다.
	// "/static/"으로 시작하는 모든 URL 경로. 일치하는 경로의 경우 요청이 파일 서버에 도달하기 전에 "/static" 접두사를 제거합니다.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
