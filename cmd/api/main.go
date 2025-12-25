package main

import (
	_ "github.com/toanuitt/bookmark_service/docs"
	"github.com/toanuitt/bookmark_service/internal/api"
)

//	@title			BookMark_Service API
//	@version		1.0
//	@description	Password Generator
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	app := api.New(cfg)
	app.Start()
}
