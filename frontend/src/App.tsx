import { Header } from "@/components/Header"
import { ProjectCard } from "@/components/ProjectCard"
import { ProjectForm } from "@/components/ProjectForm"
import { Dialog } from "@/components/Dialog"
import { Button } from "@/components/ui/button"
import { Loader2 } from "lucide-react"
import { ThemeProvider } from "@/context/ThemeContext"
import { AuthProvider, useAuth } from "@/context/AuthContext"
import { AuthForm } from "@/components/AuthForm"
import { useProjects } from "@/hooks/useProjects"
import { ErrorAlert } from "@/components/ErrorAlert"

function AppContent() {
  const { isAuthenticated, logout, loading: authLoading } = useAuth()
  const {
    projects,
    loading,
    error,
    setError,
    handleAddProject,
    handleEditProject,
    handleDeleteProject,
    handleApiKeyCreated,
    handleSubmit,
    isDialogOpen,
    setIsDialogOpen,
    editingProject,
  } = useProjects({ isAuthenticated })

  if (authLoading || loading) {
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
          <ErrorAlert error={error} onDismiss={() => setError(null)} />
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
