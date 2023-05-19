<template>
	<view class="myCenter-whole">
		<view class="container">
			<view class="myCenter-top">
				<view class="avatar">
					<img src="https://avatar.cdn.deepin.com/public/default.png" alt="avatar">
				</view>
				<view class="title">
					<h4>{{userInfo.account ?? 'luck user'}}</h4>
					<h4>{{userInfo.email ?? ''}}</h4>
				</view>

			</view>
			<view class="myCenter-bottom">
				<button type="primary" @click="loginOut">退出登录</button>
			</view>
		</view>
	</view>
</template>

<script>
	export default {
		data() {
			return {

			};
		},
		methods: {
			loginOut() {
				uni.removeStorageSync('userInfo');
				this.$store.commit('setUserInfo', {})
				uni.showModal({
					content: '退出登录成功',
					showCancel: false
				});
				uni.redirectTo({
					url: "/pages/login/login"
				})
			}
		},
		computed: {
			userInfo() {
				return this.$store.state.userInfo
			}
		}
	}
</script>

<style lang="scss">
	.myCenter-whole {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;

		.container {
			padding: 0 20rpx;
			width: 100%;
			position: relative;
			top: -200rpx;


			.myCenter-top {

				.avatar {
					text-align: center;

					img {
						width: 200rpx;
						height: 200rpx;
					}
				}

				.title {
					margin: 20rpx 0;

					h4 {
						text-align: center;
						margin: 20rpx 0;
					}
				}

			}

			.myCenter-bottom {}
		}
	}
</style>