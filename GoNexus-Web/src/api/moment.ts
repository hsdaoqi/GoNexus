import request from '@/utils/request'

export interface CreatePostRequest {
  content: string
  media?: string[]
  visibility: number // 0-公开, 1-好友, 2-私密
  location?: string
  mood?: string
}

export interface Comment {
  id: number
  post_id: number
  user_id: number
  user_nickname: string
  user_avatar: string
  content: string
  parent_id?: number
  created_at: string
}

export interface Post {
  id: number
  user_id: number
  user_nickname: string
  user_avatar: string
  content: string
  media: string[]
  visibility: number
  location: string
  mood: string
  like_count: number
  comment_count: number
  is_liked: boolean
  created_at: string
  comments?: Comment[]
}

// 发布动态
export function createPost(data: CreatePostRequest) {
  return request({
    url: '/moment/post',
    method: 'post',
    data
  })
}

// 获取动态列表
export function getMoments(params: { page: number; page_size: number; type: string; user_id?: number }) {
  return request({
    url: '/moment/list',
    method: 'get',
    params
  })
}

// 发表评论
export function createComment(data: { post_id: number; content: string; parent_id?: number }) {
  return request({
    url: '/moment/comment',
    method: 'post',
    data
  })
}

// 获取评论
export function getComments(post_id: number) {
  return request({
    url: '/moment/comments',
    method: 'get',
    params: { post_id }
  })
}

// 点赞/取消点赞
export function toggleLike(post_id: number) {
  return request({
    url: '/moment/like',
    method: 'post',
    data: { post_id }
  })
}
