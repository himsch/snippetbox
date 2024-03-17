package validator

import (
	"strings"
	"unicode/utf8"
)

// 양식 필드에 대한 유효성 검사 오류 맵을 포함하는 새로운 유효성 검사기 유형을 정의합니다.
type Validator struct {
	FieldErrors map[string]string
}

// Valid()는 FieldErrors 맵에 항목이 포함되어 있지 않으면 true를 반환합니다.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError()는 FieldErrors 맵에 오류 메시지를 추가합니다(주어진 키에 대한 항목이 이미 존재하지 않는 한).
func (v *Validator) AddFieldError(key, message string) {
	// 참고: map이 아직 초기화되지 않은 경우 먼저 map를 초기화해야 합니다.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField()는 유효성 검사가 'ok'가 아닌 경우에만 FieldErrors 맵에 오류 메시지를 추가합니다.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank()는 값이 빈 문자열이 아닌 경우 true를 반환합니다.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars()는 값에 n자 이하가 포함된 경우 true를 반환합니다.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt()는 값이 허용되는 정수 목록에 있으면 true를 반환합니다.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
