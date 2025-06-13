import type { MaybeRef } from 'vue'
import { router } from '@inertiajs/vue3'
import { unref } from 'vue'

export function useLink(target: MaybeRef<string>) {
  const navigate = (e: Event) => {
    e.preventDefault()
    e.stopPropagation()

    router.visit(unref(target))
  }

  return {
    link: true,
    href: unref(target),
    onClick: navigate,
  }
}
