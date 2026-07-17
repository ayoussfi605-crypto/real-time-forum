import { navigate } from "./navigate.js";
import { checkSession } from "./checkSession.js";

async function initApp() {
    const loggedIn = await checkSession();

    if (loggedIn) {
        navigate("feed");
    } else {
        navigate("signin");
    }
}

initApp();