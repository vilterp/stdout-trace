import { DateTime } from "luxon";

export type TraceEvent =
  | {
      trace_evt: "start_span";
      id: number;
      parent_id: number;
      ts: string;
      op: string;
    }
  | {
      trace_evt: "log";
      id: number;
      ts: string;
      line: string;
    }
  | {
      trace_evt: "finish_span";
      id: number;
      ts: string;
    };

export interface Span extends NormalizedSpan {
  children: Span[];
}

export interface LogLine {
  line: string;
  timestamp: DateTime;
}

export interface NormalizedSpan {
  id: number;
  parentID?: number;
  op: string;
  startedAt: DateTime;
  finishedAt?: DateTime;
  logLines: LogLine[];
}

export type TraceDB = {
  byID: { [id: number]: NormalizedSpan };
  byParentID: { [parentID: number]: number[] };
  roots: { [id: number]: true };
};

export const EMPTY_TRACE_DB: TraceDB = {
  byID: {},
  byParentID: {},
  roots: {}
};

function insert(db: TraceDB, span: NormalizedSpan): TraceDB {
  const newByID = { ...db.byID, [span.id]: span };
  const newByParentID = span.parentID
    ? {
        ...db.byParentID,
        [span.parentID]: [...(db.byParentID[span.parentID] || []), span.id]
      }
    : db.byParentID;
  // TODO: should probably use strings for ids and stop using these sentinel values
  const isRoot = !span.parentID;
  console.log({ isRoot });
  return {
    byID: newByID,
    byParentID: newByParentID,
    roots: isRoot ? { ...db.roots, [span.id]: true } : db.roots
  };
}

function update(
  db: TraceDB,
  id: number,
  up: (s: NormalizedSpan) => NormalizedSpan
): TraceDB {
  // assumes that id and parent id don't change
  const span = db.byID[id];
  if (!span) {
    throw new Error(`span not found with id ${id}`);
  }
  const newSpan = up(span);
  if (newSpan.id !== span.id) {
    throw new Error("can't change id");
  }
  if (newSpan.parentID !== span.parentID) {
    throw new Error("can't change parentID");
  }

  return {
    ...db,
    byID: {
      ...db.byID,
      [id]: newSpan
    }
  };
}

export function saveEvent(db: TraceDB, evt: TraceEvent): TraceDB {
  switch (evt.trace_evt) {
    case "start_span": {
      return insert(db, {
        id: evt.id,
        op: evt.op,
        parentID: evt.parent_id === -1 ? undefined : evt.parent_id,
        startedAt: parseTimestamp(evt.ts), // TODO: parse date???
        logLines: []
      });
    }
    case "finish_span": {
      return update(db, evt.id, span => ({
        ...span,
        finishedAt: parseTimestamp(evt.ts)
      }));
    }
    case "log": {
      return update(db, evt.id, span => ({
        ...span,
        logLines: [
          ...span.logLines,
          { timestamp: parseTimestamp(evt.ts), line: evt.line }
        ]
      }));
    }
  }
}

export function denormalize(db: TraceDB): Span[] {
  return Object.keys(db.roots).map(rootID => getSpan(db, parseInt(rootID)));
}

function getSpan(db: TraceDB, id: number): Span {
  return {
    ...db.byID[id],
    children: (db.byParentID[id] || []).map(childID => {
      if (childID === id) {
        throw new Error(`child id same as id: ${id}`);
      }
      return getSpan(db, childID);
    })
  };
}

function parseTimestamp(ts: string): DateTime {
  return DateTime.fromISO(ts);
}
