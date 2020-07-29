import React from "react";
import { NormalizedSpan } from "./trace";
import { Action } from "./TraceView";
import "./Sidebar.css";

export function Sidebar(props: {
  span: NormalizedSpan;
  highlightedLogIdx: number;
  handleAction: (a: Action) => void;
}) {
  const { span, highlightedLogIdx } = props;
  return (
    <div style={{ padding: 10 }}>
      <table>
        <tbody style={{ textAlign: "left" }}>
          <tr>
            <th>Operation:</th>
            <td>{span.op}</td>
          </tr>
          <tr>
            <th>Start:</th>
            <td>{span.startedAt.toString()}</td>
          </tr>
          <tr>
            <th>Duration:</th>
            <td>
              {span.finishedAt
                ? span.finishedAt.diff(span.startedAt).milliseconds
                : span.startedAt.diffNow().milliseconds}
              ms
            </td>
          </tr>
        </tbody>
      </table>
      <h3>Log Messages</h3>
      <table className="table table-sm">
        <thead>
          <tr>
            <th>Age</th>
            <th>Message</th>
          </tr>
        </thead>
        <tbody>
          {span.logLines.map((logLine, idx) => {
            const highlighted = idx === highlightedLogIdx;
            return (
              <tr
                key={idx}
                className={`log-line ${
                  highlighted ? "log-line--highlighted" : ""
                }`}
                onMouseOver={() =>
                  props.handleAction({
                    type: "HOVER_LOG_LINE",
                    spanID: span.id,
                    logIdx: idx,
                  })
                }
              >
                <td style={{ whiteSpace: "nowrap", fontFamily: "monospace" }}>
                  {logLine.timestamp.toString()}
                </td>
                <td style={{ fontFamily: "monospace", whiteSpace: "pre" }}>
                  {logLine.line}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
