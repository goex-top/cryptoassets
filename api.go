package main

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

var ErrorMsg = map[int]string{
	3000: "nickname duplicate",
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func SendErrorMsg(c echo.Context, errorCode int) error {
	msg, ok := ErrorMsg[errorCode]
	if ok {
		return c.JSON(http.StatusOK, Response{
			Code: 30000,
			Data: map[string]interface{}{
				"error_code": errorCode,
				"error_msg":  msg,
			},
		})
	} else {
		return c.JSON(http.StatusOK, Response{
			Code: 30000,
			Data: map[string]interface{}{
				"error_code": errorCode,
				"error_msg":  "unknown error",
			},
		})
	}
}

func SendOK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code: 20000,
		Data: data,
	})
}

// setting - GET
func GetSetting(c echo.Context) error {
	type AccountSummary struct {
		gorm.Model
		NickName     string `json:"nick_name"`
		ExchangeName string `json:"exchange_name"`
		ApiKey       string `json:"api_key"`
	}
	time.Now().Format(time.StampMicro)
	acc := make([]AccountSummary, 0)
	for _, v := range accounts {
		acc = append(acc, AccountSummary{
			Model:        v.Model,
			NickName:     v.NickName,
			ExchangeName: v.ExchangeName,
			ApiKey:       v.ApiKey,
		})
	}
	return SendOK(c, acc)
}

// setting - POST
func AddSetting(c echo.Context) error {
	type Setting struct {
		ExchangeName string `json:"exchange_name"`
		NickName     string `json:"nick_name"`
		ApiKey       string `json:"api_key"`
		SecKey       string `json:"sec_key"`
		PassKey      string `json:"pass_key"`
	}
	setting := new(Setting)
	if err := c.Bind(setting); err != nil {
		return err
	}

	if orm.HasNickName(setting.NickName) {
		return SendErrorMsg(c, 3000)
	}
	enApiSecretKey := ""
	enApiPassphrase := ""
	if setting.SecKey != "" {
		enApiSecretKey = string(AESECBEncrypt([]byte(setting.SecKey), []byte(conf.User.Password)))
	}
	if setting.PassKey != "" {
		enApiPassphrase = string(AESECBEncrypt([]byte(setting.PassKey), []byte(conf.User.Password)))
	}
	account := Account{
		NickName:      setting.NickName,
		ExchangeName:  setting.ExchangeName,
		ApiKey:        setting.ApiKey,
		ApiSecretKey:  enApiSecretKey,
		ApiPassphrase: enApiPassphrase,
	}
	acc, err := orm.AddAccount(account)
	if err != nil {
		return SendErrorMsg(c, -1)
	}
	accounts = append(accounts, acc)

	return SendOK(c, acc)
}

// setting - DELETE
func DeleteSetting(c echo.Context) error {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	err := orm.DeleteAccount(uint(id))
	if err != nil {
		return SendErrorMsg(c, -2)
	}
	for k, v := range accounts {
		if v.ID == uint(id) {
			accounts = append(accounts[:k], accounts[k+1:]...)
			break
		}
	}
	return SendOK(c, "{}")
}

func GetAssetHistory(c echo.Context) error {
	all := make([][]Asset, 0)
	maxSize := 0
	for _, v := range accounts {
		acc := orm.GetAssetsFromNickname(v.NickName)
		all = append(all, acc)
		if len(acc) > maxSize {
			maxSize = len(acc)
		}
	}
	history := make([]Asset, 0)
	for i := 0; i < maxSize; i++ {
		total := Asset{}
		for _, v := range all {
			index := i - (maxSize - len(v))
			if index < 0 {
				continue
			}
			total.Model = v[index].Model
			total.Btc += v[index].Btc
			total.Btc_Usdt += v[index].Btc_Usdt
			total.Usdt += v[index].Usdt
			total.Usdt_Usd += v[index].Usdt_Usd
			total.Usd_Cny += v[index].Usd_Cny
		}
		history = append(history, Asset{
			Model:    total.Model,
			Btc:      total.Btc,
			Usdt:     total.Usdt,
			Btc_Usdt: total.Usdt,
			Usdt_Usd: total.Usdt_Usd,
			Usd_Cny:  total.Usd_Cny,
		})
	}
	return SendOK(c, history)
}

func GetCurrentAsset(c echo.Context) error {
	type NicknameAsset struct {
		ID       uint
		NickName string
		Btc      float64
		Usdt     float64
	}
	all := make([]NicknameAsset, 0)
	for _, v := range accounts {
		acc := orm.GetAssetsFromNickname(v.NickName)
		all = append(all, NicknameAsset{
			ID:       acc[len(acc)-1].ID,
			NickName: v.NickName,
			Btc:      acc[len(acc)-1].Btc,
			Usdt:     acc[len(acc)-1].Usdt,
		})
	}

	return SendOK(c, all)
}

func GetCurrentCoins(c echo.Context) error {
	//nick_name := c.QueryParam("nick_name")
	ids := c.QueryParam("ID")
	id, _ := strconv.Atoi(ids)
	return SendOK(c, orm.GetCoinsFromAssetId(uint(id)))
}

func UserLogin(c echo.Context) error {

	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if u.UserName == conf.User.UserName && u.Password == conf.User.Password {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = u.UserName
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 6).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(tokenSecKey))
		if err != nil {
			return err
		}

		return SendOK(c, map[string]interface{}{
			"token": t,
		})
	} else {
		return echo.ErrUnauthorized
	}
}

// Token parses and validates a token and return the logged in user
func parseToken(tokenString string) (interface{}, error) {
	if tokenString == "" {
		return nil, nil // unauthorized
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenSecKey), nil
	})
	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("token invalid")
		}
		// In a real authentication, here we should actually validate that the token is valid
		return map[string]interface{}{
			"name":  claims["name"].(string),
			"admin": claims["admin"].(bool),
			"exp":   int64(claims["exp"].(float64)),
		}, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
			return nil, errors.New("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
			return nil, errors.New("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
			return nil, errors.New("Couldn't handle this token: " + err.Error())
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return nil, errors.New("Couldn't handle this token" + err.Error())
	}

	return nil, err
}

func GetUserInfo(c echo.Context) error {
	token := c.QueryParam("token")
	info, err := parseToken(token)
	if err != nil {
		fmt.Println(err)
		return SendErrorMsg(c, -2)
	}
	user := info.(map[string]interface{})
	return SendOK(c, map[string]interface{}{
		"roles":        []string{"admin"},
		"introduction": "I am a super administrator",
		"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		"name":         user["name"].(string),
	})
}

func UserLogout(c echo.Context) error {
	return SendOK(c, "{}")
}

func GetSupport(c echo.Context) error {
	list := make([]string, 0)
	for k := range List {
		list = append(list, k)
	}
	return SendOK(c, list)
}
