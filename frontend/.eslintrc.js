module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true // 如果需要支持 Node.js 环境
  },
  extends: [
    'eslint:recommended', // 这一句话解决了props报错？
    'plugin:vue/vue3-recommended', // 使用 Vue 3 的推荐配置
    '@vue/eslint-config-standard' // 使用标准的 Vue 配置
  ],
  parserOptions: {
    ecmaVersion: 12, // 你使用的是 ES2021 语法
    sourceType: 'module' // 确保支持 ES 模块
  },
  plugins: [
    'vue' // 添加 Vue 插件
  ],
  rules: {
    // 在这里添加你需要的自定义规则
    // "vue/html-self-closing": "off", // 关闭自闭合标签警告
    'vue/component-definition-name-casing': 'off', // 关闭组件命名风格警告
    'vue/attribute-hyphenation': ['off'],
    // 你可以根据需求关闭其他规则
    'vue/html-self-closing': [
      'error',
      {
        html: {
          void: 'always', // 自闭合 void 元素，如 <img>
          normal: 'never', // 正常元素不要求自闭合
          component: 'always' // Vue 组件要求自闭合
        }
      }
    ],
    'vue/max-attributes-per-line': ['error', {
      singleline: 10000, // 单行最多允许 5 个属性
      multiline: 10000 // 多行每个属性放一行
    }],
    'vue/singleline-html-element-content-newline': 'off',
    'vue/require-default-prop': 'off' // 禁用此规则
  }
}
