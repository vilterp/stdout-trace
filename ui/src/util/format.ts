export function formatNanos(nanos: number): string {
  const justNS = Math.round(nanos % 1000);
  const us = nanos / 1000;
  const justUS = Math.round(us % 1000);
  const ms = us / 1000;
  const justMS = Math.round(ms % 1000);
  const sec = ms / 1000;
  const justSec = Math.round(sec % 1000);

  if (justSec > 0) {
    // s to ns
    if (justNS > 0) {
      return `${justSec}s${justMS}ms${justUS}µs${justNS}ns`;
    }
    // s to µs
    if (justNS > 0) {
      return `${justSec}s${justMS}ms${justUS}µs`;
    }
    // s to ms
    if (justMS > 0) {
      return `${justSec}s${justMS}ms`;
    }
    // s only
    return `${justSec}s`;
  }

  if (justMS > 0) {
    // ms to ns
    if (justNS > 0) {
      return `${justMS}ms${justUS}µs${justNS}ns`;
    }
    // ms to µs
    if (justNS > 0) {
      return `${justMS}ms${justUS}µs`;
    }
    // ms only
    return `${justMS}ms`;
  }

  if (justUS > 0) {
    // µs to ns
    if (justNS > 0) {
      return `${justUS}µs${justNS}ns`;
    }
    // µs only
    return `${justUS}µs`;
  }

  // ns only
  if (justNS > 0) {
    return `${justNS}ns`;
  }

  // The default value for 0 is 0µs.
  return `0µs`;
}
