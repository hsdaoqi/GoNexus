import request from '../utils/request'

// 登录接口
export const login = (data: any) => {
    return request({
        url: '/user/login',
        method: 'post',
        data
    })
}

// 注册接口
export const register = (data: any) => {
    console.log(data)
    return request({
        url: '/user/register',
        method: 'post',
        data
    })
}