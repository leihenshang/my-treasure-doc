import {
	defineConfig
} from 'vite';
import uni from '@dcloudio/vite-plugin-uni';

export default defineConfig({
	plugins: [uni()],
	server: {
		proxy: {
			// 使用 proxy 实例
			'/test': {
				target: 'http://127.0.0.1:2021',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/test/, ''),
			},
		},
	},
});