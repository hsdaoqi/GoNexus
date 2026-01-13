import request from '../utils/request'

// 获取好友列表
export const getFriendList = () => {
  return request({
    url: '/friend/list',
    method: 'get'
  })
}

// 发送好友申请
export const addFriend = (data: { receiver_id: number, verify_msg: string }) => {
  return request({
    url: '/friend/request',
    method: 'post',
    data
  })
}

// 处理好友申请
export const handleRequest = (data: { request_id: number, action: number }) => {
  return request({
    url: '/friend/process',
    method: 'post',
    data
  })
}

// 获取待处理的好友请求
export const getPendingRequests = () => {
  return request({
    url: '/friend/pending',
    method: 'get'
  })
}

// 删除好友
export const deleteFriend = (data: { friend_id: number }) => {
  return request({
    url: '/friend/delete', // 对应后端的路由
    method: 'post',
    data
  })
}