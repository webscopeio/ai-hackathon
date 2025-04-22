import { cn } from "@/lib/utils";
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
    <Card
      className={cn("w-[380px] overflow-hidden", className)}
      active={active}
      {...props}
    >
      <CardHeader>
        <AgentTitle active={active}>Analyzer</AgentTitle>
        <CardDescription>
          Analyzes the website and generates scenarios.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        <div className="h-[300px] overflow-hidden relative">
          <div className="absolute bottom-0 w-full flex flex-col gap-1">
            {messages.map((message, index) => (
              <div
                key={index}
                className={cn(
                  "mb-4 flex flex-col gap-1 items-start pb-4 last:mb-0 last:pb-0 transition-all duration-500 opacity-0",
                  "animate-fade-in blur-sm"
                )}
              >
                <p className="text-sm font-medium leading-none">
                  {message.title === "ANALYZER" ? "Agent" : "Toolcall"}
                </p>
                <p className="text-sm text-muted-foreground break-all">
                  {message.description.length > 150
                    ? `${message.description.slice(0, 150)}...`
                    : message.description}
                </p>
              </div>
            ))}
          </div>
          <div className="absolute top-0 w-full h-[52px] bg-gradient-to-b from-background to-transparent" />
        </div>
        <ShinyText
          className={cn(
            "justify-self-end",
            active ? "opacity-100" : "opacity-0"
          )}
          text="Generating..."
          speed={2}
        />
      </CardContent>
    </Card>
  );
}
