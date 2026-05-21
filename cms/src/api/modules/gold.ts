import {request} from '../request'

export const goldApi = {
    add: (data: { userId: string | number; amount: number }) => {
        return request.post<{ gold: number }>('/gold/add', data)
    },

    sub: (data: { userId: string | number; amount: number }) => {
        return request.post<{ gold: number }>('/gold/sub', data)
    },
}
