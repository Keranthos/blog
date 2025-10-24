/* eslint-disable brace-style */
/* eslint-disable camelcase */
/* eslint-disable no-unused-vars */
/* eslint-disable spaced-comment */
/* eslint-disable no-undef */
/* eslint-disable prefer-const */
/* eslint-disable no-return-assign */
/* eslint-disable quotes */
/* eslint-disable semi */
/* eslint-disable no-tabs */
/* eslint-disable space-before-function-paren */
/* eslint-disable one-var */
/* eslint-disable indent */

// 注意：live2d_path 参数应使用绝对路径
const live2d_path = "/live2d-widget-master/";

// 添加全局错误处理器来捕获 Live2D 错误
window.addEventListener('error', function(event) {
	if (event.filename && event.filename.includes('live2d')) {
		console.warn('Live2D 错误已被捕获:', event.message);
		// 不阻止错误冒泡，只记录警告
	}
});

// 添加未处理的 Promise 拒绝处理器
window.addEventListener('unhandledrejection', function(event) {
	if (event.reason && event.reason.toString().includes('live2d')) {
		console.warn('Live2D Promise 错误已被捕获:', event.reason);
		// 不阻止错误冒泡，只记录警告
	}
});

// 封装异步加载资源的方法
function loadExternalResource(url, type) {
	return new Promise((resolve, reject) => {
		let tag;

		if (type === "css") {
			tag = document.createElement("link");
			tag.rel = "stylesheet";
			tag.href = url;
		}
		else if (type === "js") {
			tag = document.createElement("script");
			tag.src = url;
		}
		if (tag) {
			tag.onload = () => resolve(url);
			tag.onerror = () => reject(url);
			document.head.appendChild(tag);
		}
	});
}

// 加载 waifu.css live2d.min.js waifu-tips.js
// 移除屏幕宽度限制，让看板娘在所有设备上都能显示
if (true) {
	Promise.all([
		loadExternalResource(live2d_path + "waifu.css", "css"),
		loadExternalResource(live2d_path + "live2d.min.js", "js"),
		loadExternalResource(live2d_path + "waifu-tips.js", "js"),
		loadExternalResource(live2d_path + "asteroids.js", "js")
	]).then(() => {
		console.log('Live2D 资源加载成功，正在初始化看板娘...');
		initWidget({
			waifuPath: live2d_path + "waifu-tips.json",
			cdnPath: "https://fastly.jsdelivr.net/gh/fghrsh/live2d_api/"
		});
	}).catch((error) => {
		console.warn('Live2D 资源加载失败，将禁用看板娘功能:', error);
		// 不显示错误信息给用户，静默处理
	});
}

// initWidget 第一个参数为 waifu-tips.json 的路径，第二个参数为 API 地址
// API 后端可自行搭建，参考 https://github.com/fghrsh/live2d_api
// 初始化看板娘会自动加载指定目录下的 waifu-tips.json

console.log(`
  く__,.ヘヽ.        /  ,ー､ 〉
           ＼ ', !-─‐-i  /  /´
           ／｀ｰ'       L/／｀ヽ､
         /   ／,   /|   ,   ,       ',
       ｲ   / /-‐/  ｉ  L_ ﾊ ヽ!   i
        ﾚ ﾍ 7ｲ｀ﾄ   ﾚ'ｧ-ﾄ､!ハ|   |
          !,/7 '0'     ´0iソ|    |
          |.从"    _     ,,,, / |./    |
          ﾚ'| i＞.､,,__  _,.イ /   .i   |
            ﾚ'| | / k_７_/ﾚ'ヽ,  ﾊ.  |
              | |/i 〈|/   i  ,.ﾍ |  i  |
             .|/ /  ｉ：    ﾍ!    ＼  |
              kヽ>､ﾊ    _,.ﾍ､    /､!
              !'〈//｀Ｔ´', ＼ ｀'7'ｰr'
              ﾚ'ヽL__|___i,___,ンﾚ|ノ
                  ﾄ-,/  |___./
                  'ｰ'    !_,.:
`);
