import { AuthPanel } from "@/components/auth-panel";

export default function AuthPage() {
  return (
    <main className="grid min-h-screen items-center gap-10 px-6 py-12 lg:grid-cols-[1fr_520px] lg:px-16">
      <section className="max-w-3xl">
        <p className="mb-4 text-sm font-medium uppercase tracking-[0.3em] text-violet-300">Nexus identity</p>
        <h1 className="text-4xl font-semibold tracking-tight text-white md:text-6xl">Create the first organization workspace.</h1>
        <p className="mt-6 text-lg leading-8 text-zinc-300">
          This slice introduces the first real Nexus product boundary: users, organizations, workspaces, sessions, and PostgreSQL-backed API persistence.
        </p>
      </section>
      <AuthPanel />
    </main>
  );
}

