package utils

/*
	webUtils.go = 웹 관련 공통 함수들을 모아놓은 파일
*/

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
)

/*
RenderHTML html 렌더링 시 필수값을 포함하여야 해서 만든 렌더 함수 입니다.
c.HTML() 대신 사용 시 현재 url 정보와, 세션 정보를 포함하므로, 모든 html 처리가 가능합니다.
*/
func RenderHTML(c *gin.Context, code int, name string, data gin.H) {

	// 현재 URL 얻어옵니다.
	currentURL := c.Request.URL.String()

	// 세션에서 유저 정보를 불러옵니다.
	//session := sessions.Default(c)

	// 파싱하여 파라미터 제외한 실제 url만 획득합니다.
	parsedURL, err := url.Parse(currentURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// 공통 데이터를 추가합니다.
	if data == nil {
		data = gin.H{}
	}
	// 현재 URL 정보
	data["currentURL"] = parsedURL.Path

	// HTML을 렌더링합니다.
	c.HTML(code, name, data)
}

/*
CustomAbortWithCode View에서 임의로 에러를 발생시킬 때 사용 하는 함수입니다.
c와 함께 http 코드를 보내면, HTML과 JSON을 자동으로 판단해 반환합니다.
*/
func CustomAbortWithCode(c *gin.Context, code int) {
	c.Set("error_code", code)
	c.Abort()
}
