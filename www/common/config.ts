if (!process.env.NEXT_PUBLIC_API_URL) {
  throw new Error("NEXT_PUBLIC_API_URL is not defined")
}

const AppConfig = {
  API_URL: process.env.NEXT_PUBLIC_API_URL,
}

export default AppConfig
