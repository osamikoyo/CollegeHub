package api

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"el-diary/el-diary/database"
	"el-diary/el-diary/database/models"
)

func Home(c echo.Context) error {
	return c.String(http.StatusOK, "Server work succesfuly")
}
func Login(c echo.Context) error {
	var u models.User
	loger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := c.Bind(&u); err != nil{
		loger.Error(err.Error())
		return err
	}

	tkn := database.Login_User(u)
	return c.JSON(http.StatusOK, models.TKN{
		Token: tkn,
	})
}
func Profile(c echo.Context) error {
	jwtplayload, ok := database.JwtPayloadFromRequest(c)
	if !ok {
		return c.String(http.StatusNotFound, "error")
	}
	user := database.Get_Prof(jwtplayload["sub"].(string))

	return c.JSON(http.StatusOK, models.User{Email: user.Email})
}
