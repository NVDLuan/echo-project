package auth

import (
	"my-project/configs/app"
	"my-project/pkg/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetUsersHandler lấy danh sách user
// @Summary Lấy danh sách người dùng
// @Description API này trả về danh sách tất cả người dùng
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func GetUsersHandler(c echo.Context) error {
	users, err := GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// GetUserHandler lấy thông tin user theo ID
// @Summary Lấy user theo ID
// @Description API này trả về thông tin một user dựa trên ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUserHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUserHandler xóa user
// @Summary Xóa user theo ID
// @Description API này xóa một user từ database
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUserHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
}

// RegisterHandler
// @Summary Đăng ký user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Thông tin user"
// @Success 201 {object} map[string]string
// @Router /auth/register [post]
func RegisterHandler(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	req.Password = hashedPassword
	if err := CreateUser(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered"})
}

// LoginHandler
// @Summary Đăng nhập user
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Email và Password"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func LoginHandler(c echo.Context) error {
	var req LoginRequest
	log := logger.GetLogger()

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		log.Error("user not found")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	if !CheckPassword(user.PasswordHash, req.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	log.Info("Login success")
	config := app.GetConfig()
	if config.JWTSetCookie {
		return SetCookie(config.JWTAuthCookie, token, c)
	}
	return c.JSON(http.StatusOK, map[string]string{"access_token": token})
}
