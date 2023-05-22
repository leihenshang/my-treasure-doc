<script>
	import
	store
	from "@/store/index.js"

	export default {
		onLaunch: function() {
			let token = uni.getStorageSync('userInfo');
			if (token) {
				//存在则关闭启动页进入首页
				uni.switchTab({
					url: "/pages/index/index"
				})
				store.commit("setUserInfo", token)
			} else {
				//不存在则跳转至登录页
				uni.reLaunch({
					url: "/pages/login/login",
				})
			}

			uni.addInterceptor('request', {
				invoke(args) {
					// request 触发前拼接 url 
					// args.url = args.url

				},
				success(args) {
					// 请求成功后
					console.log('interceptor-success', args)
					// 检查是否登录出现问题
					if (args.statusCode == 401) {
						console.log('response status code is 401')
						uni.removeStorageSync('userInfo')
						store.commit('setUserInfo', {})
						uni.redirectTo({
							url: "/pages/login/login"
						})
					}
				},
				fail(err) {
					console.log('interceptor-fail', err)
				},
				complete(res) {
					// console.log('interceptor-complete', res)
				}
			})



			console.warn('当前组件仅支持 uni_modules 目录结构 ，请升级 HBuilderX 到 3.1.0 版本以上！')
			console.log('App Launch')
		},
		onShow: function() {
			console.log('App Show')
		},
		onHide: function() {
			console.log('App Hide')
		}
	}
</script>

<style lang="scss">
	/*每个页面公共css */
	@import '@/uni_modules/uni-scss/index.scss';
	/* #ifndef APP-NVUE */
	@import '@/static/customicons.css';

	// 设置整个项目的背景色
	page {
		background-color: #f5f5f5;
		width: 100%;
		height: 100%;
	}

	uni-page-body,
	#app {
		height: 100%;
	}

	/* #endif */
	.example-info {
		font-size: 14px;
		color: #333;
		padding: 10px;
	}
</style>