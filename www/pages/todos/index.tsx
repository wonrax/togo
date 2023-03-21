// TODO prevent flash when mutate optimistically
// sort the data server side so the client doesn't have to do it

import useSWR, { useSWRConfig } from "swr"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import { Textarea } from "@/components/ui/textarea"

const fetcher = (url: string) =>
  fetch(url, {
    method: "GET",
    credentials: "include",
  }).then(async (res) => {
    await new Promise((r) => setTimeout(r, 1000))
    return res.json()
  })

export default function TodosPage() {
  const {
    data: response,
    error,
    isLoading,
  } = useSWR("http://localhost:3000/todos", fetcher, {
    keepPreviousData: true,
  })

  const { mutate } = useSWRConfig()

  return (
    <div className="min-w-full max-w-[100vw] min-h-[100vh] flex flex-col items-center px-6 py-16">
      <div className="w-full sm:w-[400px] flex flex-col gap-8">
        <form
          onSubmit={async (e) => {
            const newData = {
              ...response,
              data: [
                ...response.data,
                {
                  id: response.data.length + 2,
                  title: e.currentTarget["todo-title"].value,
                  description: e.currentTarget.description.value,
                  created_at: new Date().toISOString(),
                  updated_at: new Date().toISOString(),
                },
              ],
            }
            mutate(
              "http://localhost:3000/todos",
              async () => {
                await handleTodoSubmit(e)
                return fetcher
              },
              {
                optimisticData: newData,
              }
            )
          }}
          className=" flex flex-col gap-5"
        >
          <h3 className="text-xl font-extrabold leading-tight tracking-tight md:text-2xl">
            Add new todo
          </h3>
          <div className="flex flex-col gap-4">
            <div className="grid w-full items-center gap-2">
              <Label htmlFor="todo-title">Title</Label>
              <Input type="text" id="todo-title" placeholder="Title" />
            </div>
            <div className="grid w-full items-center gap-2">
              <Label htmlFor="description">Description</Label>
              <Textarea id="description" placeholder="Description" />
            </div>
          </div>
          <Button type="submit">Add todo</Button>
        </form>
        <div className="flex flex-col gap-4">
          <h3 className="text-xl font-extrabold leading-tight tracking-tight md:text-2xl">
            Your todos
          </h3>
          <Todos isLoading={isLoading} error={error} todos={response?.data} />
        </div>
      </div>
    </div>
  )
}

async function handleTodoSubmit(
  e: React.FormEvent<HTMLFormElement>
): Promise<Response> {
  e.preventDefault()
  const title = e.currentTarget["todo-title"].value
  const description = e.currentTarget.description.value
  await new Promise((r) => setTimeout(r, 1000))
  return await fetch("http://localhost:3000/todos", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ title, description }),
  })
}

function Todos({ todos, error, isLoading }) {
  if (isLoading) return <div>We're fetching your todos...</div>
  if (error)
    return (
      <div className="text-red-500">
        Failed to fetch todos. {JSON.stringify(error)}
      </div>
    )
  if (!todos) return null
  todos.sort(
    (a, b) =>
      new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()
  )
  return (
    <div className="flex flex-col gap-4">
      {todos.map((todo) => (
        <Todo key={todo.id} todo={todo} />
      ))}
    </div>
  )
}

function Todo({ todo }) {
  return (
    <div className="p-4 rounded-lg border shadow-sm flex flex-col gap-2">
      {todo.title && <h5 className="font-bold">{todo.title}</h5>}
      <p>
        {todo.updated_at && (
          <span className="text-sm text-gray-300">
            {new Date(todo.updated_at).toLocaleDateString("vi-VN") + " "}
          </span>
        )}
        {todo.description && (
          <span className="break-words break-all">{todo.description}</span>
        )}
      </p>
    </div>
  )
}
