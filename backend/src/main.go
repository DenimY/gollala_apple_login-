package main

import (
	"fmt"
	"github.com/BillSJC/appleLogin"
	"github.com/labstack/echo"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"

	jwt2 "gopkg.in/square/go-jose.v2/jwt"
)

type AppleConfig struct {
	TeamID      string      //Your Apple Team ID
	ClientID    string      //Your Service which enable sign-in-with-apple service
	RedirectURI string      //Your RedirectURI config in apple website
	KeyID       string      //Your Secret Key ID
	AESCert     interface{} //Your Secret Key Created By X509 package
}

func main() {

	e := echo.New()
	e.GET("/get-apple-callback", func(c echo.Context) error {
		fmt.Println("In get-apple: ", c)

		return c.String(http.StatusOK, "Hello World!")
	})

	e.POST("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})

	e.POST("/redirect", appleSignInLogic)

	e.Logger.Fatal(e.StartTLS(":1323", "ssl/localhost.crt", "ssl/localhost.key")) // localhost:1323
}

// callbackurl 생성
func getAppleCallbackUrl() string {

	type AppleConfig struct {
		TeamID       string      //Your Apple Team ID
		ClientID     string      //Your Service which enable sign-in-with-apple service
		RedirectURI  string      //Your RedirectURI config in apple website
		KeyID        string      //Your Secret Key ID
		AESCert      interface{} //Your Secret Key Created By X509 package
		ResponseMode string
	}

	a := &AppleConfig{
		"5NX4697WZG",                              //Team ID
		"com.ch2ho.hybridshop.gollala.services",   //Client ID (Service ID)
		"https://local.example.com:1323/redirect", //Callback URL
		"T8HKCBDAH9",
		nil,
		"form_post",
	} //Key ID

	u := url.Values{}
	u.Add("response_type", "code id_token")
	u.Add("redirect_uri", a.RedirectURI)
	u.Add("client_id", a.ClientID)
	u.Add("state", "")
	u.Add("scope", "name email")
	u.Add("response_mode", a.ResponseMode)

	callbackURL := "https://appleid.apple.com/auth/authorize?" + u.Encode()

	return callbackURL

}

func appleSignInLogic(c echo.Context) error {
	fmt.Println("IN SIGN IN ")
	fmt.Println("c :" + c.FormValue("user"))

	params := make(map[string]string)

	if err := c.Bind(&params); err != nil {
		panic(err)
	}

	//https://local.example.com:1323/redirect
	a := appleLogin.InitAppleConfig(
		"5NX4697WZG",                              //Team ID
		"com.ch2ho.hybridshop.gollala.services",   //Client ID (Service ID)
		"https://local.example.com:1323/redirect", //Callback URL
		"T8HKCBDAH9")                              //Key ID

	//import cert
	//err := a.LoadP8CertByFile("auth/AuthKey_T8HKCBDAH9.p8") //path to cert file
	//or you can load cert from a string
	if err := a.LoadP8CertByByte([]byte("-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQg9AHHvAJtPf0NAo1Q\n7AYrXlHxFZgnSW31FG4SM6RQBkmgCgYIKoZIzj0DAQehRANCAASum095DhuAkJ+z\nC+BNNdZ5pHuGRKnh/CxIvwvpIP13rWmwF1zmi5yaihcGTKFx+2TgkTd75IC7JgV0\nT4Lj5nkn\n-----END PRIVATE KEY-----")); err != nil {
		panic(err)
	}

	// ... some code to get Apple`s AuthorizationCode
	code := params["code"]
	token, err := a.GetAppleToken(code, 3600)
	if err != nil {
		panic(err)
	}
	fmt.Println("RESULT :   ")
	fmt.Println(token)

	claims := jwt.MapClaims{}

	tokenId := token.IDToken

	_token2, _ := jwt2.ParseSigned(tokenId)
	_ = _token2.UnsafeClaimsWithoutVerification(&claims)

	//payload 값 확인
	fmt.Println(claims)

	return c.String(http.StatusOK, "success")
}
