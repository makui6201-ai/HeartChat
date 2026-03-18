const api = require('./utils/api')

App({
  globalData: {
    role: '',
    love: 0,
    openId: '',   // stable WeChat user identifier obtained via login
  },

  onLaunch() {
    this.globalData.role   = wx.getStorageSync('role')   || ''
    this.globalData.love   = wx.getStorageSync('love')   || 0
    this.globalData.openId = wx.getStorageSync('openId') || ''

    // If we don't have an openId yet, perform WeChat login now so the chat
    // page can use it immediately on first load.
    if (!this.globalData.openId) {
      this._login()
    }
  },

  // Perform WeChat jscode2session and persist the resulting openId.
  _login() {
    const self = this
    wx.login({
      success(res) {
        if (!res.code) return
        api.login(res.code)
          .then((data) => {
            if (data.openId) {
              self.globalData.openId = data.openId
              wx.setStorageSync('openId', data.openId)
            }
          })
          .catch((err) => {
            console.warn('HeartChat login failed:', err)
          })
      },
    })
  },

  saveData() {
    wx.setStorageSync('role',   this.globalData.role)
    wx.setStorageSync('love',   this.globalData.love)
    wx.setStorageSync('openId', this.globalData.openId)
  },
})
