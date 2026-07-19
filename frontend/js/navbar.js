import { getCurrentUser, clearCurrentUser } from "./state.js";
import { navigate } from "./navigate.js";

export function hideNavbar() {
  const navbar = document.getElementById("navbar");

  if (navbar) {
    navbar.classList.add("hidden");
  document.getElementById("navbar").style.display = "none";}
}

export function updateNavbar() {
  const navbar = document.getElementById("navbar");
  const user = getCurrentUser();
  const username = document.getElementById("nav-username");

  if (!navbar) {
    return;
  }

  if (user) {
    navbar.classList.remove("hidden");
    navbar.style.display = "flex";
    if (username) {
      username.textContent = user.nickname;
    }
  } else {
    navbar.classList.add("hidden");
    navbar.style.display = "none";
  }
}

export function setupLogout() {
  document.getElementById("logout-btn").addEventListener("click", async () => {
    try {
      const response = await fetch("/logout", { credentials: "include", method: "POST" });
      if (!response.ok) {
        console.error("Logout failed on server, clearing local state anyway");
      }
    } catch (err) {
      console.error(err);
    } finally {
      clearCurrentUser();
      updateNavbar();
      navigate("signin");
    }
  });
}