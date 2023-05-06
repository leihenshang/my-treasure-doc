<template>
	<view class="uni-common-mt">
		<view class="btn-group">
			<button type="primary" @tap="save">保存</button>
			<button type="warn" @tap="undo">取消</button>
		</view>
		<view class="uni-form-item uni-column">
			<input class="uni-input" focus placeholder="输入你的标题" />
		</view>
		<view class="uni-form-item uni-column container">
			<editor id="editor" class="ql-container" :placeholder="placeholder" @ready="onEditorReady"></editor>
		</view>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				placeholder: '挥洒你的创意吧...'
			};
		},
		methods: {
			onEditorReady() {
				// #ifdef MP-BAIDU
				this.editorCtx = requireDynamicLib('editorLib').createEditorContext('editor');
				// #endif

				// #ifdef APP-PLUS || H5 ||MP-WEIXIN
				uni.createSelectorQuery().select('#editor').context((res) => {
					this.editorCtx = res.context
				}).exec()
				// #endif
			},
			undo() {
				this.editorCtx.undo()
			},
			save() {
				this.editorCtx.getContents({
					success: (data) => {
						uni.showModal({
							content: "表单数据为:" + data.text,
							showCancel: false
						})
					}
				})
			},
		}
	}
</script>

<style lang="scss">
	.container {
		padding: 10rpx;
	}

	#editor {
		width: 100%;
		height: 800rpx;
		background-color: #CCCCCC;
	}

	.btn-group {
		display: flex;
		justify-content: space-between;
		padding: 0 10rpx;

		button {
			margin: 10rpx 0;
			width: 200rpx;
		}
	}
</style>