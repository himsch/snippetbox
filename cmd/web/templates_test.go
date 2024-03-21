package main

import (
	"snippetbox/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// 테스트 사례 이름, humanDate() 함수에 대한 입력(tm 필드) 및
	// 예상 출력(wan 필드)을 포함하는 익명 구조체 조각을 만듭니다.
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	for _, tt := range tests {
		// t.Run() 함수를 사용하여 각 테스트 사례에 대한 하위 테스트를 실행합니다.
		// 첫 번째 매개변수는 테스트 이름(로그 출력에서 하위 테스트를 식별하는 데 사용됨)이고
		// 두 번째 매개변수는 각 사례에 대한 실제 테스트를 포함하는 익명 함수입니다.
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
