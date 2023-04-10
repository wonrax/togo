import Head from "next/head"
import Link from "next/link"

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
          <h3 className="text-2xl font-extrabold leading-tight tracking-tighter sm:text-2xl md:text-3xl lg:text-4xl">
            A bleeding edge todo manager
          </h3>
          <p className="max-w-[700px] text-lg text-slate-700 dark:text-slate-400 sm:text-xl">
            Add and manage your todo efficiently.
          </p>
        </div>
        <div className="flex gap-4">
          <Link href="/signup" className={buttonVariants({ size: "lg" })}>
            Sign up
          </Link>
          <Link
            href="/login"
            className={buttonVariants({ variant: "outline", size: "lg" })}
          >
            Log in
          </Link>
        </div>
      </section>
    </Layout>
  )
}
