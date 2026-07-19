import { navigate } from "./navigate.js";
import { checkSession } from "./checkSession.js";
import { updateNavbar, setupLogout, hideNavbar } from "./navbar.js";

async function initApp() {
    setupLogout();
    hideNavbar();

    const loggedIn = await checkSession();

    if (loggedIn) {
        updateNavbar()
        navigate("feed");
    } else {
        navigate("signin");
        updateNavbar()
    }
}

initApp();