import {
  usePipewaveStatus,
  usePipewaveSend,
  usePipewaveResetConnection,
  usePipewaveMessage,
} from "@pipewave/reactpkg";
import { useState } from "react";

const Encoder = new TextEncoder();
const Decoder = new TextDecoder();

interface InputAreaProps {
  isConnected: boolean;
}

function InputArea({ isConnected }: InputAreaProps) {
  const [input, setInput] = useState("");
  const { send } = usePipewaveSend();
  const handleSend = (text: string) => {
    send({
      id: crypto.randomUUID(),
      msgType: "ECHO_REQ",
      data: Encoder.encode(text),
    });
    setInput("");
  };

  return (
    <div>
      <input
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && handleSend(input)}
        placeholder="Type a message..."
      />
      <button
        onClick={() => handleSend(input)}
        disabled={!isConnected || !input.trim()}
      >
        Send
      </button>
    </div>
  );
}

export function EchoApp() {
  const [messages, setMessages] = useState<string[]>([]);

  const { status, isConnected, isSuspended } = usePipewaveStatus();
  const { resetRetryCount } = usePipewaveResetConnection();

  usePipewaveMessage("ECHO_RES", async (data: Uint8Array) => {
    setMessages((prev) => [...prev, Decoder.decode(data)]);
  });

  return (
    <div>
      <InputArea isConnected={isConnected} />

      <p>
        Status:{" "}
        <span style={{ color: isConnected ? "green" : "red" }}>{status}</span>
      </p>
      {isSuspended && <button onClick={resetRetryCount}>Retry</button>}
      {messages.map((msg, i) => (
        <p key={i}>{msg}</p>
      ))}
    </div>
  );
}
