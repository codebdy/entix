package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"rxdrag.com/entify/common/errorx"
	"rxdrag.com/entify/common/middlewares"
	"rxdrag.com/entify/config"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/db"
	"rxdrag.com/entify/handler"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/resolve"
	"rxdrag.com/entify/schema"
)

const PORT = 4000

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logFile, err := os.OpenFile("./debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开日志文件异常")
	}
	log.SetOutput(logFile)
}

func checkParams() {
	dbConfig := config.GetDbConfig()
	if dbConfig.Driver == "" ||
		dbConfig.Host == "" ||
		dbConfig.Database == "" ||
		dbConfig.User == "" ||
		dbConfig.Port == "" ||
		dbConfig.Password == "" {
		panic("Params is not enough, please set")
	}
}

func checkMetaInstall() {
	if !repository.IsEntityExists(consts.META_ENTITY_NAME) {
		schema.Installed = false
	} else {
		schema.Installed = true
	}
}

func main() {
	defer db.Close()
	log.Println("启动应用")
	checkParams()
	checkMetaInstall()

	h := handler.New(&handler.Config{
		Pretty:         true,
		GraphiQLConfig: &handler.GraphiQLConfig{},
		FormatErrorFn:  errorx.Format,
	})

	http.Handle("/graphql",
		middlewares.CorsMiddleware(
			middlewares.ContextMiddleware(
				resolve.LoadersMiddleware(h),
			),
		),
	)

	if config.Storage() == consts.LOCAL {
		prefix := "/" + consts.STATIC_PATH + "/"
		fmt.Println(fmt.Sprintf("Running a file server at http://localhost:%d/static/", PORT))
		http.Handle(prefix,
			http.StripPrefix(
				prefix,
				middlewares.CorsMiddleware(http.FileServer(http.Dir("./"+consts.STATIC_PATH)))),
		)
	}

	fmt.Println(fmt.Sprintf("🚀 Graphql server ready at http://localhost:%d/graphql", PORT))
	err2 := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err2 != nil {
		fmt.Printf("启动失败:%s", err2)
	}
}
