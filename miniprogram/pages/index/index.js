const { CHARACTERS } = require('../../utils/config')
const app = getApp()

Page({
  data: {
    boyfriendChar: CHARACTERS.boyfriend,
    girlfriendChar: CHARACTERS.girlfriend,
  },

  selectBoyfriend() {
    this._selectRole('boyfriend')
  },

  selectGirlfriend() {
    this._selectRole('girlfriend')
  },

  _selectRole(role) {
    app.globalData.role = role
    app.saveData()
    wx.navigateTo({ url: '/pages/chat/chat' })
  },
})
