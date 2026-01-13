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