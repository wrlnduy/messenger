let ws = null;
let chat = null;
const authUI = document.getElementById("auth-ui");
const authError = document.getElementById("auth-error");

// --- AUTH ---

function login() {
  const usernameInput = document.getElementById("login-username");
  const passwordInput = document.getElementById("login-password");

  const username = usernameInput.value;
  const password = passwordInput.value;

  passwordInput.value = "";
  passwordInput.type = "password";

  fetch("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  })
    .then(checkStatus)
    .then(() => {
      window.location.href = "/logged";
    })
    .catch(err => showError(err));
}

function register() {
  const usernameInput = document.getElementById("reg-username");
  const passwordInput = document.getElementById("reg-password");

  const username = usernameInput.value;
  const password = passwordInput.value;

  passwordInput.value = "";
  passwordInput.type = "password";

  fetch("/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  })
    .then(checkStatus)
    .then(() => {
      alert("Registered!");
    })
    .catch(err => showError(err));
}

function togglePassword(id) {
  const input = document.getElementById(id);
  input.type = input.type === "password" ? "text" : "password";
}

function logout() {
  fetch("/logout", { method: "POST" })
    .finally(() => {
      ws?.close();
      ws = null;
      chat.innerHTML = "";
      chatUI.style.display = "none";
      authUI.style.display = "block";
      window.location.href = "/";
    });
}

// --- CHAT ---

function startChat() {
  ws = new WebSocket(`wss://${window.location.host}/logged/ws?chat_id=00000000-0000-0000-0000-000000000001`);

  ws.onmessage = e => {
    const msg = JSON.parse(e.data);
    printMessage(msg);
  };

  fetch("/logged/history?chat_id=00000000-0000-0000-0000-000000000001")
    .then(res => res.json())
    .then(data => {
      const users = data.mapping;

      data.messages.forEach(msg => {
        msg.username = users[msg.userId];
        printMessage(msg);
      });
    })
    .catch(err => console.error(err));
}

function send() {
  const input = document.getElementById("input");
  const text = input.value.trim();
  if (!text) return;

  fetch("/logged/message?chat_id=00000000-0000-0000-0000-000000000001", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ text })
  });

  input.value = "";
}

function printMessage(msg) {
  const ts = new Date(msg.timestamp * 1000);
  const li = document.createElement("li");
  li.innerText = `${msg.username}: ${msg.text} (${ts.toLocaleString()})`;
  chat.appendChild(li);
  li.scrollIntoView();
}

// --- HELPERS ---

function checkStatus(res) {
  if (!res.ok) {
    return res.text().then(text => { throw new Error(text || "Error"); });
  }
  return res;
}

function showError(err) {
  authError.innerText = err.message;
}
