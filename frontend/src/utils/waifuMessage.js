// 看板娘消息工具
// 用于在应用中显示看板娘的可爱提示

/**
 * 显示看板娘消息
 * @param {string} text - 要显示的文本
 * @param {number} timeout - 显示时长（毫秒）
 * @param {number} priority - 优先级（数字越大优先级越高）
 */
export function showWaifuMessage (text, timeout = 6000, priority = 9) {
  if (typeof window !== 'undefined' && window.showMessage) {
    window.showMessage(text, timeout, priority)
  }
}

/**
 * 错误消息映射
 */
const errorMessages = {
  // 网络错误
  network: '(｡•́︿•̀｡)<br>呜呜…网络好像出问题了，刷新试试吧？',
  timeout: '等太久啦！网络可能有点慢哦～',

  // HTTP 状态码错误
  400: '请求格式不对呢…检查一下输入吧～',
  401: '(｡•́︿•̀｡)<br>登录已过期，需要重新登录哦～',
  403: '(｡•́︿•̀｡)<br>这个操作没有权限呢…',
  404: '哎呀，找不到这个内容了…要不回首页看看？',
  500: '服务器小哥哥累倒了…稍后再试试吧～',
  502: '网关出问题了…服务器可能在维护呢～',
  503: '服务暂时不可用…过会儿再来吧！',

  // 业务错误
  load_failed: '(๑•̀ㅂ•́)و✧<br>加载失败了呢…再试一次吧！',
  submit_failed: '提交失败了…检查一下网络？',
  upload_failed: '上传失败了…文件太大了吗？',
  empty_input: '还没有输入内容哦～',
  invalid_format: '格式好像不太对呢…',
  invalid_email: '邮箱格式不对哦～',
  invalid_password: '密码格式不符合要求呢～',
  password_mismatch: '两次密码不一样哦！',
  username_taken: '这个用户名已经被占用了呢…',
  login_failed: '登录失败了…用户名或密码错误？',
  register_failed: '注册失败了…可能是用户名重复了？',

  // 数据库错误
  database_error: '数据库小姐姐罢工了…稍后再试吧～',

  // 默认错误
  unknown: '出了点小问题…不要担心，刷新一下就好啦！',

  // API相关错误
  api_error: '(｡•́︿•̀｡)<br>API请求失败了呢…检查一下网络？',
  server_error: '(｡•́︿•̀｡)<br>服务器小哥哥累倒了…稍后再试试吧～',

  // 数据相关错误
  data_parse_error: '(｡•́︿•̀｡)<br>数据解析失败了呢…',
  validation_error: '(｡•́︿•̀｡)<br>输入的数据格式不对哦～',

  // 权限相关错误
  insufficient_permissions: '(｡•́︿•̀｡)<br>权限不够呢…联系管理员吧～',
  resource_not_found: '(｡•́︿•̀｡)<br>找不到这个资源了呢…',

  // 界面相关错误
  presentation_error: '(｡•́︿•̀｡)<br>讲演加载失败了呢…再试试看？',
  media_error: '(｡•́︿•̀｡)<br>媒体内容加载失败…检查一下网络？',
  comment_error: '(｡•́︿•̀｡)<br>评论提交失败…再试一次吧～',
  timeline_error: '(｡•́︿•̀｡)<br>时间线加载失败了呢…',
  search_error: '(｡•́︿•̀｡)<br>搜索功能出问题了…'
}

/**
 * 成功消息映射
 */
const successMessages = {
  submit: 'ヾ(✿ﾟ▽ﾟ)ノ<br>提交成功啦！',
  save: '保存成功！✨',
  delete: '删除成功～',
  update: '更新成功！',
  upload: '上传成功！✨',
  login: '(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧<br>欢迎回来！',
  logout: '再见啦～记得常回来看看！',
  register: '注册成功！欢迎加入～',
  comment: '评论发表成功！主人会很开心的～',
  copy: '复制成功！记得注明出处哦～',
  presentation: '讲演内容加载成功！',
  media: '媒体内容加载成功！',
  search: '搜索完成！找到相关内容啦～',
  timeline: '时间线加载成功！',
  profile: '个人资料加载成功！'
}

/**
 * 警告消息映射
 */
const warningMessages = {
  unsaved: '还有内容没保存哦，确定要离开吗？',
  delete_confirm: '真的要删除吗？删了就找不回来了…',
  network_slow: '网络有点慢呢…请耐心等待～',
  file_too_large: '文件太大了！换个小一点的吧～',
  invalid_file_type: '文件类型不对哦～',
  empty_search: '搜索框是空的哦～输入点什么吧！',
  long_content: '内容有点长呢…慢慢看哦～',
  unsaved_changes: '还有内容没保存哦，确定要离开吗？'
}

/**
 * 显示错误消息
 * @param {Error|string|number} error - 错误对象、错误类型或 HTTP 状态码
 */
export function showErrorMessage (error) {
  let message = errorMessages.unknown

  if (typeof error === 'number') {
    // HTTP 状态码
    message = errorMessages[error] || errorMessages.unknown
  } else if (typeof error === 'string') {
    // 错误类型字符串
    message = errorMessages[error] || errorMessages.unknown
  } else if (error && error.response) {
    // Axios 错误对象
    const status = error.response.status
    message = errorMessages[status] || errorMessages.unknown
  } else if (error && error.message) {
    // 标准 Error 对象
    if (error.message.includes('Network')) {
      message = errorMessages.network
    } else if (error.message.includes('timeout')) {
      message = errorMessages.timeout
    }
  }

  showWaifuMessage(message, 6000, 10)
}

/**
 * 显示成功消息
 * @param {string} type - 成功类型
 */
export function showSuccessMessage (type = 'submit') {
  const message = successMessages[type] || successMessages.submit
  showWaifuMessage(message, 4000, 9)
}

/**
 * 显示警告消息
 * @param {string} type - 警告类型
 */
export function showWarningMessage (type) {
  const message = warningMessages[type] || type
  showWaifuMessage(message, 5000, 9)
}

/**
 * 显示自定义消息
 * @param {string} text - 自定义文本
 * @param {number} timeout - 显示时长
 */
export function showCustomMessage (text, timeout = 6000) {
  showWaifuMessage(text, timeout, 9)
}

// 导出所有消息映射，方便扩展
export const messages = {
  error: errorMessages,
  success: successMessages,
  warning: warningMessages
}
