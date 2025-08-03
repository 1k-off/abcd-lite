import { useState, useEffect, useCallback } from "react"
import { Project, APIKey } from "@/types/project"
import {
  getProjects,
  createProject,
  updateProject,
  deleteProject,
} from "@/services/api"

interface UseProjectsOptions {
  isAuthenticated: boolean
}

export function useProjects({ isAuthenticated }: UseProjectsOptions) {
  const [projects, setProjects] = useState<Project[]>([])
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  const fetchProjects = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getProjects()
      setProjects(data)
    } catch (err: any) {
      const errorMessage = err?.message || err?.error || err?.toString() || "Failed to fetch projects."
      setError(errorMessage)
      // Do not call logout here
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (isAuthenticated) {
      fetchProjects()
    } else {
      setLoading(false)
    }
  }, [fetchProjects, isAuthenticated])

  const handleAddProject = () => {
    setEditingProject(null)
    setIsDialogOpen(true)
  }

  const handleEditProject = (project: Project) => {
    setEditingProject(project)
    setIsDialogOpen(true)
  }

  const handleDeleteProject = async (id: string) => {
    try {
      setError(null)
      await deleteProject(id)
      setProjects((prev) => prev.filter((project) => project.id !== id))
    } catch (err: any) {
      const errorMessage = err?.message || err?.error || err?.toString() || "Failed to delete project. Please try again later."
      setError(errorMessage)
      // Optionally log error
    }
  }

  const handleApiKeyCreated = (projectId: string, apiKeyMeta: APIKey) => {
    setProjects((prev) =>
      prev.map((proj) =>
        proj.id === projectId
          ? {
              ...proj,
              apiKeys: [...(proj.apiKeys || []), apiKeyMeta],
              updatedAt: new Date().toISOString(),
            }
          : proj
      )
    )
    if (editingProject && editingProject.id === projectId) {
      setEditingProject({
        ...editingProject,
        apiKeys: [...(editingProject.apiKeys || []), apiKeyMeta],
        updatedAt: new Date().toISOString(),
      })
    }
  }

  const handleUpdateProject = async (
    id: string,
    data: Omit<Project, "id" | "createdAt" | "updatedAt">
  ) => {
    try {
      const updated = await updateProject(id, data)
      setProjects((prev) =>
        prev.map((proj) =>
          proj.id === id
            ? {
                ...proj,
                ...updated,
                createdAt: updated.createdAt || proj.createdAt || new Date().toISOString(),
                updatedAt: updated.updatedAt || new Date().toISOString(),
              }
            : proj
        )
      )
      setEditingProject(null)
      setIsDialogOpen(false)
    } catch (err: any) {
      const errorMessage = err?.message || err?.error || err?.toString() || "Failed to update project. Please try again later."
      setError(errorMessage)
    }
  }

  const handleSubmit = async (data: Omit<Project, "id" | "createdAt" | "updatedAt">) => {
    try {
      setError(null)
      if (editingProject) {
        await handleUpdateProject(editingProject.id, data)
      } else {
        const newProject = await createProject(data)
        if (newProject && newProject.id) {
          setProjects((prev) => [...prev, newProject])
          setIsDialogOpen(false)
          setEditingProject(null)
        } else {
          throw new Error("Invalid project response from server")
        }
      }
    } catch (err: any) {
      const errorMessage = err?.message || err?.error || err?.toString() || "Failed to save project. Please try again later."
      setError(errorMessage)
    }
  }

  return {
    projects,
    loading,
    error,
    setError,
    handleAddProject,
    handleEditProject,
    handleDeleteProject,
    handleApiKeyCreated,
    handleUpdateProject,
    handleSubmit,
    isDialogOpen,
    setIsDialogOpen,
    editingProject,
    setEditingProject,
  }
} 