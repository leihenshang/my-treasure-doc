<template>
	<view class="container">
		<ul>
			<li v-for="item in list.list">
				<p class="listItem" @click="detail(item.id)">{{item.id}}-{{item.title}}</p>
				<p>创建时间:{{item.createdAt}}</p>
				<p>{{item.content.substring(0,20)}}...</p>
			</li>
		</ul>
	</view>
</template>

<script>
	import {
		docList
	} from "@/request/api.js"

	let docList1 = {
		"total": 6,
		"list": [{
				"content": "sssssssss",
				"id": 1,
				"title": "hhhhhhhhh"
			},
			{
				"content": "sssssssss",
				"id": 2,
				"title": "hhhhhhhhh"
			}, {
				"content": "sssssssss",
				"id": 3,
				"title": "hhhhhhhhh"
			},
			{
				"content": "sssssssss",
				"id": 4,
				"title": "hhhhhhhhh"
			}, {
				"content": "sssssssss",
				"id": 5,
				"title": "hhhhhhhhh"
			}, {
				"content": "sssssssss",
				"id": 6,
				"title": "hhhhhhhhh"
			}
		]
	};

	export default {
		data() {
			return {
				href: 'https://uniapp.dcloud.io/component/README?id=uniui',
				list: docList1
			}
		},
		methods: {
			async getDocList() {
				await docList({
					"page": 1,
					"pageSize": 10
				}).then(res => {
					console.log(res)
					this.list = res
				}).catch(res => {
					uni.showModal({
						content: res
					})
				})
			},
			detail(id) {
				console.log(id)
				uni.showToast({
					icon: "none",
					title: String(id)
				})
				uni.navigateTo({
					url: "/pages/docDetail/docDetail?id=" + id
				})
			}
		},
		beforeMount() {
			this.getDocList()
		}
	}
</script>
<style lang="scss">
	.container {
		padding: 20px;
		font-size: 14px;
		line-height: 24px;

		ul {
			list-style-type: none;
			margin: 0;
			padding: 0;

			li {
				height: 150rpx;
				width: 100%;
				background-color: pink;
				margin: 10rpx 0;
			}
		}
	}
</style>