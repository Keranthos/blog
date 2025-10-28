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
/*
 * Live2D Widget
 * https://github.com/stevenjoezhang/live2d-widget
 */

function loadWidget(config) {
	let { waifuPath, apiPath, cdnPath } = config;
	let useCDN = false, modelList;
	if (typeof cdnPath === "string") {
		useCDN = true;
		if (!cdnPath.endsWith("/")) cdnPath += "/";
	} else if (typeof apiPath === "string") {
		if (!apiPath.endsWith("/")) apiPath += "/";
	} else {
		console.error("Invalid initWidget argument!");
		return;
	}
	localStorage.removeItem("waifu-display");
	sessionStorage.removeItem("waifu-text");
	console.log('正在创建看板娘DOM元素...');
		document.body.insertAdjacentHTML("beforeend", `<div id="waifu">
			<div id="waifu-tips"></div>
			<canvas id="live2d" width="800" height="800"></canvas>
			<div id="waifu-tool">
				<span class="tool-icon" title="小游戏">🎮</span>
				<span class="tool-icon" title="换装">👗</span>
				<span class="tool-icon" title="关闭">❌</span>
			</div>
		</div>`);
	console.log('看板娘DOM元素已创建');
	// https://stackoverflow.com/questions/24148403/trigger-css-transition-on-appended-element
	setTimeout(() => {
		document.getElementById("waifu").style.bottom = 0;
	}, 0);

	function randomSelection(obj) {
		return Array.isArray(obj) ? obj[Math.floor(Math.random() * obj.length)] : obj;
	}
	// 检测用户活动状态，并在空闲时显示消息
	let userAction = false,
		userActionTimer,
		messageTimer,
		messageArray = [
			"好久不见，日子过得好快呢……",
			"大坏蛋！你都多久没理人家了呀，嘤嘤嘤～",
			"嗨～快来逗我玩吧！",
			"(｡•́︿•̀｡)<br>不要冷落我嘛～",
			"我在这里等你好久啦！",
			"要不要休息一下呢？",
			"记得多喝水哦～",
			"坐久了要起来活动活动呢！",
			"眼睛累了吗？看看远处吧～",
			"今天过得开心吗？"
		];
	window.addEventListener("mousemove", () => userAction = true);
	window.addEventListener("keydown", () => userAction = true);
	setInterval(() => {
		if (userAction) {
			userAction = false;
			clearInterval(userActionTimer);
			userActionTimer = null;
		} else if (!userActionTimer) {
			userActionTimer = setInterval(() => {
				showMessage(randomSelection(messageArray), 6000, 9);
			}, 20000);
		}
	}, 1000);

	(function registerEventListener() {
		// 等待DOM元素创建后再绑定事件
		setTimeout(() => {
			const toolIcons = document.querySelectorAll("#waifu-tool .tool-icon");
			if (toolIcons.length >= 3) {
				toolIcons[0].addEventListener("click", () => { // 🎮 小游戏
					try {
						console.log('=== 游戏启动调试信息 ===');
						console.log('检查Asteroids函数:', typeof Asteroids);
						console.log('检查window.Asteroids:', typeof window.Asteroids);
						console.log('检查window对象:', window);
						console.log('检查脚本是否加载:', document.querySelector('script[src*="asteroids.js"]'));

						if (typeof window.Asteroids === 'function') {
							showMessage("🎮 飞机大战启动！<br>使用方向键移动，空格键射击", 3000, 10);
							// 延迟启动游戏，确保所有资源都已加载
							setTimeout(() => {
								console.log('开始创建Asteroids游戏实例...');

								// 确保ASTEROIDSPLAYERS数组存在
								if (!window.ASTEROIDSPLAYERS) {
									window.ASTEROIDSPLAYERS = [];
								}

								const game = new window.Asteroids();
								console.log('Asteroids游戏实例创建完成:', game);

								// 手动将游戏实例添加到数组中（模拟自动启动的行为）
								window.ASTEROIDSPLAYERS.push(game);
								console.log('游戏实例已添加到ASTEROIDSPLAYERS数组:', window.ASTEROIDSPLAYERS);

								// 暂停看板娘的定时器活动，防止与游戏冲突
								if (userActionTimer) {
									clearInterval(userActionTimer);
									userActionTimer = null;
									console.log('游戏启动：已暂停看板娘定时器');
								}

								// 监听ESC键退出游戏，自动刷新页面复原界面
								const handleEscKey = (event) => {
									if (event.keyCode === 27) { // ESC键
										console.log('检测到ESC键，准备退出游戏并刷新页面...');
										// 移除事件监听器
										document.removeEventListener('keydown', handleEscKey);
										// 清理游戏实例数组，恢复看板娘活动
										window.ASTEROIDSPLAYERS = [];
										console.log('游戏退出：已清理游戏实例，看板娘活动将恢复');
										// 延迟刷新页面，确保游戏完全退出
										setTimeout(() => {
											window.location.reload();
										}, 500);
									}
								};
								document.addEventListener('keydown', handleEscKey);

								// 检查游戏是否正常启动
								setTimeout(() => {
									const canvas = document.querySelector('#ASTEROIDS-CANVAS');
									const gameContainer = document.querySelector('#ASTEROIDS-GAMECONTAINER');
									console.log('游戏Canvas:', canvas);
									console.log('游戏容器:', gameContainer);
									console.log('游戏是否运行:', window.ASTEROIDS);
									console.log('ASTEROIDSPLAYERS数组:', window.ASTEROIDSPLAYERS);
								}, 1000);
							}, 500);
						} else {
							console.error('Asteroids函数不存在，尝试重新加载脚本...');
							// 尝试重新加载脚本
							const script = document.createElement('script');
							script.src = '/live2d-widget-master/asteroids.js';
							script.onload = () => {
								console.log('脚本重新加载完成，检查Asteroids:', typeof window.Asteroids);
								if (typeof window.Asteroids === 'function') {
									showMessage("🎮 飞机大战启动！<br>使用方向键移动，空格键射击", 3000, 10);
									setTimeout(() => {
										// 确保ASTEROIDSPLAYERS数组存在
										if (!window.ASTEROIDSPLAYERS) {
											window.ASTEROIDSPLAYERS = [];
										}

										const game = new window.Asteroids();
										// 手动将游戏实例添加到数组中
										window.ASTEROIDSPLAYERS.push(game);
										console.log('重新加载后游戏实例已添加到数组:', window.ASTEROIDSPLAYERS);

										// 暂停看板娘的定时器活动，防止与游戏冲突
										if (userActionTimer) {
											clearInterval(userActionTimer);
											userActionTimer = null;
											console.log('游戏启动：已暂停看板娘定时器');
										}

										// 监听ESC键退出游戏，自动刷新页面复原界面
										const handleEscKey = (event) => {
											if (event.keyCode === 27) { // ESC键
												console.log('检测到ESC键，准备退出游戏并刷新页面...');
												// 移除事件监听器
												document.removeEventListener('keydown', handleEscKey);
												// 清理游戏实例数组，恢复看板娘活动
												window.ASTEROIDSPLAYERS = [];
												console.log('游戏退出：已清理游戏实例，看板娘活动将恢复');
												// 延迟刷新页面，确保游戏完全退出
												setTimeout(() => {
													window.location.reload();
												}, 500);
											}
										};
										document.addEventListener('keydown', handleEscKey);
									}, 500);
								} else {
									showMessage("游戏脚本加载失败，请刷新页面重试～", 2000, 9);
								}
							};
							script.onerror = () => {
								showMessage("游戏脚本加载失败，请刷新页面重试～", 2000, 9);
							};
							document.head.appendChild(script);
						}
					} catch (error) {
						console.error('游戏启动失败:', error);
						showMessage("游戏启动失败，请稍后再试～", 2000, 9);
					}
				});
				toolIcons[1].addEventListener("click", () => { // 👗 换装
					if (window.loadRandModel) {
						window.loadRandModel();
					}
				});
				toolIcons[2].addEventListener("click", () => { // ❌ 关闭
					localStorage.setItem("waifu-display", Date.now());
					showMessage("(｡•́︿•̀｡)<br>呜呜…记得要回来看我哦！", 2000, 11);
					document.getElementById("waifu").style.bottom = "-500px";
					setTimeout(() => {
						document.getElementById("waifu").style.display = "none";
						const toggle = document.getElementById("waifu-toggle");
						if (toggle) toggle.classList.add("waifu-toggle-active");
					}, 3000);
				});
				console.log('看板娘工具栏事件已绑定');
			}
		}, 100);
		const devtools = () => {};
		console.log("%c", devtools);
		devtools.toString = () => {
			showMessage("Σ(っ °Д °;)っ<br>你打开了控制台！是想搞事情吗？", 6000, 9);
		};
		window.addEventListener("copy", () => {
			showMessage("你复制了什么？记得注明出处哦～", 6000, 9);
		});
		window.addEventListener("visibilitychange", () => {
			if (!document.hidden) showMessage("(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧<br>哇，你终于回来了～好想你呀！", 6000, 9);
		});
	})();

	(function welcomeMessage() {
		let text;
		if (location.pathname === "/") { // 如果是主页
			const now = new Date().getHours();
			if (now > 5 && now <= 7) text = "早上好呀！☀️ 新的一天开始啦，要元气满满哦！";
			else if (now > 7 && now <= 11) text = "上午好～工作顺利吗？记得多喝水，不要久坐哦！";
			else if (now > 11 && now <= 13) text = "中午啦！🍱 该吃午饭了，要好好吃饭才有力气呢～";
			else if (now > 13 && now <= 17) text = "下午好～☕ 午后容易犯困，要不要起来活动一下？";
			else if (now > 17 && now <= 19) text = "傍晚了！🌆 今天辛苦啦，窗外的夕阳很美呢～";
			else if (now > 19 && now <= 21) text = "晚上好呀！🌙 今天过得开心吗？";
			else if (now > 21 && now <= 23) text = ["😴<br>已经很晚了呢，早点休息吧，晚安～", "夜深了，要爱护眼睛哦！"];
			else text = "🦉<br>哇！你是夜猫子吗？这么晚还不睡，明天起得来嘛？";
		} else if (location.pathname.includes("/blog")) {
			text = "来看博客啦！📝 主人写的文章都很用心呢～";
		} else if (location.pathname.includes("/moments")) {
			text = "碎碎念～💭 主人的日常小心情都在这里啦！";
		} else if (location.pathname.includes("/fragments/books")) {
			text = "书单推荐！📚 这些书都值得一读哦～";
		} else if (location.pathname.includes("/fragments/novels")) {
			text = "小说推荐！📖 主人的品味很不错呢～";
		} else if (location.pathname.includes("/fragments/movies")) {
			text = "电影推荐！🎬 一起来看好电影吧～";
		} else if (location.pathname.includes("/questionbox")) {
			text = "💌<br>悄悄话箱～有什么想问的吗？不要害羞哦！";
		} else if (location.pathname.includes("/timeline")) {
			text = "时间树来啦！⏰ 看看主人的成长轨迹～";
		} else if (location.pathname.includes("/presentation")) {
			text = "讲演展示！🎤 主人的演讲很精彩呢～";
		} else if (location.pathname.includes("/profile")) {
			text = "个人资料～👤 来看看主人的自我介绍吧！";
		} else if (location.pathname.includes("/search")) {
			text = "搜索功能！🔍 在找什么呢？让我帮你找找～";
		} else {
			text = "ヾ(◍°∇°◍)ﾉﾞ<br>欢迎来到主人的小站～";
		}
		showMessage(text, 7000, 8);
	})();

	// 监听路由变化，显示相应页面的欢迎消息
	(function routeChangeListener() {
		let currentPath = location.pathname;
		
		// 监听popstate事件（浏览器前进后退）
		window.addEventListener('popstate', () => {
			if (location.pathname !== currentPath) {
				currentPath = location.pathname;
				setTimeout(() => {
					showPageWelcomeMessage();
				}, 500); // 延迟500ms确保页面加载完成
			}
		});
		
		// 监听pushState和replaceState（程序化导航）
		const originalPushState = history.pushState;
		const originalReplaceState = history.replaceState;
		
		history.pushState = function(...args) {
			originalPushState.apply(this, args);
			if (location.pathname !== currentPath) {
				currentPath = location.pathname;
				setTimeout(() => {
					showPageWelcomeMessage();
				}, 500);
			}
		};
		
		history.replaceState = function(...args) {
			originalReplaceState.apply(this, args);
			if (location.pathname !== currentPath) {
				currentPath = location.pathname;
				setTimeout(() => {
					showPageWelcomeMessage();
				}, 500);
			}
		};
		
		function showPageWelcomeMessage() {
			let text;
			if (location.pathname === "/") {
				text = "回到首页啦～欢迎回来！";
			} else if (location.pathname.includes("/blog")) {
				text = "来看博客啦！📝 主人写的文章都很用心呢～";
			} else if (location.pathname.includes("/moments")) {
				text = "碎碎念～💭 主人的日常小心情都在这里啦！";
			} else if (location.pathname.includes("/fragments/books")) {
				text = "书单推荐！📚 这些书都值得一读哦～";
			} else if (location.pathname.includes("/fragments/novels")) {
				text = "小说推荐！📖 主人的品味很不错呢～";
			} else if (location.pathname.includes("/fragments/movies")) {
				text = "电影推荐！🎬 一起来看好电影吧～";
			} else if (location.pathname.includes("/questionbox")) {
				text = "💌<br>悄悄话箱～有什么想问的吗？不要害羞哦！";
			} else if (location.pathname.includes("/timeline")) {
				text = "时间树来啦！⏰ 看看主人的成长轨迹～";
			} else if (location.pathname.includes("/presentation")) {
				text = "讲演展示！🎤 主人的演讲很精彩呢～";
			} else if (location.pathname.includes("/profile")) {
				text = "个人资料～👤 来看看主人的自我介绍吧！";
			} else if (location.pathname.includes("/search")) {
				text = "搜索功能！🔍 在找什么呢？让我帮你找找～";
			} else {
				text = "ヾ(◍°∇°◍)ﾉﾞ<br>欢迎来到主人的小站～";
			}
			showMessage(text, 5000, 7); // 优先级稍低，避免与初始欢迎消息冲突
		}
	})();

	function showMessage(text, timeout, priority) {
		if (!text || (sessionStorage.getItem("waifu-text") && sessionStorage.getItem("waifu-text") > priority)) return;

		// 检查游戏是否正在运行，如果是则暂停消息显示
		if (window.ASTEROIDSPLAYERS && window.ASTEROIDSPLAYERS.length > 0) {
			console.log('游戏进行中，暂停看板娘消息显示');
			return;
		}

		// 检查DOM元素是否存在
		const tips = document.getElementById("waifu-tips");
		if (!tips) {
			console.log('看板娘DOM元素不存在，跳过消息显示');
			return;
		}

		if (messageTimer) {
			clearTimeout(messageTimer);
			messageTimer = null;
		}
		text = randomSelection(text);
		sessionStorage.setItem("waifu-text", priority);
		tips.innerHTML = text;
		tips.classList.add("waifu-tips-active");
		messageTimer = setTimeout(() => {
			sessionStorage.removeItem("waifu-text");
			if (tips) { // 再次检查DOM元素是否存在
				tips.classList.remove("waifu-tips-active");
			}
		}, timeout);
	}

	// 将 showMessage 暴露到全局，供 Vue 应用使用
	window.showMessage = showMessage;

	(function initModel() {
		let modelId = localStorage.getItem("modelId"),
			modelTexturesId = localStorage.getItem("modelTexturesId");
		if (modelId === null) {
			// 首次访问加载 指定模型 的 指定材质
			modelId = 1; // 模型 ID
			modelTexturesId = 53; // 材质 ID
		}
		console.log('准备加载Live2D模型:', modelId, modelTexturesId);
		loadModel(modelId, modelTexturesId);
		fetch(waifuPath)
			.then(response => response.json())
			.then(result => {
				window.addEventListener("mouseover", event => {
					for (let { selector, text } of result.mouseover) {
						if (!event.target.matches(selector)) continue;
						text = randomSelection(text);
						text = text.replace("{text}", event.target.innerText);
						showMessage(text, 4000, 8);
						return;
					}
				});
				window.addEventListener("click", event => {
					for (let { selector, text } of result.click) {
						if (!event.target.matches(selector)) continue;
						text = randomSelection(text);
						text = text.replace("{text}", event.target.innerText);
						showMessage(text, 4000, 8);
						return;
					}
				});
				result.seasons.forEach(({ date, text }) => {
					const now = new Date(),
						after = date.split("-")[0],
						before = date.split("-")[1] || after;
					if ((after.split("/")[0] <= now.getMonth() + 1 && now.getMonth() + 1 <= before.split("/")[0]) && (after.split("/")[1] <= now.getDate() && now.getDate() <= before.split("/")[1])) {
						text = randomSelection(text);
						text = text.replace("{year}", now.getFullYear());
						//showMessage(text, 7000, true);
						messageArray.push(text);
					}
				});
			});
	})();

	async function loadModelList() {
		try {
			const response = await fetch(`${cdnPath}model_list.json`);
			if (!response.ok) {
				throw new Error(`HTTP ${response.status}`);
			}
			modelList = await response.json();
		} catch (error) {
			console.warn('Live2D 模型列表加载失败:', error);
			// 设置一个默认的模型列表
			modelList = {
				models: [
					["chitose"],
					["haruto"],
					["hibiki"],
					["hijiki"],
					["izumi"],
					["koharu"],
					["shizuku"],
					["tororo"],
					["tsumiki"],
					["unitychan"],
					["wanko"]
				],
				messages: [
					"嗨，很高兴见到你！",
					"你好呀！",
					"很高兴认识你！",
					"欢迎来到我的世界！",
					"很高兴见到你！",
					"你好！",
					"欢迎！",
					"很高兴认识你！",
					"你好呀！",
					"欢迎来到我的世界！",
					"很高兴见到你！"
				]
			};
		}
	}

	async function loadModel(modelId, modelTexturesId, message) {
		localStorage.setItem("modelId", modelId);
		localStorage.setItem("modelTexturesId", modelTexturesId);

		// 显示当前模型名称
		if (!modelList) await loadModelList();

		try {
			// 根据modelId动态选择模型
			console.log('开始加载Live2D模型...');
			console.log('CDN路径:', cdnPath);
			console.log('loadlive2d函数存在:', typeof loadlive2d);

			// 完整的模型映射表（根据CDN实际可用的模型）
			const modelGroups = {
				0: { // Potion-Maker系列
					base: "Potion-Maker/Pio",
					textures: ["Pio"],
					names: ["Pio"]
				},
				1: { // Potion-Maker系列
					base: "Potion-Maker/Tia",
					textures: ["Tia"],
					names: ["Tia"]
				},
				2: { // bilibili-live系列
					base: "bilibili-live/22",
					textures: ["22"],
					names: ["22"]
				},
				3: { // bilibili-live系列
					base: "bilibili-live/33",
					textures: ["33"],
					names: ["33"]
				},
				4: { // ShizukuTalk系列
					base: "ShizukuTalk",
					textures: ["shizuku-48", "shizuku-pajama"],
					names: ["Shizuku", "Shizuku睡衣"]
				},
				5: { // HyperdimensionNeptunia系列 - Neptune
					base: "HyperdimensionNeptunia",
					textures: ["neptune_classic", "nepnep", "neptune_santa", "nepmaid", "nepswim"],
					names: ["Neptune经典", "Neptune普通", "Neptune圣诞", "Neptune女仆", "Neptune泳装"]
				},
				6: { // HyperdimensionNeptunia系列 - Noir
					base: "HyperdimensionNeptunia",
					textures: ["noir_classic", "noir", "noir_santa", "noireswim"],
					names: ["Noir经典", "Noir普通", "Noir圣诞", "Noir泳装"]
				},
				7: { // HyperdimensionNeptunia系列 - Blanc
					base: "HyperdimensionNeptunia",
					textures: ["blanc_classic", "blanc_normal", "blanc_swimwear"],
					names: ["Blanc经典", "Blanc普通", "Blanc泳装"]
				},
				8: { // HyperdimensionNeptunia系列 - Vert
					base: "HyperdimensionNeptunia",
					textures: ["vert_classic", "vert_normal", "vert_swimwear"],
					names: ["Vert经典", "Vert普通", "Vert泳装"]
				},
				9: { // HyperdimensionNeptunia系列 - Nepgear
					base: "HyperdimensionNeptunia",
					textures: ["nepgear", "nepgear_extra", "nepgearswim"],
					names: ["Nepgear", "Nepgear特别", "Nepgear泳装"]
				},
				10: { // HyperdimensionNeptunia系列 - Histoire
					base: "HyperdimensionNeptunia",
					textures: ["histoire", "histoirenohover"],
					names: ["Histoire", "Histoire无悬停"]
				},
				11: { // KantaiCollection系列
					base: "KantaiCollection/murakumo",
					textures: ["murakumo"],
					names: ["Murakumo"]
				}
			};

			// 根据modelId和modelTexturesId选择模型
			const modelGroup = modelGroups[modelId];
			let targetModel = "Potion-Maker/Pio"; // 默认模型
			let currentModelName = `模型${modelId}`;

			if (modelGroup) {
				const textureIndex = Math.min(modelTexturesId, modelGroup.textures.length - 1);
				const textureName = modelGroup.textures[textureIndex];
				const displayName = modelGroup.names[textureIndex];

				// 优先使用本地模型
				if (modelId === 9) {
					// Nepgear本地模型系列
					const nepgearLocalModels = [
						"/live2d-widget-master/models/HyperdimensionNeptunia/nepgear", // 普通版
						"/live2d-widget-master/models/HyperdimensionNeptunia/nepgear_extra", // 特别版
						"/live2d-widget-master/models/HyperdimensionNeptunia/nepgearswim" // 泳装版
					];
					const nepgearNames = ["Nepgear (本地)", "Nepgear特别 (本地)", "Nepgear泳装 (本地)"];

					targetModel = nepgearLocalModels[textureIndex] || nepgearLocalModels[0];
					currentModelName = nepgearNames[textureIndex] || nepgearNames[0];
				} else {
					// 构建完整的模型路径（CDN）
					if (modelGroup.base.includes("/")) {
						// 直接路径模型（如 Potion-Maker/Pio）
						targetModel = modelGroup.base;
					} else {
						// 需要组合路径的模型（如 HyperdimensionNeptunia）
						targetModel = `${modelGroup.base}/${textureName}`;
					}
					currentModelName = displayName;
				}
			}

			// 显示当前模型名称
			showMessage(`${message || '欢迎来到主人的小站～'}<br>当前模型: ${currentModelName}`, 4000, 10);

			// 构建模型URL
			let modelUrl;
			if (targetModel.startsWith("/")) {
				// 本地模型路径
				modelUrl = `${targetModel}/index.json`;
			} else {
				// CDN模型路径
				modelUrl = `${cdnPath}model/${targetModel}/index.json`;
			}

			console.log('完整URL:', modelUrl);

			if (typeof loadlive2d === 'function') {
				loadlive2d("live2d", modelUrl);
				console.log(`正在加载模型: ${currentModelName}`);
				console.log(`模型路径: ${targetModel}`);
				window.modelLoaded = true; // 在真正加载后才标记
			} else {
				console.error('loadlive2d函数不存在！');
				showMessage("Live2D库未加载", 3000, 9);
			}
		} catch (error) {
			console.error('Live2D 模型加载失败:', error);
			// 显示错误信息帮助调试
			showMessage("模型加载失败，请检查控制台", 3000, 9);
		}
	}

	// 将函数暴露到全局作用域
	window.loadModel = loadModel;
	window.loadRandModel = loadRandModel;
	window.loadOtherModel = loadOtherModel;

	async function loadRandModel() {
		const modelId = localStorage.getItem("modelId"),
			modelTexturesId = localStorage.getItem("modelTexturesId");

		// 只对Nepgear模型(9)进行本地换装
		if (modelId == 9) {
			const currentTextureId = parseInt(modelTexturesId) || 0;
			const nextTextureId = (currentTextureId + 1) % 3; // 0, 1, 2 循环

			const textureNames = ["普通版", "特别版", "泳装版"];
			showMessage(`✨<br>换上${textureNames[nextTextureId]}！`, 3000, 10);
			loadModel(9, nextTextureId, "✨<br>我的新衣服好看吗？");
		} else {
			// 其他模型使用CDN换装
			if (useCDN) {
				if (!modelList) await loadModelList();
				const target = randomSelection(modelList.models[modelId]);
				loadlive2d("live2d", `${cdnPath}model/${target}/index.json`);
				showMessage("我的新衣服好看嘛？", 4000, 10);
			} else {
				// 可选 "rand"(随机), "switch"(顺序)
				fetch(`${apiPath}rand_textures/?id=${modelId}-${modelTexturesId}`)
					.then(response => response.json())
					.then(result => {
						if (result.textures.id === 1 && (modelTexturesId === 1 || modelTexturesId === 0)) showMessage("(｡•́︿•̀｡)<br>我还没有其他衣服呢！", 4000, 10);
						else loadModel(modelId, result.textures.id, "✨<br>我的新衣服好看吗？");
					});
			}
		}
	}

	async function loadOtherModel() {
		let modelId = localStorage.getItem("modelId");
		if (useCDN) {
			if (!modelList) await loadModelList();
			const index = (++modelId >= modelList.models.length) ? 0 : modelId;
			loadModel(index, 0, modelList.messages[index]);
		} else {
			fetch(`${apiPath}switch/?id=${modelId}`)
				.then(response => response.json())
				.then(result => {
					loadModel(result.model.id, 0, result.model.message);
				});
		}
	}
}

