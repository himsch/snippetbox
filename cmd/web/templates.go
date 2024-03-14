package main

import "snippetbox/internal/models"

// 보유 구조 역할을 할 templateData 유형을 정의합니다.
// HTML 템플릿에 전달하려는 동적 데이터입니다.
// 현재는 필드가 하나만 포함되어 있지만 더 추가하겠습니다.
// 빌드가 진행됨에 따라.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
