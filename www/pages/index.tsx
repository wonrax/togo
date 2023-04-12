import Head from "next/head"
import Link from "next/link"
import AppConfig from "@/common/config"
import { fetcher } from "@/common/fetcher"
import { ArrowRight } from "lucide-react"
import useSWR from "swr"

import { Layout } from "@/components/layout"
import { buttonVariants } from "@/components/ui/button"

export default function IndexPage() {
  return (
    <Layout>
      <Head>
        <title>Togo</title>
        <meta name="description" content="A bleeding edge todo manager" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <section className="container grid items-center gap-6 pt-6 pb-8 md:py-10">
        <div className="flex max-w-[980px] flex-col items-start gap-2">
          <h3 className="text-3xl font-extrabold leading-tight tracking-tight sm:text-3xl md:text-4xl lg:text-5xl">
            <span className="text-transparent bg-clip-text bg-gradient-to-br from-cyan-500 to-indigo-500 dark:from-cyan-400 dark:to-indigo-400">
              A bleeding edge
            </span>
            <br />
            todo manager
          </h3>
          <p className="max-w-[700px] text-lg text-slate-700 dark:text-slate-400 sm:text-xl">
            Add and manage your todo efficiently.
          </p>
          <p className="max-w-[700px] text-xs text-slate-300 dark:text-slate-700 sm:text-sm">
            Tiny unimportant note: Your data are not guaranteed to be kept
            forever.&nbsp;
            <br className="hidden sm:block" />
            Do not store critical data here.
          </p>
        </div>
        <CTA />
      </section>
    </Layout>
  )
}

function CTA() {
  const {
    data: response,
    error,
    isLoading,
  } = useSWR(`${AppConfig.API_URL}/me`, fetcher, {
    revalidateIfStale: false,
    revalidateOnFocus: false,
    shouldRetryOnError: false,
  })

  let ctaDestination = "/signup"

  if (!isLoading && !error && response?.data?.username) {
    ctaDestination = "/todos"
  }

  return (
    <Link
      href={ctaDestination}
      className={buttonVariants({ size: "lg", variant: "outline" }) + " w-fit"}
    >
      Free – Forever
      <ArrowRight className="ml-2" size={18} strokeWidth={2} />
    </Link>
  )
}
