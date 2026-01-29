const ws = new WebSocket("ws://localhost:8080/ws");
const chat = document.getElementById("chat");

ws.onmessage = (e) => {
  const msg = JSON.parse(e.data);
  const li = document.createElement("li");
  li.innerText = `${msg.user_id}: ${msg.text} | ${new Date(msg.timestamp).toDateString()}/${new Date(msg.timestamp).toTimeString()}`;
  chat.appendChild(li);
};

function send() {
  const input = document.getElementById("input");
  fetch("/message", {
    method: "POST",
    body: JSON.stringify({ text: input.value }),
  });
  input.value = "";
}