<template>
	<view class="myCenter-whole">
		<view class="myCenter-top">
			<img src="https://avatar.cdn.deepin.com/public/default.png" alt="avatar">
			<h1>我的名字</h1>
		</view>
		<view class="myCenter-bottom">
			<button type="primary">退出登录</button>
		</view>
	</view>
</template>

<script>
	export default {
		data() {
			return {

			};
		}
	}
</script>

<style lang="scss">
	.myCenter-whole {
		margin: 20rpx;

		.myCenter-top {
			img {
				width: 200rpx;
				height: 200rpx;
			}

			h1 {
				color: pink;
			}
		}
	}
</style>