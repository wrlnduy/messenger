const authUI = document.getElementById("auth-ui");
const authError = document.getElementById("auth-error");
const loginPasswordInput = document.getElementById("login-password");
const regPasswordInput = document.getElementById("reg-password");


function login() {
  const usernameInput = document.getElementById("login-username");

  const username = usernameInput.value.trim();
  const password = loginPasswordInput.value;

  loginPasswordInput.value = "";
  loginPasswordInput.type = "password";

  fetch("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  })
    .then(checkStatus)
    .then(() => {
      window.location.href = "/logged";
    })
    .catch(showError);
}

loginPasswordInput.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    login();
  }
});

function register() {
  const usernameInput = document.getElementById("reg-username");
  
  const username = usernameInput.value.trim();
  const password = regPasswordInput.value;

  regPasswordInput.value = "";
  regPasswordInput.type = "password";

  fetch("/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  })
    .then(checkStatus)
    .then(() => {
      alert("Registered!");
    })
    .catch(showError);
}

regPasswordInput.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    register();
  }
});

function togglePassword(id) {
  const input = document.getElementById(id);
  input.type = input.type === "password" ? "text" : "password";
}

function checkStatus(res) {
  if (!res.ok) {
    return res.text().then(text => {
      throw new Error(text || "Error");
    });
  }
  return res;
}

function showError(err) {
  authError.innerText = err.message;
}
