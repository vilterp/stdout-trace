import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";

const wsAddr =
  process.env.NODE_ENV === "development"
    ? "localhost:8888"
    : window.location.host;
const wsURL = `ws://${wsAddr}/ws`;

const ws = new WebSocket(wsURL);
ws.addEventListener("message", evt => {
  console.log("message", evt);
});
ws.addEventListener("error", evt => {
  console.log("error", evt);
});
ws.addEventListener("close", evt => {
  console.log("close", evt);
});
ws.addEventListener("open", evt => {
  console.log("open", evt);
});

ReactDOM.render(<App />, document.getElementById("root"));
