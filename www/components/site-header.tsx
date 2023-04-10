import { useRouter } from "next/router"
import { fetcher } from "@/common/fetcher"
import { Avatar, AvatarFallback } from "@radix-ui/react-avatar"
import {
  Cloud,
  CreditCard,
  Github,
  Keyboard,
  LifeBuoy,
  LogOut,
  Mail,
  MessageSquare,
  Plus,
  PlusCircle,
  Settings,
  User,
  UserPlus,
  Users,
} from "lucide-react"
import useSWR, { useSWRConfig } from "swr"

import { siteConfig } from "@/config/site"
import { MainNav } from "@/components/main-nav"
import { ThemeToggle } from "@/components/theme-toggle"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-40 w-full border-b border-b-slate-200 bg-white dark:border-b-slate-700 dark:bg-slate-900">
      <div className="container flex h-16 items-center space-x-4 sm:justify-between sm:space-x-0">
        <MainNav items={siteConfig.mainNav} />
        <div className="flex flex-1 items-center justify-end space-x-4">
          <nav className="flex items-center space-x-1">
            <ThemeToggle />
            <Header />
          </nav>
        </div>
      </div>
    </header>
  )
}

function Header() {
  const router = useRouter()
  const {
    data: response,
    error,
    isLoading,
  } = useSWR("http://localhost:3000/me", fetcher, {
    revalidateIfStale: false,
    revalidateOnFocus: false,
    shouldRetryOnError: false,
  })
  const { mutate } = useSWRConfig()
  return (
    !(error || isLoading) && (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Avatar className="rounded-full aspect-square bg-gray-100 h-10 flex items-center justify-center border border-gray-200 cursor-pointer">
            {response?.data?.username
              ? response.data.username[0].toUpperCase()
              : "U"}
            <AvatarFallback>{}</AvatarFallback>
          </Avatar>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56" align="end">
          <DropdownMenuLabel>
            <div className="flex flex-col">
              <p className="text-base">
                {response?.data?.username
                  ? response.data.username
                  : "My Account"}
              </p>
              <p className="font-normal text-gray-500">
                {response?.data?.created_at
                  ? `Member since ${new Date(
                      response?.data?.created_at
                    ).toLocaleDateString("en-UK")}`
                  : "Member"}
              </p>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem onSelect={() => handleUserLogout(router, mutate)}>
              <LogOut className="mr-2 h-4 w-4" />
              <span>Logout</span>
            </DropdownMenuItem>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    )
  )
}

async function handleUserLogout(router, mutate) {
  await fetch("http://localhost:3000/logout", {
    method: "GET",
    credentials: "include",
  })
  await mutate(() => true, undefined, { revalidate: false })
  router.push("/login")
}
