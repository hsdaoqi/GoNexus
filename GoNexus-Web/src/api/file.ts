import request from '@/utils/request'

export const uploadFile = (formData: FormData) => {
  return request({
    url: '/file/upload',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

