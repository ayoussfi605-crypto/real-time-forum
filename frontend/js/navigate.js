import { renderSignIn } from "./signin.js";
import { renderSignUp } from "./signup.js";
import { renderfeed } from "./feed.js";
import { hideNavbar, updateNavbar } from "./navbar.js";
import { initChat, closeChatAndWS } from "./chat.js";

export function navigate(route) {
    switch (route) {
        case "signin":
            hideNavbar();
            closeChatAndWS();
            renderSignIn();
            break;

        case "signup":
            hideNavbar();
            closeChatAndWS();
            renderSignUp();
            break;
        
        case "feed":
            updateNavbar();
            initChat();
            renderfeed();
            break;
    }
}