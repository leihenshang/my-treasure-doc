<template>
	<view class="edit-box">
		<uni-forms ref="form" label-position="top" :model="docData">
			<view class="btn-group">
				<button type="primary" @click="save">保存</button>
			</view>
			<uni-forms-item label="标题" name="">
				<uni-easyinput focus class="title" v-model="docData.title"></uni-easyinput>
			</uni-forms-item>
			<uni-forms-item label="内容" name="" class="">
				<editor id="editor" class="ql-container" @ready="onEditorReady" @input="editorInput"></editor>
			</uni-forms-item>
		</uni-forms>
	</view>
</template>

<script>
	import {
		ApiDocUpdate,docDetail
	} from "@/request/api.js"

	export default {
		data() {
			return {
				docData: {
					id: 0,
					createdAt: "0000-00-00 00:00:00",
					updatedAt: "0000-00-00 00:00:00",
					deletedAt: null,
					userId: 0,
					title: "没有标题",
					content: "没有内容",
					docStatus: 0,
					groupId: 0,
					viewCount: 0,
					likeCount: 0,
					isTop: 0,
					priority: 0
				}
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
				this.docData.content = e.detail.text
			},
			undo() {
				this.editorCtx.undo()
			},
			save() {
				console.log(this.docData)
				this.editorCtx.getContents({
					success: (data) => {
						ApiDocUpdate(this.docData).then((res) => {
							uni.showModal({
								content: "更新成功",
								showCancel: false
							})
						}).catch(res => {
							console.log(res)
						})

					}
				})
			}
		},
		onLoad: function(option) {
			this.docData.id = Number(option.id)
			docDetail({
				id: this.docData.id
			}).then(res => {
				console.log(res)
				this.docData = res
				this.editorCtx.setContents({
					html: this.docData.content
				})
			}).catch(err => {
				uni.showToast({
					icon: "none",
					title: String(err)
				})

				uni.navigateBack()
			})
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