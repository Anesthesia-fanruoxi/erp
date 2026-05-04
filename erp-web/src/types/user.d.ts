interface UserInfo {
  id: number
  userName: string
  realName: string
  email: string
  phone: string
  status: number
  createdAt: number
}

interface LoginForm {
  userName: string
  password: string
}

interface RegisterForm {
  userName: string
  password: string
  realName: string
  email?: string
  phone?: string
}

interface MenuItem {
  id: number
  name: string
  path: string
  icon: string
  sort: number
  children?: MenuItem[]
}
