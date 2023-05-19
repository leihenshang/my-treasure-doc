<template>
	<view class="container">
		<uni-list :border="true">
			<uni-list-item :to="'/pages/docDetail/docDetail?id='+item.id" :clickable=true :title="item.title"
				:right-text="item.createdAt" :note="item.content.substring(0,20)+'...'"
				v-for="item in list.list"></uni-list-item>
		</uni-list>
		<view class="bottom-fill">
			<text class="bottom-fill-text">暂时没有更多了...</text>
		</view>
		<uni-fab horizontal="right" @fabClick="fabClick"></uni-fab>
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
				list: {
					"total": 0,
					"list": []
				},
				status: 'more',
				docPage: {
					page: 1,
					pageSize: 20
				},
				lastPage: 0
			}
		},
		methods: {
			async getDocList() {
				if (this.lastPage == this.docPage.page) {
					return
				}

				await docList({
					page: this.docPage.page,
					pageSize: this.docPage.pageSize
				}).then(res => {
					this.lastPage = this.docPage.page
					if (Math.ceil(res.total / this.docPage.pageSize) > this.docPage.page) {
						this.docPage.page++
					}
					this.list.list.push(...res.list)
				}).catch(err => {
					console.log(err)
					uni.showModal({
						content: '出现异常了！'
					})
				})
			},
			detail(id) {
				uni.showToast({
					icon: "none",
					title: String(id)
				})
				uni.navigateTo({
					url: "/pages/docDetail/docDetail?id=" + id
				})
			},
			clickLoadMore(e) {
				console.log(e.detail)
			},
			fabClick() {
				uni.navigateTo({
					url: '/pages/write/write',
					success: res => {},
					fail: () => {},
					complete: () => {}
				});
			}
		},
		beforeMount() {
			this.getDocList()
		},
		onLoad: function(options) {
			setTimeout(function() {
				console.log('start pulldown');
			}, 1000);
			uni.startPullDownRefresh();
		},
		onPullDownRefresh() {
			console.log('refresh');
			setTimeout(function() {
				uni.stopPullDownRefresh();
			}, 1000);
		},
		onReachBottom() {
			uni.showToast({
				icon: "none",
				title: "触底了"
			})
			this.getDocList()

		}
	}
</script>
<style lang="scss">
	.container {
		.bottom-fill {
			height: 180rpx;
			width: 100%;
			text-align: center;
			font-size: 10rpx;
			color: gray;

			.bottom-fill-text {
				display: block;
				margin: 20rpx;
			}
		}
	}
</style>