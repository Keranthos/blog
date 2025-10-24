<template>
  <div class="test-guide">
    <div class="guide-container">
      <h1>🖼️ 图片API测试工具使用指南</h1>

      <div class="guide-section">
        <h2>📋 功能概述</h2>
        <p>这个测试工具可以帮助您：</p>
        <ul>
          <li>✅ 测试各个图片API的可用性和稳定性</li>
          <li>⚡ 测量API的响应速度</li>
          <li>🖼️ 预览获取到的图片质量</li>
          <li>📊 生成详细的测试报告</li>
          <li>💾 导出测试结果</li>
        </ul>
      </div>

      <div class="guide-section">
        <h2>🚀 如何使用</h2>
        <div class="step">
          <h3>步骤 1: 访问测试页面</h3>
          <p>在浏览器中访问：<code>/image-api-test</code></p>
        </div>

        <div class="step">
          <h3>步骤 2: 配置测试选项</h3>
          <ul>
            <li><strong>自动测试下一个</strong>：启用后会自动依次测试所有API</li>
            <li><strong>显示详细信息</strong>：显示每个API的详细测试信息</li>
          </ul>
        </div>

        <div class="step">
          <h3>步骤 3: 开始测试</h3>
          <p>点击 <strong>"开始全部测试"</strong> 按钮开始测试所有API</p>
        </div>

        <div class="step">
          <h3>步骤 4: 查看结果</h3>
          <p>测试完成后可以：</p>
          <ul>
            <li>查看每个API的测试状态和响应时间</li>
            <li>预览获取到的图片</li>
            <li>查看测试统计信息</li>
            <li>导出测试结果</li>
          </ul>
        </div>
      </div>

      <div class="guide-section">
        <h2>📊 测试状态说明</h2>
        <div class="status-examples">
          <div class="status-example success">
            <div class="status-icon">✅</div>
            <div class="status-info">
              <strong>成功</strong>
              <p>API响应正常，图片加载成功，响应时间 &lt; 3秒</p>
            </div>
          </div>

          <div class="status-example slow">
            <div class="status-icon">⏰</div>
            <div class="status-info">
              <strong>慢速</strong>
              <p>API响应正常，但响应时间 &gt; 3秒</p>
            </div>
          </div>

          <div class="status-example error">
            <div class="status-icon">❌</div>
            <div class="status-info">
              <strong>失败</strong>
              <p>API无法访问或图片加载失败</p>
            </div>
          </div>

          <div class="status-example testing">
            <div class="status-icon">🔄</div>
            <div class="status-info">
              <strong>测试中</strong>
              <p>正在测试该API</p>
            </div>
          </div>
        </div>
      </div>

      <div class="guide-section">
        <h2>🔧 测试的API列表</h2>
        <div class="api-categories">
          <div class="category">
            <h3>🏞️ 风景类</h3>
            <ul>
              <li>随机风景 (api.btstu.cn)</li>
              <li>随机图片 (api.btstu.cn)</li>
              <li>随机壁纸 (api.wetab.link)</li>
            </ul>
          </div>

          <div class="category">
            <h3>🎨 动漫类</h3>
            <ul>
              <li>随机动漫 (api.btstu.cn)</li>
              <li>萌图 (api.yimian.xyz)</li>
              <li>动漫图片 (api.gmit.vip)</li>
              <li>Lolimi动漫 (api.lolimi.cn)</li>
            </ul>
          </div>

          <div class="category">
            <h3>🎭 综合类</h3>
            <ul>
              <li>综合随机 (api.btstu.cn)</li>
              <li>小歪API (api.ixiaowai.cn)</li>
              <li>墨天逸 (api.mtyqx.cn)</li>
              <li>国际源随机 (picsum.photos)</li>
            </ul>
          </div>
        </div>
      </div>

      <div class="guide-section">
        <h2>📈 测试报告</h2>
        <p>测试完成后，您会看到：</p>
        <ul>
          <li><strong>总体统计</strong>：成功率、平均响应时间等</li>
          <li><strong>推荐API</strong>：响应快且稳定的API</li>
          <li><strong>问题API</strong>：测试失败的API及原因</li>
        </ul>
      </div>

      <div class="guide-section">
        <h2>💡 使用建议</h2>
        <ul>
          <li>建议在网络状况良好时进行测试</li>
          <li>测试结果会因网络环境而异</li>
          <li>可以多次测试以获得更准确的结果</li>
          <li>导出测试结果可以用于后续分析</li>
        </ul>
      </div>

      <div class="guide-section">
        <h2>🔗 快速访问</h2>
        <div class="quick-links">
          <a href="/image-api-test" class="link-btn primary">
            <font-awesome-icon icon="play" />
            开始测试
          </a>
          <a href="/editarticle" class="link-btn secondary">
            <font-awesome-icon icon="edit" />
            返回编辑页面
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
</script>

