import React from "react";
import { denormalize, EMPTY_TRACE_DB, saveEvent, TraceDB } from "./trace";
import "./App.css";
import PanelLayout from "./util/PanelLayout";
import TraceView, {
  Action,
  EMPTY_TRACE_VIEW_STATE,
  TraceViewState,
  update
} from "./TraceView";
import { Sidebar } from "./Sidebar";

const wsAddr =
  process.env.NODE_ENV === "development"
    ? "localhost:8888"
    : window.location.host;
const wsURL = `ws://${wsAddr}/ws`;

type WebSocketState = "CONNECTING" | "OPEN" | "CLOSED";

interface AppState {
  db: TraceDB;
  wsState: WebSocketState;
  traceState: TraceViewState;
}

class App extends React.Component<{}, AppState> {
  constructor(props: {}) {
    super(props);
    this.state = {
      db: EMPTY_TRACE_DB,
      wsState: "CONNECTING",
      traceState: EMPTY_TRACE_VIEW_STATE
    };
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

  handleAction = (action: Action) => {
    console.log("trace view action:", action);
    this.setState(p => ({
      ...p,
      traceState: update(this.state.traceState, action)
    }));
  };

  render() {
    const denormalized = denormalize(this.state.db);
    return (
      <PanelLayout
        titleArea={<h1 style={{ fontSize: 20, margin: 5 }}>Trace</h1>}
        sidebar={
          this.state.traceState.hoveredSpanID ? (
            <Sidebar
              span={this.state.db.byID[this.state.traceState.hoveredSpanID]}
              handleAction={this.handleAction}
              highlightedLogIdx={
                this.state.traceState.highlightedLogLineBySpanID[
                  this.state.traceState.hoveredSpanID
                ]
              }
            />
          ) : null
        }
        mainContent={
          <>
            <p>WS State: {this.state.wsState}</p>
            {denormalized ? (
              <TraceView
                traces={denormalized}
                width={800}
                traceState={this.state.traceState}
                handleAction={this.handleAction}
              />
            ) : (
              "no trace yet"
            )}
          </>
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
