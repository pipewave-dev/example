import { useState } from "react";
import { useDebugLogger } from "@pipewave/reactpkg/hooks";

export function DebugLoggerExample() {
  const [enabled, setEnabled] = useState(true);
  useDebugLogger(enabled);

  return (
    <div>
      <p>Open browser DevTools → Console to see debug logs.</p>
      <label>
        <input
          type="checkbox"
          checked={enabled}
          onChange={(e) => setEnabled(e.target.checked)}
        />{" "}
        Enable debug logger
      </label>
    </div>
  );
}
