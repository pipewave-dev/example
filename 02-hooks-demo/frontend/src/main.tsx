import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import PipewaveWrapper from "./PipewaveWrapper.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <PipewaveWrapper>
      <App />
    </PipewaveWrapper>
  </StrictMode>,
);
