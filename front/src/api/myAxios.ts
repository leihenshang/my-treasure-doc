import axios from 'axios'

const defaultOptions = {
    // request interceptor handler
    reqHandleFunc: config => config,
    reqErrorFunc: error => Promise.reject(error),
    // response interceptor handler
    resHandleFunc: response => response,
    resErrorFunc: error => Promise.reject(error)
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
                params: data
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
                data: data
            }
        }
        return myAxios(axiosOpt)
    }
}

export {
    myAxios,
    myHttp
}