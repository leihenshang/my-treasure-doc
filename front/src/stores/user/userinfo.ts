import {defineStore} from 'pinia'

export const useUserinfoStore = defineStore('userinfo',{
    state:() => {
        return {
           userinfo: {
            id : 0,
            createdAt : "",
            updatedAt : "",
            deletedAt : null,
            nickname : "",
            account : "",
            email : "",
            password : "",
            userType : 0,
            userStatus : 0,
            mobile : "",
            avatar : "",
            bio : "",
            token : "",
            tokenExpire : "",
            lastLoginIp : "",
            lastLoginTime : ""}
        }
    },
    actions: {
        updateUserinfo(newUserinfo) {
            this.userinfo = newUserinfo
        }
    },
    getters: {
        userId:(state) => state.userinfo.id
    }
})