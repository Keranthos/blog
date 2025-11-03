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
            reuseExistingChunk: true,
            minChunks: 1,
            minSize: 0
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
          },
          // Highlight.js 单独打包（代码高亮，按需加载）
          highlightjs: {
            name: 'chunk-highlightjs',
            test: /[\\/]node_modules[\\/](highlight\.js|@highlightjs)[\\/]/,
            priority: 18,
            reuseExistingChunk: true
          }
        }
      },
      // 运行时 chunk 单独提取，便于长期缓存
      runtimeChunk: {
        name: 'runtime'
      }
    },
    // 生产环境优化
    ...(process.env.NODE_ENV === 'production' ? {
      performance: {
        hints: 'warning',
        maxEntrypointSize: 512000,  // 500KB
        maxAssetSize: 512000
      }
    } : {})
  },
  // 生产环境优化
  productionSourceMap: false,  // 禁用 source map 减小构建体积
  // 图片优化（使用 imagemin-webpack-plugin 需要额外安装）
  chainWebpack: config => {
    // 压缩图片（需要安装 imagemin-webpack-plugin）
    // config.module
    //   .rule('images')
    //   .use('image-webpack-loader')
    //   .loader('image-webpack-loader')
    //   .options({
    //     mozjpeg: { progressive: true, quality: 85 },
    //     optipng: { enabled: false },
    //     pngquant: { quality: [0.65, 0.90], speed: 4 },
    //     gifsicle: { interlaced: false },
    //     webp: { quality: 85 }
    //   })
  }
})
