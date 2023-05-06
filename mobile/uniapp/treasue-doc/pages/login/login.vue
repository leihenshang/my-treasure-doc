<template>
	<view class="login-box">

		<view class="input-box">
			<uni-forms @submit="formSubmit" @reset="formReset">
				<view class="uni-form-item uni-column">
					<uni-forms-item label="用户名" required>
						<uni-easyinput placeholder="请输入姓名" focus name="username" />
					</uni-forms-item>
				</view>
				<view class="uni-form-item uni-column">
					<uni-forms-item label="密码" required>
						<uni-easyinput placeholder="请输入密码" focus name="password" />
					</uni-forms-item>
				</view>
				<view class="uni-form-item uni-column">
					<checkbox-group>
						<label>
							<checkbox value="cb" />记住密码
						</label>
					</checkbox-group>

				</view>
				<view class="uni-btn-v">
					<button form-type="submit" type="primary">登录</button>
					<button type="warn">注册</button>
				</view>
			</uni-forms>
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
			formSubmit: function(e) {
				console.log('form发生了submit事件，携带数据为：' + JSON.stringify(e.detail.value))
				var formdata = e.detail.value
				uni.showModal({
					content: '表单数据内容：' + JSON.stringify(formdata),
					showCancel: false
				});
			},
			formReset: function(e) {
				console.log('清空数据')
			}
		}
	}
</script>

<style lang="scss">
	.login-box {
		padding: 50rpx 20rpx;

		.uni-btn-v {
			padding: 0 10rpx;
			width: 100%;

			button {
				margin: 10rpx 0;
				width: 100%;
			}
		}
	}
</style>