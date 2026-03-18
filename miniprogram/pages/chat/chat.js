const app = getApp()
const api = require('../../utils/api')
const { CHARACTERS } = require('../../utils/config')

// Love score level definitions
const LOVE_LEVELS = [
  { min: 0,  max: 9,        name: '普通朋友' },
  { min: 10, max: 29,       name: '熟悉' },
  { min: 30, max: 59,       name: '亲密' },
  { min: 60, max: Infinity, name: '暧昧' },
]

function getLoveLevel(love) {
  for (const level of LOVE_LEVELS) {
    if (love >= level.min && love <= level.max) return level.name
  }
  return '暧昧'
}

function formatTime() {
  const now = new Date()
  const h = String(now.getHours()).padStart(2, '0')
  const m = String(now.getMinutes()).padStart(2, '0')
  return `${h}:${m}`
}

Page({
  data: {
    char: {},            // current character config
    love: 0,
    loveLevel: '普通朋友',
    messages: [],
    inputText: '',
    sending: false,
    scrollToId: '',
    _msgCounter: 0,     // private counter for unique IDs
  },

  onLoad() {
    const role   = app.globalData.role   || 'girlfriend'
    const love   = app.globalData.love   || 0
    const openId = app.globalData.openId || ''
    const char   = CHARACTERS[role] || CHARACTERS.girlfriend

    this.setData({
      char,
      love,
      loveLevel: getLoveLevel(love),
    })

    // If we have a server-side user identity, load their stored history first
    // so returning users see their previous conversation.
    if (openId) {
      this._loadServerHistory(openId, role, char, love)
    } else {
      this._loadGreeting()
    }
  },

  // ------------------------------------------------------------------ history

  /** Load persisted history from the server for returning users. */
  _loadServerHistory(openId, role, char, love) {
    api.getHistory(openId)
      .then((data) => {
        // Restore role and love from the server if available.
        const serverRole = data.role || role
        const serverLove = data.love != null ? data.love : love
        const serverChar = CHARACTERS[serverRole] || char

        // Sync global state with server state.
        app.globalData.role = serverRole
        app.globalData.love = serverLove

        this.setData({
          char: serverChar,
          love: serverLove,
          loveLevel: getLoveLevel(serverLove),
        })

        // Render persisted history messages.
        if (data.history && data.history.length > 0) {
          const msgs = data.history.map((m, i) => ({
            id: i + 1,
            role: m.role === 'assistant' ? 'ai' : 'user',
            content: m.content,
            time: '',  // historical messages don't have a stored time
          }))
          this.setData({
            messages: msgs,
            _msgCounter: msgs.length,
            scrollToId: `msg-${msgs.length}`,
          })
        }
        // Always show a greeting after history is loaded.
        this._loadGreeting()
      })
      .catch(() => {
        // History load failed – just show the greeting as normal.
        this._loadGreeting()
      })
  },

  // ------------------------------------------------------------------ greeting

  _loadGreeting() {
    const { char, love } = this.data
    // Show typing indicator while greeting loads
    this._addTyping()
    // 1-second natural delay before the greeting appears
    setTimeout(() => {
      api
        .greeting(char.key, love)
        .then((res) => {
          this._removeTyping()
          this._typewriter(res.greeting)
        })
        .catch(() => {
          this._removeTyping()
          const fallback = love > 0 ? '你终于来了！' : '你好'
          this._addMessage('ai', fallback)
        })
    }, 1000)
  },

  // ------------------------------------------------------------------ send

  onInput(e) {
    this.setData({ inputText: e.detail.value })
  },

  sendMessage() {
    const text = (this.data.inputText || '').trim()
    if (!text || this.data.sending) return

    this.setData({ inputText: '', sending: true })

    // Add user message
    this._addMessage('user', text)

    // Show typing indicator
    this._addTyping()

    // Gather last 10 non-typing turns for context
    const history = this.data.messages
      .filter((m) => m.role === 'user' || m.role === 'ai')
      .slice(-10)
      .map((m) => ({
        role: m.role === 'ai' ? 'assistant' : 'user',
        content: m.content,
      }))

    const { char, love } = this.data
    const openId = app.globalData.openId || ''

    api
      .chat(openId, char.key, love, text, history)
      .then((res) => {
        this._removeTyping()
        return this._typewriter(res.reply)
      })
      .then(() => {
        const newLove = this.data.love + 1
        app.globalData.love = newLove
        app.saveData()
        this.setData({
          love: newLove,
          loveLevel: getLoveLevel(newLove),
          sending: false,
        })
      })
      .catch(() => {
        this._removeTyping()
        this._addMessage('ai', '抱歉，我现在不太方便说话，稍后再聊~')
        this.setData({ sending: false })
      })
  },

  goBack() {
    wx.navigateBack({ delta: 1 })
  },

  // ------------------------------------------------------------------ helpers

  /** Append a regular message bubble and scroll to it. */
  _addMessage(role, content) {
    const counter = this.data._msgCounter + 1
    const messages = [
      ...this.data.messages,
      { id: counter, role, content, time: formatTime() },
    ]
    this.setData({
      messages,
      _msgCounter: counter,
      scrollToId: `msg-${counter}`,
    })
    return counter
  },

  /** Append a typing-indicator bubble. */
  _addTyping() {
    const counter = this.data._msgCounter + 1
    const messages = [
      ...this.data.messages,
      { id: counter, role: 'typing', content: '', time: '' },
    ]
    this.setData({
      messages,
      _msgCounter: counter,
      scrollToId: `msg-${counter}`,
    })
  },

  /** Remove the last message (used to remove typing indicator). */
  _removeTyping() {
    const messages = this.data.messages.slice(0, -1)
    this.setData({ messages })
  },

  /**
   * Typewriter effect: adds an empty AI bubble then progressively reveals text.
   * Uses path-based setData for efficiency (no full-array clone per frame).
   */
  _typewriter(text) {
    return new Promise((resolve) => {
      const counter = this.data._msgCounter + 1
      const msgIndex = this.data.messages.length  // index of the new message
      const messages = [
        ...this.data.messages,
        { id: counter, role: 'ai', content: '', time: formatTime() },
      ]
      this.setData({
        messages,
        _msgCounter: counter,
        scrollToId: `msg-${counter}`,
      })

      let i = 0
      const interval = setInterval(() => {
        i++
        const partial = text.slice(0, i)
        this.setData({
          [`messages[${msgIndex}].content`]: partial,
          scrollToId: `msg-${counter}`,
        })
        if (i >= text.length) {
          clearInterval(interval)
          resolve()
        }
      }, 50)
    })
  },
})
