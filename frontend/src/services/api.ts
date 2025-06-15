import { Project } from "@/types/project"

// In development, use the full URL. In production, use relative paths
const API_BASE_URL = import.meta.env.DEV ? "http://localhost:8900/api" : "/api"

interface ApiError {
  error: string
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: "Unknown error occurred" })) as ApiError
    throw new Error(errorData.error || `HTTP error! status: ${response.status}`)
  }
  return response.json()
}

export async function getProjects(): Promise<Project[]> {
  const response = await fetch(`${API_BASE_URL}/projects`)
  return handleResponse<Project[]>(response)
}

export async function createProject(data: Omit<Project, "id" | "createdAt" | "updatedAt">): Promise<Project> {
  const response = await fetch(`${API_BASE_URL}/projects`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
  return handleResponse<Project>(response)
}

export async function updateProject(id: string, data: Omit<Project, "id" | "createdAt" | "updatedAt">): Promise<Project> {
  const response = await fetch(`${API_BASE_URL}/projects/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
  return handleResponse<Project>(response)
}

export async function deleteProject(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/projects/${id}`, {
    method: "DELETE",
  })
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: "Unknown error occurred" })) as ApiError
    throw new Error(errorData.error || `Failed to delete project: ${response.status}`)
  }
} 