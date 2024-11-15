import type { ArtalkPlugin } from '@/types'

export const UnreadBadge: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  let $unreadBadge: HTMLElement | null = null

  const showUnreadBadge = (count: number) => {
    if (!$unreadBadge) return

    if (count > 0) {
      $unreadBadge.innerText = `${Number(count || 0)}`
      $unreadBadge.style.display = 'block'
    } else {
      $unreadBadge.style.display = 'none'
    }
  }

  ctx.on('mounted', () => {
    $unreadBadge = list.getEl().querySelector<HTMLElement>('.atk-unread-badge')
  })

  ctx.on('notifies-updated', (notifies) => {
    showUnreadBadge(notifies.length || 0)
  })
}
