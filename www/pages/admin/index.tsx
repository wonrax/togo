import { useMemo, useReducer, useState } from "react"
import { useRouter } from "next/router"
import AppConfig from "@/common/config"
import { fetcher } from "@/common/fetcher"
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table"
import useSWR from "swr"

import { Layout } from "@/components/layout"

export default function AdminPage() {
  const router = useRouter()
  const {
    data: response,
    error,
    isLoading,
  } = useSWR(
    `${AppConfig.API_URL}/admin/users`,
    (url: string) => fetcher(url, router, true),
    {
      revalidateIfStale: false,
      revalidateOnFocus: false,
    }
  )

  return (
    <Layout>
      <div className="min-w-full max-w-[100vw] min-h-[100vh] flex flex-col items-center px-6 py-16">
        <div className="w-full sm:w-[800px] flex flex-col gap-8">
          {!error && !isLoading && <Data data={response?.data} />}
        </div>
      </div>
    </Layout>
  )
}
type User = {
  username: string
  created_at: string
}

const columnHelper = createColumnHelper<User>()

const columns = [
  columnHelper.accessor("username", {
    header: "Username",
    cell: (info) => info.getValue(),
  }),
  columnHelper.accessor("created_at", {
    header: () => "Created at",
    cell: (info) => new Date(info.getValue()).toLocaleDateString("en-UK"),
  }),
]

function Data({ data }) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  return (
    <table className="border-collapse w-full  bg-white dark:bg-slate-800 text-sm">
      <thead>
        {table.getHeaderGroups().map((headerGroup) => (
          <tr key={headerGroup.id}>
            {headerGroup.headers.map((header) => (
              <th
                key={header.id}
                className="w-1/2 border-b  border-slate-300 dark:border-slate-600 font-semibold p-4 text-slate-900 dark:text-slate-200 text-left"
              >
                {header.isPlaceholder
                  ? null
                  : flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
              </th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody>
        {table.getRowModel().rows.map((row) => (
          <tr key={row.id}>
            {row.getVisibleCells().map((cell) => (
              <td
                key={cell.id}
                className="border-t border-slate-300 dark:border-slate-700 p-4 text-slate-500 dark:text-slate-400"
              >
                {flexRender(cell.column.columnDef.cell, cell.getContext())}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  )
}
