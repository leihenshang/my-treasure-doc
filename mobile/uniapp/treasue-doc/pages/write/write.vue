<template>
	<view class="write-box">
		<!-- <richTextEditor></richTextEditor> -->
		<uni-forms ref="form" label-position="top" :model="formData">
			<view class="btn-group">
				<button type="primary" @click="save">保存</button>
			</view>
			<uni-forms-item label="标题" name="">
				<uni-easyinput focus placeholder="输入你的标题" class="title" v-model="formData.title"></uni-easyinput>
			</uni-forms-item>
			<uni-forms-item label="内容" name="" class="">
				<richTextEditor @input="editorInput"></richTextEditor>
			</uni-forms-item>
		</uni-forms>
	</view>
</template>

<script>
	import {
		docCreate
	} from "@/request/api.js"
	import richTextEditor from '@/component/editor/editor.vue'

	export default {
		components: {
			richTextEditor
		},
		data() {
			return {
				placeholder: '挥洒你的创意吧...',
				formData: {
					title: "",
					content: ""
				}
			};
		},
		methods: {
			editorInput(e) {
				this.formData.content = e
			},
			save() {
				console.log(this.formData)
				docCreate(this.formData).then((res) => {
					uni.reLaunch({
						url: "/pages/index/index",
					})
				}).catch(res => {
					console.log(res)
				})
			}
		}

	}
</script>

<style lang="scss">
	.write-box {
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