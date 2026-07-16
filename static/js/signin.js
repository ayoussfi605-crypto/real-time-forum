

function renderSignIn(){

 document.getElementById("app").innerHTML =`
 <div class="header"><h1>Real-Time-Forum</h1></div>
 <form id="signin-form">
 <label for="identifier">Nickname or Email:</label>
 <input  type="text" id="identifier" name="identifier" required>

 <label for="password">Password:</label>
 <input  type="password" id="password" name="password" required>

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
    
   let x = document.getElementById("identifier").value.trim()
   let y = document.getElementById("password").value
   const data = {
       identifier: x,
       password: y
    };
 
   
}

// renderSignIn();