<template>
	<view class="login-box">
		<uni-forms label-position="top" :border="false" :modelValue="formData" ref="form" :rules="rules">
			<view class="uni-form-item uni-column column-img">
				<view class="image-box">
					<image src="../../static/notebook.png"></image>
				</view>
			</view>
			<view class="uni-form-item uni-column">
				<uni-forms-item label="用户名" required name="account">
					<uni-easyinput placeholder="请输入用户名" v-model="formData.account" focus name="account" />
				</uni-forms-item>
			</view>
			<view class="uni-form-item uni-column">
				<uni-forms-item label="密码" required name="password">
					<uni-easyinput type="password" placeholder="请输入密码" v-model="formData.password" name="password" />
				</uni-forms-item>
			</view>
			<view class="uni-form-item uni-column login-btn">
				<uni-forms-item>
					<button form-type="submit" @click="formSubmit" type="primary">登录</button>
				</uni-forms-item>
			</view>
		</uni-forms>
	</view>
</template>

<script>
	import {
		login as apiLogin
	} from "@/request/api.js"

	export default {
		data() {
			return {
				formData: {
					account: '',
					password: '',
					verifyCode: "123456789",
				},
				rules: {
					account: {
						rules: [{
								required: true,
								errorMessage: '请输入用户名',
							},
							{
								minLength: 3,
								maxLength: 200,
								errorMessage: '姓名长度在 {minLength} 到 {maxLength} 个字符',
							}
						]
					},
					password: {
						rules: [{
								required: true,
								errorMessage: '请输入密码',
							},
							{
								minLength: 3,
								maxLength: 200,
								errorMessage: '密码长度在 {minLength} 到 {maxLength} 个字符',
							}
						]
					}
				}
			};
		},
		methods: {
			formSubmit: function(e) {
				this.$refs.form.validate().then(res => {
					console.log('form发生了submit事件,表单数据信息：', res);
					apiLogin(this.formData).then(res => {
						console.log(res)
						uni.setStorageSync('userInfo', res)
						this.$store.commit('setUserInfo', res)
						console.log(this.$store.state.userInfo)
						uni.switchTab({
							url: "/pages/index/index"
						})
					}).catch(err => {
						console.log(err)
						uni.showToast({
							icon: "none",
							title: '登录失败:' + err
						});
					})


				}).catch(err => {
					console.log('表单错误信息：', err);
				})

			},

		},
		onLoad() {
			console.log('login page onload', this.$store.state.userInfo)
		}
	}
</script>

<style lang="scss">
	.login-box {
		display: flex;
		flex: auto;
		// min-height: 1080rpx;
		height: 100%;
		flex-direction: row;
		align-items: flex-start;


		.uni-forms {
			width: 95%;
			margin: 20% auto;

			.column-img {
				.image-box {
					margin: 0 auto;
					border-radius: 50%;
					overflow: hidden;
					width: 200rpx;
					height: 200rpx;
					text-align: center;
					background-color: pink;
					display: flex;
					flex: auto;
					flex-direction: row;
					align-items: center;
					justify-content: center;

					image {
						height: 150rpx;
						width: 150rpx;
					}
				}
			}

		}
	}
</style>