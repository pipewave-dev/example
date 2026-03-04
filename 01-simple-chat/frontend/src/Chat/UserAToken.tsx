

import { PipewaveProvider } from '@pipewave/reactpkg'
import { PipewaveModuleConfig } from '@pipewave/reactpkg'
import { Chat } from './ChatComponent'



const eventHandler = {
    onOpen: async () => {
        console.log('WebSocket connected')
    },
    onClose: async () => {
        console.log('WebSocket disconnected')
    },
    onError: async (error: Event) => {
        console.log('WebSocket error', error)
    },
}

// --- Config placeholder ---
const config = new PipewaveModuleConfig(
    {
        backendEndpoint: 'localhost:8080/websocket',
        insecure: true,
        debugMode: true,
        getAccessToken: async () => "UserA",
    }
)
// --- Page wrapper: should wrap your app with PipewaveProvider ---
export function UserA() {
    return (
        <PipewaveProvider config={config} eventHandler={eventHandler}>
            <Chat toUserId="UserB" />
        </PipewaveProvider>
    )
}

