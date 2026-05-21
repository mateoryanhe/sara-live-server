import {request} from '../request'

export const diamondApi = {
    add: (data: { userId: string | number; amount: number }) => {
        return request.post<{ diamond: number }>('/diamond/add', data)
    },

    sub: (data: { userId: string | number; amount: number }) => {
        return request.post<{ diamond: number }>('/diamond/sub', data)
    },
}
