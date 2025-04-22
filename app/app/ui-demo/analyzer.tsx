import { BellRing, Check } from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
} from "@/components/ui/card";
import { AgentTitle } from "@/components/agent";
import { Message } from "@/lib/types";
import ShinyText from "@/components/ShinyText";

const notifications = [
  {
    title: "Your call has been confirmed.",
    description: "1 hour ago",
  },
  {
    title: "You have a new message!",
    description: "1 hour ago",
  },
  {
    title: "Your subscription is expiring soon!",
    description: "2 hours ago",
  },
];

type AnalyzerProps = React.ComponentProps<typeof Card> & {
  active: boolean;
  messages: Message[];
};

export function Analyzer({
  className,
  active,
  messages,
  ...props
}: AnalyzerProps) {
  return (
    <Card className={cn("w-[380px]", className)} active={active} {...props}>
      <CardHeader>
        <AgentTitle active={active}>Analyzer</AgentTitle>
        <CardDescription>
          Analyzes the website and generates scenarios.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-4 h-[350px]">
        <div>
          {messages.slice(-3).map((message, index) => (
            <div
              key={index}
              className="mb-4 flex flex-row items-start pb-4 last:mb-0 last:pb-0"
            >
              <div className="space-y-1">
                <p className="text-sm font-medium leading-none">
                  {message.title === "ANALYZER" ? "Agent" : "Toolcall"}
                </p>
                <p className="text-sm text-muted-foreground">
                  {message.description}
                </p>
              </div>
            </div>
          ))}
        </div>
        {active && <ShinyText text="Generating..." speed={2} />}
      </CardContent>
    </Card>
  );
}
