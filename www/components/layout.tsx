import { SiteHeader } from "@/components/site-header"

interface LayoutProps {
  children: React.ReactNode
}

export function Layout({ children }: LayoutProps) {
  return (
    <>
      <SiteHeader />
      <main className="min-h-[calc(100vh_+_1px)]">{children}</main>
    </>
  )
}
