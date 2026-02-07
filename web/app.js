let ws = null;
const chat = document.getElementById("chat");
const authUI = document.getElementById("auth-ui");
const authError = document.getElementById("auth-error");

// --- AUTH ---

function login() {
  const username = document.getElementById("login-username").value;
  const password = document.getElementById("login-password").value;

  fetch("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  })
    .then(checkStatus)
    .then(() => {
      window.location.href = "/logged.html";
    })
    .catch(err => showError(err));
}

function register() {
  const username = document.getElementById("reg-username").value;
  const password = document.getElementById("reg-password").value;

  fetch("/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  })
    .then(checkStatus)
    .then(() => {
      alert("Registered! Wait for admin approval.");
    })
    .catch(err => showError(err));
}

function logout() {
  fetch("/logout", { method: "POST" })
    .finally(() => {
      ws?.close();
      ws = null;
      chat.innerHTML = "";
      chatUI.style.display = "none";
      authUI.style.display = "block";
    });
}

// --- CHAT ---

function startChat() {
  ws = new WebSocket(`ws://${window.location.host}/logged/ws`);

  ws.onmessage = e => {
    const msg = JSON.parse(e.data);
    printMessage(msg);
  };

  fetch("/logged/history")
    .then(res => res.json())
    .then(data => data.messages.forEach(printMessage))
    .catch(err => console.error(err));
}

function send() {
  const input = document.getElementById("input");
  const text = input.value.trim();
  if (!text) return;

  fetch("/logged/message", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ text })
  });

  input.value = "";
}

function printMessage(msg) {
  const ts = new Date(msg.timestamp * 1000);
  const li = document.createElement("li");
  li.innerText = `${msg.userId}: ${msg.text} (${ts.toLocaleString()})`;
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
