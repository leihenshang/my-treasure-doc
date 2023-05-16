import request from '@/request/request.js'

export function login(data) {
	return request({
		url: '/user/login',
		method: 'POST',
		data
	})
}

export function docList(data) {
	return request({
		url: '/test/doc/list',
		method: 'POST',
		data
	})
}

export function docCreate(data) {
	return request({
		url: '/test/doc/create',
		method: 'POST',
		data
	})
}