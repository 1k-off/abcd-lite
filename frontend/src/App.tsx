import { useState, useEffect } from "react"
import { Header } from "@/components/Header"
import { ProjectCard } from "@/components/ProjectCard"
import { ProjectForm } from "@/components/ProjectForm"
import { Dialog } from "@/components/Dialog"
import { Project, APIKey } from "@/types/project"
import { getProjects, createProject, updateProject, deleteProject } from "@/services/api"
import { Button } from "@/components/ui/button"
import { AlertCircle, Loader2 } from "lucide-react"
import { ThemeProvider } from "@/context/ThemeContext"
import { AuthProvider, useAuth } from "@/context/AuthContext"
import { AuthForm } from "@/components/AuthForm"

function AppContent() {
  const { isAuthenticated, logout, loading } = useAuth()
  const [projects, setProjects] = useState<Project[]>([])
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (isAuthenticated) {
      fetchProjects()
    }
  }, [isAuthenticated])

  const fetchProjects = async () => {
    try {
      setError(null)
      const data = await getProjects()
      setProjects(data)
    } catch (err) {
      await logout()
    }
  }

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
    } catch (err) {
      setError("Failed to delete project. Please try again later.")
      console.error("Error deleting project:", err)
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

  const handleUpdateProject = async (id: string, data: Omit<Project, "id" | "createdAt" | "updatedAt">) => {
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
    } catch (err) {
      setError("Failed to update project. Please try again later.")
      console.error("Error updating project:", err)
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
    } catch (err) {
      setError("Failed to save project. Please try again later.")
      console.error("Error saving project:", err)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!isAuthenticated) {
    return <AuthForm />
  }

  return (
    <div className="min-h-screen bg-background">
      <Header onAddProject={handleAddProject} />
      <button onClick={logout} className="absolute top-4 right-4 z-50 bg-secondary px-4 py-2 rounded shadow">Logout</button>
      <main className="container mx-auto py-8 px-4">
        {error && (
          <div className="mb-6 p-4 bg-destructive/10 border border-destructive rounded-lg flex items-center gap-2 text-destructive">
            <AlertCircle className="h-5 w-5" />
            <p>{error}</p>
            <Button
              variant="ghost"
              size="sm"
              className="ml-auto"
              onClick={() => setError(null)}
            >
              Dismiss
            </Button>
          </div>
        )}

        {projects.length === 0 ? (
          <div className="text-center py-12">
            <h2 className="text-2xl font-semibold mb-2">No projects yet</h2>
            <p className="text-muted-foreground mb-4">
              Get started by creating your first project
            </p>
            <Button onClick={handleAddProject}>Create Project</Button>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {(projects || []).map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                onEdit={handleEditProject}
                onDelete={handleDeleteProject}
              />
            ))}
          </div>
        )}
      </main>

      <Dialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        title={editingProject ? "Edit Project" : "Create Project"}
        description={
          editingProject
            ? "Update your project details"
            : "Add a new project to your list"
        }
      >
        <ProjectForm
          project={editingProject ?? undefined}
          onSubmit={handleSubmit}
          onCancel={() => setIsDialogOpen(false)}
          onApiKeyCreated={handleApiKeyCreated}
        />
      </Dialog>
    </div>
  )
}

export default function AppWithProviders() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </ThemeProvider>
  )
}
