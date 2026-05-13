import { getPublicConfig } from "@nexus/config";
import type { HealthStatus } from "@nexus/types";

const status: HealthStatus = {
  service: "web",
  status: "ok"
};

export default function HomePage() {
  const config = getPublicConfig();

  return (
    <main className="flex min-h-screen items-center justify-center px-6 py-16">
      <section className="w-full max-w-4xl rounded-3xl border border-white/10 bg-white/[0.04] p-10 shadow-2xl shadow-violet-950/30 backdrop-blur">
        <p className="mb-4 text-sm font-medium uppercase tracking-[0.3em] text-violet-300">
          {config.appName}
        </p>
        <h1 className="max-w-3xl text-4xl font-semibold tracking-tight text-white md:text-6xl">
          Operational intelligence for AI-native engineering teams.
        </h1>
        <p className="mt-6 max-w-2xl text-lg leading-8 text-zinc-300">
          Nexus is initialized with a production-grade engineering foundation: monorepo tooling,
          service boundaries, local infrastructure, documentation workflows, and CI-ready checks.
        </p>
        <dl className="mt-10 grid gap-4 text-sm text-zinc-300 sm:grid-cols-3">
          <div className="rounded-2xl border border-white/10 bg-black/20 p-4">
            <dt className="text-zinc-500">Workspace</dt>
            <dd className="mt-2 font-medium text-white">Modular monolith</dd>
          </div>
          <div className="rounded-2xl border border-white/10 bg-black/20 p-4">
            <dt className="text-zinc-500">Health</dt>
            <dd className="mt-2 font-medium text-emerald-300">{status.status}</dd>
          </div>
          <div className="rounded-2xl border border-white/10 bg-black/20 p-4">
            <dt className="text-zinc-500">Service</dt>
            <dd className="mt-2 font-medium text-white">{status.service}</dd>
          </div>
        </dl>
      </section>
    </main>
  );
}

