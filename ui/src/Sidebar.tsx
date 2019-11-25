import React from "react";
import { NormalizedSpan } from "./trace";

export function Sidebar(props: { span: NormalizedSpan }) {
  const span = props.span;
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
                ? span.finishedAt.diff(span.startedAt).toISO()
                : span.startedAt.diffNow().toISO()}
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
          {span.logLines.map((logLine, idx) => (
            <tr key={idx}>
              <td>{logLine.timestamp.toString()}</td>
              <td style={{ fontFamily: "monospace", whiteSpace: "pre" }}>
                {logLine.line}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
