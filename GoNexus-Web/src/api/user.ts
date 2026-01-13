import request from '../utils/request'

export const getUserInfo = () => {
    return request({
        url:'user/info',
        method:'get'
    })
}

export const updateAvatar = (data: any) => {
    return request({
        url:'user/avatar',
        method:'post',
        data
    })
}

export const updateUserInfo = (data: any) => {
    return request({
        url:'user/info',
        method:'post',
        data
    })
}