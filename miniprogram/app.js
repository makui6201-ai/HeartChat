App({
  globalData: {
    role: '',
    love: 0,
  },

  onLaunch() {
    this.globalData.role = wx.getStorageSync('role') || ''
    this.globalData.love = wx.getStorageSync('love') || 0
  },

  saveData() {
    wx.setStorageSync('role', this.globalData.role)
    wx.setStorageSync('love', this.globalData.love)
  },
})
