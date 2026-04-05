interface BlogPost {
  id: number
  slug: string
  title: string
  category: string
  tags: string[]
  date: string
  excerpt: string
  views: number
  likes: number
  thumbnail?: string
  is_pinned: boolean
  pinned_order: number
  user: ListUser
}

interface ApiPost {
  id: number
  slug: string
  title: string
  category: string
  tags: string[]
  date?: string
  published_at?: string
  excerpt?: string
  views?: number
  likes?: number
  thumbnail?: string
  is_pinned?: boolean
  pinned_order?: number
  user?: ListUser
}

interface ListUser {
  avatar?: string
  is_admin?: boolean
  nickname?: string
  username?: string
}

export type {BlogPost, ApiPost, ListUser}
