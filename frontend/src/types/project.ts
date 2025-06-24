export interface Project {
  id: string
  name: string
  iisSites: string[]
  apiKeys: APIKey[]
  createdAt: string
  updatedAt: string
} 

export interface APIKey {
  id: string
  hash: string
  createdAt: string
  prefix: string
  suffix: string
}