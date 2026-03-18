/**
 * Character configuration and API base URL.
 * Update BASE_URL to your deployed server address.
 * Image domains must be whitelisted in the WeChat Developer Console for production.
 */
// BASE_URL is the deployed server address.
// For local development the default points to localhost:8080.
// Change this to your production server URL before uploading the mini-program.
const BASE_URL = 'http://localhost:8080'

const CHARACTERS = {
  boyfriend: {
    key: 'boyfriend',
    name: '浩然',
    tag: '温柔幽默',
    desc: '温润如玉的灵魂，总能在疲惫的一天后给你温暖的笑颜...',
    cardImage:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuC41mo-RNHOGwmnUwEYlQe3w-teJtL0veEgnJuRxsN-VGEHkI71niyt9puXcR5e8fmK3Pw2hzpZaodHiWygZh4pLe4uBTnhtc-tFT9X5FdMNRD44T7VpJ-O0LzhboS5fkAXy0k2Rwh2rgfAsOhXA26Y4kx-yfaSqx1-P_LM3gaI3EBWKG0zziUwcilLwimrH-0CGtosbr-9or6Me2Zt3BnDj8NB1jA_dEE8FL9vCoR92RQAKDP61D0P9iNrSYiOM4bO9b652S-xdA',
    avatarImage:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuC41mo-RNHOGwmnUwEYlQe3w-teJtL0veEgnJuRxsN-VGEHkI71niyt9puXcR5e8fmK3Pw2hzpZaodHiWygZh4pLe4uBTnhtc-tFT9X5FdMNRD44T7VpJ-O0LzhboS5fkAXy0k2Rwh2rgfAsOhXA26Y4kx-yfaSqx1-P_LM3gaI3EBWKG0zziUwcilLwimrH-0CGtosbr-9or6Me2Zt3BnDj8NB1jA_dEE8FL9vCoR92RQAKDP61D0P9iNrSYiOM4bO9b652S-xdA',
  },
  girlfriend: {
    key: 'girlfriend',
    name: '米娜',
    tag: '可爱粘人',
    desc: '总是渴望和你在一起，分享每一天的点点滴滴...',
    cardImage:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuCIyXmkiIQgswh7TWnqfpuP1iA3myAroeNBOBATC67O2ReMNZhHoyUD-XGlzaLzRVEY1gHnLeLGgRvLTytHuqiPIuZQCjwnL11-FBNrMZQcmxCTbyhipRztBsJJcWRnm-WWCC2pP8l2ZVFKictc5UD2INt_wHvsQSwLJGasGHXOpiE5PImpXdk5jnwF3W5W0oosmB0oHehNU4vjxB2Mr0vSPJZKyRFRYv-VVP1ktF3q-nE9J41WY3o8uxT_TDszhNbDprI5rTw',
    avatarImage:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuCGnEdJq1vgreTdY6gyXOQfDI4af86cbC9pYcM8bQUT2BmC5aDCJT_Ys_ZN7dSG6s-v3aWcVJ953eoCSZNTundRmfUwbqKz_cr8OQ75rSvdi7MvJ_16Or2_5g3JkzxYTQwg_2ctYl9iVVqqFX7-1WDOINW32tfNGbuVjlNPHNLYxdLSx4nPxP4zokAEscn6W_c5wBG74xh0iHxy1RF3XPjIDXeXrpG5Wrsi7aj-tKkrsZNELZs-pf45aZkdinLfol7V_urTngAcAA',
  },
}

module.exports = { BASE_URL, CHARACTERS }
