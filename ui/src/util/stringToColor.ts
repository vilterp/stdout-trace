import * as d3chromatic from "d3-scale-chromatic";

function hashString(s: string) {
  var hash = 0;
  for (var i = 0; i < s.length; i++) {
    hash = (hash << 5) - hash + s.charCodeAt(i);
    hash |= 0;
  }
  return Math.abs(hash);
}

const colors = d3chromatic.schemeSpectral[11];

export function stringToColor(s: string) {
  return colors[hashString(s) % colors.length];
}
