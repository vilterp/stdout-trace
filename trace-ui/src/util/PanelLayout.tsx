import React from "react";

import "./PanelLayout.css";

function PanelLayout(props: {
  mainContent: React.ReactNode;
  sidebar: React.ReactNode;
  titleArea: React.ReactNode;
}) {
  return (
    <div className="panel-layout">
      <div className="panel-layout__title-area">{props.titleArea}</div>
      <div className="panel-layout__main-content">{props.mainContent}</div>
      <div className="panel-layout__sidebar">{props.sidebar}</div>
    </div>
  );
}

export default PanelLayout;
