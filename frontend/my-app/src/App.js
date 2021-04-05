import logo from './logo.svg';
import './App.css';

// import AppleLogin from './apple-login'

import MyAppleSigninButton from "./AppleSignInButton";
import AppleLogin from "react-apple-login";

function App() {
    return (
        <div className="App">
            <h2>App</h2>
            {/*<MyAppleSigninButton/>*/}

            <AppleLogin
                clientId="com.ch2ho.hybridshop.gollala.services"
                redirectURI="https://local.example.com:1323/redirect"
                responseType="code id_token"
                state="test"
                scope="name email"
                nonce=""
                responseMode="form_post"
                usePopup={false}
                designProp={
                    {
                        height: 30,
                        width: 140,
                        color: "black",
                        border: false,
                        type: "sign-in",
                        border_radius: 15,
                        scale: 1,
                        locale: "en_US",
                    }
                }

            />

            {/*<AppleLogin/>*/}
        </div>
    );
}

export default App;
