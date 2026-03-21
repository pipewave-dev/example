import {
  PipewaveProvider,
  PipewaveModuleConfig,
  PipewaveDebugger,
} from "@pipewave/reactpkg";
import { Chat } from "./ChatComponent";

const config = new PipewaveModuleConfig({
  backendEndpoint: "localhost:8080/pipewave",
  insecure: true,
  getAccessToken: async () => "UserA",
  enableLongPollingFallback: true,
  heartbeatInterval: 30000,
});

export function UserA() {
  return (
    <PipewaveProvider config={config}>
      <Chat toUserId="UserB" />
      <PipewaveDebugger
        buttonLabel="PwDebuggerUserA"
        buttonPosition={{ bottom: 20, left: 20 }}
      />
    </PipewaveProvider>
  );
}