<style scoped>
.test-guide {
  min-height: 100vh;
  padding-top: 80px;
  padding-bottom: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
}

.test-guide::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  z-index: -1;
}

.guide-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.guide-container h1 {
  color: #333;
  text-align: center;
  margin-bottom: 30px;
  font-size: 2.5rem;
}

.guide-section {
  margin-bottom: 40px;
}

.guide-section h2 {
  color: #667eea;
  margin-bottom: 20px;
  font-size: 1.5rem;
  border-bottom: 2px solid #e9ecef;
  padding-bottom: 10px;
}

.guide-section h3 {
  color: #495057;
  margin-bottom: 10px;
  font-size: 1.2rem;
}

.guide-section p {
  color: #6c757d;
  line-height: 1.6;
  margin-bottom: 15px;
}

.guide-section ul {
  color: #6c757d;
  line-height: 1.8;
  margin-left: 20px;
}

.guide-section li {
  margin-bottom: 8px;
}

.step {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  border-left: 4px solid #667eea;
}

.step h3 {
  color: #667eea;
  margin-bottom: 15px;
}

code {
  background: #e9ecef;
  color: #e83e8c;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
}

.status-examples {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.status-example {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border-radius: 12px;
  border: 2px solid transparent;
}

.status-example.success {
  background: rgba(40, 167, 69, 0.1);
  border-color: #28a745;
}

.status-example.slow {
  background: rgba(255, 193, 7, 0.1);
  border-color: #ffc107;
}

.status-example.error {
  background: rgba(220, 53, 69, 0.1);
  border-color: #dc3545;
}

.status-example.testing {
  background: rgba(23, 162, 184, 0.1);
  border-color: #17a2b8;
}

.status-icon {
  font-size: 2rem;
  width: 50px;
  text-align: center;
}

.status-info strong {
  display: block;
  margin-bottom: 5px;
  font-size: 1.1rem;
}

.status-info p {
  margin: 0;
  font-size: 0.9rem;
  color: #6c757d;
}

.api-categories {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.category {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 20px;
  border-left: 4px solid #667eea;
}

.category h3 {
  color: #667eea;
  margin-bottom: 15px;
  font-size: 1.2rem;
}

.category ul {
  margin: 0;
  padding-left: 20px;
}

.category li {
  margin-bottom: 8px;
  font-size: 0.9rem;
}

.quick-links {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-top: 30px;
}

.link-btn {
  padding: 15px 30px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.link-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.link-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.link-btn.secondary {
  background: #f8f9fa;
  color: #6c757d;
  border: 2px solid #e9ecef;
}

.link-btn.secondary:hover {
  background: #e9ecef;
  border-color: #dee2e6;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .guide-container {
    padding: 20px 15px;
  }

  .guide-container h1 {
    font-size: 2rem;
  }

  .status-examples {
    grid-template-columns: 1fr;
  }

  .api-categories {
    grid-template-columns: 1fr;
  }

  .quick-links {
    flex-direction: column;
    align-items: center;
  }
}
</style>
