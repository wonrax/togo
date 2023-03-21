import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export default function LoginPage() {
  return (
    <div className="min-w-full min-h-[100vh] flex flex-col items-center justify-center px-6 py-16">
      <form
        onSubmit={handleLogin}
        className="min-w-full flex flex-col gap-5 sm:min-w-[320px] -mt-[20vh]"
      >
        <h3 className="text-xl font-extrabold leading-tight tracking-tight md:text-2xl">
          Access to your Togo
        </h3>
        <div className="flex flex-col gap-4">
          <div className="grid w-full items-center gap-2">
            <Label htmlFor="username">Username</Label>
            <Input type="text" id="username" placeholder="Username" />
          </div>
          <div className="grid w-full items-center gap-2">
            <Label htmlFor="password">Password</Label>
            <Input type="password" id="password" placeholder="Password" />
          </div>
        </div>
        <Button type="submit">Login</Button>
      </form>
    </div>
  )
}

async function handleLogin(e: React.FormEvent<HTMLFormElement>) {
  e.preventDefault()
  const username = e.currentTarget.username.value
  const password = e.currentTarget.password.value
  const response = await fetch("http://localhost:3000/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  })
  if (response.status === 200) {
    location.href = "/todos"
  }
}
