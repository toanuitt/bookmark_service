package main

import (
	_ "github.com/toanuitt/bookmark_service/docs"
	"github.com/toanuitt/bookmark_service/internal/infrastructure"
)

//	@title			BookMark_Service API
//	@version		1.2
//	@description	API documentation for bookmark service
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your Bearer token in the format: Bearer {token}
func main() {
	app := infrastructure.CreateAPI()
	app.Start()
}
