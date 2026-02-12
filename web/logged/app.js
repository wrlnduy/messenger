let ws = null;
let currentChatId = null;

const chatListEl = document.getElementById("chatList");
const messagesEl = document.getElementById("messages");
const inputEl = document.getElementById("messageInput");
const directInput = document.getElementById("directInput");
const directBtn = document.getElementById("directBtn");

async function loadChats() {
  const res = await fetch("/logged/chats");
  const data = await res.json();

  chatListEl.innerHTML = "";

  data.chats.forEach(chat => {
    const li = document.createElement("li");

    li.textContent = chat.title || chat.chatId;
    li.dataset.id = chat.chatId;

    li.onclick = () => selectChat(chat.chatId, li);

    chatListEl.appendChild(li);
  });

}

async function selectChat(chatId, element) {
  currentChatId = chatId;

  document.querySelectorAll("#chatList li")
    .forEach(li => li.classList.remove("active"));

  element.classList.add("active");

  if (ws) {
    ws.close();
    ws = null;
  }

  await loadHistory(chatId);
  connectWebSocket(chatId);
}

async function loadHistory(chatId) {
  const res = await fetch(`/logged/history?chat_id=${chatId}`);
  const data = await res.json();

  messagesEl.innerHTML = "";

  const users = data.mapping || {};

  data.messages.forEach(msg => {
    const username = msg.username || users[msg.userId] || "Unknown";
    addMessage(username, msg.text, msg.timestamp);
  });
}

function connectWebSocket(chatId) {
  const protocol = location.protocol === "https:" ? "wss" : "ws";

  ws = new WebSocket(`${protocol}://${location.host}/logged/ws?chat_id=${chatId}`);

  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    addMessage(msg.username, msg.text, msg.timestamp);
  };

  ws.onclose = () => {
    console.log("WebSocket closed");
  };
}

async function sendMessage() {
  if (!currentChatId) return;

  const text = inputEl.value.trim();
  if (!text) return;

  await fetch(`/logged/message?chat_id=${currentChatId}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ text })
  });

  inputEl.value = "";
}

inputEl.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    sendMessage();
  }
});

function addMessage(username, text, timestamp) {
  const div = document.createElement("div");

  let timeStr = "";
  if (timestamp) {
    const ts = new Date(timestamp);
    timeStr = ` <span style="color:gray;font-size:12px;">(${ts.toLocaleString()})</span>`;
  }

  div.innerHTML = `<strong>${username}:</strong> ${text}${timeStr}`;

  messagesEl.appendChild(div);
  messagesEl.scrollTop = messagesEl.scrollHeight;
}

async function createDirect() {
  const username = directInput.value.trim();
  if (!username) return;

  const res = await fetch("/logged/direct", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username })
  });

  if (!res.ok) {
    alert("User not found");
    return;
  }

  const chat = await res.json();

  directInput.value = "";

  await loadChats();

  setTimeout(() => {
    const items = document.querySelectorAll("#chatList li");
    items.forEach(li => {
      if (li.dataset.id === chat.chatId) {
        li.click();
      }
    });
  }, 50);
}

directBtn.onclick = createDirect;

directInput.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    createDirect();
  }
});

loadChats();
