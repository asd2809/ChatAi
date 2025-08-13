package main
import (
	"github.com/go-playground/validator/v10"
	"demt/mysql/jiaoyan/utils"
	 "github.com/gin-gonic/gin/binding"
)

func main() {
	// r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		utils.RegisterCustomValidators(v)
	}

	// r.POST("/register", RegisterHandler)
	// r.Run()
}
