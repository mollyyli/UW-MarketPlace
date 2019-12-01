addSession = async () => {
  const email = document.querySelector("#email-input").value;
  const pass = document.querySelector("#password-input").value;
  const passConf = document.querySelector("#passwordconf-input").value;
  const userName = document.querySelector("#username-input").value;
  const firstName = document.querySelector("#firstname-input").value;
  const lastName = document.querySelector("#lastname-input").value;

  await fetch(`https://api.briando.me/v1/users`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: {
      "Email": email,
      "Password": pass,
      "PasswordConf": passConf,
      "UserName": userName,
      "FirstName": firstName,
      "LastName": lastName
    }
  });

  const response = await fetch(`https://api.briando.me/v1/sessions`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: {
      "Email": email,
      "Password": pass
    }
  });

  return response.headers.Authorization.split("Bearer ")[0]
};



const token = addSession()
let socket = new WebSocket(`https://api.briando.me/ws?auth=${token}`);
socket.onopen = () => {
  console.log("Websocket connection open")
}
socket.onclose = () => {
  console.log("Websocket connection closed")
}
socket.onmessage = () => {

}


async function getData() {
  url = document.querySelector(".site-input").value;
  // Default options are marked with *
  const response = await fetch(
    `https://api.briando.me/v1/summary?url=${url}`
  )
  if (response.status == 200) {
    const data = await response.json();
    document.querySelector(".data").innerHTML = "";
    Object.keys(data).forEach(key => {
      if (key == "images") {
        data[key].forEach(image => {
          Object.keys(image).forEach(key => {
            if (key == "url") {
              const node = document.createElement("li");
              const text = document.createTextNode("IMAGE: ");
              const img = document.createElement("img");
              img.setAttribute("src", image[key]);
              node.appendChild(text);
              node.appendChild(img);
              document.querySelector(".data").appendChild(node);
            }
          });
        });
      } else if (key == "icon") {
        const node = document.createElement("li");
        const text = document.createTextNode("ICON: ");
        const img = document.createElement("img");
        img.setAttribute("src", data.icon.url);
        node.appendChild(text);
        node.appendChild(img);
        document.querySelector(".data").appendChild(node);
      } else {
        const node = document.createElement("li");
        const text = document.createTextNode(
          key.toUpperCase() + ": " + data[key]
        );
        node.appendChild(text);
        document.querySelector(".data").appendChild(node);
      }
    });
  } else {
    document.querySelector(".data").innerHTML = ""
    document.querySelector(".data").textContent = `${response.status}: ${response.statusText}`
  }
}
