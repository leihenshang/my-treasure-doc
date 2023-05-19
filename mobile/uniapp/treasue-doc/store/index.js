// 页面路径：store/index.js
import {
	createStore
} from 'vuex'
const store = createStore({
	state: { //存放状态
		userInfo: {}
	},
	mutations: {
		setUserInfo(state, userInfo) {
			state.userInfo = userInfo
		},
	}
})

export default store