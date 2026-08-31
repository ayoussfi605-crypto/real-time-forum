import { navigate } from "./navigate.js";
import { validateRegister } from "./validateRegister.js";

export function renderSignUp() {
  
  document.getElementById("app").innerHTML=`
  <div class="auth-container">
    <div class="auth-form-side">
      <div class="auth-card">
        <h2>Create your account</h2>
        <p class="auth-subtitle">Join the community and start your journey</p>
        <form id="signup-form">
          <div class="input-group">
            <label for="nickname">Nickname</label>
            <input type="text" id="nickname" name="nickname" placeholder="Choose a nickname" required>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
          </div>
          <div class="input-row">
            <div class="input-group">
              <label for="first_name">First Name</label>
              <input type="text" id="first_name" name="first_name" placeholder="First name" required>
            </div>
            <div class="input-group">
              <label for="last_name">Last Name</label>
              <input type="text" id="last_name" name="last_name" placeholder="Last name" required>
            </div>
          </div>
          <div class="input-row">
            <div class="input-group">
              <label for="age">Age</label>
              <input type="number" id="age" name="age" min="1" placeholder="Age" required>
            </div>
            <div class="input-group">
              <label for="gender">Gender</label>
              <select id="gender" name="gender" required>
                <option value="">Select</option>
                <option value="Male">Male</option>
                <option value="Female">Female</option>
              </select>
            </div>
          </div>
          <div class="input-group">
            <label for="email">Email</label>
            <input type="email" id="email" name="email" placeholder="Enter your email" required>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
          </div>
          <div class="input-group" id="password-group">
            <label for="password">Password</label>
            <input type="password" id="password" name="password" placeholder="Create a password" required>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
          </div>
          <div class="input-group" id="confirm-password-group">
            <label for="confirm_password">Confirm Password</label>
            <input type="password" id="confirm_password" name="confirm_password" placeholder="Confirm your password" required>
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
          </div>
          <button type="submit" class="auth-submit-btn">Create Account</button>
          <p id="signup-error"></p>
        </form>
        <p class="auth-switch">Already have an account? <a href="#" id="go-signin">Sign In</a></p>
      </div>
    </div>
  </div>
  `;
  
  document.getElementById("go-signin").addEventListener("click", (e) =>{
    e.preventDefault();
     navigate("signin");
});

document.getElementById("signup-form").addEventListener("submit", handleSignUp);
}

async function handleSignUp(e){
    e.preventDefault();
  const errorBox = document.getElementById("signup-error");
  errorBox.textContent = "";
 

  const data = {
    nickname: document.getElementById("nickname").value.trim(),
    first_name: document.getElementById("first_name").value.trim(),
    last_name: document.getElementById("last_name").value.trim(),
    age: parseInt(document.getElementById("age").value),
    gender: document.getElementById("gender").value,
    email: document.getElementById("email").value.trim(),
    password: document.getElementById("password").value,
    confirm_password: document.getElementById("confirm_password").value,
  };
 
  const error = validateRegister(data);

  if (error) {
    errorBox.textContent = error;
    return;
  }

  try{
    const response = await fetch("/register",{
        method: "POST",
        headers: {"Content-Type": "application/json"},
        credentials: "include",
        body: JSON.stringify(data),
        
    })
     
    
    const result = await response.json();

    if (!response.ok){
        errorBox.textContent = result.message ;
        
        return;
    }
    navigate("signin");
  } catch(err){
    errorBox.textContent = "internal server error"
    console.error(err);
    
     

  }

}
