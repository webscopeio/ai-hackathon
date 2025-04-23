import { cn } from "@/lib/utils";

interface TestFilesProps {
  count: number;
  names: string[];
  scenario: string;
  className?: string;
}

export function TestFiles({
  count,
  names,
  scenario,
  className,
}: TestFilesProps) {
  if (count === 0) {
    return null;
  }

  return (
    <div className={cn("flex flex-row gap-16", className)}>
      <div className="flex flex-col">
        <h2 className="text-lg font-bold pb-1">Scenario {names.length + 1}:</h2>
        <p className="text-sm whitespace-pre-wrap">{scenario}</p>
      </div>
      <div className="flex flex-col">
        <h2 className="text-lg font-bold pb-1">Generated test files:</h2>
        {Array.from({ length: count }).map((_, index) => {
          const fileName = names[index];
          const isKnown = !!fileName;

          return (
            <div key={index} className="flex items-center gap-2 pl-1.5 pb-1">
              <div
                className={cn(
                  "h-2 w-2 rounded-full",
                  isKnown
                    ? "bg-blue-500 shadow-[0_0_12px_4px_rgba(59,130,246,0.7)]"
                    : "bg-gray-400"
                )}
              />
              <span className="font-mono">
                {isKnown ? fileName : "Pending..."}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
