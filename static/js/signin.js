
let signinForm = document.getElementById("signin-forum");

document.getElementById("signin-forum").innerHTML = `
  <div class="header">
    <h1>Real-Time-Forum</h1>
  </div>
  <form id="signin-form">
    <label for="username">Username:</label>
    <input type="text" id="username" name="username" required>
    <label for="password">Password:</label>
    <input type="password" id="password" name="password" required>
    <button type="submit">Sign In</button>
  </form>
  <p>Don't have an account? <a href="/signup">Sign Up</a></p>
`;  