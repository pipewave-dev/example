import {
  usePipewaveStatus,
  usePipewaveResetConnection,
} from "@pipewave/reactpkg/hooks";

export function ResetConnectionExample() {
  const { status, isSuspended } = usePipewaveStatus();
  const { resetRetryCount } = usePipewaveResetConnection();

  return (
    <div>
      <p>
        When status is <code>SUSPEND</code> (max retries hit), click Reset to
        allow reconnection attempts again.
      </p>
      <p>
        Current status: <strong>{status}</strong>
      </p>
      <button onClick={resetRetryCount} disabled={!isSuspended}>
        Reset Retry Count
      </button>
      {!isSuspended && (
        <p style={{ color: "#aaa", fontSize: 13 }}>
          (button active only when suspended)
        </p>
      )}
    </div>
  );
}
