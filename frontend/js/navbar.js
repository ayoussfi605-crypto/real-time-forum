import { getCurrentUser, clearCurrentUser } from "./state.js";
import { navigate } from "./navigate.js";

export function updateNavbar() {
  const navbar = document.getElementById("navbar");
  const user = getCurrentUser();

  if (user) {
    navbar.classList.remove("hidden")

    document.getElementById("nav-username").textContent = user.nickname

  } else {
    navbar.classList.add("hidden")
  }
}

export function setupLogout() {
document.getElementById("logout-btn").addEventListener("click", async () => {
  try {

   const response = await fetch("/logout",{ 
     credentials: "include" , 
           method: "POST"
    })


   if(!response.ok){
    console.error("Logout failed");
    return
   }
    clearCurrentUser();
    updateNavbar();
    navigate("signin");

  } catch (err) {
     console.error(err);
    return
  }
});
}