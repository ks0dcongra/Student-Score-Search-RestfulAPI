package controller

import (
	"example1/app/model"
	"example1/app/model/responses"
	"example1/app/service"
	"example1/utils/global"
	"example1/utils/token"
	"net/http"
	"strconv"
	"fmt"
	
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type UserController struct {
	UserService service.UserServiceInterface 
}

func NewUserController(UserService service.UserServiceInterface) *UserController {
	return &UserController{
		UserService,
	}
}

// Login
func (h *UserController) LoginUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		requestData := new(model.LoginStudent)
		// var login model.LoginStudent
		if err := c.ShouldBindJSON(requestData); err != nil {
			c.JSON(http.StatusNotFound, responses.Status(responses.ParameterErr, nil))
			return
		}
		student, status := h.UserService.Login(requestData)
		// student, status:= service.NewUserService().Login(requestData)
		if status == responses.Success{
			c.JSON(http.StatusOK, responses.Status(responses.Success, gin.H{
				"Student": student,
				// [Session用]:拿到上面session暫存
				// [Session用]:用id存至session暫存
				// middleware.SaveSession(c, student.Id)
				// "Sessions": middleware.GetSession(c),
				// [Token用]:回傳的參數
				// "Token": tokenResult,
				}))
			return
		}
		c.JSON(http.StatusNotFound, responses.Status(status, nil))
	}
}

// Logout
func (h *UserController) LogoutUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// [Session用]:清除目前Session
		// middleware.ClearSession(c)
		// [Token用]:取得Header
		tokenString := c.GetHeader("Authorization")
		global.Blacklist[tokenString] = true // 将 Token 加入黑名单
		c.JSON(http.StatusOK, responses.Status(responses.Success, gin.H{
			"message": "Logout Successfully.",
		}))
	}
}

// Create User
func (h *UserController) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestData := new(model.CreateStudent)
		if err := c.ShouldBindJSON(requestData); err != nil {
			c.JSON(http.StatusOK, responses.Status(responses.ParameterErr, nil))
			return
		}
		student_id, status := service.NewUserService().CreateUser(requestData)
		if status != responses.Success {
			c.JSON(http.StatusOK, responses.Status(responses.Error, nil))
			return
		}
		c.JSON(http.StatusOK, responses.Status(responses.Success, student_id))
	}
}

// ScoreSearch
func (h *UserController) ScoreSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestData := c.Param("id")
		if requestData == "0" || requestData == "" {
			c.JSON(http.StatusOK, responses.Status(responses.ParameterErr, nil))
			return
		}

		// 創建 JwtFactory 實例
		JwtFactory := token.Newjwt()

		// [Token用]:先將uint轉換成int再運用strconv轉換成string。
		user_id, err := JwtFactory.ExtractTokenID(c)
		str_user_id := strconv.Itoa(int(user_id))
		// [Token用]:限制只有本人能查詢分數，如果Token login時所暫存的user_id與傳入c的user_id不相符，則回傳只限本人查詢分數。
		if str_user_id != requestData {
			c.JSON(http.StatusOK, responses.Status(responses.ScoreTokenErr, nil))
			return
		}
		// [Token用]:Token那邊出錯了!
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Status(responses.TokenErr, nil))
			return
		}

		redisKey := fmt.Sprintf("user_%s", requestData)
		student, status := service.NewUserService().ScoreSearch(requestData, redisKey)
		if status == responses.Error {
			// 失敗
			c.JSON(http.StatusOK, responses.Status(responses.Error, nil))
		} else if status == responses.SuccessDb {
			// 成功但來自DB
			c.JSON(http.StatusOK, responses.Status(responses.SuccessDb, student))
		} else {
			// 成功但來自Redis
			c.JSON(http.StatusOK, responses.Status(responses.SuccessRedis, student))
		}
	}
}
