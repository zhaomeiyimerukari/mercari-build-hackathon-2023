import { useCookies } from "react-cookie";
import "./Header.css";

export const Header: React.FC = () => {
  const [cookies, _, removeCookie] = useCookies(["userID", "token"]);

  const onLogout = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    removeCookie("userID");
    removeCookie("token");
  };

  return (
    <>
      <header>
        <div className="topLeft">
          <i className="topIcon fab fa-facebook-square"></i>
          <i className="topIcon fab fa-instagram-square"></i>
          <i className="topIcon fab fa-pinterest-square"></i>
          <i className="topIcon fab fa-twitter-square"></i>
        </div>
        <div className="headerTitle">
          <span>Simple Mercari</span>
          <span>-Team 12</span>
        </div>
        <div className="topRight">
          <img
              className="headerImg"
              src="./images/me.jpg"
          />
          <img
              className="headerImg"
              src="./images/Tan.jpg"
          />
          <img
              className="headerImg"
              src="./images/wang.jpg"
          />
        </div>
        <div className="LogoutButtonContainer">
          <button onClick={onLogout} id="MerButton">
            Logout
          </button>
        </div>
      </header>
    </>
  );
}
