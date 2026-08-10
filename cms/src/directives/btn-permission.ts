import type {Directive, DirectiveBinding} from 'vue'
import {hasButtonPermission, resolveButtonPermission} from '@/utils/permission'

function applyBtnPermission(el: HTMLElement, binding: DirectiveBinding<string>) {
    const resolved = resolveButtonPermission(binding.value, binding.arg)
    if (!resolved) {
        return
    }
    const allowed = hasButtonPermission(resolved.page, resolved.action)
    if (!allowed) {
        el.parentNode?.removeChild(el)
    }
}

export const btnPermission: Directive<HTMLElement, string> = {
    mounted(el, binding) {
        applyBtnPermission(el, binding)
    },
    updated(el, binding) {
        applyBtnPermission(el, binding)
    },
}
