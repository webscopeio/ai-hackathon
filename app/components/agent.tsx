import * as React from "react";
import { cn } from "@/lib/utils";
import { CardTitle } from "@/components/ui/card";

type AgentTitleProps = React.ComponentProps<typeof CardTitle> & {
  active?: boolean;
};

export function AgentTitle({
  className,
  active,
  children,
  ...props
}: AgentTitleProps) {
  return (
    <div className="flex items-center gap-3">
      <span
        className={cn(
          "block h-2 w-2 rounded-full transition-all duration-200",
          active
            ? ["bg-blue-500", "shadow-[0_0_12px_4px_rgba(59,130,246,0.7)]"]
            : "bg-gray-400"
        )}
      />
      <CardTitle className={className} {...props}>
        {children}
      </CardTitle>
    </div>
  );
}
