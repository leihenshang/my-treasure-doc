import request from '@/request/request.js'

export function login(data) {
	return request({
		url: '/test/user/login',
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

export function docDetail(data) {
	return request({
		url: '/test/doc/detail',
		method: 'POST',
		data
	})
}