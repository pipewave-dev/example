import {
  PipewaveProvider,
  PipewaveModuleConfig,
  PipewaveDebugger,
} from "@pipewave/reactpkg";
import { Chat } from "./ChatComponent";

const config = new PipewaveModuleConfig({
  backendEndpoint: "localhost:8080/pipewave",
  insecure: true,
  getAccessToken: async () => "UserB",
  enableLongPollingFallback: true,
  heartbeatInterval: 30000,
});

export function UserB() {
  return (
    <PipewaveProvider config={config}>
      <Chat toUserId="UserA" />
      <PipewaveDebugger
        buttonLabel="PwDebuggerUserB"
        buttonPosition={{ bottom: 20, right: 20 }}
        panelSide="right"
      />
    </PipewaveProvider>
  );
}
