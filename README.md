# gollala_apple_login-
go, apple login 



golang apple_login git link <br>
https://pkg.go.dev/github.com/kohge4/appleLogin


참고 java link <br>
https://whitepaek.tistory.com/60
<br>
https://whitepaek.tistory.com/61

payload 값 확인 <br>
https://jwt.io/


apple doc <br>
https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_js


# API 

url:
<br>
(GET)  <br>
[ /getAppleUrl ] : Apple 로그인 페이지 얻기 API <br> 
-> response type: string <br>
-> parameter: state(string) 'signIn, signUp' <br>
-> state : signIn(로그인), signUp(회원가입)   <br>
ex) https://appleid.apple.com/auth/authorize?response_type=code%20id_token&response_mode=form_post&client_id=com.ch2ho.hybridshop.gollala.services&redirect_uri=https%3A%2F%2Flocal.example.com%3A1323%2Fredirect&state=signUp&nonce=nonce&scope=email



(POST)  <br>
[ /redirect ] : Apple 로그인 redirection API <br>
-> response type: json

회원가입:
https://dev.clayful.io/ko/http/apis/customer/create 
<br>
로그인 :
https://dev.clayful.io/ko/http/apis/customer/authenticate#payload-user-id



(GET)  <br>
[ / ] : Test용 API 
