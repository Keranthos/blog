import katex from 'katex'
import 'katex/dist/katex.min.css'

/**
 * 在 Markdown 渲染前保护 LaTeX 公式，避免被 Markdown 解析器处理
 * @param {string} markdown - 原始 Markdown 内容
 * @returns {Object} - { protected: 保护后的内容, placeholders: 占位符映射 }
 */
export function protectLatex (markdown) {
  if (!markdown) return { protected: '', placeholders: [] }

  const placeholders = []
  let placeholderIndex = 0
  let protectedContent = markdown

  // 先保护块级公式 $$...$$
  protectedContent = protectedContent.replace(/\$\$([\s\S]*?)\$\$/g, (match, formula) => {
    // 使用 HTML 注释作为占位符，这样 Markdown 不会处理它
    const placeholder = `<!--LATEX_BLOCK_${placeholderIndex}-->`
    placeholders.push({
      placeholder,
      formula: formula.trim(),
      type: 'block'
    })
    placeholderIndex++
    return placeholder
  })

  // 再保护行内公式 $...$（避免匹配 $$...$$）
  // 使用更兼容的方式：先标记所有 $$，然后处理单个 $
  const doubleDollarMarkers = []
  let markerIndex = 0
  protectedContent = protectedContent.replace(/\$\$/g, () => {
    const marker = `<!--DOUBLE_DOLLAR_${markerIndex}-->`
    doubleDollarMarkers.push(marker)
    markerIndex++
    return marker
  })

  // 现在处理行内公式 $...$（此时已经没有 $$ 了）
  protectedContent = protectedContent.replace(/\$([^$\n]+?)\$/g, (match, formula) => {
    // 使用 HTML 注释作为占位符
    const placeholder = `<!--LATEX_INLINE_${placeholderIndex}-->`
    placeholders.push({
      placeholder,
      formula: formula.trim(),
      type: 'inline'
    })
    placeholderIndex++
    return placeholder
  })

  // 恢复 $$ 标记
  doubleDollarMarkers.forEach((marker, index) => {
    protectedContent = protectedContent.replace(marker, '$$')
  })

  return { protected: protectedContent, placeholders }
}

/**
 * 恢复并渲染 LaTeX 公式
 * @param {string} html - 已渲染的 HTML 内容
 * @param {Array} placeholders - 占位符映射
 * @returns {string} - 处理后的 HTML 内容
 */
export function restoreAndRenderLatex (html, placeholders) {
  if (!html || !placeholders || placeholders.length === 0) return html

  let processedHtml = html

  // 恢复并渲染每个占位符
  placeholders.forEach(({ placeholder, formula, type }) => {
    try {
      if (!formula) {
        processedHtml = processedHtml.replace(placeholder, '')
        return
      }

      // 检测是否包含对齐语法（& 和 \\），如果是块级公式且没有包裹在 aligned 环境中，自动添加
      let processedFormula = formula
      if (type === 'block') {
        const hasAmpersand = formula.includes('&')
        const hasLineBreak = formula.includes('\\')
        const hasAlignedEnv = formula.includes('\\begin{aligned}') || formula.includes('\\begin{align}')

        // 如果使用了 equation 环境（会自动编号），替换为 equation* 环境（不会编号）
        if (formula.includes('\\begin{equation}')) {
          processedFormula = formula.replace(/\\begin\{equation\}/g, '\\begin{equation*}')
          processedFormula = processedFormula.replace(/\\end\{equation\}/g, '\\end{equation*}')
        }
        // 如果使用了 align 环境（会自动编号），替换为 aligned 环境（不会编号）
        if (formula.includes('\\begin{align}')) {
          processedFormula = processedFormula.replace(/\\begin\{align\}/g, '\\begin{aligned}')
          processedFormula = processedFormula.replace(/\\end\{align\}/g, '\\end{aligned}')
        } else if ((hasAmpersand || hasLineBreak) && !hasAlignedEnv) {
          // 如果包含对齐语法但没有 aligned 环境，自动包裹
          // 移除首尾空白，然后包裹在 aligned 环境中（不会自动编号）
          processedFormula = `\\begin{aligned}\n${processedFormula.trim()}\n\\end{aligned}`
        }
      }

      const rendered = katex.renderToString(processedFormula, {
        displayMode: type === 'block',
        throwOnError: false,
        errorColor: '#cc0000',
        strict: false,
        trust: false,
        output: 'html'
      })

      // KaTeX 返回的 HTML 已经包含 <span class="katex">，我们只需要包装 div
      const wrappedHtml = type === 'block'
        ? `<div class="katex-block" style="text-align: center !important; display: block !important; width: 100% !important; margin: 24px 0 !important; max-width: 100% !important;">${rendered}</div>`
        : rendered

      processedHtml = processedHtml.replace(placeholder, wrappedHtml)
    } catch (error) {
      console.error('LaTeX rendering error:', error, formula)
      const errorHtml = type === 'block'
        ? `<div class="katex-error">公式渲染错误: ${formula}</div>`
        : `<span class="katex-error">公式错误: ${formula}</span>`
      processedHtml = processedHtml.replace(placeholder, errorHtml)
    }
  })

  return processedHtml
}

/**
 * 处理 Markdown 内容中的 LaTeX 公式（旧版本，保留以兼容）
 * 支持行内公式 $...$ 和块级公式 $$...$$
 * @param {string} html - 已渲染的 HTML 内容
 * @returns {string} - 处理后的 HTML 内容
 */
export function renderLatex (html) {
  if (!html) return ''

  let processedHtml = html

  // 先处理块级公式 $$...$$，避免与行内公式冲突
  const blockFormulaPlaceholders = []
  let placeholderIndex = 0

  processedHtml = processedHtml.replace(/\$\$([\s\S]*?)\$\$/g, (match, formula) => {
    try {
      const trimmedFormula = formula.trim()
      if (!trimmedFormula) return match
      const rendered = katex.renderToString(trimmedFormula, {
        displayMode: true,
        throwOnError: false,
        errorColor: '#cc0000'
      })
      const placeholder = `__KATEX_BLOCK_${placeholderIndex}__`
      blockFormulaPlaceholders.push({ placeholder, html: `<div class="katex-block">${rendered}</div>` })
      placeholderIndex++
      return placeholder
    } catch (error) {
      console.error('LaTeX block rendering error:', error, formula)
      const placeholder = `__KATEX_BLOCK_${placeholderIndex}__`
      blockFormulaPlaceholders.push({ placeholder, html: `<div class="katex-error">公式渲染错误: ${formula}</div>` })
      placeholderIndex++
      return placeholder
    }
  })

  // 处理行内公式 $...$（此时已没有 $$...$$）
  processedHtml = processedHtml.replace(/\$([^$\n]+?)\$/g, (match, formula) => {
    try {
      const trimmedFormula = formula.trim()
      if (!trimmedFormula) return match
      const rendered = katex.renderToString(trimmedFormula, {
        displayMode: false,
        throwOnError: false,
        errorColor: '#cc0000'
      })
      return rendered
    } catch (error) {
      console.error('LaTeX inline rendering error:', error, formula)
      return `<span class="katex-error">公式错误: ${formula}</span>`
    }
  })

  // 恢复块级公式占位符
  blockFormulaPlaceholders.forEach(({ placeholder, html }) => {
    processedHtml = processedHtml.replace(placeholder, html)
  })

  return processedHtml
}
