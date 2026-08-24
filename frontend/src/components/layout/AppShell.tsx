import type { ReactNode } from 'react';

type AppShellProps = {
  children: ReactNode;
};

export function AppShell({ children }: AppShellProps) {
  return (
    <main className="app-shell">
      <h1>Calculator</h1>
      <div className="app-shell__content">{children}</div>
    </main>
  );
}
