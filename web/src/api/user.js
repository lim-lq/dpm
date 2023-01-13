import request from '@/utils/request'

export function getUserList (postData) {
  return request({
    url: '/accounts/search',
    method: 'post',
    data: postData
  })
}

export function createUserList (postData) {
  return request({
    url: '/accounts',
    method: 'post',
    data: postData
  })
}
