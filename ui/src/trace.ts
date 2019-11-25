import { DateTime } from "luxon";

export type TraceEvent =
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

export interface Span extends NormalizedSpan {
  children: Span[];
}

export interface LogLine {
  line: string;
  timestamp: DateTime;
}

interface NormalizedSpan {
  id: number;
  parentID?: number;
  op: string;
  startedAt: DateTime;
  finishedAt?: DateTime;
  logLines: LogLine[];
}

export type TraceDB = {
  byID: { [id: string]: NormalizedSpan };
  byParentID: { [parentID: string]: number[] };
};

export const EMPTY_TRACE_DB: TraceDB = {
  byID: {},
  byParentID: {}
};

function insert(db: TraceDB, span: NormalizedSpan): TraceDB {
  const newByID = { ...db.byID, [span.id]: span };
  const newByParentID = span.parentID
    ? {
        ...db.byParentID,
        [span.parentID]: [...(db.byParentID[span.parentID] || []), span.id]
      }
    : db.byParentID;
  return {
    byID: newByID,
    byParentID: newByParentID
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
  switch (evt.evt) {
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

export function denormalize(db: TraceDB): Span {
  return getSpan(db, 1); // TODO: assuming root span id always 1...
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

// TODO: could index by whether finished or not
export function allFinished(db: TraceDB): boolean {
  return !Object.values(db.byID).some(v => !v.finishedAt);
}

function parseTimestamp(ts: string): DateTime {
  return DateTime.fromISO(ts);
}
