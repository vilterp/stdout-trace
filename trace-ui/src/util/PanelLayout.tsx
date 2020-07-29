import React, { useState } from "react";

import "./PanelLayout.css";

function PanelLayout(props: {
  mainContent: React.ReactNode;
  sidebar: React.ReactNode;
  titleArea: React.ReactNode;
}) {
  const [sidebarWidthPx, setSidebarWidthPx] = useState(300);
  const [dragging, setDragging] = useState(false);

  return (
    <div
      className="panel-layout"
      style={{
        gridTemplateColumns: `auto 5px ${sidebarWidthPx}px`,
        userSelect: dragging ? "none" : "inherit",
        cursor: dragging ? "grabbing" : "default",
      }}
      onMouseMove={(evt) => {
        if (dragging) {
          setSidebarWidthPx(window.innerWidth - evt.clientX);
        }
      }}
      onMouseUp={() => {
        if (dragging) {
          setDragging(false);
        }
      }}
    >
      <div className="panel-layout__title-area">{props.titleArea}</div>
      <div className="panel-layout__main-content">{props.mainContent}</div>
      <div
        className="panel-layout__splitter"
        onMouseDown={() => setDragging(true)}
        style={{ cursor: dragging ? "inherit" : "grab" }}
      />
      <div className="panel-layout__sidebar">{props.sidebar}</div>
    </div>
  );
}

export default PanelLayout;
