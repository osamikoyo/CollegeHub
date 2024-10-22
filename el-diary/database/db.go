package database

import (
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"el-diary/el-diary/database/models"
)

var Key string = "fjb13jb2hgjbh35bhbjbjfcj1klkj4tr932490gvj9290vj1o3finhb13iulrfh89i42ghfvn9023ruj"

const contextKeyUser = "user"

// jwtPayloadFromRequest извлекает JWT токен из контекста и возвращает его указанные заявленные, если это возможно.
func JwtPayloadFromRequest(c echo.Context) (jwt.MapClaims, bool) {
	//тут мы по ключу который в константу скинули ("user") получаем от фронтенда токен, тоесть фронт отправляет {"user" : "(тут типа токен, без скобочек)"}
	jwtToken, ok := c.Get(contextKeyUser).(*jwt.Token)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_context_value": c.Get(contextKeyUser),
		}).Error("wrong type of JWT token in context")
		return nil, false
	}
	//тут мы его в нормальный вид приводим, так скажем и возвращаем мапу с email по которому мы ищем юзера в базе данных
	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_claims": jwtToken.Claims,
		}).Error("wrong type of JWT token claims")
		return nil, false
	}

	return payload, true
}
func GenerateToken(u models.User) string {
	//тут мы генерим токен для пользователя, "sub" это то по какому параметру мы будем искать юзера в базе данных
	payload := jwt.MapClaims{
		"sub": u.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	//генерирую подспорье для токена
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := tkn.SignedString(Key)
	//тут токен превращаем в строку, и если есть ошибка выводим её
	if err != nil {
		//вывожу логером ошибку
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Error(err.Error())
	}
	return t
}
func Get_Prof(email string) models.User {
	var u models.User

	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! ЗАМЕНИ DNS НА НУЖНЫЙ !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	dsn := "host=localhost user=osami password= dbname=users port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Error(err.Error())
	}

	if err := db.Where("Email = ?", email).Find(&u); err != nil {
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Error(err.Error.Error())
	}

	return u
}
func Login_User(u models.User) (string) {
	// подключение gorm к базе данных
	dsn := "host=localhost user=osami password= dbname=users port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn))
	//инициализация логера
	loger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err != nil {
		loger.Error(err.Error())
	}
	tkn := GenerateToken(u)
	//тут очень сложная запись, понимаю, что тут происходит, во первых строков if errs := мы сначало иницализируем переменую с ошибкой, а в теле цикла её выводим
	//функцией db.where мы ищем юзера в базе данных где email равна u.Email, с password такая же логика, а функция update обновляет нулевой токен на сгенереный функцией
	if errs := db.Where("Email = ?", u.Email).Where("Password = ?", u.Password).Update("Token", tkn).Error; errs != nil {
		loger.Error(errs.Error())
		return ""
	}
	return tkn
}
