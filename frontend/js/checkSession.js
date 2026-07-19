import { setCurrentUser, clearCurrentUser } from "./state.js";
import { updateNavbar } from "./navbar.js";

export async function checkSession() {
    try {
        const response = await fetch("/me", {
            credentials: "include",
        });

        if (!response.ok) {
            clearCurrentUser();
            updateNavbar();
            return false;
        }

        const data = await response.json();
        setCurrentUser(data);
        updateNavbar();
        return true;

    } catch (err) {
        clearCurrentUser();
        updateNavbar();
        console.error(err);
        return false;
    }
}