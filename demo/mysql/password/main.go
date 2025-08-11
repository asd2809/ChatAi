package main
import (
	"fmt"
	"log"
	"demt/mysql/password/utils"
)

func main() {
	password := "myPassword123"

	// 在注册的时候使用存放入数据库中
	// 生成哈希
	hash, err := utils.GeneratePasswordHash(password, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash:", hash)

	// 登录的时候使用，通过同样的Argon2算法和参数生成哈希，跟数据库中的哈希做安全比较
	// 验证密码
	ok, err := utils.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Password valid:", ok)
}
