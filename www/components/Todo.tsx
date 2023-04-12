import { motion } from "framer-motion"
import { Edit, Trash2 } from "lucide-react"

import { Button } from "@/components/ui/button"

export default function Todo({
  todo,
  isProcessing,
  handleRemoveTodo,
  className,
}: {
  todo: any
  isProcessing: boolean
  handleRemoveTodo: (id: number) => void
  className?: string
}) {
  return (
    <motion.div
      initial={{ height: 0 }}
      animate={{ height: "auto" }}
      transition={{ type: "spring", stiffness: 1000, damping: 40 }}
      exit={{ opacity: 0 }}
      style={{ overflow: "hidden" }}
      className={`rounded-xl border border-gray-200 dark:border-slate-800 shadow-sm transition-opacity duration-1000 bg-white dark:bg-slate-800${
        className ? " " + className : ""
      }`}
    >
      <div className={`p-6 flex flex-col`}>
        {todo.title && <h5 className="font-medium text-lg">{todo.title}</h5>}
        <p>
          {todo.updated_at && (
            <span className="text-gray-500">
              {new Date(todo.updated_at).toLocaleDateString("vi-VN") + " – "}
            </span>
          )}
          {todo.description ? (
            <span className="text-gray-500 break-words">
              {todo.description}
            </span>
          ) : (
            <span className="text-gray-500">No content</span>
          )}
        </p>
        <div className="flex gap-2 mt-3">
          <Button
            onClick={() => handleRemoveTodo(todo.id)}
            variant="subtle"
            className="w-fit h-fit p-2"
            disabled={isProcessing}
          >
            <Trash2 size={18} />
          </Button>
          <Button
            onClick={() => {}}
            variant="subtle"
            className="w-fit h-fit p-2"
            disabled={isProcessing}
          >
            <Edit size={18} />
          </Button>
        </div>
      </div>
    </motion.div>
  )
}
