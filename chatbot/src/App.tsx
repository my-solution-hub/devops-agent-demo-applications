import { useState, useRef, useEffect, FormEvent } from 'react'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

interface Message {
  id: number
  sender: 'user' | 'bot'
  text: string
  timestamp: Date
}

function formatJson(data: unknown): string {
  return JSON.stringify(data, null, 2)
}

async function handleCommand(input: string): Promise<string> {
  const lower = input.toLowerCase().trim()

  try {
    if (lower.includes('list cats')) {
      const res = await fetch(`${API_URL}/api/cats`)
      if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      const data = await res.json()
      if (Array.isArray(data) && data.length === 0) {
        return '📋 No cats found. Try "add cat" to create one!'
      }
      return `📋 Cats:\n\`\`\`\n${formatJson(data)}\n\`\`\``
    }

    if (lower.includes('list devices')) {
      const res = await fetch(`${API_URL}/api/devices`)
      if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      const data = await res.json()
      if (Array.isArray(data) && data.length === 0) {
        return '📋 No devices found. Try "add device" to register one!'
      }
      return `📋 Devices:\n\`\`\`\n${formatJson(data)}\n\`\`\``
    }

    if (lower.includes('add cat')) {
      const res = await fetch(`${API_URL}/api/cats`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: 'Whiskers',
          breed: 'Maine Coon',
          ageMonths: 24,
          weightKg: 5.5,
          dietaryRestrictions: ['grain-free'],
          ownerId: 'demo-user',
        }),
      })
      if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      const data = await res.json()
      return `✅ Cat created!\n\`\`\`\n${formatJson(data)}\n\`\`\``
    }

    if (lower.includes('add device')) {
      const res = await fetch(`${API_URL}/api/devices`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          device_type: 'feeder',
          name: 'Kitchen Auto-Feeder',
          status: 'online',
          config: { portion_grams: 50 },
        }),
      })
      if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      const data = await res.json()
      return `✅ Device registered!\n\`\`\`\n${formatJson(data)}\n\`\`\``
    }

    return (
      "🤔 I don't understand that yet. Try one of these commands:\n\n" +
      '• **list cats** — show all cat profiles\n' +
      '• **list devices** — show all registered devices\n' +
      '• **add cat** — create a sample cat profile\n' +
      '• **add device** — register a sample device'
    )
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err)
    return `❌ Error: ${message}\n\nMake sure the API Gateway is running at ${API_URL}`
  }
}

function renderMessageText(text: string) {
  const parts = text.split(/(`{3}[\s\S]*?`{3}|\*\*.*?\*\*)/g)
  return parts.map((part, i) => {
    if (part.startsWith('```') && part.endsWith('```')) {
      const code = part.slice(3, -3).replace(/^\n/, '')
      return (
        <pre key={i} className="code-block">
          <code>{code}</code>
        </pre>
      )
    }
    if (part.startsWith('**') && part.endsWith('**')) {
      return <strong key={i}>{part.slice(2, -2)}</strong>
    }
    return <span key={i}>{part}</span>
  })
}

export default function App() {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 0,
      sender: 'bot',
      text:
        "👋 Welcome to the Smart Home Cat Demo!\n\n" +
        "I'm a simple command interpreter that talks to the API Gateway. Try:\n\n" +
        '• **list cats** — show all cat profiles\n' +
        '• **list devices** — show all registered devices\n' +
        '• **add cat** — create a sample cat profile\n' +
        '• **add device** — register a sample device',
      timestamp: new Date(),
    },
  ])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const sendMessage = async (e: FormEvent) => {
    e.preventDefault()
    const text = input.trim()
    if (!text || loading) return

    const userMsg: Message = {
      id: Date.now(),
      sender: 'user',
      text,
      timestamp: new Date(),
    }
    setMessages((prev) => [...prev, userMsg])
    setInput('')
    setLoading(true)

    const response = await handleCommand(text)

    const botMsg: Message = {
      id: Date.now() + 1,
      sender: 'bot',
      text: response,
      timestamp: new Date(),
    }
    setMessages((prev) => [...prev, botMsg])
    setLoading(false)
  }

  return (
    <div className="app">
      <header className="header">
        <h1>🐱 Smart Home Cat Demo</h1>
        <span className="api-badge" title={API_URL}>
          API: {API_URL}
        </span>
      </header>

      <div className="messages">
        {messages.map((msg) => (
          <div key={msg.id} className={`message ${msg.sender}`}>
            <div className="message-bubble">
              <div className="message-text">{renderMessageText(msg.text)}</div>
              <div className="message-time">
                {msg.timestamp.toLocaleTimeString()}
              </div>
            </div>
          </div>
        ))}
        {loading && (
          <div className="message bot">
            <div className="message-bubble">
              <div className="typing-indicator">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      <form className="input-area" onSubmit={sendMessage}>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a command... (e.g. list cats)"
          disabled={loading}
          aria-label="Chat input"
        />
        <button type="submit" disabled={loading || !input.trim()}>
          Send
        </button>
      </form>
    </div>
  )
}
