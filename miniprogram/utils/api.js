const { BASE_URL } = require('./config')

/**
 * Exchange a WeChat login code for the user's stable openId.
 * @param {string} code - the code from wx.login()
 */
function login(code) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${BASE_URL}/api/login`,
      method: 'POST',
      header: { 'Content-Type': 'application/json' },
      data: { code },
      success(res) {
        if (res.statusCode === 200) {
          resolve(res.data)
        } else {
          reject(new Error((res.data && res.data.error) || '登录失败'))
        }
      },
      fail(err) { reject(err) },
    })
  })
}

/**
 * Load the server-side chat history (role, love score, message turns) for a
 * specific user.
 * @param {string} openId - the user's WeChat openId
 */
function getHistory(openId) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${BASE_URL}/api/history?openId=${encodeURIComponent(openId)}`,
      method: 'GET',
      header: { 'Content-Type': 'application/json' },
      success(res) {
        if (res.statusCode === 200) {
          resolve(res.data)
        } else {
          reject(new Error((res.data && res.data.error) || '获取历史失败'))
        }
      },
      fail(err) { reject(err) },
    })
  })
}

/**
 * Clear the conversation history for a user on the server.
 * @param {string} openId - the user's WeChat openId
 */
function deleteHistory(openId) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${BASE_URL}/api/history`,
      method: 'DELETE',
      header: { 'Content-Type': 'application/json' },
      data: { openId },
      success(res) {
        if (res.statusCode === 200) {
          resolve(res.data)
        } else {
          reject(new Error((res.data && res.data.error) || '清除历史失败'))
        }
      },
      fail(err) { reject(err) },
    })
  })
}

/**
 * Send a chat message and receive an AI reply.
 * @param {string} openId  - the user's WeChat openId (may be empty string)
 * @param {string} role    - 'boyfriend' | 'girlfriend'
 * @param {number} love    - current love score
 * @param {string} message - user message text
 * @param {Array}  history - last N {role, content} turns for context
 */
function chat(openId, role, love, message, history) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${BASE_URL}/api/chat`,
      method: 'POST',
      header: { 'Content-Type': 'application/json' },
      data: { openId, role, love, message, history },
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

module.exports = { login, chat, greeting, getHistory, deleteHistory }
