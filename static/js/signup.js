// signup.js
// Registration form + fetch l /api/register
 
function renderSignUp() {

document.getElementById("app").innerHTML=`
    <div class="header">
      <h1>Real-Time-Forum</h1>
    </div>
    <form id="signup-form">
      <label for="nickname">Nickname:</label>
      <input type="text" id="nickname" name="nickname" required>
 
      <label for="first_name">First Name:</label>
      <input type="text" id="first_name" name="first_name" required>
 
      <label for="last_name">Last Name:</label>
      <input type="text" id="last_name" name="last_name" required>
 
      <label for="age">Age:</label>
      <input type="number" id="age" name="age" min="1" required>
 
      <label for="gender">Gender:</label>
      <select id="gender" name="gender" required>
        <option value=""> select </option>
        <option value="Male">Male</option>
        <option value="Female">Female</option>
      </select>
 
      <label for="email">Email:</label>
      <input type="email" id="email" name="email" required>
 
      <label for="password">Password:</label>
      <input type="password" id="password" name="password" required>
 
      <button type="submit">Sign Up</button>
      <p id="signup-error" style="color:red;"></p>
    </form>
    <p>Already have an account? <a href="#" id="go-signin">Sign In</a></p>
  `;

  document.getElementById("go-signin").addEventListener("click", (e) =>{
    e.preventDefault();
     navigate("signup");
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
  };
 

  try{
    const response = await fetch("/register",{
        method: "POST",
        headers: {"content-type":"applicaton/json"},
        credentials: "include",
        body: JSON.stringify(data),
        
    })
     
    
    const result = await response.json();

    if (!response.ok){
        errorBox.textContent = "error fetch data";
        
        return;
    }
    alert("Tregister succec! you can login.");
    navigate("signup");
  } catch(err){
    errorBox.textContent = "internal server error"
    console.error(err);
     
    
  }
 

}
renderSignUp()