import request from '@/utils/request'
// 创建群组
export const createGroup = (data: { name: string, avatar?: string, notice?: string }) => {
    return request({
        url: '/group/create',
        method: 'post',
        data
    })
}

// 获取我加入的群组
export const getMyGroups = () => {
    return request({
        url: '/group/mine',
        method: 'get'
    })
}

// 更新群资料（群主）
export const updateGroup = (data: { id: number, name?: string, avatar?: string, notice?: string }) => {
    return request({
        url: '/group/update',
        method: 'post',
        data
    })
}

// 获取群成员
export const getGroupMembers = (groupId: number) => {
    return request({
        url: '/group/members',
        method: 'get',
        params: { group_id: groupId }
    })
}

// 邀请好友入群
export const inviteMember = (data: { group_id: number, friend_id: number }) => {
    return request({
        url: '/group/invite',
        method: 'post',
        data
    })
}

// 踢人
export const kickMember = (data: { group_id: number, member_id: number }) => {
    return request({
        url: '/group/kick',
        method: 'post',
        data
    })
}

// 禁言
export const muteMember = (data: { group_id: number, member_id: number, mute: number }) => {
    return request({
        url: '/group/mute',
        method: 'post',
        data
    })
}

// 设置/取消管理员
export const setAdmin = (data: { group_id: number, member_id: number, is_admin: boolean }) => {
    return request({
        url: '/group/admin',
        method: 'post',
        data
    })
}

// 转让群主
export const transferGroup = (data: { group_id: number, member_id: number }) => {
    return request({
        url: '/group/transfer',
        method: 'post',
        data
    })
}
