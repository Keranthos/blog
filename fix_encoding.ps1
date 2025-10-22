# 修复EditArticleView.vue文件中的字符编码问题
$filePath = "D:\blog\frontend\src\views\EditArticleView.vue"
$content = Get-Content $filePath -Raw -Encoding UTF8

# 修复所有损坏的字符
$content = $content.Replace('内容类型?/label>', '内容类型：</label>')
$content = $content.Replace('文章类型选择（仅文章模式显示?-->', '文章类型选择（仅文章模式显示） -->')
$content = $content.Replace('文章类型?/label>', '文章类型：</label>')
$content = $content.Replace('媒体类型选择（仅媒体模式显示?-->', '媒体类型选择（仅媒体模式显示） -->')
$content = $content.Replace('媒体类型?/label>', '媒体类型：</label>')
$content = $content.Replace('文章标题?/label>', '文章标题：</label>')
$content = $content.Replace('请输入文章标?', '请输入文章标题')
$content = $content.Replace('媒体名称?/label>', '媒体名称：</label>')
$content = $content.Replace('请输入媒体名?', '请输入媒体名称')
$content = $content.Replace('评分?-10）', '评分（1-10）')
$content = $content.Replace('请输入评?', '请输入评分')
$content = $content.Replace('技?前端', '技术,前端')
$content = $content.Replace('封面图片?/label>', '封面图片：</label>')
$content = $content.Replace('图片加载?..', '图片加载中...')
$content = $content.Replace('请输入图片链?', '请输入图片链接')
$content = $content.Replace('媒体封面?/label>', '媒体封面：</label>')
$content = $content.Replace('图片选择?', '图片选择器')

# 保存修复后的文件
Set-Content $filePath -Value $content -Encoding UTF8
Write-Host "文件修复完成"
