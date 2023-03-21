import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export default function LoginPage() {
  return (
    <div className="min-w-full min-h-[100vh] flex flex-col items-center justify-center p-6">
      <form
        onSubmit={handleLogin}
        className="min-w-full flex flex-col gap-5 sm:min-w-[320px]"
      >
        <h3 className="text-2xl font-extrabold leading-tight tracking-tight md:text-3xl">
          Access to your Togo
        </h3>
        <div className="flex flex-col gap-4">
          <div className="grid w-full items-center gap-1.5">
            <Label htmlFor="username">Username</Label>
            <Input type="text" id="username" placeholder="Username" />
          </div>
          <div className="grid w-full items-center gap-1.5">
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
    method: "POST", // *GET, POST, PUT, DELETE, etc.
    cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
    headers: {
      "Content-Type": "application/json",
    },
    redirect: "follow", // manual, *follow, error
    referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: JSON.stringify({ username, password }), // body data type must match "Content-Type" header
  })
  console.log(await response.json())
}
