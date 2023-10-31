package main

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"fmt"
	"github.com/On-Jun9/onjung/config"
	"github.com/On-Jun9/onjung/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func init() {
	//설정 yaml 파일 세팅
	profile := "secrets"
	setRuntimeConfig(profile)
	profile = "properties"
	setRuntimeConfig(profile)
}

// setRuntimeConfig
/*
	PROFILE을 기반으로
	config파일을 읽고 전역변수에 언마샬링 해줍니다.
*/
func setRuntimeConfig(profile string) {
	viper.AddConfigPath(".")

	// 환경변수에서 읽어온 profile이름의 yaml파일을 configPath로 설정합니다.
	viper.SetConfigName(profile)
	viper.SetConfigType("yaml")

	// 기본값을 설정합니다.
	viper.SetDefault("Server.Port", 8080)
	viper.SetDefault("Server.SessionTimeOut", 600)
	viper.SetDefault("Server.DBLogLevel", 2)     //ERROR
	viper.SetDefault("Server.ServerLogLevel", 2) //ERROR
	viper.SetDefault("Datasource.SslMode", "prefer")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// viper는 읽어온 설정파일의 정보를 가지고있으니, 전역변수에 언마샬링해
	// 애플리케이션의 원하는곳에서 사용하도록 합니다.
	err = viper.Unmarshal(&config.RuntimeConf)
	if err != nil {
		panic(err)
	}

	// viper는 설정파일이 변경된 이벤트를 핸들링할수 있습니다.
	// 설정파일이 변경되면 다시 언마샬링해줍니다.
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		var err error
		err = viper.ReadInConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = viper.Unmarshal(&config.RuntimeConf)
		if err != nil {
			fmt.Println(err)
			return
		}
	})
	viper.WatchConfig()
}

// 템플릿 파일을 포함하여 빌드하도록 세팅합니다.
//
//go:embed common/templates/*
var templateFiles embed.FS

// 정적 파일을 포함하여 빌드하도록 세팅합니다.
//
//go:embed common/static/*
var staticFiles embed.FS

// loadTemplates
/*
	템플릿(html) 위치를 지정합니다.
	템플릿의 이름을 Path에 맞게 지정하여, 라우터에 세팅합니다.
	템플릿 파일의 이름이 중복 되어도, 다른 경로라면 사용 가능합니다.
*/
func loadTemplates(router *gin.Engine, templateDirs ...string) {
	tmpl := template.New("").Funcs(router.FuncMap)
	for _, templateDir := range templateDirs {
		err := fs.WalkDir(templateFiles, templateDir, func(path string, d fs.DirEntry, err error) error {
			if d == nil || d.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".html") {
				templateName := strings.Replace(path, string(os.PathSeparator), "/", -1)
				content, err := fs.ReadFile(templateFiles, path)
				if err != nil {
					panic(err)
				}
				tmpl.New(templateName).Parse(string(content))
			}
			return nil
		})

		if err != nil {
			panic(err)
		}
	}
	router.SetHTMLTemplate(tmpl)
}

// registerRoutes
/*
	url 추가합니다.
*/
func registerRoutes(r *gin.Engine) {
	// "/" 기본 경로 접속 시 메인으로 redirect

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/main")
	})
}

// setupRouter
/*
	gin Router 세팅
*/
func setupRouter() *gin.Engine {

	// gin mode 설정
	ginMode := "release"
	if utils.IsInSlice(config.RuntimeConf.Server.Mode, []string{"debug", "release", "test"}) {
		ginMode = config.RuntimeConf.Server.Mode
	}
	gin.SetMode(ginMode)

	// gin 생성, 세팅 합니다.
	r := gin.Default()

	// 로거를 사용합니다.
	r.Use(gin.Logger())

	// panic 발생 시 500 에러로 리턴합니다.
	r.Use(gin.Recovery())

	// 템플릿(html) 위치를 지정합니다.
	//r.LoadHTMLGlob("templates/**/*.html",...)
	loadTemplates(r, "common/templates")

	// static 위치를 지정합니다.
	// Embed static files
	staticFS, err := fs.Sub(staticFiles, "common/static")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/static", http.FS(staticFS))

	//DB 연결-> 전역변수 생성
	err = config.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	//세션 설정
	key := make([]byte, 66)
	_, err = rand.Read(key)
	if err != nil {
		panic(err)
	}
	store := memstore.NewStore([]byte(hex.EncodeToString(key)))
	// 세션 시간 설정 (초) - 설정파일에서 불러옵니다 (기본값 600)
	store.Options(sessions.Options{
		MaxAge: config.RuntimeConf.Server.SessionTimeOut,
		Path:   "/",
	})
	r.Use(sessions.Sessions("onjung-session", store))

	// ---- URL 추가 ----
	registerRoutes(r)

	return r
}

func main() {
	r := setupRouter()
	port := ":" + strconv.Itoa(config.RuntimeConf.Server.Port)
	r.Run(port)
}
