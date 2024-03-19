package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX
// regexp.MustCompile() 함수를 사용하여 정규식을 구문 분석합니다.
// 이메일 주소의 형식을 확인하는 패턴입니다.
// 이는 '컴파일된' regexp.Regexp 유형에 대한 포인터를 반환하거나 오류 발생 시 패닉을 발생시킵니다.
// 시작 시 이 패턴을 한 번 구문 분석하고 컴파일된 *regexp.Regexp를 변수에 저장하는 것이
// 필요할 때마다 패턴을 다시 구문 분석하는 것보다 성능이 더 좋습니다.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
