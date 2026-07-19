import { renderSignIn } from "./signin.js";
import { renderSignUp } from "./signup.js";
import { renderfeed } from "./feed.js";
import { hideNavbar, updateNavbar } from "./navbar.js";

export function navigate(route) {
    switch (route) {
        case "signin":
            hideNavbar();
            renderSignIn();
            break;

        case "signup":
            hideNavbar();
            renderSignUp();
            break;
        
        case "feed":
            updateNavbar();
            renderfeed();
            break;
    }
}