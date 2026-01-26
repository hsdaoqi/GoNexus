import request from '../utils/request'

// 语义搜索
export const semanticSearch = (params: { query: string, target_id: number, chat_type: number }) => {
    return request({
        url: '/ai/search',
        method: 'get',
        params
    })
}

// 聊天总结
export const getChatSummary = (params: { target_id: number, chat_type: number }) => {
    return request({
        url: '/ai/summary',
        method: 'get',
        params
    })
}

// 回复建议
export const getReplySuggestions = (params: { target_id: number, chat_type: number }) => {
    return request({
        url: '/ai/suggest',
        method: 'get',
        params
    })
}
