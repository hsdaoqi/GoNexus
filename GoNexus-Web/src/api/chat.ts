import request from '../utils/request'

// 获取历史消息
export const getChatHistory = (params: any) => {
  return request({
    url: '/chat/history',
    method: 'get',
    params
  })
}

// 询问 AI (RAG)
export const askAI = (query: string,target_id: number,chat_type: number) => {
  return request({
    url: '/ai/search',
    method: 'get',
    params: { query,target_id,chat_type }
  })
}

// 撤回消息
export const revokeMessage = (data: { msg_id: number, chat_type: number, target_id: number }) => {
  return request({
    url: '/chat/revoke',
    method: 'post',
    data
  })
}

// 消息已读上报
export const readMessage = (data: { target_id: number, chat_type: number }) => {
  return request({
    url: '/chat/read',
    method: 'post',
    data
  })
}