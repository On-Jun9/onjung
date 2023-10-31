// Package utils 공통 Utils 패키지
package utils

/*
	utils.go = 공통 사용 함수, 구조체를 모아놓은 파일
*/

// IsInSlice
/*
	IN 함수 []string 안에서 특정 문자열 유무 판단
*/
func IsInSlice(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
