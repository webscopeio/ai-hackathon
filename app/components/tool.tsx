import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Message } from "@/lib/types";
import ShinyText from "@/components/ShinyText";

type ToolProps = React.ComponentProps<typeof Card> & {
  active: boolean;
  messages: Message[];
  title: string;
  description: string;
};

export function Tool({
  className,
  active,
  messages,
  title,
  description,
  ...props
}: ToolProps) {
  return (
    <Card
      className={cn("w-[300px] overflow-hidden", className)}
      active={active}
      {...props}
    >
      <CardHeader className="pb-0">
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-1">
        <div className="h-[150px] overflow-hidden relative">
          <div className="absolute bottom-0 w-full flex flex-col gap-0.5">
            {messages.map((message, index) => (
              <div
                key={index}
                className={cn(
                  "mb-1 flex flex-col items-start pb-1 last:mb-0 last:pb-0 transition-all duration-500 opacity-0",
                  "animate-fade-in blur-sm"
                )}
              >
                <p className="text-sm font-medium leading-none capitalize">
                  {message.title}
                </p>
                <p className="text-sm text-muted-foreground break-all">
                  {message.description.length > 100
                    ? `${message.description.slice(0, 100)}...`
                    : message.description}
                </p>
              </div>
            ))}
          </div>
          <div
            className={cn(
              "absolute top-0 w-full h-[52px] bg-gradient-to-b from-background to-transparent",
              active && "from-blue-50/30"
            )}
          />
        </div>
      </CardContent>
    </Card>
  );
}
