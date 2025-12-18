/**
 * Font Awesome 图标按需加载工具
 * 首屏只加载必要的图标，其他图标按需动态加载
 */

import { library } from '@fortawesome/fontawesome-svg-core'

// 图标加载缓存，避免重复加载
const loadedIcons = new Set()

/**
 * 动态加载单个图标
 * @param {string} iconName - 图标名称（如 'diagramProject', 'flask'，注意驼峰命名）
 * @param {string} style - 图标样式（'solid', 'regular', 'brands'）
 * @returns {Promise}
 */
export async function loadIcon (iconName, style = 'solid') {
  const cacheKey = `${style}-${iconName}`

  // 如果已加载，直接返回
  if (loadedIcons.has(cacheKey)) {
    return Promise.resolve()
  }

  try {
    let iconModule

    switch (style) {
      case 'solid':
        iconModule = await import('@fortawesome/free-solid-svg-icons')
        break
      case 'regular':
        iconModule = await import('@fortawesome/free-regular-svg-icons')
        break
      case 'brands':
        iconModule = await import('@fortawesome/free-brands-svg-icons')
        break
      default:
        throw new Error(`Unknown icon style: ${style}`)
    }

    // Font Awesome 图标名称格式：faHouse, faBlog 等
    // iconName 应该是驼峰格式：diagramProject -> faDiagramProject
    const iconKey = `fa${iconName.charAt(0).toUpperCase() + iconName.slice(1)}`
    const icon = iconModule[iconKey]

    if (icon) {
      library.add(icon)
      loadedIcons.add(cacheKey)
    } else {
      console.warn(`Icon ${iconKey} not found in ${style} style`)
    }
  } catch (error) {
    console.error(`Failed to load icon ${iconName} from ${style}:`, error)
  }
}

/**
 * 批量加载图标
 * @param {Array} icons - 图标数组，格式：[{name: 'diagramProject', style: 'solid'}, ...]
 * @returns {Promise}
 */
export async function loadIcons (icons) {
  const promises = icons.map(({ name, style = 'solid' }) => loadIcon(name, style))
  await Promise.all(promises)
}

/**
 * 预加载常用图标（在需要时调用）
 * @param {Array} iconNames - 图标名称数组（驼峰格式）
 * @param {string} style - 图标样式
 */
export async function preloadIcons (iconNames, style = 'solid') {
  const icons = iconNames.map(name => ({ name, style }))
  await loadIcons(icons)
}
