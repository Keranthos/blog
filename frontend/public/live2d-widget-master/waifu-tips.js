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
			const live2dCanvas = document.getElementById("live2d");
			if (live2dCanvas) {
				// 左键点击看板娘 - 换装
				live2dCanvas.addEventListener("click", (event) => {
					if (event.button === 0) { // 左键
						if (window.loadRandModel) {
							window.loadRandModel();
						}
					}
				});

				// 右键点击看板娘 - 游戏
				live2dCanvas.addEventListener("contextmenu", (event) => {
					event.preventDefault(); // 阻止默认右键菜单
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
			}

			// ESC键消失功能
			let isHoveringWaifu = false;
			const waifuElement = document.getElementById("waifu");
			if (waifuElement) {
				waifuElement.addEventListener("mouseenter", () => {
					isHoveringWaifu = true;
				});
				waifuElement.addEventListener("mouseleave", () => {
					isHoveringWaifu = false;
				});
			}

			// 监听ESC键
			document.addEventListener("keydown", (event) => {
				if (event.key === "Escape" && isHoveringWaifu) {
					localStorage.setItem("waifu-display", Date.now());
					showMessage("(｡•́︿•̀｡)<br>呜呜…记得要回来看我哦！", 2000, 11);
					document.getElementById("waifu").style.bottom = "-500px";
					setTimeout(() => {
						document.getElementById("waifu").style.display = "none";
						const toggle = document.getElementById("waifu-toggle");
						if (toggle) toggle.classList.add("waifu-toggle-active");
					}, 3000);
				}
				});
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
			if (now > 5 && now <= 7) text = "早上好，一日之计在于晨，很高兴在这时候看到你，朋友";
			else if (now > 7 && now <= 11) text = "上午好，希望你有快乐的一天，而我现在一般在考虑中午吃什么";
			else if (now > 11 && now <= 13) text = "中午到了，午饭和午觉对我都是不可或缺的";
			else if (now > 13 && now <= 17) text = "下午好，很快就可以到令人愉悦的晚上了，我的朋友";
			else if (now > 17 && now <= 19) text = "傍晚了，今天很辛苦，也许可以考虑喝一杯奶茶？";
			else if (now > 19 && now <= 21) text = "晚上好呀，现在可以稍微放松一下了，宁静的夜晚总是让我惬意";
			else if (now > 21 && now <= 23) text = ["😴已经很晚了呢，明天又是崭新的一天，晚安玛卡巴卡", "夜深了，要爱护眼睛哦"];
			else text = "🦉注意身体，请尽快休息，我的朋友";
		} else if (location.pathname.includes("/blog")) {
			const blogMessages = [
				"我努力想探寻世界上所有好玩的东西，这里留着一点足迹",
				"这里除了代码，还有很多有趣的事情",
				"这是我自己一点点搭建的博客，也许只是想在虚拟的网络中拥有一点属于自己的东西"
			];
			text = blogMessages[Math.floor(Math.random() * blogMessages.length)];
		} else if (location.pathname.includes("/moments")) {
			const momentsMessages = [
				"一些散落在过去时光中的故事与感想，好故事可是昂贵的，我的朋友",
				"你说有没有可能，正在看这个网站的你，已经或将要为我带来某一篇随笔的些许想法呢？",
				"文章只是一张地图，找到回忆还需要跋山涉水"
			];
			text = momentsMessages[Math.floor(Math.random() * momentsMessages.length)];
		} else if (location.pathname === "/fragments/books") {
			const booksMessages = [
				"您有多久没有完整读过一本书？那里面有与我们共鸣的灵魂",
				"我记不清上一次打开非专业的纸质书是什么时候了，但它们一定组成了我的一部分",
				"也许您有好书要推荐给我，朋友？",
				"我会永远怀念小时候和亲人朋友在书店里面打发的时光，哪怕那时候的书店朴实的只有书"
			];
			text = booksMessages[Math.floor(Math.random() * booksMessages.length)];
		} else if (location.pathname === "/fragments/novels") {
			const novelsMessages = [
				"仔细回想，小说已经贯穿了我到现在为止超过一半的人生",
				"虚构的故事里藏着最真实的情感，欢迎来到想象的世界",
				"我们在小说中寻找想象里的自己，但不要忘了自己不是找到的，而是创造的"
			];
			text = novelsMessages[Math.floor(Math.random() * novelsMessages.length)];
		} else if (location.pathname === "/fragments/movies") {
			const moviesMessages = [
				"无关乎题材与风格，我总是痴迷于所有的好电影",
				"我喜欢有年代感的电影，它里面有时代的影子",
				"好的电影不厌百回看，也许现在是时候翻出你最爱的电影了，朋友"
			];
			text = moviesMessages[Math.floor(Math.random() * moviesMessages.length)];
		} else if (location.pathname.includes("/questionbox")) {
			const questionMessages = [
				"知无不言，言无不尽",
				"每一个问题都是成长的契机，欢迎交流"
			];
			text = questionMessages[Math.floor(Math.random() * questionMessages.length)];
		} else if (location.pathname.includes("/timeline")) {
			const timelineMessages = [
				"朋友，你觉得时间轴应该是一支箭，还是一棵树？",
				"现在的我依然是菜菜的，但是三十年河东三十年河西，莫欺少年穷、莫欺中年穷、莫欺老年穷、似者为大，朋友"
			];
			text = timelineMessages[Math.floor(Math.random() * timelineMessages.length)];
		} else if (location.pathname.includes("/presentation")) {
			const presentationMessages = [
				"作为一个很讨厌麻烦的i人，我必须承认讲演时常困扰我",
				"我想尝试一下Slidev,做PPT实在太麻烦了"
			];
			text = presentationMessages[Math.floor(Math.random() * presentationMessages.length)];
		} else if (location.pathname.includes("/profile")) {
			const profileMessages = [
				"想要知道我的故事？也许您可以直接问问我",
				"真正了解一个人就像开包抽卡，谁也不知道开出来的是金色传说还是白色普通，但只要开的够多总能开出来的。你觉得呢，我的朋友？"
			];
			text = profileMessages[Math.floor(Math.random() * profileMessages.length)];
		} else if (location.pathname.includes("/search")) {
			const searchMessages = [
				"如果在这里也没有的话，可以试试催更",
				"现在这个搜索支持了标题、标签与内容搜索，如果依然没找到那也许就是我的知识盲区"
			];
			text = searchMessages[Math.floor(Math.random() * searchMessages.length)];
		} else {
			text = "ヾ(◍°∇°◍)ﾉﾞ<br>可以随便逛逛，也许有不一样的小彩蛋或者新发现呢";
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
			const homeMessages = [
				"回到首页啦～欢迎回来，朋友"
			];
			text = homeMessages[Math.floor(Math.random() * homeMessages.length)];
			} else if (location.pathname.includes("/blog")) {
				const blogMessages = [
					"我努力想探寻世界上所有好玩的东西，这里留着一点足迹",
					"这里除了代码，还有很多有趣的事情",
					"这是我自己一点点搭建的博客，也许只是想在虚拟的网络中拥有一点属于自己的东西"
				];
				text = blogMessages[Math.floor(Math.random() * blogMessages.length)];
			} else if (location.pathname.includes("/moments")) {
				const momentsMessages = [
					"一些散落在过去时光中的故事与感想，好故事可是昂贵的，我的朋友",
					"你说有没有可能，正在看这个网站的你，已经或将要为我带来某一篇随笔的些许想法呢？",
					"文章只是一张地图，找到回忆还需要跋山涉水",
				];
				text = momentsMessages[Math.floor(Math.random() * momentsMessages.length)];
			} else if (location.pathname === "/fragments/books") {
				const booksMessages = [
					"您有多久没有完整读过一本书？那里面有与我们共鸣的灵魂",
					"我记不清上一次打开非专业的纸质书是什么时候了，但它们一定组成了我的一部分",
					"也许您有好书要推荐给我，朋友？",
					"我会永远怀念小时候和亲人朋友在书店里面打发的时光，哪怕那时候的书店朴实的只有书"
				];
				text = booksMessages[Math.floor(Math.random() * booksMessages.length)];
			} else if (location.pathname === "/fragments/novels") {
				const novelsMessages = [
					"仔细回想，小说已经贯穿了我到现在为止超过一半的人生",
					"虚构的故事里藏着最真实的情感，欢迎来到想象的世界",
					"我们在小说中寻找想象里的自己，但不要忘了自己不是找到的，而是创造的"
				];
				text = novelsMessages[Math.floor(Math.random() * novelsMessages.length)];
			} else if (location.pathname === "/fragments/movies") {
				const moviesMessages = [
					"无关乎题材与风格，我总是痴迷于所有的好电影",
					"我喜欢有年代感的电影，它里面有时代的影子",
					"好的电影不厌百回看，也许现在是时候翻出你最爱的电影了，朋友"
				];
				text = moviesMessages[Math.floor(Math.random() * moviesMessages.length)];
			} else if (location.pathname.includes("/questionbox")) {
				const questionMessages = [
					"知无不言，言无不尽",
					"每一个问题都是成长的契机，欢迎交流"
				];
				text = questionMessages[Math.floor(Math.random() * questionMessages.length)];
			} else if (location.pathname.includes("/timeline")) {
				const timelineMessages = [
					"朋友，你觉得时间轴应该是一支箭，还是一棵树？",
					"现在的我依然是菜菜的，但是三十年河东三十年河西，莫欺少年穷、莫欺中年穷、莫欺老年穷、似者为大，朋友"
				];
				text = timelineMessages[Math.floor(Math.random() * timelineMessages.length)];
			} else if (location.pathname.includes("/presentation")) {
				const presentationMessages = [
					"作为一个很讨厌麻烦的i人，我必须承认讲演时常困扰我",
					"我想尝试一下Slidev,做PPT实在太麻烦了"
				];
				text = presentationMessages[Math.floor(Math.random() * presentationMessages.length)];
			} else if (location.pathname.includes("/profile")) {
				const profileMessages = [
					"想要知道我的故事？也许您可以直接问问我",
					"真正了解一个人就像开包抽卡，谁也不知道开出来的是金色传说还是白色普通，但只要开的够多总能开出来的。你觉得呢，我的朋友？"
				];
				text = profileMessages[Math.floor(Math.random() * profileMessages.length)];
			} else if (location.pathname.includes("/search")) {
				const searchMessages = [
					"如果在这里也没有的话，可以试试催更",
					"现在这个搜索支持了标题、标签与内容搜索，如果依然没找到那也许就是我的知识盲区"
				];
				text = searchMessages[Math.floor(Math.random() * searchMessages.length)];
			} else {
				text = "ヾ(◍°∇°◍)ﾉﾞ<br>可以随便逛逛，也许有不一样的小彩蛋或者新发现呢";
			}
			showMessage(text, 5000, 9); // 提高优先级到9，确保路由变化消息能显示
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

	async function loadModel(modelId, modelTexturesId, message) {
		localStorage.setItem("modelId", modelId);
		localStorage.setItem("modelTexturesId", modelTexturesId);

		try {
			console.log('开始加载Live2D模型...');
			console.log('loadlive2d函数存在:', typeof loadlive2d);

			// 只支持本地Nepgear模型系列
					const nepgearLocalModels = [
						"/live2d-widget-master/models/HyperdimensionNeptunia/nepgear", // 普通版
						"/live2d-widget-master/models/HyperdimensionNeptunia/nepgear_extra", // 特别版
						"/live2d-widget-master/models/HyperdimensionNeptunia/nepgearswim" // 泳装版
					];
					const nepgearNames = ["Nepgear (本地)", "Nepgear特别 (本地)", "Nepgear泳装 (本地)"];

			const textureIndex = Math.min(modelTexturesId, nepgearLocalModels.length - 1);
			const targetModel = nepgearLocalModels[textureIndex] || nepgearLocalModels[0];
			const currentModelName = nepgearNames[textureIndex] || nepgearNames[0];

			// 显示当前模型名称
			showMessage(`${message || '欢迎来到我的小站～'}`, 4000, 10);

			// 构建模型URL
			const modelUrl = `${targetModel}/index.json`;

			console.log('完整URL:', modelUrl);

			if (typeof loadlive2d === 'function') {
				loadlive2d("live2d", modelUrl);
				console.log(`正在加载模型: ${currentModelName}`);
				console.log(`模型路径: ${targetModel}`);
				window.modelLoaded = true;
			} else {
				console.error('loadlive2d函数不存在！');
				showMessage("Live2D库未加载", 3000, 9);
			}
		} catch (error) {
			console.error('Live2D 模型加载失败:', error);
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

		// 只支持本地Nepgear模型换装
			const currentTextureId = parseInt(modelTexturesId) || 0;
			const nextTextureId = (currentTextureId + 1) % 3; // 0, 1, 2 循环

			const textureNames = ["普通版", "特别版", "泳装版"];
			showMessage(`✨<br>换上${textureNames[nextTextureId]}！`, 3000, 10);
			loadModel(9, nextTextureId, "✨<br>我的新衣服好看吗？");
	}

	async function loadOtherModel() {
		// 只支持本地Nepgear模型，不需要切换其他模型
		showMessage("(｡•́︿•̀｡)<br>我只有这一套衣服呢！", 4000, 10);
	}

	// 将函数暴露到全局作用域
	window.loadModel = loadModel;
	window.loadRandModel = loadRandModel;
	window.loadOtherModel = loadOtherModel;
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
				window.loadModel(9, 0, "ヾ(◍°∇°◍)ﾉﾞ<br>欢迎来到我的小站～");
			}
		}, 1500);
	}
})();
