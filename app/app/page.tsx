"use client";

import { Skeleton } from "@/components/ui/skeleton";
import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";

const StatusDemo: React.FC = () => {
  const { status, data, error } = useQuery(api.getStatus);

  return (
    <div className="flex-1">
      <h2 className="mb-6 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors">
        Status
      </h2>
      {status === "pending" ? (
        <div className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200">
          <div className="p-6">
            <Skeleton className="h-6 w-3/4 mb-4" />
            <Skeleton className="h-4 w-full mb-2" />
          </div>
        </div>
      ) : status === "error" ? (
        <div className="bg-red-50 border-l-4 border-red-500 p-6 rounded shadow">
          <div className="ml-3">
            <h3 className="text-xl font-medium text-red-800 mb-2">
              Error loading status
            </h3>
            <p className="leading-7 text-red-700">
              {error?.message ?? "Unknown error"}
            </p>
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition-shadow duration-300">
          <div className="p-6">
            <h3 className="text-xl font-semibold tracking-tight mb-2">
              Status
            </h3>
            <p className="leading-7 text-gray-600">{JSON.stringify(data)}</p>
          </div>
        </div>
      )}
    </div>
  );
};

export default function Home() {
  return <StatusDemo />;
}
