const host = "http://127.0.0.1:2021"

async function getDocList() {
	let finallRes = null;
	 uni.request({
		url: host + "/doc/list",
		data: {

		},
		method: "POST",
		header: {
			"X-Token": "3b5b3d702a9637860ac351550859cd19"
		},
		success: (res) => {
			// console.log(res.data);
			if (!res.data) {
				finallRes = 'data is empty'
			}
			if (res.data.code > 0) {
				finallRes = res.data.msg;
			}

			finallRes = res.data.data;
			console.log(finallRes)
		}
	})

	return finallRes
}

export {
	getDocList
}