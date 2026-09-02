import { navigate } from "./navigate.js";
import { setCurrentUser } from "./state.js";
import { checkSession } from "./checkSession.js";
import { updateNavbar } from "./navbar.js";
import { initChat } from "./chat.js";

export function renderSignIn(){

 document.getElementById("app").innerHTML =`
 <div class="auth-container">
   <div class="auth-form-side">
     <div class="auth-card">
       <h2>Welcome back 👋</h2>
       <p class="auth-subtitle">Sign in to your account to continue</p>
       <form id="signin-form">
         <div class="input-group">
           <label for="identifier">Username or Email</label>
           <input type="text" id="identifier" name="identifier" placeholder="Enter your username or email" required>
           <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
         </div>
         <div class="input-group" id="password-group">
           <label for="password">Password</label>
           <input type="password" id="password" name="password" placeholder="Enter your password" required>
           <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
         </div>
         <button type="submit" class="auth-submit-btn">Sign In</button>
         <p id="signin-error"></p>
       </form>
       <p class="auth-switch">Don't have an account? <a href="#" id="go-signup">Sign Up</a></p>
     </div>
   </div>
 </div>
 `;

document.getElementById("go-signup").addEventListener("click", (e) =>{
    e.preventDefault();
     navigate("signup");

});

 document.getElementById("signin-form").addEventListener("submit", handleSignIn);
}

async function handleSignIn(e){
    e.preventDefault();

    const errorBox = document.getElementById("signin-error");
    errorBox.textContent = "";
    
   const x = document.getElementById("identifier").value.trim()
   const y = document.getElementById("password").value
   const data = {
       identifier: x,
       password: y
    };

  try {

    const response = await fetch("/login", {  
        method: "POST",
        headers: {"Content-Type": "application/json"},
        credentials: "include",
        body: JSON.stringify(data),
    })
  
    const result = await response.json();

    if (!response.ok) {
      errorBox.textContent = result.message ;
      return;
    }

    await checkSession()
    updateNavbar()
    navigate("feed");
    initChat();

  } catch (err) {
    errorBox.textContent = "Internal Server Error";
    console.error(err);
  }
   
}
