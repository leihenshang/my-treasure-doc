<template>
	<page-meta>
		<navigation-bar :title="nbTitle" :title-icon="titleIcon" :title-icon-radius="titleIconRadius"
			:subtitle-text="subtitleText" :subtitle-color="nbFrontColor" :loading="nbLoading"
			:front-color="nbFrontColor" :background-color="nbBackgroundColor" :color-animation-duration="2000"
			color-animation-timing-func="easeIn" />
	</page-meta>
	<view class="edit-box">
		<uni-forms ref="form" label-position="top" :model="formData">
			<view class="btn-group">
				<button type="primary" @click="save">保存</button>
				<button type="warn" @click="undo">上一步</button>
			</view>
			<uni-forms-item label="标题" name="">
				<uni-easyinput focus placeholder="输入你的标题" class="title" v-model="formData.title"></uni-easyinput>
			</uni-forms-item>
			<uni-forms-item label="内容" name="" class="">

				<editor id="editor" class="ql-container" :placeholder="placeholder" @ready="onEditorReady"
					@input="editorInput"></editor>
			</uni-forms-item>
		</uni-forms>
	</view>
</template>

<script>
	import {
		docCreate
	} from "@/request/api.js"

	export default {
		data() {
			return {
				placeholder: '挥洒你的创意吧...',
				formData: {
					title: "",
					content: ""
				},
				nbTitle: '编辑',
				titleIcon: '/static/logo.png',
				titleIconRadius: '20px',
				subtitleText: 'subtitleText',
				nbLoading: false,
				nbFrontColor: '#000000',
				nbBackgroundColor: '#ffffff',

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
			editorInput(e) {
				this.formData.content = e.detail.text
			},
			undo() {
				this.editorCtx.undo()
			},
			save() {
				console.log(this.formData)
				this.editorCtx.getContents({
					success: (data) => {
						docCreate(this.formData).then((res) => {
							uni.reLaunch({
								url: "/pages/index/index",
							})
						}).catch(res => {
							console.log(res)
						})

					}
				})
			}
		}

	}
</script>

<style lang="scss">
	.edit-box {
		padding: 10rpx 10rpx 10rpx;

		#editor {
			width: 100%;
			height: 800rpx;
			background-color: #FFF;
			padding: 10rpx;
		}

		.btn-group {
			display: flex;
			justify-content: space-between;

			button {
				margin: 10rpx 0;
			}
		}
	}
</style>