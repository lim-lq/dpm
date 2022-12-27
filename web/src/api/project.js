import request from '@/utils/request'

export function getProjectList (postData) {
  return request({
    url: '/projects/search',
    method: 'post',
    data: postData
  })
}
