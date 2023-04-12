// TODO prevent flash when mutate optimistically
// sort the data server side so the client doesn't have to do it

import { FormEvent } from "react"
import { useRouter } from "next/router"
import AppConfig from "@/common/config"
import { fetcher } from "@/common/fetcher"
import { AnimatePresence } from "framer-motion"
import useSWR from "swr"

import Todo from "@/components/Todo"
import { Layout } from "@/components/layout"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { Textarea } from "@/components/ui/textarea"

const NEW_TODO_ANIMATION_DURATION = 0.5

export default function TodosPage() {
  const router = useRouter()
  const {
    data: response,
    error,
    isLoading,
    mutate,
  } = useSWR(
    `${AppConfig.API_URL}/todos`,
    (url: string) => fetcher(url, router, true),
    {
      revalidateIfStale: false,
      revalidateOnFocus: false,
    }
  )

  const handleAddTodo = async (e: FormEvent<HTMLFormElement>) => {
    const optimisticData = {
      ...response,
      data: [
        ...response.data,
        {
          id: Date.now(),
          title: e.currentTarget["todo-title"].value,
          description: e.currentTarget.description.value,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          isProcessing: true,
        },
      ],
    }
    mutate(
      async () => {
        try {
          const r = addTodo(e)
          await new Promise((resolve) =>
            setTimeout(resolve, NEW_TODO_ANIMATION_DURATION * 1000)
          ) // wait for the animation
          return await (await r).json()
        } catch (e) {
          console.log(e)
        }
      },
      {
        optimisticData: optimisticData,
        rollbackOnError: true,
        populateCache: (addedTodo, currentTodos) => {
          if (!addedTodo) return currentTodos
          return {
            ...currentTodos,
            data: [...currentTodos.data, addedTodo.data],
          }
        },
        revalidate: false,
      }
    )
  }

  const handleRemoveTodo = async (id: number) => {
    const optimisticData = {
      ...response,
      data: response.data.map((todo) => {
        if (todo.id === id) {
          return {
            ...todo,
            isProcessing: true,
          }
        }
        return todo
      }),
    }
    mutate(removeTodo(id), {
      optimisticData: optimisticData,
      rollbackOnError: true,
      populateCache: (_, currentTodos) => ({
        ...currentTodos,
        data: currentTodos.data.filter((todo) => todo.id !== id),
      }),
      revalidate: false,
    })
  }

  return (
    <Layout>
      <div className="min-w-full max-w-[100vw] min-h-[100vh] flex flex-col items-center px-6 py-16">
        <div className="w-full sm:w-[400px] flex flex-col gap-8">
          <form onSubmit={handleAddTodo} className="flex flex-col gap-4">
            <h4 className="text-md font-bold leading-tight tracking-tight md:text-lg">
              Add new todo
            </h4>
            <div className="flex flex-col gap-2">
              <Input type="text" id="todo-title" placeholder="Title" />
              <Textarea id="description" placeholder="Description" />
              <Button type="submit">Add todo</Button>
            </div>
          </form>
          <Separator />
          <div className="flex flex-col gap-4">
            <h3 className="text-xl font-bold leading-tight tracking-tight md:text-2xl">
              Your todos
            </h3>
            <Todos
              isLoading={isLoading}
              error={error}
              todos={response?.data}
              handleRemoveTodo={handleRemoveTodo}
            />
          </div>
        </div>
      </div>
    </Layout>
  )
}

async function addTodo(e: React.FormEvent<HTMLFormElement>): Promise<Response> {
  e.preventDefault()
  const title = e.currentTarget["todo-title"].value
  const description = e.currentTarget.description.value
  return await fetch(`${AppConfig.API_URL}/todos`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ title, description }),
  })
}

async function removeTodo(id: number) {
  return await fetch(`${AppConfig.API_URL}/todos/${id}`, {
    method: "DELETE",
    credentials: "include",
  })
}

function Todos({ todos, error, isLoading, handleRemoveTodo }) {
  if (isLoading) return <div>We`re fetching your todos...</div>
  if (error)
    return (
      <div className="text-red-500">
        Failed to fetch todos. {JSON.stringify(error)}
      </div>
    )
  if (!todos) return <p>You don`t have any todo! Create one!</p>
  if (!todos.length) return <p>You don`t have any todo! Create one!</p>
  todos.sort(
    (a, b) =>
      new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()
  )
  return (
    <div className="flex flex-col gap-4">
      {todos.map((todo) => {
        return (
          <AnimatePresence
            key={todo.id}
            initial={todo.isProcessing ? true : false}
          >
            <Todo
              isProcessing={todo.isProcessing}
              todo={todo}
              handleRemoveTodo={handleRemoveTodo}
            />
          </AnimatePresence>
        )
      })}
    </div>
  )
}
