package main

import (
	"fmt"
	"github.com/BillSJC/appleLogin"
	"github.com/denimY/gollala_apple_login-/backend/src/config"
	url_conf "github.com/denimY/gollala_apple_login-/backend/src/config/url"
	apiModel "github.com/denimY/gollala_apple_login-/backend/src/model"
	math_util "github.com/denimY/gollala_apple_login-/backend/src/util"
	"github.com/labstack/echo"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"math/rand"
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

	//http.HandleFunc("/hello", HelloServer)
	//err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

	e.GET("/getAppleUrl", func(c echo.Context) error {
		state := c.QueryParam("state")
		return c.String(http.StatusOK, getAppleCallbackUrl(state))
	})

	e.POST("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})

	e.POST("/redirect", appleSignInLogic)

	e.Logger.Fatal(e.StartTLS(":1323", "ssl/localhost.crt", "ssl/localhost.key")) // localhost:1323
	//e.Logger.Fatal(e.StartTLS(":1323", "../../ssl/local.example.com/cert.pem", "../../ssl/local.example.com/key.pem")) // localhost:1323
	//e.Logger.Fatal(e.StartTLS(":1323", "../../ssl/minica.pem", "../../ssl/minica-key.pem")) // localhost:1323
	//e.Logger.Fatal(e.StartTLS(":1323", "../../ssl/127.0.0.1/cert.pem", "../../ssl/127.0.0.1/key.pem")) // localhost:1323
	//e.Logger.Fatal(e.StartTLS(":1323", "../../ssl/local.example.com/cert.pem", "../../ssl/local.example.com/key.pem")) // localhost:1323
	//e.Logger.Fatal(e.StartTLS(":1323", "../../ssl/certs/local.example.com/cert.pem", "../../ssl/certs/local.example.com/key.pem")) // localhost:1323
	//e.Logger.Fatal(e.StartTLS(":1323", "../../ssl/server.crt", "../../ssl/server.key"))                                            // localhost:1323
}

func getAppleCallbackUrl(state string) string {

	u := url.Values{}
	u.Add("response_type", "code id_token") // code와 id_token 모두 받기
	u.Add("redirect_uri", config.RedirectURI)
	u.Add("client_id", config.ClientID)
	u.Add("state", state)               // signIn, singUp 둘 중 하나만 넣도록 해야함
	u.Add("scope", "name email")        // name email 모두 요청
	u.Add("response_mode", "form_post") //scope 쓰기 위해서는 항시 필요

	// url 생성
	callbackURL := config.AppleCallbackURL + u.Encode()

	return callbackURL

}

func appleSignInLogic(c echo.Context) error {
	fmt.Println("IN SIGN IN ")
	fmt.Println("c :" + c.FormValue("user"))

	params := make(map[string]interface{})

	if err := c.Bind(&params); err != nil {
		panic(err)
	}

	state := fmt.Sprint(params["state"])
	if state == "" || (state != "signIn" && state != "signUp") {
		return c.JSON(http.StatusBadRequest, `"{"message": "wrong state value. set 'signIn' or 'signUp'"}"`)
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

	//callbackURL := a.CreateCallbackURL(state)
	//fmt.Println("callbackurl: {}", callbackURL)

	// ... some code to get Apple`s AuthorizationCode
	code := fmt.Sprint(params["code"])
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

	fmt.Println(claims)
	fmt.Println("success")

	//state := params["state"]
	if state == "signIn" {
		resp, _ := signInClayfulApi(fmt.Sprintf("%v", claims["sub"]))
		return c.JSON(http.StatusOK, resp)
	} else if state == "signUp" {
		resp, _ := signUpClayfulApi(fmt.Sprintf("%v", claims["sub"]))
		return c.JSON(http.StatusOK, resp)
	} else {
		fmt.Printf("wrong state value. state :[%v]\n", state)
		return c.JSON(http.StatusBadRequest, `"{message: wrong state value. set 'signIn' or 'signUp' }"`)
	}

}

// createClayfulApi is sign in Claypul api
func signInClayfulApi(sub string) (*http.Response, error) {

	// Generate 8-digit random string variables
	// ex) Guest00000001
	randomGuest := "Guest" + math_util.Lpad(fmt.Sprintf("%v", rand.Intn(100000000)), "0", 2)

	// api 사용 이름 임시지정
	nameFirst := "temp_first"
	nameLast := "temp_last"

	jsonBody := apiModel.SignInBody{
		Connect: true, // empty pwd option
		UserId:  sub,
		Alias:   randomGuest,
		Name: apiModel.UserName{
			First: nameFirst,
			Last:  nameLast,
			Full:  randomGuest,
		},
	}

	resp, body, errs := gorequest.New().Post(url_conf.SignInClayPulUserId).
		Set("Accept", "application/json").
		//Set("Accept-Encoding", "gzip").
		//Set("Accept-Encoding", "gzip, deflate, br").
		Set("Authorization", config.ClaypulApiToken).
		//Set("Access-Control-Allow-Origin", "*").
		//Set("Access-Control-Allow-Methods", "*").
		Send(jsonBody).
		End()

	if errs != nil {
		return nil, fmt.Errorf(fmt.Sprint(errs))
	}

	fmt.Printf("resp=%v body=%v err=%v", resp, body, errs)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		fmt.Println(str)
	}

	if err != nil {
		panic(err)
	}

	return resp, nil

}

func signUpClayfulApi(sub string) (*http.Response, error) {

	jsonBody := apiModel.SignUpBody{
		Connect: true,
		UserId:  sub,
	}

	resp, body, errs := gorequest.New().Post(url_conf.SignUpClayPulUserId).
		Set("Accept", "application/json").
		//Set("Accept-Encoding", "gzip").
		Set("Authorization", config.ClaypulApiToken).
		Send(jsonBody).
		End()

	if errs != nil {
		return nil, fmt.Errorf(fmt.Sprint(errs))
	}

	fmt.Printf("resp=%v body=%v err=%v", resp, body, errs)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		fmt.Println(str)
	}

	if err != nil {
		panic(err)
	}

	return resp, nil

}
