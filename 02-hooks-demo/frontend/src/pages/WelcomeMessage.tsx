import { usePipewave } from "@pipewave/reactpkg/hooks";
import { useState } from "react";

const WELCOME = "WELCOME";
const decoder = new TextDecoder();

export function WelcomeMessage() {
    const [msg, setMsg] = useState<string[]>([]);

    usePipewave(
        {
            [WELCOME]: async (data) => {
                setMsg((prev) => [...prev, `← WELCOME: ${decoder.decode(data)}`]);
            },
        },
    );
    return (
        <div>
            {msg.map((m, i) => (
                <div key={i}>{m}</div>
            ))}
        </div>
    );
}
