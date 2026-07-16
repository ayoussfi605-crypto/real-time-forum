let currentUser = null;

function setCurrentUser(user) {
  currentUser = user;
}

function getCurrentUser() {
  return currentUser;
}

function clearCurrentUser() {
  currentUser = null;
}

function isLoggedIn() {
  return currentUser !== null;
}