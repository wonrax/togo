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
        <div className="w-full sm:w-[800px] flex flex-col gap-4">
          <p className="text-xl md:text-2xl font-bold">Users</p>
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
    <table className="border-collapse w-full dark:bg-slate-800 text-sm">
      <thead>
        {table.getHeaderGroups().map((headerGroup) => (
          <tr className="bg-gray-100" key={headerGroup.id}>
            {headerGroup.headers.map((header, index) => {
              let borderRadius = "border-y"
              if (index === 0) borderRadius = " rounded-l border-l border-y"
              else if (index === headerGroup.headers.length - 1)
                borderRadius = " rounded-r border-r border-y"
              return (
                <th
                  key={header.id}
                  className={
                    "w-1/2 font-semibold p-4 text-slate-900 dark:text-slate-200 text-left border-separate" +
                    borderRadius
                  }
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </th>
              )
            })}
          </tr>
        ))}
      </thead>
      <tbody>
        {table.getRowModel().rows.map((row, index) => {
          const border =
            index === 0 ? "" : "border-t border-slate-300 dark:border-slate-700"
          return (
            <tr className={border} key={row.id}>
              {row.getVisibleCells().map((cell) => (
                <td
                  key={cell.id}
                  className="p-4 text-slate-500 dark:text-slate-400"
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          )
        })}
      </tbody>
    </table>
  )
}
