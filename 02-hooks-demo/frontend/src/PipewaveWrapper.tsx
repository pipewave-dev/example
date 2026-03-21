import {
  PipewaveProvider,
  PipewaveModuleConfig,
  PipewaveDebugger,
} from "@pipewave/reactpkg";

const accessToken = { value: "default" };
const config = new PipewaveModuleConfig({
  backendEndpoint: "localhost:8080/pipewave",
  enableLongPollingFallback: true,
  insecure: true,
  getAccessToken: async () => accessToken.value,
  retry: {
    maxRetry: 5,
    initialRetryDelay: 1000,
    maxRetryDelay: 3000,
  },
});

export default function PipewaveWrapper(props: React.PropsWithChildren) {
  return (
    <PipewaveProvider config={config}>
      {props.children}
      {/* Optional: Include the debugger for development purposes */}
      <PipewaveDebugger />
    </PipewaveProvider>
  );
}
