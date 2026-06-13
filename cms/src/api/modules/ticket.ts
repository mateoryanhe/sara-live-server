import {request} from '../request'
import type {PageResponse, Ticket, TicketQuery} from '@/types/api'

export const ticketApi = {
    getTicketList: (params: TicketQuery) => {
        return request.post<PageResponse<Ticket>>('/ticket/ticketList', params)
    },

    createTicket: (data: {
        price: number
        sort: number
    }) => {
        return request.post<{ id: string }>('/ticket/createTicket', data)
    },

    updateTicket: (data: {
        id: string | number
        price: number
        sort: number
    }) => {
        return request.post<boolean>('/ticket/updateTicket', data)
    },

    deleteTicket: (id: string | number) => {
        return request.post<boolean>('/ticket/deleteTicket', {id})
    },

    onShelfTicket: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/ticket/onShelfTicket', {id})
    },

    offShelfTicket: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/ticket/offShelfTicket', {id})
    },
}
