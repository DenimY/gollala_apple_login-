import logo from './logo.svg';
import './App.css';

// import AppleLogin from './apple-login'

import MyAppleSigninButton from "./AppleSignInButton";
import AppleLogin from "react-apple-login";

function App() {
    return (
        <div className="App">
            {/*<header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>*/}

            <h2>App</h2>
            <MyAppleSigninButton/>
            <AppleLogin clientId="com.react.apple.login" redirectURI="https://redirectUrl.com"/>

            <AppleLogin
                clientId={"com.react.apple.login"}
                redirectURI={"https://redirectUrl.com"}
                responseType={"code"}
                responseMode={"query"}
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
                        // locale: "KR",
                    }
                }

            />
            <AppleLogin
                clientId={"com.react.apple.login"}
                redirectURI={"https://redirectUrl.com"}
                responseType={"code"}
                responseMode={"query"}
                usePopup={false}
                designProp={
                    {
                        height: 30,
                        width: 140,
                        color: "black",
                        border: false,
                        type: "continue",
                        border_radius: 15,
                        scale: 1,
                        // locale: "KR",
                    }
                }

            />


            {/*<AppleLogin/>*/}
        </div>
    );
}

export default App;
