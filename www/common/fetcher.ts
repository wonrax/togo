import { NextRouter } from "next/router"

async function fetcher(
  url: string,
  router: NextRouter,
  redirectIfUnauthorize = false
) {
  const res = await fetch(url, {
    method: "GET",
    credentials: "include",
  })
  if (res.status === 401 && redirectIfUnauthorize && router) {
    router.push("/login")
  }
  if (!res.ok) {
    throw {
      status: res.status,
      message: res.statusText,
    }
  }
  return await res.json()
}

export { fetcher }
