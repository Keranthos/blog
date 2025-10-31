const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 3002, // 设置端口为3002
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        secure: false
      },
      '/uploads': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        secure: false
      }
    }
  },
  configureWebpack: {
    resolve: {
      alias: {
        '@': require('path').resolve(__dirname, 'src') // 确保别名 '@' 指向 src
      }
    },
    // 代码分割优化
    optimization: {
      splitChunks: {
        chunks: 'all',
        cacheGroups: {
          // 第三方库单独打包
          vendor: {
            name: 'chunk-vendors',
            test: /[\\/]node_modules[\\/]/,
            priority: 10,
            reuseExistingChunk: true
          },
          // Markdown 相关库（marked, DOMPurify）单独打包
          markdown: {
            name: 'chunk-markdown',
            test: /[\\/]node_modules[\\/](marked|dompurify)[\\/]/,
            priority: 20,
            reuseExistingChunk: true
          },
          // Font Awesome 单独打包（体积较大）
          fontAwesome: {
            name: 'chunk-fontawesome',
            test: /[\\/]node_modules[\\/]@fortawesome[\\/]/,
            priority: 20,
            reuseExistingChunk: true
          },
          // Vue Router 和 Vuex 单独打包
          vueCommon: {
            name: 'chunk-vue-common',
            test: /[\\/]node_modules[\\/](vue-router|vuex)[\\/]/,
            priority: 15,
            reuseExistingChunk: true
          }
        }
      }
    }
  }
})
