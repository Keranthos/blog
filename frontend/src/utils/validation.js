/**
 * 前端验证工具函数
 * 与后端验证逻辑保持一致，提供即时反馈
 */

/**
 * 验证用户名
 * @param {string} username - 用户名
 * @returns {{valid: boolean, error: string}} 验证结果
 */
export function validateUsername (username) {
  // 去除首尾空格
  const trimmed = username.trim()

  if (!trimmed) {
    return { valid: false, error: '用户名不能为空' }
  }

  // 计算字符数（支持中文）
  const charCount = [...trimmed].length
  if (charCount > 12) {
    return { valid: false, error: '用户名最多12个字符' }
  }

  if (charCount < 2) {
    return { valid: false, error: '用户名至少2个字符' }
  }

  // 检查是否包含非法字符（字母、数字、下划线、中文）
  const usernameRegex = /^[a-zA-Z0-9_\u4e00-\u9fa5]+$/
  if (!usernameRegex.test(trimmed)) {
    return { valid: false, error: '用户名只能包含字母、数字、下划线和中文' }
  }

  return { valid: true, error: '' }
}

/**
 * 验证密码
 * @param {string} password - 密码
 * @returns {{valid: boolean, error: string}} 验证结果
 */
export function validatePassword (password) {
  if (!password) {
    return { valid: false, error: '密码不能为空' }
  }

  if (password.length < 6) {
    return { valid: false, error: '密码至少6位' }
  }

  if (password.length > 50) {
    return { valid: false, error: '密码最多50位' }
  }

  return { valid: true, error: '' }
}

/**
 * 验证两次密码是否一致
 * @param {string} password - 密码
 * @param {string} confirmPassword - 确认密码
 * @returns {{valid: boolean, error: string}} 验证结果
 */
export function validatePasswordMatch (password, confirmPassword) {
  if (password !== confirmPassword) {
    return { valid: false, error: '两次密码不一致' }
  }
  return { valid: true, error: '' }
}

/**
 * 验证注册表单
 * @param {string} username - 用户名
 * @param {string} password - 密码
 * @param {string} confirmPassword - 确认密码
 * @returns {{valid: boolean, error: string}} 验证结果
 */
export function validateRegisterForm (username, password, confirmPassword) {
  // 验证用户名
  const usernameResult = validateUsername(username)
  if (!usernameResult.valid) {
    return usernameResult
  }

  // 验证密码
  const passwordResult = validatePassword(password)
  if (!passwordResult.valid) {
    return passwordResult
  }

  // 验证密码一致性
  const matchResult = validatePasswordMatch(password, confirmPassword)
  if (!matchResult.valid) {
    return matchResult
  }

  return { valid: true, error: '' }
}
