import { navigate } from "./navigate.js";
import { setCurrentUser } from "./state.js";
import { checkSession } from "./checkSession.js";
import { updateNavbar } from "./navbar.js";

export function renderSignIn(){

 document.getElementById("app").innerHTML =`
 <div class="header"><h1>01_Forum</h1></div>
 <form id="signin-form">
 <label for="identifier">Nickname or Email:</label>
 <input  type="text" id="identifier" name="identifier" placeholder="Enter your identifier" required>

 <label for="password">Password:</label>
 <input  type="password" id="password" name="password" placeholder="Enter your password" required>

 <button type="submit">Sign In</button>
 <p id="signin-error" style="color:red;"></p>
 </form>
 <p>Don't have an account? <a href="#" id="go-signup">Sign Up</a></p>
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

    const response = await fetch("/login",{
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
