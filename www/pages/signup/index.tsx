import Link from "next/link"
import { useRouter } from "next/router"

import { Layout } from "@/components/layout"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export default function SignUpPage() {
  const router = useRouter()
  return (
    <Layout>
      <div className="min-w-full min-h-[100vh] flex flex-col items-center justify-center px-6 py-16">
        <form
          onSubmit={(e) => handleSignup(e, router)}
          className="min-w-full flex flex-col gap-5 sm:min-w-[320px] -mt-[20vh]"
        >
          <h3 className="text-xl font-extrabold leading-tight tracking-tight md:text-2xl">
            Create your account
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
            <div className="grid w-full items-center gap-2">
              <Label htmlFor="confirm">Confirm Password</Label>
              <Input
                type="password"
                id="confirm"
                placeholder="Confirm Password"
              />
            </div>
          </div>
          <Button type="submit">Create my account</Button>
        </form>
        <p className="mt-4 text-sm">
          Already have an account?{" "}
          <Link className="font-bold" href="/login">
            Log in.
          </Link>
        </p>
      </div>
    </Layout>
  )
}

async function handleSignup(e: React.FormEvent<HTMLFormElement>, router) {
  e.preventDefault()
  const username = e.currentTarget.username.value
  const password = e.currentTarget.password.value
  const confirm = e.currentTarget.confirm.value
  if (password !== confirm) {
    console.log("not match")
    return
  }
  const response = await fetch("http://localhost:3000/signup", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  })
  if (response.status === 201) {
    router.push("/todos")
  }
}
