import { Project, APIKey } from "@/types/project"

// In development, use the full URL. In production, use relative paths
const API_BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:8900/api" : "/api"

interface ApiError {
  error?: string
  message?: string
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ message: "Unknown error occurred" })) as ApiError
    throw new Error(errorData.message || errorData.error || `HTTP error! status: ${response.status}`)
  }

  return response.json()
}

export async function getProjects(): Promise<Project[]> {
  const response = await fetch(`${API_BASE_URL}/projects`, { credentials: "include" })
  const data = await handleResponse<{ projects: Project[] }>(response)
  return data.projects
}

export async function createProject(data: Omit<Project, "id" | "createdAt" | "updatedAt">): Promise<Project> {
  const response = await fetch(`${API_BASE_URL}/projects`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
    credentials: "include",
  })
  const result = await handleResponse<{ project: Project }>(response)
  
  if (!result?.project) {
    throw new Error("Invalid project response from server")
  }

  const project = result.project
  if (!project.id || !project.createdAt || !project.updatedAt) {
    throw new Error("Project response missing required fields")
  }

  return project
}

export async function updateProject(id: string, data: Omit<Project, "id" | "createdAt" | "updatedAt">): Promise<Project> {
  const response = await fetch(`${API_BASE_URL}/projects/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
    credentials: "include",
  })
  const result = await handleResponse<{ project: Project }>(response)
  return result.project
}

export async function deleteProject(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/projects/${id}`, {
    method: "DELETE",
    credentials: "include",
  })
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ message: "Unknown error occurred" })) as ApiError
    throw new Error(errorData.message || errorData.error || `Failed to delete project: ${response.status}`)
  }
}

export async function generateApiKey(projectId: string): Promise<{ apiKey: string; apiKeyMeta: APIKey }> {
  const response = await fetch(`${API_BASE_URL}/projects/${projectId}/api-keys`, {
    method: "POST",
    credentials: "include",
  })
  return handleResponse<{ apiKey: string; apiKeyMeta: APIKey }>(response)
}

export async function deleteApiKey(projectId: string, keyId: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/projects/${projectId}/api-keys/${keyId}`, {
    method: "DELETE",
    credentials: "include",
  })
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ message: "Failed to delete API key" })) as ApiError
    throw new Error(errorData.message || errorData.error || "Failed to delete API key")
  }
} 