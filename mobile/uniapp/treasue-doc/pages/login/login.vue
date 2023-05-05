<template>
	<view class="login-box">
		<view>
			<form @submit="formSubmit" @reset="formReset">
				<view class="uni-form-item uni-column">
					<view class="title">用户名：</view>
					<view>
						<input class="uni-input" focus name="username" placeholder="用户名" />
					</view>
				</view>
				<view class="uni-form-item uni-column">
					<view class="title">密码：</view>
					<view>
						<input class="uni-input" name="password" placeholder="密码" />
					</view>
				</view>
				<view class="uni-form-item uni-column">
					<view>
						<checkbox-group>
							<label>
								<checkbox value="cb" />记住密码
							</label>
						</checkbox-group>
					</view>
				</view>
				<view class="uni-btn-v">
					<button form-type="submit">Submit</button>
					<button type="default" form-type="reset">Reset</button>
				</view>
			</form>
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
		background-color: pink;

		.uni-form-item,
		.title {
			padding: 20rpx 0;
		}
	}
</style>