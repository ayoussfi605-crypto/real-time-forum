import { renderSignIn } from "./signin.js";
import { renderSignUp } from "./signup.js";
import { renderfeed } from "./feed.js";
export function navigate(route) {
    switch (route) {
        case "signin":
            renderSignIn();
            break;

        case "signup":
            renderSignUp();
            break;
        
        case "feed":
            renderfeed();
            break;
    }
}