const nameField = document.getElementById("nameField");
const nameBtn = document.getElementById("nameBtn");
const resetBtn = document.getElementById("reset");
const showBtn = document.getElementById("show");
const cards = document.getElementById("cards");
const choices = document.getElementById("choices");

const vals = [1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144];
vals.forEach((v) => {
  el = document.createElement("btn");
  el.classList.add("card");
  el.innerText = v;
  el.addEventListener("click", () => {
    socket.send(JSON.stringify({ cmd: "pick", choice: v }));
  });
  cards.appendChild(el);
});

let socket = new WebSocket(
  `${location.protocol === "https:" ? "wss" : "ws"}://${location.host}/ws`
);

let data = {};

socket.addEventListener("error", (e) => {
  alert("The websocket ran into an error. Please refresh and try again.");
});

socket.addEventListener("message", (e) => {
  data = JSON.parse(e.data);
  console.log(data);

  if (data.show) {
    showBtn.hidden = true;
    resetBtn.hidden = false;
  } else {
    showBtn.hidden = false;
    resetBtn.hidden = true;
  }

  choices.innerHTML = null;
  data.clients.forEach((c) => {
    el = document.createElement("div");
    el.classList.add("choice");
    el.innerHTML = `
      <div class="choice">
        <div class="card ${c.choice == 0 ? "" : "selected"}">
          ${data.show ? c.choice : "?"}
        </div>
        <div>${c.name}</div>
      </div>
    `;
    choices.appendChild(el);

    if (c.name === nameField.value) {
      [...cards.children].forEach((v) => {
        if (v.innerText == c.choice) {
          v.classList.add("selected");
        } else {
          v.classList.remove("selected");
        }
      });
    }
  });
});

nameBtn.addEventListener("click", () => {
  socket.send(JSON.stringify({ cmd: "name", name: nameField.value }));
});

showBtn.addEventListener("click", () => {
  socket.send(JSON.stringify({ cmd: "show" }));
});

resetBtn.addEventListener("click", () => {
  socket.send(JSON.stringify({ cmd: "reset" }));
});
