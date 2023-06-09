import axios from 'axios'

let VueAxiosPlugin:any = {}

VueAxiosPlugin.install = (app, options) => {
  const defaultOptions = {
    // request interceptor handler
    reqHandleFunc: config => config,
    reqErrorFunc: error => Promise.reject(error),
    // response interceptor handler
    resHandleFunc: response => response,
    resErrorFunc: error => Promise.reject(error)
  }

  const initOptions = {
    ...defaultOptions,
    ...options
  }

  const service = axios.create(initOptions)

  // Add a request interceptor
  service.interceptors.request.use(
    config => initOptions.reqHandleFunc(config),
    error => initOptions.reqErrorFunc(error)
  )
  // Add a response interceptor
  service.interceptors.response.use(
    response => initOptions.resHandleFunc(response),
    error => initOptions.resErrorFunc(error)
  )

  app.config.globalProperties.$axios = service
  app.config.globalProperties.$http = {
    get: (url, data, options) => {
      let axiosOpt = {
        ...options,
        ...{
          method: 'get',
          url: url,
          params: data
        }
      }
      return service(axiosOpt)
    },
    post: (url, data, options) => {
      let axiosOpt = {
        ...options,
        ...{
          method: 'post',
          url: url,
          data: data
        }
      }
      return service(axiosOpt)
    }
  }
}

// Auto-install
let GlobalVue = null
if (typeof window !== 'undefined') {
  GlobalVue = window.Vue
} else if (typeof global !== 'undefined') {
  GlobalVue = global.Vue
}
if (GlobalVue) {
  GlobalVue.use(VueAxiosPlugin)
}

export default VueAxiosPlugin
