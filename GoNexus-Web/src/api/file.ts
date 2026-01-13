import request from '../utils/request'

export const upload = (data: any) => {
    return request({
        url:'/file/upload',
        method:'post',
        data
    })
}