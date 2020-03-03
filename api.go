package main

import (
	"errors"
	"fmt"
	"github.com/beaquant/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func SendErrorMsg(c echo.Context, errorCode int, msg string) error {
	return c.JSON(http.StatusOK, Response{
		Code: 30000,
		Data: map[string]interface{}{
			"error_code": errorCode,
			"error_msg":  msg,
		},
	})
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
		return SendErrorMsg(c, 3000, "个性名重复")
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
	err := verifyAccount(account)
	if err != nil {
		return SendErrorMsg(c, 3001, err.Error())
	}
	acc, err := orm.AddAccount(account)
	if err != nil {
		return SendErrorMsg(c, 3002, err.Error())
	}
	addAccount(acc)
	return SendOK(c, fmt.Sprintf(`{"id":%d}`, acc.ID))
}

// setting - DELETE
func DeleteSetting(c echo.Context) error {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	err := orm.DeleteAccount(uint(id))
	if err != nil {
		return SendErrorMsg(c, 3001, err.Error())
	}
	deleteAccount(uint(id))
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
			total.Btc += v[index].Btc
			total.Usdt += v[index].Usdt
			total.Usd += v[index].Usd
			total.Cny += v[index].Cny

			total.Model = v[index].Model
			total.Btc_Usdt = v[index].Btc_Usdt
			total.Btc_Usd = v[index].Btc_Usd
			total.Btc_Cny = v[index].Btc_Cny
			total.Usdt_Usd = v[index].Usdt_Usd
			total.Usdt_Cny = v[index].Usdt_Cny
			total.Usd_Cny = v[index].Usd_Cny
		}
		history = append(history, Asset{
			Btc:      total.Btc,
			Usdt:     total.Usdt,
			Usd:      total.Usd,
			Cny:      total.Cny,
			Model:    total.Model,
			Btc_Usdt: total.Btc_Usdt,
			Btc_Usd:  total.Btc_Usd,
			Btc_Cny:  total.Btc_Cny,
			Usdt_Usd: total.Usdt_Usd,
			Usdt_Cny: total.Usdt_Cny,
			Usd_Cny:  total.Usd_Cny,
		})
	}
	return SendOK(c, history)
}

func GetCurrentAsset(c echo.Context) error {
	type NicknameAsset struct {
		ID           uint    `json:"id"`
		ExchangeName string  `json:"exchange_name"`
		NickName     string  `json:"nick_name"`
		Btc          float64 `json:"btc"`
		Usdt         float64 `json:"usdt"`
	}
	all := make([]NicknameAsset, 0)
	for _, v := range accounts {
		acc := orm.GetAssetsFromNickname(v.NickName)
		all = append(all, NicknameAsset{
			ID:           acc[len(acc)-1].ID,
			ExchangeName: v.ExchangeName,
			NickName:     v.NickName,
			Btc:          utils.Float64Round(acc[len(acc)-1].Btc, 8),
			Usdt:         utils.Float64Round(acc[len(acc)-1].Usdt, 4),
		})
	}

	return SendOK(c, all)
}

func GetCurrentCoinList(c echo.Context) error {
	type Coin struct {
		CoinName string  `json:"name"`
		Usdt     float64 `json:"value"`
	}
	all := make([]Coin, 0)
	for _, v := range accounts {
		assets := orm.GetAssetsFromAccountId(v.ID)
		asset := assets[len(assets)-1]
		coins := orm.GetCoinsFromAssetId(asset.ID)
		for _, vv := range coins {
			found := false
			for kkk, vvv := range all {
				if vv.CoinName == vvv.CoinName {
					found = true
					all[kkk].Usdt += vv.Usdt
					break
				}
			}
			if !found {
				all = append(all, Coin{
					CoinName: vv.CoinName,
					Usdt:     vv.Usdt,
				})
			}
		}
	}

	return SendOK(c, all)
}

func GetCurrentCoins(c echo.Context) error {
	ids := c.QueryParam("id")
	//ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	if id == 0 {
		return SendErrorMsg(c, 40000, "can't find `id` param")
	}
	return SendOK(c, orm.GetCoinsFromAssetId(uint(id)))
}

// user login
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

// user info
func GetUserInfo(c echo.Context) error {
	token := c.QueryParam("token")
	info, err := parseToken(token)
	if err != nil {
		return SendErrorMsg(c, 4000, err.Error())
	}
	user := info.(map[string]interface{})
	return SendOK(c, map[string]interface{}{
		"roles":        []string{"admin"},
		"introduction": "I am a super administrator",
		"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		"name":         user["name"].(string),
	})
}

// user logout
func UserLogout(c echo.Context) error {
	return SendOK(c, "{}")
}

// support
func GetSupport(c echo.Context) error {
	list := make([]string, 0)
	for k := range List {
		list = append(list, k)
	}
	return SendOK(c, list)
}
