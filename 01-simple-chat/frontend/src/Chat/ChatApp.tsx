
import { UserA } from "./UserAToken"
import { UserB } from "./UserBToken"

// --- Page wrapper: should wrap your app with PipewaveProvider ---
export default function ChatApp() {
    return (
        <div style={{ display: 'flex', flexDirection: 'row', width: '100%', height: '100vh' }}>
            <div style={{ flex: 1 }}>
                <UserA />
            </div>
            <div style={{ flex: 1 }}>
                <UserB />
            </div>
        </div>
    )
}