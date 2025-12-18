/**
 * 移动端优化工具
 * 提供移动端特定的优化功能
 */

/**
 * 检测是否为移动端
 * @returns {boolean}
 */
export function isMobileDevice () {
  if (typeof window === 'undefined') return false
  return window.innerWidth <= 768 || /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
}

/**
 * 移动端优化：减少首屏数据量
 * @param {number} desktopLimit - 桌面端限制
 * @param {number} mobileLimit - 移动端限制（通常更小）
 * @returns {number}
 */
export function getOptimizedLimit (desktopLimit, mobileLimit) {
  return isMobileDevice() ? mobileLimit : desktopLimit
}

/**
 * 移动端优化：简化动画
 * @param {number} desktopDuration - 桌面端动画时长（ms）
 * @param {number} mobileDuration - 移动端动画时长（ms，通常更短）
 * @returns {number}
 */
export function getOptimizedAnimationDuration (desktopDuration, mobileDuration) {
  return isMobileDevice() ? mobileDuration : desktopDuration
}

/**
 * 移动端优化：禁用非必要功能
 * @param {Function} callback - 桌面端执行的回调
 * @param {Function} mobileCallback - 移动端执行的回调（可选）
 */
export function runIfDesktop (callback, mobileCallback = null) {
  if (isMobileDevice()) {
    if (mobileCallback) {
      mobileCallback()
    }
    return
  }
  callback()
}

/**
 * 移动端优化：预加载策略
 * @param {Array} items - 要预加载的项目
 * @param {number} desktopCount - 桌面端预加载数量
 * @param {number} mobileCount - 移动端预加载数量
 * @returns {Array}
 */
export function getOptimizedPreloadItems (items, desktopCount, mobileCount) {
  const limit = isMobileDevice() ? mobileCount : desktopCount
  return items.slice(0, limit)
}
