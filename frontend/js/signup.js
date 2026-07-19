
import { navigate } from "./navigate.js";
import { validateRegister } from "./validateRegister.js";

export function renderSignUp() {
  
  document.getElementById("app").innerHTML=`
  <div class="header">
  <h1>01_Forum</h1>
  </div>
  <form id="signup-form">
  <label for="nickname">Nickname:</label>
  <input type="text" id="nickname" name="nickname" placeholder="Enter your nickname"  required>
  
  <label for="first_name">First Name:</label>
  <input type="text" id="first_name" name="first_name" placeholder="Enter your first name"required>
  
  <label for="last_name">Last Name:</label>
  <input type="text" id="last_name" name="last_name"  placeholder="Enter your last name" required>
  
  <label for="age">Age:</label>
  <input type="number" id="age" name="age" min="1" placeholder="Enter your age" min="1" required>
  
  <label for="gender">Gender:</label>
  <select id="gender" name="gender" required>
  <option value=""> select </option>
  <option value="Male">Male</option>
  <option value="Female">Female</option>
  </select>
  
  <label for="email">Email:</label>
  <input type="email" id="email" name="email" placeholder="Enter your email" required>
  
  <label for="password">Password:</label>
  <input type="password" id="password" name="password" placeholder="Enter your password" required>
  
  <label for="confirmpassword">confirmpassword:</label>
  <input type="password" id="confirmpassword" name="confirmpassword" placeholder="confirm your password" required>
  
  <button type="submit">Sign Up</button>
  <p id="signup-error" style="color:red;"></p>
  </form>
  <p>Already have an account? <a href="#" id="go-signin">Sign In</a></p>
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
    confirmpassword: document.getElementById("confirmpassword").value,
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
    alert("Tregister succec! you can login.");
    navigate("signin");
  } catch(err){
    errorBox.textContent = "internal server error"
    console.error(err);
     
    
  }

}