module.exports = {
  devServer: {
    disableHostCheck: true,
    clientLogLevel: 'info'
  },
  configureWebpack: {
    resolve: {
      alias: {
        'vue$': 'vue/dist/vue.esm.js'
      }
    }
  }
}