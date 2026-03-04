
import { useMemo, useState, useRef } from 'react'
import { usePipewave, type OnMessage } from '@pipewave/reactpkg'
import { encode, decode } from '@msgpack/msgpack'

// --- Message types from backend ---
const MSG_TYPE_CHAT_SEND_MSG = "CHAT_SEND_MSG"
const MSG_TYPE_CHAT_TYPING = "CHAT_TYPING"

const MSG_TYPE_CHAT_INCOMING_MSG = "CHAT_INCOMING_MSG"
const MSG_TYPE_CHAT_USER_TYPING = "CHAT_USER_TYPING"
const MSG_TYPE_CHAT_ACK = "CHAT_ACK"
const MSG_TYPE_ECHO_RESPONSE = "ECHO_RESPONSE"
const MSG_TYPE_CHAT_FAIL = "CHAT_FAIL"

interface Message {
    id: string
    from: string
    text: string
    timestamp: number
    isMe: boolean
}

// --- Main component using websocket ---
export function Chat({ toUserId }: { toUserId: string }) {
    const [messages, setMessages] = useState<Message[]>([])
    const [input, setInput] = useState('')
    const [typingUser, setTypingUser] = useState<string | null>(null)
    const [toUserOffline, setToUserOffline] = useState(false)
    const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    const onMessage: OnMessage = useMemo(() => ({
        [MSG_TYPE_CHAT_INCOMING_MSG]: async (data: Uint8Array, id: string) => {
            const payload = decode(data) as { from_user_id: string; content: string; timestamp: number }
            setMessages(prev => [...prev, {
                id,
                from: payload.from_user_id,
                text: payload.content,
                timestamp: payload.timestamp,
                isMe: false
            }])
            setTypingUser(null) // Stop typing indicator when message arrives
            setToUserOffline(false)
        },
        [MSG_TYPE_CHAT_USER_TYPING]: async (data: Uint8Array) => {
            const payload = decode(data) as { from_user_id: string }
            setTypingUser(payload.from_user_id)

            // Clear indicator after 3 seconds of no typing updates
            if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current)
            typingTimeoutRef.current = setTimeout(() => {
                setTypingUser(null)
            }, 3000)
        },
        [MSG_TYPE_CHAT_ACK]: async () => {
            setToUserOffline(false)
        },
        [MSG_TYPE_ECHO_RESPONSE]: async (data: Uint8Array, id: string) => {
            const text = new TextDecoder().decode(data)
            setMessages(prev => [...prev, {
                id,
                from: 'System',
                text,
                timestamp: Date.now() / 1000,
                isMe: false
            }])
        },
        [MSG_TYPE_CHAT_FAIL]: async () => {
            setToUserOffline(true)
        },
    }), [])

    const { status, send, resetRetryCount } = usePipewave(onMessage)

    const handleSend = () => {
        if (!input.trim()) return

        const payload = {
            to_user_id: toUserId,
            content: input,
        }

        const msgId = crypto.randomUUID()
        send({
            id: msgId,
            msgType: MSG_TYPE_CHAT_SEND_MSG,
            data: encode(payload),
        })

        setMessages(prev => [...prev, {
            id: msgId,
            from: 'Me',
            text: input,
            timestamp: Date.now() / 1000,
            isMe: true,
        }])

        setInput('')

    }

    const sendTyping = () => {
        send({
            id: crypto.randomUUID(),
            msgType: MSG_TYPE_CHAT_TYPING,
            data: encode({
                to_user_id: toUserId,
            }),
        })
    }

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInput(e.target.value)
        if (e.target.value.length > 0) {
            sendTyping()
        }
    }

    return (
        <div style={{ padding: 24, maxWidth: 600, margin: '0 auto', fontFamily: 'sans-serif' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #eee', paddingBottom: 12 }}>
                <div>
                    <h2 style={{ margin: 0 }}>Simple Chat</h2>
                    <div style={{ fontSize: '0.85rem', color: '#666' }}>
                        Status: <span style={{ color: status === 'READY' ? 'green' : 'red', fontWeight: 'bold' }}>{status}</span>
                        {status === 'SUSPEND' && <button onClick={resetRetryCount} style={{ marginLeft: 8 }}>Retry</button>}
                    </div>
                </div>
            </div>

            <div style={{
                height: 400,
                overflowY: 'auto',
                border: '1px solid #eee',
                borderRadius: 8,
                padding: 16,
                margin: '16px 0',
                backgroundColor: '#f9f9f9',
                display: 'flex',
                flexDirection: 'column',
                gap: 12
            }}>
                {messages.length === 0 && <p style={{ textAlign: 'center', color: '#999' }}>No messages yet.</p>}
                {messages.map((msg) => (
                    <div key={msg.id} style={{
                        alignSelf: msg.isMe ? 'flex-end' : 'flex-start',
                        maxWidth: '80%',
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: msg.isMe ? 'flex-end' : 'flex-start'
                    }}>
                        <div style={{
                            padding: '8px 12px',
                            borderRadius: 12,
                            backgroundColor: msg.isMe ? '#007AFF' : '#E9E9EB',
                            color: msg.isMe ? 'white' : 'black',
                            fontSize: '0.95rem'
                        }}>
                            {msg.text}
                        </div>
                        <small style={{ color: '#888', fontSize: '0.7rem', marginTop: 4 }}>
                            {msg.from} • {new Date(msg.timestamp * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                        </small>
                    </div>
                ))}
                {typingUser && (
                    <div style={{ color: '#888', fontSize: '0.8rem', fontStyle: 'italic' }}>
                        {typingUser} is typing...
                    </div>
                )}
            </div>

            {toUserOffline && (
                <div style={{ marginBottom: 8, padding: '6px 12px', borderRadius: 8, backgroundColor: '#fff3cd', color: '#856404', fontSize: '0.85rem' }}>
                    {toUserId} đang offline
                </div>
            )}
            <div style={{ display: 'flex', gap: 8 }}>
                <input
                    value={input}
                    onChange={handleInputChange}
                    onKeyDown={e => e.key === 'Enter' && handleSend()}
                    placeholder="Type a message..."
                    style={{ flex: 1, padding: '10px 16px', borderRadius: 20, border: '1px solid #ccc', outline: 'none' }}
                />
                <button
                    onClick={handleSend}
                    disabled={status !== 'READY' || !input.trim()}
                    style={{
                        padding: '10px 20px',
                        borderRadius: 20,
                        border: 'none',
                        backgroundColor: '#007AFF',
                        color: 'white',
                        fontWeight: 'bold',
                        cursor: status === 'READY' ? 'pointer' : 'default',
                        opacity: status === 'READY' ? 1 : 0.5
                    }}
                >
                    Send
                </button>
            </div>
        </div>
    )
}
