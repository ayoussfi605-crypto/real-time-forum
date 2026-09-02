import { navigate } from "./navigate.js";
import { checkSession } from "./checkSession.js";
import { updateNavbar, setupLogout, hideNavbar } from "./navbar.js";
import { initChat } from "./chat.js";
import { renderfeed } from "./feed.js";

async function initApp() {
    // Setup logout button and hide navbar initially
    setupLogout();
    hideNavbar();
    // Check if the user is logged in by checking the session
    const loggedIn = await checkSession();

    // User is logged in, navigate to feed and update navbar
    if (loggedIn) {
        updateNavbar()
        initChat();
        renderfeed()
        // navigate("feed");
    // User is not logged in, navigate to signin and update navbar
    } else {
        navigate("signin");
        updateNavbar()
    }
}

initApp();