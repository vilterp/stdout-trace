import React, { useEffect, useState } from "react";
import "./App.css";

const wsAddr =
  process.env.NODE_ENV === "development"
    ? "localhost:8888"
    : window.location.host;
const wsURL = `ws://${wsAddr}/ws`;

type TraceEvent =
  | {
      evt: "start_span";
      id: number;
      parent_id: number;
      ts: string;
      op: string;
    }
  | {
      evt: "log";
      id: number;
      ts: string;
      line: string;
    }
  | {
      evt: "finish_span";
      id: number;
      ts: string;
    };

class Appc extends React.Component<{}, { events: TraceEvent[] }> {
  constructor(props: {}) {
    super(props);
    this.state = { events: [] };
  }

  componentDidMount() {
    const ws = new WebSocket(wsURL);
    ws.addEventListener("message", evt => {
      this.setState({ events: [...this.state.events, JSON.parse(evt.data)] });
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
  }

  render() {
    return (
      <div>
        <h1>Messages</h1>
        <ul>
          {this.state.events.map((msg, idx) => (
            <li key={idx.toString()}>
              <pre>{JSON.stringify(msg)}</pre>
            </li>
          ))}
        </ul>
      </div>
    );
  }
}

export default Appc;

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
