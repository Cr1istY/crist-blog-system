export interface User {
  id: string
  username: string
  displayName: string
  avatar: string
  verified: boolean
  bio: string
  email: string
}

export interface Tweet {
  id: string
  user: User
  content: string
  timestamp: Date
  likes: number
  images?: string[]
}

export interface TweetListProps {
  tweets: Tweet[]
  loading?: boolean
  currentUser?: User
}

export interface TweetEmits {
  (event: 'submit', content: string): void
  (event: 'like', tweetId: string): void
  (event: 'retweet', tweetId: string): void
  (event: 'reply', tweetId: string): void
}

export interface TweetImage {
  id: string
  url: string
  thumbnailUrl?: string
  width?: number
  height?: number
}

export interface UploadedImage {
  id: string
  file: File
  previewUrl: string
  uploading: boolean
  progress: number
  error?: string
  uploadedUrl?: string // 上传成功后的服务器 URL
}

export interface TweetComposerProps {
  placeholder?: string
  submitLabel?: string
  maxLength?: number
  maxImages?: number
  allowedImageTypes?: string[]
  maxImageSize?: number
  currentUser?: User
}
