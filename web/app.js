const ws = new WebSocket("ws://localhost:8080/ws");
const chat = document.getElementById("chat");
const myUserID = getCookies().user_id;

ws.onmessage = (e) => {
  const msg = JSON.parse(e.data);
  printMessage(msg);
};

fetch("/history")
  .then(res => {
    if (!res.ok) throw new Error("HTTP error " + res.status);
    return res.json();
  })
  .then(msgs => {
    msgs.messages.forEach(msg => {
      printMessage(msg);
    });
  })
  .catch(err => console.error(err));

function send() {
  const input = document.getElementById("input");
  fetch("/message", {
    method: "POST",
    body: JSON.stringify({ text: input.value }),
  });
  input.value = "";
}

function getCookies() {
  return document.cookie.split(';').reduce((acc, c) => {
    const [key, value] = c.trim().split('=');
    acc[key] = decodeURIComponent(value);
    return acc;
  }, {});
}

function printMessage(msg) {
  const ts = new Date(msg.timestamp * 1000);
  const li = document.createElement("li");
  li.innerText = `${msg.userId}: ${msg.text}\t ${ts.toDateString()} / ${ts.toTimeString()}`;
  if (msg.user_id === myUserID) {
    li.style.fontWeight = "bold";
  }
  chat.appendChild(li);
}