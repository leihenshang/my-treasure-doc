<template>
	<view class="c-container">
		<view class="c-title">
			<uni-section :title="docData.title" :sub-title="docData.createdAt" title-font-size="14px"></uni-section>
		</view>
		<hr class="c-hr">
		<view class="c-content">
			<rich-text :nodes="docData.content"></rich-text>
		</view>
		<uni-fab :pattern="{icon:'compose'}" horizontal="right" @fabClick="fabClick"></uni-fab>
	</view>
</template>

<script>
	import {
		docDetail
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
			fabClick() {
				uni.redirectTo({
					url: "/pages/edit/edit?id=" + this.docData.id
				})
			}
		},
		onLoad: function(option) { //option为object类型，会序列化上个页面传递的参数
			docDetail({
				id: Number(option.id)
			}).then(res => {
				console.log(res)
				this.docData = res
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
	.c-container {
		height: 100%;
		width: 100%;

		.c-hr {
			color: darkgrey;
			width: 95%;
			margin: 0 auto;
		}

		.c-content {
			background-color: #fff;
			padding: 20rpx 10rpx;
			height: calc(100% - 40rpx);
			width: calc(100% - 20rpx);
		}
	}
</style>