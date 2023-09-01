import axios from 'axios'
import {useUserinfoStore} from '@/stores/user/userinfo'

const userStore = useUserinfoStore()

interface DefaultOptions {
    reqHandleFunc: Function,
    reqErrorFunc: Function,
    resHandleFunc: Function,
    resErrorFunc: Function
}

const defaultOptions: DefaultOptions = {
    // request interceptor handler
    reqHandleFunc: (config: any) => config,
    reqErrorFunc: (error: any) => Promise.reject(error),
    // response interceptor handler
    resHandleFunc: (response: any) => response,
    resErrorFunc: (error: any) => Promise.reject(error)
}

const options = {}

const initOptions: Object = {
    ...defaultOptions,
    ...options
}

const myAxios = axios.create(initOptions)

// Add a request interceptor
myAxios.interceptors.request.use(
    config => initOptions.reqHandleFunc(config),
    error => initOptions.reqErrorFunc(error)
)
// Add a response interceptor
myAxios.interceptors.response.use(
    response => initOptions.resHandleFunc(response),
    error => initOptions.resErrorFunc(error)
)

interface MyHttp {
    get: Function,
    post: Function
}

const myHttp: MyHttp = {
    get: (url: string, data: any, options: any) => {
        let axiosOpt = {
            ...options,
            ...{
                method: 'get',
                url: url,
                params: data,
                headers:{
                    'X-Token': userStore.userinfo.token
                }
            }
        }
        return myAxios(axiosOpt)
    },
    post: (url: string, data: any, options: any) => {
        let axiosOpt = {
            ...options,
            ...{
                method: 'post',
                url: url,
                data: data,
                headers:{
                    'X-Token': userStore.userinfo.token
                }
            }
        }
        return myAxios(axiosOpt)
    }
}

export {
    myAxios,
    myHttp
}