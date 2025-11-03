# Git提交信息乱码修复指南

## 问题原因

在Windows PowerShell中，Git提交信息中的中文字符可能因为编码问题显示为乱码。这是因为：
1. PowerShell的默认编码可能与Git不一致
2. Git在保存提交信息时使用了错误的编码

## 解决方案

### 方案1：使用Git Bash（推荐）

在Git Bash中执行（Git Bash对中文支持更好）：

```bash
# 修改最近的提交
git commit --amend -m "完善部署配置：添加部署文档、配置模板和CORS更新"

# 修改倒数第二个提交
git rebase -i HEAD~2
# 在编辑器中，将第二个提交的 'pick' 改为 'reword'
# 然后修改提交信息为：
# "部署准备：清理废弃代码、优化配置、添加部署文档"

# 强制推送
git push origin master --force-with-lease
```

### 方案2：使用英文提交信息（更安全）

如果不想修改历史，后续提交使用英文：

```bash
git commit -m "Deploy preparation: add deployment docs and config templates"
```

### 方案3：配置Git使用UTF-8（全局修复）

在Git Bash或PowerShell中执行：

```bash
# 全局配置Git编码
git config --global i18n.commitencoding utf-8
git config --global i18n.logoutputencoding utf-8
git config --global core.quotepath false

# 在PowerShell中设置编码
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$env:LANG = 'zh_CN.UTF-8'
```

## 当前状态

- 最近的提交信息已经是乱码（已在GitHub上）
- 需要修复最近的2个提交信息

## 建议

1. **立即修复**：使用Git Bash修改最近的提交信息并强制推送
2. **后续提交**：使用英文提交信息，避免编码问题
3. **或者**：保持现状，后续提交使用正确编码配置

