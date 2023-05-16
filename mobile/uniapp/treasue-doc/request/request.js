//把配置项单独处理
import store from '@/store/index.js'; //vuex  
let server_url = ' '; //请求地址
let token = ' '; // 凭证
server_url = process.env.NODE_ENV === 'development' ? '' : 'http://***/api'; //环境配置
token = store.state.xToken

function service(options = {}) {
	store.state.token && (token = store.state.token); //从vuex中获取登录凭证
	options.url = `${server_url}${options.url}`;
	//配置请求头
	options.header = {
		'Content-Type': 'application/json',
		'X-Token': `${token}` //Bearer 
	};
	return new Promise((resolved, rejected) => {
		//成功
		options.success = (res) => {
			if (Number(res.data.code) == 0) { //请求成功
				resolved(res.data.data);
			} else {
				uni.showToast({
					icon: 'none',
					duration: 3000,
					title: `${res.data.msg}`
				});
				rejected(res.data.msg); //错误
			}
		}
		//错误
		options.fail = (err) => {
			rejected(err); //错误
		}
		uni.request(options);
	});
}
export default service;