function initWidget(config, apiPath) {
	if (typeof config === "string") {
		config = {
			waifuPath: config,
			apiPath
		};
	}
	document.body.insertAdjacentHTML("beforeend", `<div id="waifu-toggle">
			<span>看板娘</span>
		</div>`);
	const toggle = document.getElementById("waifu-toggle");
	toggle.addEventListener("click", () => {
		toggle.classList.remove("waifu-toggle-active");
		if (toggle.getAttribute("first-time")) {
			loadWidget(config);
			toggle.removeAttribute("first-time");
		} else {
			localStorage.removeItem("waifu-display");
			document.getElementById("waifu").style.display = "";
			setTimeout(() => {
				document.getElementById("waifu").style.bottom = 0;
			}, 0);
		}
	});
	// 简化逻辑，直接加载看板娘，不检查localStorage
	console.log('正在加载看板娘组件...');
	loadWidget(config);
}

// 初始化模型 - 固定加载模型
(function initModel() {
	if (!window.ASTEROIDS) {
		// 固定使用本地Nepgear模型
		localStorage.setItem("modelId", 9);
		localStorage.setItem("modelTexturesId", 0);

		console.log('准备加载本地Nepgear模型...');
		// 延迟加载模型，等待loadWidget完成
		setTimeout(() => {
			if (window.loadModel) {
				window.loadModel(9, 0, "ヾ(◍°∇°◍)ﾉﾞ<br>欢迎来到主人的小站～");
			}
		}, 1500);
	}
})();
