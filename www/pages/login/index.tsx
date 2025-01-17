import Link from "next/link"
import { useRouter } from "next/router"
import AppConfig from "@/common/config"
import { useSWRConfig } from "swr"

import { Layout } from "@/components/layout"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export default function LoginPage() {
  const router = useRouter()
  const { mutate } = useSWRConfig()
  return (
    <Layout>
      <div className="min-w-full min-h-[100vh] flex flex-col items-center justify-center px-6 py-16">
        <form
          onSubmit={(e) => handleLogin(e, router, mutate)}
          className="min-w-full flex flex-col gap-5 sm:min-w-[320px] -mt-[20vh]"
        >
          <h3 className="text-xl font-extrabold leading-tight tracking-tight md:text-2xl">
            Access to your Togo
          </h3>
          <div className="flex flex-col gap-4">
            <div className="grid w-full items-center gap-2">
              <Label htmlFor="username">Username</Label>
              <Input
                type="text"
                id="username"
                placeholder="Username"
                autoFocus
              />
            </div>
            <div className="grid w-full items-center gap-2">
              <Label htmlFor="password">Password</Label>
              <Input type="password" id="password" placeholder="Password" />
            </div>
          </div>
          <Button type="submit">Login</Button>
        </form>
        <p className="mt-4 text-sm">
          Don`t have an account?{" "}
          <Link className="font-bold" href="/signup">
            Sign up.
          </Link>
        </p>
      </div>
    </Layout>
  )
}

async function handleLogin(
  e: React.FormEvent<HTMLFormElement>,
  router,
  mutate
) {
  e.preventDefault()
  const username = e.currentTarget.username.value
  const password = e.currentTarget.password.value
  const response = await fetch(`${AppConfig.API_URL}/login`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  })
  if (response.status === 200) {
    router.push("/todos")
    mutate(`${AppConfig.API_URL}/me`) // refresh current user
  }
}
