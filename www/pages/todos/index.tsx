// TODO prevent flash when mutate optimistically
// sort the data server side so the client doesn't have to do it

import { FormEvent } from "react"
import { useRouter } from "next/router"
import { AnimatePresence, motion } from "framer-motion"
import useSWR, { useSWRConfig } from "swr"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { Textarea } from "@/components/ui/textarea"

const NEW_TODO_ANIMATION_DURATION = 0.5

const fetcher = (url: string) =>
  fetch(url, {
    method: "GET",
    credentials: "include",
  }).then(async (res) => {
    if (res.status === 401) {
      location.href = "/login"
    }
    if (!res.ok) {
      throw {
        status: res.status,
        message: res.statusText,
      }
    }
    return res.json()
  })

export default function TodosPage() {
  const {
    data: response,
    error,
    isLoading,
    mutate,
  } = useSWR("http://localhost:3000/todos", fetcher)

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
    <div className="min-w-full max-w-[100vw] min-h-[100vh] flex flex-col items-center px-6 py-16">
      <div className="w-full sm:w-[400px] flex flex-col gap-8">
        <Header />
        <Separator />
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
  )
}

async function addTodo(e: React.FormEvent<HTMLFormElement>): Promise<Response> {
  e.preventDefault()
  const title = e.currentTarget["todo-title"].value
  const description = e.currentTarget.description.value
  return await fetch("http://localhost:3000/todos", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ title, description }),
  })
}

async function removeTodo(id: number) {
  return await fetch(`http://localhost:3000/todos/${id}`, {
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

function Todo({ todo, isProcessing, handleRemoveTodo }) {
  const bgColor = isProcessing ? "#fafafa" : "#ffffff"
  return (
    <motion.div
      initial={{ height: 0, backgroundColor: bgColor, opacity: 0 }}
      animate={{ height: "auto", backgroundColor: bgColor, opacity: 1 }}
      transition={{ type: "spring", stiffness: 1000, damping: 40 }}
      exit={{ opacity: 0 }}
      style={{ overflow: "hidden" }}
      className={`rounded-lg border shadow-sm transition-opacity duration-1000`}
    >
      <div className={`p-4 flex flex-col`}>
        {todo.title && <h5 className="font-medium">{todo.title}</h5>}
        <p>
          {todo.updated_at && (
            <span className="text-gray-600 text-sm">
              {new Date(todo.updated_at).toLocaleDateString("vi-VN") + " – "}
            </span>
          )}
          {todo.description ? (
            <span className="text-gray-600 text-sm break-words">
              {todo.description}
            </span>
          ) : (
            <span className="text-gray-400 text-sm">No content</span>
          )}
        </p>
        <Button
          onClick={() => handleRemoveTodo(todo.id)}
          variant="link"
          className="w-fit mt-3 px-0 text-red-500"
          disabled={isProcessing}
        >
          Delete
        </Button>
      </div>
    </motion.div>
  )
}

function Header() {
  const router = useRouter()
  const {
    data: response,
    error,
    isLoading,
  } = useSWR("http://localhost:3000/me", fetcher)
  const { mutate } = useSWRConfig()
  return (
    <div className="flex flex-row gap-3 w-full items-center py-2 px-3 rounded-lg bg-gray-50 border">
      <div className="flex flex-col w-full">
        {!error &&
          (isLoading ? (
            <>
              <div className="animate-pulse h-4 mt-1 max-w-[6rem] bg-gray-200 rounded-md" />
              <div className="animate-pulse h-4 mt-2 max-w-[12rem] bg-gray-200 rounded-md" />
            </>
          ) : (
            <>
              <p className="font-medium">{response?.data?.username}</p>
              <p className="text-sm text-gray-400">{`Member since ${new Date(
                response?.data?.created_at
              ).toLocaleDateString("en-UK")}`}</p>
            </>
          ))}
      </div>
      <Button onClick={() => handleUserLogout(router, mutate)}>Logout</Button>
    </div>
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
