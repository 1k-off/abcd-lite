import { useState, useEffect } from "react"
import { Header } from "@/components/Header"
import { ProjectCard } from "@/components/ProjectCard"
import { ProjectForm } from "@/components/ProjectForm"
import { Dialog } from "@/components/Dialog"
import { Project } from "@/types/project"
import { getProjects, createProject, updateProject, deleteProject } from "@/services/api"
import { Button } from "@/components/ui/button"
import { AlertCircle, Loader2 } from "lucide-react"
import { ThemeProvider } from "@/context/ThemeContext"

function App() {
  const [projects, setProjects] = useState<Project[]>([])
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | undefined>()
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchProjects()
  }, [])

  const fetchProjects = async () => {
    try {
      setIsLoading(true)
      setError(null)
      const data = await getProjects()
      setProjects(data)
    } catch (err) {
      setError("Failed to load projects. Please try again later.")
      console.error("Error fetching projects:", err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleAddProject = () => {
    setEditingProject(undefined)
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

  const handleSubmit = async (data: Omit<Project, "id" | "createdAt" | "updatedAt">) => {
    try {
      setError(null)
      if (editingProject) {
        const updatedProject = await updateProject(editingProject.id, data)
        setProjects((prev) =>
          prev.map((project) =>
            project.id === editingProject.id ? updatedProject : project
          )
        )
      } else {
        const newProject = await createProject(data)
        if (newProject && newProject.id) {
          setProjects((prev) => [...prev, newProject])
        } else {
          throw new Error("Invalid project response from server")
        }
      }
      setIsDialogOpen(false)
    } catch (err) {
      setError("Failed to save project. Please try again later.")
      console.error("Error saving project:", err)
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <Header onAddProject={handleAddProject} />
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
          project={editingProject}
          onSubmit={handleSubmit}
          onCancel={() => setIsDialogOpen(false)}
        />
      </Dialog>
    </div>
  )
}

export default function AppWithProviders() {
  return (
    <ThemeProvider>
      <App />
    </ThemeProvider>
  )
}
