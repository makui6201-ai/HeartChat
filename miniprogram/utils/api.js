const { BASE_URL } = require('./config')

/**
 * Send a chat message and receive an AI reply.
 * @param {string} role - 'boyfriend' | 'girlfriend'
 * @param {number} love - current love score
 * @param {string} message - user message text
 * @param {Array}  history - last N {role, content} turns for context
 */
function chat(role, love, message, history) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${BASE_URL}/api/chat`,
      method: 'POST',
      header: { 'Content-Type': 'application/json' },
      data: { role, love, message, history },
      success(res) {
        if (res.statusCode === 200) {
          resolve(res.data)
        } else {
          reject(new Error((res.data && res.data.error) || '请求失败'))
        }
      },
      fail(err) {
        reject(err)
      },
    })
  })
}

/**
 * Fetch a context-sensitive greeting for the character.
 * @param {string} role - 'boyfriend' | 'girlfriend'
 * @param {number} love - current love score
 */
function greeting(role, love) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${BASE_URL}/api/greeting`,
      method: 'POST',
      header: { 'Content-Type': 'application/json' },
      data: { role, love },
      success(res) {
        if (res.statusCode === 200) {
          resolve(res.data)
        } else {
          reject(new Error((res.data && res.data.error) || '请求失败'))
        }
      },
      fail(err) {
        reject(err)
      },
    })
  })
}

module.exports = { chat, greeting }
