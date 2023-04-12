import { FormEvent } from "react"
import { useRouter } from "next/router"
import AppConfig from "@/common/config"
import { fetcher } from "@/common/fetcher"
import useSWR from "swr"

import { Layout } from "@/components/layout"

export default function AdminPage() {
  const router = useRouter()
  const {
    data: response,
    error,
    isLoading,
    mutate,
  } = useSWR(
    `${AppConfig.API_URL}/admin/users`,
    (url: string) => fetcher(url, router, true),
    {
      revalidateIfStale: false,
      revalidateOnFocus: false,
    }
  )

  return (
    <Layout>
      <div className="min-w-full max-w-[100vw] min-h-[100vh] flex flex-col items-center px-6 py-16">
        <div className="w-full sm:w-[400px] flex flex-col gap-8">
          {JSON.stringify(response)}
        </div>
      </div>
    </Layout>
  )
}
