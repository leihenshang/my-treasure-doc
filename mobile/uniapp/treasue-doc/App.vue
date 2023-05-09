<script>
	export default {
		onLaunch: function() {
			let token = uni.getStorageSync('token');
			if (token) {
				//存在则关闭启动页进入首页

			} else {
				//不存在则跳转至登录页
				uni.reLaunch({
					url: "/pages/login/login",
					success: () => {}
				})
			}

			uni.addInterceptor('request', {
				invoke(args) {
					// request 触发前拼接 url 
					args.url = args.url
				},
				success(args) {
					// 请求成功后，修改code值为1
					args.data.code = 1
					console.log('interceptor-success')
				},
				fail(err) {
					console.log('interceptor-fail', err)
				},
				complete(res) {
					console.log('interceptor-complete', res)
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
	}

	/* #endif */
	.example-info {
		font-size: 14px;
		color: #333;
		padding: 10px;
	}
</style>