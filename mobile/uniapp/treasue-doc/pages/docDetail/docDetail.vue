<template>
	<view>
		<h1>{{docData.title}}</h1>
		<p>最后更新时间:{{docData.createdAt}}</p>
		<p>{{docData.content}}</p>
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
					content: "没有内容没有内容没有内容没有内容没有内容\n没有内容没有内容没有内容没有内容没有内容没有内容没有内容没有内容没有内容没有内容没有内容没有内容",
					docStatus: 0,
					groupId: 0,
					viewCount: 0,
					likeCount: 0,
					isTop: 0,
					priority: 0
				}
			};
		},
		onLoad: function(option) { //option为object类型，会序列化上个页面传递的参数
			console.log(option.id); //打印出上个页面传递的参数。
			console.log(option.name); //打印出上个页面传递的参数。
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

</style>