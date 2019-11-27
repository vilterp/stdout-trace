import _ from "lodash";
import React, { Component } from "react";
import { Span } from "./trace";
import * as d3scale from "d3-scale";
import { DateTime } from "luxon";
import { stringToColor } from "./util/stringToColor";

function numDescendants(s: Span): number {
  return (
    s.children.length +
    s.children.map(c => numDescendants(c)).reduce((a, b) => a + b, 0) // TODO: .sum() would be nice
  );
}

const HEIGHT = 30;
const HEIGHT_PLUS_SPACE = HEIGHT + 5;

const DOWN_ARROW = "▼";
const SIDE_ARROW = "▶";

type Action =
  | { type: "TOGGLE_COLLAPSED"; spanID: number }
  | { type: "HOVER_SPAN"; spanID: number }
  | { type: "UN_HOVER_SPAN" };

interface TraceViewProps {
  traces: Span[];
  width: number;
  traceState: TraceViewState;
  handleAction: (action: Action) => void;
}

export interface TraceViewState {
  hoveredSpanID: number | null;
  collapsedSpanIDs: number[];
}

export const EMPTY_TRACE_VIEW_STATE: TraceViewState = {
  hoveredSpanID: null,
  collapsedSpanIDs: []
};

export function update(state: TraceViewState, action: Action): TraceViewState {
  switch (action.type) {
    case "TOGGLE_COLLAPSED": {
      const isCollapsed = _.includes(state.collapsedSpanIDs, action.spanID);
      return {
        ...state,
        collapsedSpanIDs: isCollapsed
          ? state.collapsedSpanIDs.filter(spanID => spanID !== action.spanID)
          : [...state.collapsedSpanIDs, action.spanID]
      };
    }
    case "HOVER_SPAN":
      return {
        ...state,
        hoveredSpanID: action.spanID
      };
    case "UN_HOVER_SPAN":
      return {
        ...state,
        hoveredSpanID: null
      };
    default:
      return state;
  }
}

// preorder traversal, excluding collapsed children
function flatten(tree: Span, collapsed: number[]): Span[] {
  const output: Span[] = [];
  function recur(node: Span) {
    output.push(node);
    if (_.includes(collapsed, node.id)) {
      return;
    }
    if (node.children) {
      node.children.forEach(child => {
        recur(child);
      });
    }
  }
  recur(tree);
  return output;
}

interface St {
  now: DateTime;
  nowIntervalID: NodeJS.Timeout;
}

class TraceView extends Component<TraceViewProps, St> {
  componentDidMount() {
    const nowIntervalID = setInterval(() => {
      this.setState({
        now: DateTime.local()
      });
    }, 16);
    this.setState({
      now: DateTime.local(),
      nowIntervalID: nowIntervalID
    });
  }

  componentWillUnmount() {
    clearInterval(this.state.nowIntervalID);
  }

  handleAction = (action: Action) => {
    this.props.handleAction(action);
  };

  render() {
    const { traces, width, traceState } = this.props;
    const { collapsedSpanIDs, hoveredSpanID } = traceState;
    const flattened = _.flatten(traces.map(t => flatten(t, collapsedSpanIDs)));

    if (traces.length === 0) {
      return null;
    }

    const firstTS = (
      _.minBy(traces, t => t.startedAt.toMillis()) || {
        startedAt: DateTime.local()
      }
    ).startedAt;
    const lastTS =
      _.max(
        flattened.map(span =>
          span.finishedAt
            ? span.finishedAt.toMillis()
            : DateTime.local().toMillis()
        )
      ) || DateTime.local();

    const scale = d3scale
      .scaleLinear()
      .domain([firstTS, lastTS])
      .range([0, width]);

    return (
      <svg
        height={flattened.length * HEIGHT_PLUS_SPACE}
        style={{ width: "100%", minWidth: width }}
      >
        {flattened.map((span, idx) => {
          const isFinished = !!span.finishedAt;
          const isHovered = hoveredSpanID === span.id;
          const isCollapsed = _.includes(collapsedSpanIDs, span.id);
          // TODO: is this nanos?
          const isLeaf = span.children.length === 0;
          const label = isLeaf
            ? `${span.op}`
            : isCollapsed
            ? `${SIDE_ARROW} ${span.op} (${numDescendants(span)})`
            : `${DOWN_ARROW} ${span.op}`;
          const startTS = span.startedAt.toMillis();
          const endTS = span.finishedAt
            ? span.finishedAt.toMillis()
            : DateTime.local();
          const color = stringToColor(span.op);
          return (
            <g
              key={span.id}
              style={{ cursor: "pointer" }}
              onMouseOver={() => {
                this.handleAction({ type: "HOVER_SPAN", spanID: span.id });
              }}
              onClick={
                isLeaf
                  ? () => {}
                  : () => {
                      this.handleAction({
                        type: "TOGGLE_COLLAPSED",
                        spanID: span.id
                      });
                    }
              }
            >
              <rect
                fill="white"
                x={0}
                y={idx * HEIGHT_PLUS_SPACE - 5}
                height={HEIGHT}
                width={width}
              />
              <rect
                fill={color}
                y={idx * HEIGHT_PLUS_SPACE - 5}
                x={scale(startTS)}
                stroke={isFinished ? "black" : "red"}
                strokeWidth={2}
                height={HEIGHT}
                width={Math.max(scale(endTS) - scale(startTS), 1)}
              />
              <text
                x={scale(startTS) + 5}
                y={idx * HEIGHT_PLUS_SPACE + HEIGHT / 2}
                style={{ textDecoration: isHovered ? "underline" : "none" }}
              >
                {label}
              </text>
              <g>
                {span.logLines.map((logEntry, logIdx) => (
                  <circle
                    key={logIdx}
                    cx={scale(logEntry.timestamp.toMillis())}
                    cy={idx * HEIGHT_PLUS_SPACE + 20}
                    r={3}
                    fill={"white"}
                    stroke={"black"}
                  />
                ))}
              </g>
            </g>
          );
        })}
      </svg>
    );
  }
}

export default TraceView;
