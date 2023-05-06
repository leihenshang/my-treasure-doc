const host = ""

function getDocList() {
	uni.request({
		url: host + "",
		data: {

		},
		// header: {},
		success: (res) => {
			console.log(res.data);
		}
	})
}