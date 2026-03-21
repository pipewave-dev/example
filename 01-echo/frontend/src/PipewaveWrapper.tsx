import {
  PipewaveProvider,
  PipewaveModuleConfig,
  PipewaveDebugger,
} from "@pipewave/reactpkg";

const config = new PipewaveModuleConfig({
  backendEndpoint: "localhost:8080/pipewave",
  insecure: true,
  getAccessToken: async () => {
    return "default";
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
