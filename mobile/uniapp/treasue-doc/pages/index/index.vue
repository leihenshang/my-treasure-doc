<template>
	<view class="container">
		<view class="search-box">
			<uni-search-bar @confirm="searchDoc" @input="searchInput" bgColor="#EEEEEE" clearButton="auto"
				placeholder="搜索一下" radius="100" cancelButton="none"></uni-search-bar>
		</view>
		<view class="group-box">
			<view class="group-box-tree">

			</view>
			<view class="group-box-btn">
			</view>
		</view>
		<view class="list-box">
			<uni-list :border="true">
				<uni-list-item :to="'/pages/docDetail/docDetail?id='+item.id" clickable :title="item.title"
					:right-text="item.createdAt" :note="removeHtmlTag(item.content)"
					v-for="item in list.list"></uni-list-item>
			</uni-list>
		</view>
		<view class="bottom-fill">
			<text class="bottom-fill-text">暂时没有更多了...</text>
		</view>
		<uni-fab ref="fab" :pattern="pattern" horizontal="right" @fabClick="fabClick"></uni-fab>
	</view>
</template>

<script>
	import {
		docList
	} from "@/request/api.js"

	export default {
		data() {
			return {
				pattern: {
					color: '#7A7E83',
					backgroundColor: '#fff',
					selectedColor: '#007AFF',
					buttonColor: '#007AFF',
					iconColor: '#fff'
				},
				list: {
					"total": 0,
					"list": []
				},
				status: 'more',
				docPage: {
					page: 1,
					pageSize: 20
				},
				lastPage: 0,
				classes: '1-2',
				dataTree: [{
						text: "一年级",
						value: "1-0",
						children: [{
								text: "1.1班",
								value: "1-1"
							},
							{
								text: "1.2班",
								value: "1-2"
							}
						]
					},
					{
						text: "二年级",
						value: "2-0",
						children: [{
								text: "2.1班",
								value: "2-1"
							},
							{
								text: "2.2班",
								value: "2-2"
							}
						]
					},
					{
						text: "三年级",
						value: "3-0",
						disable: true
					}
				]
			}
		},
		methods: {
			onnodeclick(e) {
				console.log(e);
			},
			onpopupopened(e) {
				console.log('popupopened');
			},
			onpopupclosed(e) {
				console.log('popupclosed');
			},
			onchange(e) {
				console.log('onchange:', e);
			},

			searchDoc() {},
			searchInput() {},
			removeHtmlTag(content) {
				let regex = /(<([^>]+)>)/ig
				let c = content.replace(regex, '')
				return c.substring(0, 20) + '...'
			},
			async getDocList() {
				if (this.lastPage == this.docPage.page) {
					return
				}

				await docList({
					page: this.docPage.page,
					pageSize: this.docPage.pageSize,
					isDesc: true
				}).then(res => {
					this.lastPage = this.docPage.page
					if (Math.ceil(res.total / this.docPage.pageSize) > this.docPage.page) {
						this.docPage.page++
					}
					this.list.list.push(...res.list)
				}).catch(err => {
					console.log(err)
					uni.showToast({
						title: '请求文档列表异常！',
						icon: "none",
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
		.search-box {
			width: 100%;
		}

		.group-box {
			display: flex;
			padding: 2rpx;

			.group-box-tree {
				width: 50%;
			}

			.group-box-btn {
				width: 50%;
			}
		}

		.list-box {}


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