import React from "react";
import {
  allFinished,
  denormalize,
  EMPTY_TRACE_DB,
  saveEvent,
  TraceDB
} from "./trace";
import "./App.css";
import PanelLayout from "./util/PanelLayout";

const wsAddr =
  process.env.NODE_ENV === "development"
    ? "localhost:8888"
    : window.location.host;
const wsURL = `ws://${wsAddr}/ws`;

type WebSocketState = "CONNECTING" | "OPEN" | "CLOSED";

interface AppState {
  db: TraceDB;
  wsState: WebSocketState;
}

class App extends React.Component<{}, AppState> {
  constructor(props: {}) {
    super(props);
    this.state = { db: EMPTY_TRACE_DB, wsState: "CONNECTING" };
  }

  componentDidMount() {
    // TODO: try to reconnect
    const ws = new WebSocket(wsURL);
    ws.addEventListener("message", evt => {
      const traceEvt = JSON.parse(evt.data);
      const newDB = saveEvent(this.state.db, traceEvt);
      console.log("UPDATE", this.state.db, traceEvt, "=>", newDB);
      this.setState({ db: newDB });
    });
    ws.addEventListener("error", evt => {
      console.log("error", evt);
      this.setState(s => ({
        ...s,
        wsState: "CLOSED"
      }));
    });
    ws.addEventListener("close", evt => {
      console.log("close", evt);
      this.setState(s => ({
        ...s,
        wsState: "CLOSED"
      }));
    });
    ws.addEventListener("open", evt => {
      console.log("open", evt);
      this.setState(s => ({
        ...s,
        wsState: "OPEN"
      }));
    });
  }

  render() {
    return (
      <PanelLayout
        titleArea={<h1 style={{ fontSize: 20, margin: 5 }}>Trace</h1>}
        sidebar={<p>Sup from the Sidebar</p>}
        mainContent={
          <div>
            <p>WS State: {this.state.wsState}</p>
            <p>
              Trace state:{" "}
              {allFinished(this.state.db) ? "finished" : "in progress"}
            </p>
            <pre>{JSON.stringify(denormalize(this.state.db), null, 2)}</pre>
          </div>
        }
      />
    );
  }
}

export default App;

// TODO(vilterp): figure out how to use hooks for this
// function App() {
//   const [messages, setMessages] = useState<TraceEvent[]>([]);
//
//   useEffect(() => {
//     console.log("connecting");
//     const ws = new WebSocket(wsURL);
//     ws.addEventListener("message", evt => {
//       setMessages([...messages, evt.data]);
//     });
//     ws.addEventListener("error", evt => {
//       console.log("error", evt);
//     });
//     ws.addEventListener("close", evt => {
//       console.log("close", evt);
//     });
//     ws.addEventListener("open", evt => {
//       console.log("open", evt);
//     });
//   });
//
//   return (
//     <div>
//       <h1>Messages</h1>
//       <ul>
//         {messages.map((idx, msg) => (
//           <li key={idx.toString()}>
//             <pre>{JSON.stringify(msg)}</pre>
//           </li>
//         ))}
//       </ul>
//     </div>
//   );
// }
