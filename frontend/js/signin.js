import { navigate } from "./navigate.js";
import { setCurrentUser } from "./state.js";
import { checkSession } from "./checkSession.js";
import { updateNavbar } from "./navbar.js";

export function renderSignIn(){

 document.getElementById("app").innerHTML =`
 <div class="auth-container">
   <div class="auth-branding">
     <div class="brand-content">
       <h1 class="brand-logo">01 Forum</h1>
       <p class="brand-tagline">Connect. Share. Learn.</p>
       <p class="brand-description">Join our developer community and start meaningful conversations with fellow programmers around the world.</p>
       <div class="brand-benefits">
         <div class="benefit-item">
           <div class="benefit-icon">
             <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="22" height="22"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
           </div>
           <span>Discussions</span>
         </div>
         <div class="benefit-item">
           <div class="benefit-icon">
             <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="22" height="22"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
           </div>
           <span>Private Messages</span>
         </div>
         <div class="benefit-item">
           <div class="benefit-icon">
             <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="22" height="22"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
           </div>
           <span>Community</span>
         </div>
       </div>
     </div>
   </div>
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
         <div class="input-group password-group">
           <label for="password">Password</label>
           <input type="password" id="password" name="password" placeholder="Enter your password" required>
           <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
           <button type="button" class="toggle-password" data-target="password" aria-label="Toggle password visibility">
             <svg class="eye-open" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
             <svg class="eye-closed hidden" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
           </button>
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

 // Setup password toggle visibility
 setupPasswordToggles();
}

function setupPasswordToggles() {
 document.querySelectorAll('.toggle-password').forEach(btn => {
   btn.addEventListener('click', () => {
     const targetId = btn.getAttribute('data-target');
     const input = document.getElementById(targetId);
     const eyeOpen = btn.querySelector('.eye-open');
     const eyeClosed = btn.querySelector('.eye-closed');

     if (input.type === 'password') {
       input.type = 'text';
       eyeOpen.classList.add('hidden');
       eyeClosed.classList.remove('hidden');
     } else {
       input.type = 'password';
       eyeOpen.classList.remove('hidden');
       eyeClosed.classList.add('hidden');
     }
   });
 });
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

  } catch (err) {
    errorBox.textContent = "Internal Server Error";
    console.error(err);
  }
   
}
